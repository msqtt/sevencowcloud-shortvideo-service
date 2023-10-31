package service

import (
	"context"
	"database/sql"
	"log"
	"path/filepath"
	"time"

	pb_prf "github.com/msqtt/sevencowcloud-shortvideo-service/api/pb/v1/profile"
	"github.com/msqtt/sevencowcloud-shortvideo-service/internal/pkg/config"
	db "github.com/msqtt/sevencowcloud-shortvideo-service/internal/repo/sqlc"
	"github.com/msqtt/sevencowcloud-shortvideo-service/internal/token"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProfileServer struct {
	pb_prf.UnimplementedProfileServiceServer
	config config.Config
	store  db.Store
}

// GetProfile implements pb_prf.ProfileServiceServer.
func (ps *ProfileServer) GetProfile(ctx context.Context, req *pb_prf.GetProfileRequest) (
	*pb_prf.GetProfileResponse, error) {
	id := req.GetUserId()
	u, err := ps.store.GetUserByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "no such user")
		}
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to find user")
	}
	
	p, err := ps.store.GetProfile(ctx, u.ProfileID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to find profile")
	}
	p2 := db2pbProfile(&p)
	p2.AvatarLink = filepath.Join(ps.config.KodoLink, p2.AvatarLink)
	gpr := &pb_prf.GetProfileResponse{Profile: p2}
	return gpr, nil
}

// UpdateProfile implements pb_prf.ProfileServiceServer.
func (ps *ProfileServer) UpdateProfile(ctx context.Context, req *pb_prf.UpdateProfileRequest) (
	*pb_prf.UpdateProfileResponse, error) {
	id := req.GetUserId()
	payl := ctx.Value("payload")
	payload := payl.(*token.Payload)

	if id != payload.UserID {
		return nil, status.Errorf(codes.PermissionDenied, "cannot update other user's profiles")
	}

	u, err := ps.store.GetUserByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "no such user")
		}
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to find user")
	}

	params := db.UpdateProfileParams{
		RealName: sql.NullString{String: req.GetRealName(), Valid: true},
		Mood: sql.NullString{String: req.GetMood(), Valid: true},
		Gender: db.NullProfilesGender{ProfilesGender: db.ProfilesGender(req.GetGender().String()), Valid: true},
		BirthDate: sql.NullTime{Time: time.Unix(req.GetBirthDate(), 0), Valid: true},
		Introduction: sql.NullString{String: req.GetIntroduction(), Valid: true},
		ID: u.ID,
	}

	err2 := ps.store.UpdateProfile(ctx, params)
	if err2 != nil {
		log.Println(err2)
		return nil, status.Errorf(codes.Internal, "failed to update profile")
	}

	p2, err := ps.store.GetProfile(ctx, u.ProfileID)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to find updated profile")		
	}

	p3 := db2pbProfile(&p2)
	p3.AvatarLink = filepath.Join(ps.config.KodoLink, p3.AvatarLink)
	res := &pb_prf.UpdateProfileResponse{
		Profile: p3,
	}
	return res, nil	
}

var _ pb_prf.ProfileServiceServer = (*ProfileServer)(nil)

// NewProfileServer creates a profile server then return it.
func NewProfileServer(store db.Store, conf config.Config) *ProfileServer {
	return &ProfileServer{
		config: conf,
		store:  store,
	}
}
