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

// SendCaptcha implements pb_usr.AuthServiceServer.
func (as *AuthServer) SendCaptcha(ctx context.Context, req *pb_usr.SendCaptchaRequest) (
	*pb_usr.SendCaptchaResponse, error) {
	email := req.GetEmail()

	cnt, err := as.store.TodayEmailCount(ctx, email)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to sum captcha times")
	}
	left := int32(as.config.ActivateTimes) - int32(cnt)
	if left <= 0 {
		return nil, status.Errorf(codes.ResourceExhausted, "the user has been reached %d times"+
			" limit of sending activate emails", cnt)
	}
	left--

	code := sample.RandomStr(8)
	params := db.AddCaptchaParams{
		Email:     email,
		Captcha:   code,
		ExpiredAt: time.Now().Add(10 * time.Minute),
	}
	_, err2 := as.store.AddCaptcha(ctx, params)
	if err2 != nil {
		return nil, status.Errorf(codes.Internal, "failed to insert captcha")
	}

	// open a goroutine to send an email
	m := extractMetadata(ctx)
	go func() {
		err := as.mail.SendActivateEmail(email, code, m.ClientIP)
		if err != nil {
			log.Println(err)
		}
	}()

	return &pb_usr.SendCaptchaResponse{
		TodayLeftTimes: left,
	}, nil
}

// LoginUser implements pb_usr.AuthServiceServer.
func (as *AuthServer) LoginUser(ctx context.Context, req *pb_usr.LoginUserRequest) (
	*pb_usr.LoginUserResponse, error) {
	email := req.GetEmail()
	passwd := req.GetPassword()
	u, err := as.store.GetUserByEmail(ctx, email)
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
	captcha := req.GetCaptcha()

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

	// for activate
	param := db.GetCaptchaByEmailAndCodeParams{
		Captcha: captcha,
		Email:   email,
	}

	capcha, err4 := as.store.GetCaptchaByEmailAndCode(ctx, param)
	if err4 != nil {
		if err4 == sql.ErrNoRows {
			return nil, status.Errorf(codes.InvalidArgument, "wrong captcha")
		}
		log.Println(err4)
		return nil, status.Errorf(codes.Internal, "failed to find captcha")
	}

	go func() {
		err2 := as.store.DeleteCaptcha(context.Background(), capcha.ID)
		if err2 != nil {
			log.Println(err2)
		}
	}()

	if capcha.ExpiredAt.Before(time.Now()) {
		return nil, status.Errorf(codes.DeadlineExceeded, "capcha was expired")
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
func NewAuthServer(conf config.Config, token token.TokenMaker, store db.Store) *AuthServer {
	ms := mail.NewMailSender(conf.SmtpHost, conf.SmtpAddr, conf.SmtpScrt, conf.SmtpPort)
	return &AuthServer{
		config: conf,
		store:  store,
		token:  token,
		mail:   ms,
	}
}
