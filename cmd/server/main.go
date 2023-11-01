package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb_fl "github.com/msqtt/sevencowcloud-shortvideo-service/api/pb/v1/follow"
	pb_prf "github.com/msqtt/sevencowcloud-shortvideo-service/api/pb/v1/profile"
	pb_usr "github.com/msqtt/sevencowcloud-shortvideo-service/api/pb/v1/user"
	pb_vid "github.com/msqtt/sevencowcloud-shortvideo-service/api/pb/v1/video"
	"github.com/msqtt/sevencowcloud-shortvideo-service/internal/pkg/config"
	db "github.com/msqtt/sevencowcloud-shortvideo-service/internal/repo/sqlc"
	"github.com/msqtt/sevencowcloud-shortvideo-service/internal/service"
	"github.com/msqtt/sevencowcloud-shortvideo-service/internal/token"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	// load config
	conf, err := config.LoadConfig("configs")
	if err != nil {
		log.Fatal(err)
	}
	// create db store
	var database *sql.DB
	if conf.PROC {
		database, err = sql.Open(conf.DBDriver, conf.DBSource)
	} else {
		database, err = sql.Open(conf.DBDriver, conf.DevDBSource)
	}
	if err != nil {
		log.Fatal(err)
	}
	store := db.NewSQLStore(database)
	pm, err2 := token.NewPasetoMaker([]byte(conf.TokenSymmetricKey))
	if err2 != nil {
		log.Fatal(err2)
	}
	// new service server
	authServer := service.NewAuthServer(conf, pm, store)
	profileServer := service.NewProfileServer(conf, store)
	fileServer := service.NewFileService(conf, pm, store)
	followServer := service.NewFollowServer(conf, pm, store)
	videoServer := service.NewVideoServer(conf, pm, store)

	// new service interceptor
	ai := service.NewAuthInterceptor(conf, pm,
		"/profile.v1.ProfileService/UpdateProfile",
		"/follow.v1.FollowService/FollowUser",
		"/follow.v1.FollowService/UnFollowUser",
	)

	// start listen server
	listenGrpc, err := net.Listen("tcp", conf.GrpcAddr)
	if err != nil {
		log.Fatal(err)
	}
	listenREST, err := net.Listen("tcp", conf.HttpAddr)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		log.Fatalln(runGRPCServer(authServer, profileServer, followServer, videoServer, ai, false, listenGrpc))
	}()
	// http gateway
	log.Fatalln(runRESTServer(fileServer, false, listenREST, listenGrpc))
}

func runGRPCServer(
	authServer pb_usr.AuthServiceServer,
	profileServer pb_prf.ProfileServiceServer,
	followServer pb_fl.FollowServiceServer,
	videoServer pb_vid.VideoServiceServer,
	authInerceptor *service.AuthInterceptor,
	enableTLS bool,
	listener net.Listener,
) error {

	server := grpc.NewServer(grpc.UnaryInterceptor(authInerceptor.Unary()))
	pb_usr.RegisterAuthServiceServer(server, authServer)
	pb_prf.RegisterProfileServiceServer(server, profileServer)
	pb_fl.RegisterFollowServiceServer(server, followServer)
	pb_vid.RegisterVideoServiceServer(server, videoServer)

	// for registering explore grpc api.
	reflection.Register(server)

	log.Printf("Start GRPC server at %s, TLS = %v", listener.Addr().String(), enableTLS)
	return server.Serve(listener)
}

func runRESTServer(
	fs *service.FileServer,
	enableTLS bool,
	restListener net.Listener,
	grpcListener net.Listener,
) error {
	mux := runtime.NewServeMux()
	err2 := mux.HandlePath("PUT", "/v1/file/avatar/{user_id}", fs.UploadAvatar)
	if err2 != nil {
		return err2
	}
	withCors := cors.New(cors.Options{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"ACCEPT", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler(mux)

	ctx, cf := context.WithCancel(context.Background())
	defer cf()

	conn, err := grpc.DialContext(
		context.Background(),
		grpcListener.Addr().String(),
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	err = pb_usr.RegisterAuthServiceHandler(ctx, mux, conn)
	if err != nil {
		return err
	}
	err = pb_prf.RegisterProfileServiceHandler(ctx, mux, conn)
	if err != nil {
		return err
	}
	err = pb_fl.RegisterFollowServiceHandler(ctx, mux, conn)
	if err != nil {
		return err
	}
	err = pb_vid.RegisterVideoServiceHandler(ctx, mux, conn)
	if err != nil {
		return err
	}

	log.Printf("Start REST server at %s, TLS = %v", restListener.Addr().String(), enableTLS)
	if enableTLS {
		return http.ServeTLS(restListener, withCors, "", "")
	}
	return http.Serve(restListener, withCors)
}
