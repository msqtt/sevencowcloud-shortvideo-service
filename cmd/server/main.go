package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb_usr "github.com/msqtt/sevencowcloud-shortvideo-service/api/pb/v1/user"
	"github.com/msqtt/sevencowcloud-shortvideo-service/internal/pkg/config"
	db "github.com/msqtt/sevencowcloud-shortvideo-service/internal/repo/sqlc"
	"github.com/msqtt/sevencowcloud-shortvideo-service/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	conf, err := config.LoadConfig("configs")
	if err != nil {
		log.Fatal(err)
	}

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
	authServer, err := service.NewAuthServer(store, conf)
	if err != nil {
		log.Fatal(err)
	}

	listenGrpc, err := net.Listen("tcp", conf.GrpcAddr)
	if err != nil {
		log.Fatal(err)
	}

	listenREST, err := net.Listen("tcp", conf.HttpAddr)
	if err != nil {
		log.Fatal(err)
	}

	errch := make(chan error)
	go func() {
		err = runGRPServer(authServer, false, listenGrpc)
		if err != nil {
			errch <- fmt.Errorf("cannot start grpc server: %w", err)
		}
	}()

	err = runRESTServer(authServer, false, listenREST)
	if err != nil {
		log.Fatalf("cannot start grpc server: %v", err)
	}

	err = <-errch
	if err != nil {
		log.Fatal(err)
	}
}

func runRESTServer(
	authServer pb_usr.AuthServiceServer,
	enableTLS bool,
	listener net.Listener,
) error {
	mux := runtime.NewServeMux()
	ctx, cf := context.WithCancel(context.Background())
	defer cf()

	err := pb_usr.RegisterAuthServiceHandlerServer(ctx, mux, authServer)
	if err != nil {
		return err
	}

	log.Printf("Start REST server at %s, TLS = %v", listener.Addr().String(), enableTLS)
	if enableTLS {
		return http.ServeTLS(listener, mux, "", "")
	}
	return http.Serve(listener, mux)
}

func runGRPServer(
	authServer pb_usr.AuthServiceServer,
	enableTLS bool,
	listener net.Listener,
) error {
	server := grpc.NewServer()
	pb_usr.RegisterAuthServiceServer(server, authServer)
	// for registering explore grpc api.
	reflection.Register(server)

	log.Printf("Start GRPC server at %s, TLS = %v", listener.Addr().String(), enableTLS)
	return server.Serve(listener)
}
