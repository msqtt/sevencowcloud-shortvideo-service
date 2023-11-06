package service

import (
	"context"
	"database/sql"
	"log"

	pb_usr "github.com/msqtt/sevencowcloud-shortvideo-service/api/pb/v1/user"
	"github.com/msqtt/sevencowcloud-shortvideo-service/internal/pkg/config"
	db "github.com/msqtt/sevencowcloud-shortvideo-service/internal/repo/sqlc"
	"github.com/msqtt/sevencowcloud-shortvideo-service/internal/token"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServer struct {
	pb_usr.UnimplementedUserServiceServer
	config config.Config
	store  db.Store
	token  token.TokenMaker
}

// GetUserProfile implements pb_prf.ProfileServiceServer.
func (us *UserServer) GetUserProfile(ctx context.Context, req *pb_usr.GetUserProfileRequest) (
	*pb_usr.GetUserProfileResponse, error) {
	id := req.GetUserId()
	u, err := us.store.GetUserByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "no such user")
		}
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to find user")
	}

	p, err := us.store.GetProfileByID(ctx, u.ProfileID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to find profile")
	}

	ing_cnt, err2 := us.store.CountFollowing(ctx, id)
	if err2 != nil {
		return nil, status.Errorf(codes.Internal, "failed to count following list")
	}

	ed_cnt, err3 := us.store.CountFollowed(ctx, id)
	if err3 != nil {
		return nil, status.Errorf(codes.Internal, "failed to count followed list")
	}

	m := extractMetadata(ctx)
	isfollowed := false
	if len(m.Token) != 0 {
		p, _ := us.token.ValidToken(m.Token)
		_, isfollowed, _ = CheckIsFollow(us.store, ctx, p.UserID, id)
	}

	gpr := &pb_usr.GetUserProfileResponse{
		UserItem: &pb_usr.UserItem{
			User:       db2pbUser(us.config, u, p),
			IsFollowed: isfollowed,
			FollowingNum: int32(ing_cnt),
			FollowedNum: int32(ed_cnt),
		},
	}
	return gpr, nil
}

var _ pb_usr.UserServiceServer = (*UserServer)(nil)

func NewUserServer(conf config.Config, token token.TokenMaker, store db.Store) *UserServer {
	return &UserServer{
		config: conf,
		token: token,
		store: store,
	}
}
