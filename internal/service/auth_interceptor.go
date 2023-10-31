package service

import (
	"context"
	"log"

	"github.com/msqtt/sevencowcloud-shortvideo-service/internal/pkg/config"
	"github.com/msqtt/sevencowcloud-shortvideo-service/internal/token"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	token        token.TokenMaker
	accessMethod map[string]struct{}
}

// NewAuthInterceptor creates an AuthInterceptor with methods which need to be authorized
func NewAuthInterceptor(conf config.Config, token token.TokenMaker, methods ...string) *AuthInterceptor {
	accessMap := make(map[string]struct{})
	for _, v := range methods {
		accessMap[v] = struct{}{}
	}
	return &AuthInterceptor{
		token:        token,
		accessMethod: accessMap,
	}
}

// Unary returns an auth interceptor.
// auth interceptor authorizes the token then wrap payload with context.
func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		payload, err := interceptor.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}
		ctx2 := context.WithValue(ctx, "payload", payload)
		return handler(ctx2, req)
	}
}

func (interceptor *AuthInterceptor) authorize(ctx context.Context, methord string) (*token.Payload, error) {
	_, ok := interceptor.accessMethod[methord]
	if !ok {
		// not login user can access
		return nil, nil
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}
	values := md["authorization"]
	if len(values) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "access token is not provied")
	}

	accessToken := values[0]
	payload, err := interceptor.token.ValidToken(accessToken)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}

	return payload, err
}
