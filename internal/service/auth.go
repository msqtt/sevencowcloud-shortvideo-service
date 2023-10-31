package service

import (
	"context"
	"database/sql"
	"log"
	"path/filepath"
	"time"

	pb_usr "github.com/msqtt/sevencowcloud-shortvideo-service/api/pb/v1/user"
	mail "github.com/msqtt/sevencowcloud-shortvideo-service/internal/email"
	"github.com/msqtt/sevencowcloud-shortvideo-service/internal/pkg/config"
	"github.com/msqtt/sevencowcloud-shortvideo-service/internal/pkg/sample"
	db "github.com/msqtt/sevencowcloud-shortvideo-service/internal/repo/sqlc"
	"github.com/msqtt/sevencowcloud-shortvideo-service/internal/token"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	pb_usr.UnimplementedAuthServiceServer
	config config.Config
	store  db.Store
	token  token.TokenMaker
	mail   *mail.MailSender
}

// SendActivateEmail implements pb_usr.AuthServiceServer.
func (as *AuthServer) SendActivateEmail(ctx context.Context, req *pb_usr.SendActivateEmailRequest) (
	*pb_usr.SendActivateEmailResponse, error) {
	id := req.GetUserId()
	email := req.GetEmail()
	
	u, err := as.store.GetUserByEmailNotActivated(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "no such user needs to be activated")
		}
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to find user")
	}

	if id != u.ID {
		return nil, status.Errorf(codes.InvalidArgument, "id and email do not match")
	}

	// check activate left times
	cnt, err2 := as.store.TodayActiationCount(ctx, id)
	if err2 != nil {
		log.Println(err2)
		return nil, status.Errorf(codes.Internal, "failed to mount today activate times")
	}
	left := int32(as.config.ActivateTimes) - int32(cnt)
	if left <= 0 {
		return nil, status.Errorf(codes.ResourceExhausted, "the user has been reached %d times"+
			" limit of sending activate emails", cnt)
	}
	left --

	// email server
	m := extractMetadata(ctx)
	code := sample.RandomStr(8)

	// add activation
	args := db.AddActivationParams{
		UserID:       u.ID,
		ActivateCode: code,
		ExpiredAt:    time.Now().Add(10 * time.Minute),
	}
	_, err4 := as.store.AddActivation(ctx, args)
	if err4 != nil {
		log.Println(err4)
		return nil, status.Errorf(codes.Internal, "failed to create activation")
	}

	// open a goroutine to send an email
	go func() {
		err4 = as.mail.SendActivateEmail(email, u.Nickname, code, m.ClientIP)
		if err4 != nil {
			log.Println(err4)
		}
	}()

	return &pb_usr.SendActivateEmailResponse{
		TodayLeftTimes: left,
	}, nil
}

// ActivateUser implements pb_usr.AuthServiceServer.
func (as *AuthServer) ActivateUser(ctx context.Context, req *pb_usr.ActivateUserRequest) (
	*pb_usr.ActivateUserResponse, error) {
	email := req.GetEmail()
	code := req.GetActivateCode()

	// check whether the email was already activated or not.
	_, err := as.store.GetUserByEmailActivated(ctx, email)
	if err == nil {
		return nil, status.Errorf(codes.AlreadyExists, "email has been activated")
	}
	if err != sql.ErrNoRows {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to find activated user")
	}

	u, err := as.store.GetUserByEmailNotActivated(ctx, email)
	if err != err {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "email was not registered")
		}
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to find activated user")
	}

	// find activate code.
	activation, err := as.store.GetActivationByUserIDAndCode(ctx,
		db.GetActivationByUserIDAndCodeParams{
			UserID:       u.ID,
			ActivateCode: code,
		})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.InvalidArgument, "activate code not found")
		}
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to find activation")
	}

	// check whether email was expired or not.
	if activation.ExpiredAt.Before(time.Now()) {
		// delete expired code.
		go func() {
			err2 := as.store.DeleteActivation(context.Background(), activation.ID)
			if err2 != nil {
				log.Println(err2)
			}
		}()
		return nil, status.Errorf(codes.DeadlineExceeded, "activate code has been expired")
	}

	// success activate
	err = as.store.ActivateUser(ctx, u.ID)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to activate user")
	}

	// whether delete action successes or not, return ok
	as.store.DeleteActivation(ctx, activation.ID)

	link := filepath.Join(as.config.KodoLink, "img/avatar/default.png")
	aur := &pb_usr.ActivateUserResponse{User: db2pbUser(u, &db.Profile{
		AvatarLink: sql.NullString{String: link, Valid: true},
		})}
	return aur, nil
}

// LoginUser implements pb_usr.AuthServiceServer.
func (as *AuthServer) LoginUser(ctx context.Context, req *pb_usr.LoginUserRequest) (
	*pb_usr.LoginUserResponse, error) {
	email := req.GetEmail()
	passwd := req.GetPassword()
	u, err := as.store.GetUserByEmailActivated(ctx, email)
	if err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to find user")
	}
	// valid passworn
	err2 := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(passwd))
	if err2 != nil {
		log.Println(err2)
		return nil, status.Errorf(codes.Unauthenticated, "incorrect password")
	}
	// create an access token.
	token, _, err3 := as.token.CreateToken(u.ID, u.Nickname, as.config.AccessDuration)
	if err3 != nil {
		log.Println(err3)
		return nil, status.Errorf(codes.Internal, "cannot access token")
	}
	// query profile.
	p, err4 := as.store.GetProfile(ctx, u.ProfileID)
	if err4 != nil {
		log.Println(err4)
		return nil, status.Errorf(codes.Internal, "failed to find profile")
	}
	u2 := db2pbUser(u, &p)
	u2.Profile.AvatarLink = filepath.Join(as.config.KodoLink, u2.Profile.AvatarLink)
	lur := &pb_usr.LoginUserResponse{Token: token, User: u2}
	return lur, nil
}

// RegisterUser implements pb_usr.AuthServiceServer.
func (as *AuthServer) RegisterUser(ctx context.Context, req *pb_usr.RegisterUserRequest) (
	*pb_usr.RegisterUserResponse, error) {
	email := req.GetEmail()
	nickname := req.GetNickname()
	passwd := req.GetPassword()

	// check for email and nickname.
	_, err := as.store.GetUserByEmail(ctx, email)
	if err == nil {
		return nil, status.Errorf(codes.AlreadyExists, "email already exists")
	}
	if err != sql.ErrNoRows {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to find email")
	}
	_, err2 := as.store.GetUserByNickName(ctx, nickname)
	if err2 == nil {
		return nil, status.Errorf(codes.AlreadyExists, "nickname already exists")
	}
	if err != sql.ErrNoRows {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to find nickname")
	}

	hashPasswd, err3 := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	if err3 != nil {
		log.Println(err3)
		return nil, status.Errorf(codes.Internal, "failed to hash password")
	}
	cutp := db.CreateUserTxParams{
		NickName: nickname,
		Email:    email,
		Password: string(hashPasswd),
	}
	cutr, err3 := as.store.CreateUserTx(ctx, cutp)
	if err3 != nil {
		log.Println(err3)
		return nil, status.Errorf(codes.Internal, "failed to create user")
	}

	link := filepath.Join(as.config.KodoLink, "img/avatar/default.png")
	u := db2pbUser(cutr.User, &db.Profile{UpdatedAt: time.Now(), AvatarLink: sql.NullString{
		String: link, Valid: true}})
	rur := &pb_usr.RegisterUserResponse{User: u}
	return rur, nil
}

var _ pb_usr.AuthServiceServer = (*AuthServer)(nil)

// NewAuthServer creates an auth server then return it and an error, if any.
func NewAuthServer(store db.Store, conf config.Config, token token.TokenMaker) *AuthServer {
	ms := mail.NewMailSender(conf.SmtpHost, conf.SmtpAddr, conf.SmtpScrt, conf.SmtpPort)
	return &AuthServer{
		config: conf,
		store:  store,
		token:  token,
		mail:   ms,
	}
}
