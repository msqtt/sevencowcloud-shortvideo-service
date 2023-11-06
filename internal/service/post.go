package service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	pb_pst "github.com/msqtt/sevencowcloud-shortvideo-service/api/pb/v1/post"
	"github.com/msqtt/sevencowcloud-shortvideo-service/internal/oss"
	"github.com/msqtt/sevencowcloud-shortvideo-service/internal/pkg/config"
	db "github.com/msqtt/sevencowcloud-shortvideo-service/internal/repo/sqlc"
	"github.com/msqtt/sevencowcloud-shortvideo-service/internal/token"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PostServer struct {
	pb_pst.UnimplementedPostServiceServer
	config config.Config
	token  token.TokenMaker
	store  db.Store
	kodo   *oss.Kodo
}

// TestGetRecommendPost implements pb_pst.PostServiceServer.
func (ps *PostServer) TestGetRecommendPost(ctx context.Context,
	req *pb_pst.TestGetRecommendPostRequest) (
	*pb_pst.TestGetRecommendPostResponse, error) {
	pIndex := req.GetPageIndex()
	pSize := req.GetPageSize()

	if pIndex == 0 {
		pIndex = 1
	}
	if pSize == 0 {
		pSize = 10
	}

	param := db.TestGetAllParams{
		Offset: (pIndex - 1) * pSize,
		Limit:  pSize,
	}

	m := extractMetadata(ctx)
	var payload *token.Payload
	if len(m.Token) > 0 {
		payload, _ = ps.token.ValidToken(m.Token)
	}

	tgar, err := ps.store.TestGetAll(ctx, param)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get recommend posts")
	}

	ret := make([]*pb_pst.PostItem, len(tgar))
	for i, row := range tgar {
		item, err := GetPostItem(ctx, ps.store, ps.config, payload, row.UserID, row.ID, row.VideoID,
			row.Title, row.Description, row.UpdatedAt, row.CreatedAt)
		if err != nil {
			log.Println(err)
		}
		ret[i] = item
	}
	var total int32 = 0
	if len(tgar) > 0 {
		total = int32(tgar[0].TotalSize)
	}
	return &pb_pst.TestGetRecommendPostResponse{
		Total:     total,
		PageSize:  pSize,
		PagePos:   pIndex,
		PostItems: ret,
	}, nil
}

func GetPostItem(ctx context.Context, q db.Querier, conf config.Config, payload *token.Payload,
	uid, pid, vid int64, title, desc string, updatedat, createdat time.Time) (*pb_pst.PostItem, error) {
	post := db.Post{
		ID:          pid,
		Title:       title,
		Description: desc,
		UserID:      uid,
		UpdatedAt:   updatedat,
		CreatedAt:   createdat,
	}
	// get user and profile
	v, err := q.GetVideoByID(ctx, vid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to find video")
	}

	tagsRow, err := q.GetTagsByPostID(ctx, pid)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to find video tags")
	}

	tags := make([]db.Tag, len(tagsRow))
	for i, row := range tagsRow {
		tags[i] = db.Tag{
			ID:          row.ID,
			Name:        row.Name,
			Description: row.Description,
		}
	}

	return &pb_pst.PostItem{
		Post:         db2pbPost(conf, post, v, tags),
		IsLiked:      false,
		IsCollected:  false,
		IsShared:     false,
		LikedNum:     0,
		CollectedNum: 0,
		SharedNum:    0,
	}, nil
}

// UploadPost implements pb_pst.PostServiceServer.
func (ps *PostServer) UploadPost(ctx context.Context, req *pb_pst.UploadPostRequest) (
	*pb_pst.UploadPostResponse, error) {
	title := req.GetTitle()
	desc := req.GetDescription()
	uid := req.GetUserId()
	tagids := req.GetTagIds()
	vid := req.GetVideoId()

	payload := ctx.Value("payload").(*token.Payload)
	if payload.UserID != uid {
		return nil, status.Errorf(codes.Unauthenticated, "user Unauthenticated")
	}

	_, err5 := ps.store.GetUserByID(ctx, uid)
	if err5 != nil {
		if err5 == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "not found user")
		}
		return nil, status.Errorf(codes.Internal, "failed to find user")
	}

	// oss: move video and cover
	video, err3 := ps.store.GetVideoByID(ctx, vid)
	if err3 != nil {
		if err3 == sql.ErrNoRows {
			return nil, status.Errorf(codes.InvalidArgument, "not found video")
		}
		return nil, status.Errorf(codes.Internal, "failed to find video")
	}
	bk := ps.config.KodoBucket

	cstrs := strings.Split(video.CoverLink, "/")
	coverName := cstrs[len(cstrs)-1]
	ckey := fmt.Sprintf("img/cover/%s", coverName)
	if err := ps.kodo.MoveFile(bk, video.CoverLink, bk, ckey); err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to move file on oss")
	}

	vstrs := strings.Split(video.SrcLink, "/")
	videoName := vstrs[len(vstrs)-1]
	tmpstrs := strings.Split(videoName, ".")
	videoName = fmt.Sprintf("%s-trans.%s", tmpstrs[0], tmpstrs[1])
	vkey := fmt.Sprintf("video/%s", videoName)
	if err := ps.kodo.MoveFile(bk, video.SrcLink, bk, vkey); err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to move file on oss")
	}

	// db: update video
	vParams := db.UpdateVideoLinkParams{
		CoverLink:   ckey,
		SrcLink:     vkey,
		ID: vid,
	}
	err := ps.store.UpdateVideoLink(ctx, vParams)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add video")
	}

	// db: insert post
	pparams := db.AddPostParams{
		Title:       title,
		Description: desc,
		UserID:      uid,
		VideoID:     vid,
	}
	r, err := ps.store.AddPost(ctx, pparams)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to add post")
	}
	pid, err2 := r.LastInsertId()
	if err2 != nil {
		return nil, status.Errorf(codes.Internal, "failed to get last insert post id")
	}

	// post tag
	alltags, err := ps.store.GetAllTags(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get all tags")
	}
	dbTag := make([]db.Tag, len(tagids))
	for i, id := range tagids {
		args := db.AddPostTagParams{
			PostID: pid,
			TagID:  id,
		}
		err7 := ps.store.AddPostTag(ctx, args)
		if err7 != nil {
			log.Println(err7)
		}
		dbTag[i] =  alltags[id-1]
	}

	postitem := &pb_pst.PostItem{
		Post: db2pbPost(ps.config, db.Post{
			ID:          pid,
			Title:       title,
			Description: desc,
			UserID:      uid,
			UpdatedAt:   time.Now(),
			CreatedAt:   time.Now(),
		}, db.Video{
			CoverLink:   ckey,
			SrcLink:     vkey,
		}, dbTag,
		),
		IsLiked:      false,
		IsCollected:  false,
		IsShared:     false,
		LikedNum:     0,
		CollectedNum: 0,
		SharedNum:    0,
	}

	return &pb_pst.UploadPostResponse{PostItem: postitem}, nil
}

var _ pb_pst.PostServiceServer = (*PostServer)(nil)

func NewPostServer(conf config.Config, token token.TokenMaker, store db.Store, kodo *oss.Kodo) *PostServer {
	return &PostServer{
		config: conf,
		token:  token,
		store:  store,
		kodo: kodo,
	}
}
