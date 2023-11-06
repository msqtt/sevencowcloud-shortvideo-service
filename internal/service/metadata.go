package service

import (
	"context"
	"strings"

	"google.golang.org/grpc/metadata"
)

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	userAgentHeader            = "user-agent"
	xForwardedForHeader        = "x-forwarded-for"
	Authorization              = "authorization"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
	Token     string
}

func extractMetadata(ctx context.Context) *Metadata {
	mtdt := &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if userAgents := md.Get(grpcGatewayUserAgentHeader); len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		}

		if authorization := md.Get(Authorization); len(authorization) > 0 {
			mtdt.Token = authorization[0]
		}

		if userAgents := md.Get(userAgentHeader); len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		}

		if clientIPs := md.Get(xForwardedForHeader); len(clientIPs) > 0 {
			mtdt.ClientIP = strings.Split(clientIPs[0], ",")[0]
		}
	}

	// if p, ok := peer.FromContext(ctx); ok {
	// 	mtdt.ClientIP = p.Addr.String()
	// }

	return mtdt
}
