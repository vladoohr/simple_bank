package gapi

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	xForwardedForHeader = "x-forwarded-for"
	userAgentHeader     = "user-agent"
	GrpcUserAgentHeader = "grpcgateway-user-agent"
)

type Metadata struct {
	UserAgent string
	ClientAPI string
}

func extractMetadata(ctx context.Context) *Metadata {
	mtdt := &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if userAgents := md.Get(userAgentHeader); len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		}

		if userAgents := md.Get(GrpcUserAgentHeader); len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		}

		if clientIPs := md.Get(xForwardedForHeader); len(clientIPs) > 0 {
			mtdt.ClientAPI = clientIPs[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		fmt.Printf("Metadata: %+v", p)
		mtdt.ClientAPI = p.Addr.String()
	}

	return mtdt
}
