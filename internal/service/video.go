package service

import (
	"context"
	"database/sql"

	pb_vid "github.com/msqtt/sevencowcloud-shortvideo-service/api/pb/v1/video"
	"github.com/msqtt/sevencowcloud-shortvideo-service/internal/pkg/config"
	db "github.com/msqtt/sevencowcloud-shortvideo-service/internal/repo/sqlc"
	"github.com/msqtt/sevencowcloud-shortvideo-service/internal/token"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type VideoServer struct {
	pb_vid.UnimplementedVideoServiceServer
	config config.Config
	token  token.TokenMaker
	store  db.Store
}

// ListVideoClass implements pb_vid.VideoServiceServer.
func (vs *VideoServer) ListVideoTag(ctx context.Context, req *pb_vid.ListVideoTagRequest) (*pb_vid.ListVideoTagResponse, error) {
	ta, err := vs.store.GetAllTags(ctx)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, status.Errorf(codes.Internal, "failed to find video calss")
		}
	}

	tas := make([]*pb_vid.Tag, len(ta))
	for i, e := range ta {
		tas[i] = &pb_vid.Tag{
			TagId:          e.ID,
			Name:        e.Name,
			Description: e.Description,
		}
	}
	return &pb_vid.ListVideoTagResponse{
		Tags: tas,
	}, nil
}

var _ pb_vid.VideoServiceServer = (*VideoServer)(nil)

func NewVideoServer(conf config.Config, token token.TokenMaker, store db.Store) *VideoServer {
	return &VideoServer{
		config: conf,
		token:  token,
		store:  store,
	}
}
