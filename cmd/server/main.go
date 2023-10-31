package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb_prf "github.com/msqtt/sevencowcloud-shortvideo-service/api/pb/v1/profile"
	pb_usr "github.com/msqtt/sevencowcloud-shortvideo-service/api/pb/v1/user"
	"github.com/msqtt/sevencowcloud-shortvideo-service/internal/pkg/config"
	db "github.com/msqtt/sevencowcloud-shortvideo-service/internal/repo/sqlc"
	"github.com/msqtt/sevencowcloud-shortvideo-service/internal/service"
	"github.com/msqtt/sevencowcloud-shortvideo-service/internal/token"
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
	authServer := service.NewAuthServer(store, conf, pm)
	profileServer := service.NewProfileServer(store, conf)
	fileServer := service.NewFileService(conf, pm, store)

	// new service interceptor
	ai := service.NewAuthInterceptor(conf, pm,
		"/profile.v1.ProfileService/UpdateProfile",
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
		log.Fatalln(runGRPCServer(authServer, profileServer, ai, false, listenGrpc))
	}()
	// http gateway
	log.Fatalln(runRESTServer(fileServer, false, listenREST, listenGrpc))
}

func runGRPCServer(
	authServer pb_usr.AuthServiceServer,
	profileServer pb_prf.ProfileServiceServer,
	authInerceptor *service.AuthInterceptor,
	enableTLS bool,
	listener net.Listener,
) error {

	server := grpc.NewServer(grpc.UnaryInterceptor(authInerceptor.Unary()))
	pb_usr.RegisterAuthServiceServer(server, authServer)
	pb_prf.RegisterProfileServiceServer(server, profileServer)

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

	log.Printf("Start REST server at %s, TLS = %v", restListener.Addr().String(), enableTLS)
	if enableTLS {
		return http.ServeTLS(restListener, mux, "", "")
	}
	return http.Serve(restListener, mux)
}

