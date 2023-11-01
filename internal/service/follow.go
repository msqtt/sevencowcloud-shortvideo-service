package service

import (
	"context"
	"database/sql"
	"log"
	"path/filepath"
	"time"

	pb_fl "github.com/msqtt/sevencowcloud-shortvideo-service/api/pb/v1/follow"
	pb_usr "github.com/msqtt/sevencowcloud-shortvideo-service/api/pb/v1/user"
	"github.com/msqtt/sevencowcloud-shortvideo-service/internal/pkg/config"
	db "github.com/msqtt/sevencowcloud-shortvideo-service/internal/repo/sqlc"
	"github.com/msqtt/sevencowcloud-shortvideo-service/internal/token"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type FollowServer struct {
	pb_fl.UnimplementedFollowServiceServer
	config config.Config
	store  db.Store
	token  token.TokenMaker
}

// FollowedList implements pb_fl.FollowServiceServer.
func (lf *FollowServer) FollowedList(ctx context.Context, req *pb_fl.FollowedListRequest) (
	*pb_fl.FollowedListResponse, error) {
	uid := req.GetUserId()
	gflr, err := lf.store.GetFollowedList(ctx, uid)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Println(err)
			return nil, status.Errorf(codes.PermissionDenied, "failed to find following list")
		}
	}

	// check if login
	m, _ := metadata.FromIncomingContext(ctx)
	token := m["authorization"]
	var loginUid int64 = -1
	if len(token) > 0 {
		if payload, err := lf.token.ValidToken(token[0]); err == nil {
			loginUid = payload.UserID
		}
	}

	// reduce data
	uis := make([]*pb_usr.UserItem, len(gflr))
	for i, row := range gflr {
		isFollow := false
		if loginUid != -1 {
			_, isFollow, err = CheckIsFollow(lf.store, ctx, loginUid, row.ID)
			if err != nil {
				log.Println(err)
				return nil, status.Errorf(codes.Internal, "failed to check if follow")
			}
		}
		ui := &pb_usr.UserItem{
			IsFollowed: isFollow,
			User: db2pbUser(db.User{ID: row.ID, Nickname: row.Nickname, Email: row.Email},
				&db.Profile{RealName: row.RealName, Mood: row.Mood, Gender: row.Gender,
					BirthDate: row.BirthDate, Introduction: row.Introduction,
					AvatarLink: sql.NullString{String: filepath.Join(lf.config.KodoLink,
						row.AvatarLink.String),
						Valid: true}},
			),
		}
		uis[i] = ui
	}
	ret := &pb_fl.FollowedListResponse{
		Users: uis,
	}
	return ret, nil
}

// FollowingList implements pb_fl.FollowServiceServer.
func (lf *FollowServer) FollowingList(ctx context.Context, req *pb_fl.FollowingListRequest) (
	*pb_fl.FollowingListResponse, error) {
	uid := req.GetUserId()
	gflr, err := lf.store.GetFollowingList(ctx, uid)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Println(err)
			return nil, status.Errorf(codes.PermissionDenied, "failed to find following list")
		}
	}

	// check if login
	m, _ := metadata.FromIncomingContext(ctx)
	token := m["authorization"]
	var loginUid int64 = -1
	if len(token) > 0 {
		if payload, err := lf.token.ValidToken(token[0]); err == nil {
			loginUid = payload.UserID
		}
		log.Println(err)
	}

	// reduce data
	uis := make([]*pb_usr.UserItem, len(gflr))
	for i, row := range gflr {
		isFollow := false
		if loginUid != -1 {
			_, isFollow, err = CheckIsFollow(lf.store, ctx, loginUid, row.ID)
			if err != nil {
				log.Println(err)
				return nil, status.Errorf(codes.Internal, "failed to check if follow")
			}
		}
		ui := &pb_usr.UserItem{
			IsFollowed: isFollow,
			User: db2pbUser(db.User{ID: row.ID, Nickname: row.Nickname, Email: row.Email},
				&db.Profile{RealName: row.RealName, Mood: row.Mood, Gender: row.Gender,
					BirthDate: row.BirthDate, Introduction: row.Introduction,
					AvatarLink: sql.NullString{String: filepath.Join(lf.config.KodoLink,
						row.AvatarLink.String),
						Valid: true}},
			),
		}
		uis[i] = ui
	}
	ret := &pb_fl.FollowingListResponse{
		Users: uis,
	}
	return ret, nil
}

// CheckIsFollow checks whether uid1 followed uid2.
func CheckIsFollow(q db.Querier, ctx context.Context, uid1, uid2 int64) (
	*db.Follow, bool, error) {
	gfp := db.GetFollowParams{FollowingUserID: uid1, FollowedUserID: uid2}
	f, err := q.GetFollow(ctx, gfp)
	if err != nil {
		// do not followed
		if err == sql.ErrNoRows {
			return nil, false, nil
		}
		log.Println(err)
		return nil, false, status.Errorf(codes.Internal, "failed to find follow")
	}
	// already followed
	return &f, true, nil
}

// CheckFollow implements pb_fl.FollowServiceServer.
func (fl *FollowServer) CheckFollow(ctx context.Context, req *pb_fl.CheckFollowRequest) (
	*pb_fl.CheckFollowResponse, error) {
	followingId := req.GetFollowingId()
	followedId := req.GetFollowedId()
	f, isFollow, err := CheckIsFollow(fl.store, ctx, followingId, followedId)
	if err != nil {
		return nil, err
	}
	var time int64 = -1
	if isFollow {
		time = f.FollowedAt.Unix()
	}
	return &pb_fl.CheckFollowResponse{IsFollowing: isFollow, FollowedAt: time}, nil
}

// FollowUser implements pb_fl.FollowServiceServer.
func (fl *FollowServer) FollowUser(ctx context.Context, req *pb_fl.FollowUserRequest) (
	*pb_fl.FollowUserResponse, error) {
	payload := ctx.Value("payload").(*token.Payload)
	followingId := payload.UserID
	if followingId != req.GetFollowingId() {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	followedId := req.GetFollowedId()
	_, isfollow, err := CheckIsFollow(fl.store, ctx, followingId, followedId)
	if err != nil {
		return nil, err
	}
	if isfollow {
		return nil, status.Errorf(codes.AlreadyExists, "already been followed")
	}
	params := db.AddFollowParams{
		FollowingUserID: followingId,
		FollowedUserID:  followedId,
	}

	_, err = fl.store.AddFollow(ctx, params)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to follow")
	}
	return &pb_fl.FollowUserResponse{
		FollowedAt: time.Now().Unix(),
	}, nil
}

// UnFollowUser implements pb_fl.FollowServiceServer.
func (fl *FollowServer) UnFollowUser(ctx context.Context, req *pb_fl.UnFollowUserRequest) (
	*pb_fl.UnFollowUserResponse, error) {
	payload := ctx.Value("payload").(*token.Payload)
	followingId := payload.UserID
	if followingId != req.GetFollowingId() {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}
	followedId := req.GetFollowedId()
	_, isfollow, err := CheckIsFollow(fl.store, ctx, followingId, followedId)
	if err != nil {
		return nil, err
	}
	if !isfollow {
		return nil, status.Errorf(codes.AlreadyExists, "have not been followed")
	}
	param := db.DeleteFollowParams{
		FollowingUserID: followingId,
		FollowedUserID:  followedId,
	}
	err2 := fl.store.DeleteFollow(ctx, param)
	if err2 != nil {
		log.Println(err2)
		return nil, status.Errorf(codes.Internal, "failed to unfollow user")
	}
	return &pb_fl.UnFollowUserResponse{Now: time.Now().Unix()}, nil
}

var _ pb_fl.FollowServiceServer = (*FollowServer)(nil)

func NewFollowServer(conf config.Config, token token.TokenMaker, store db.Store) *FollowServer {
	return &FollowServer{
		config: conf,
		store:  store,
		token:  token,
	}
}
