package peername

import (
	"os"

	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	PeerNameFlag = "doraemon-peer-name"
	PeerNameTag  = "peername"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var (
			extCtxTags = grpc_ctxtags.Extract(ctx)
		)

		md, ok := metadata.FromIncomingContext(ctx)
		if ok && len(md[PeerNameFlag]) != 0 {
			extCtxTags.Set(PeerNameTag, md[PeerNameFlag][0])
		}

		resp, err := handler(ctx, req)
		return resp, err
	}
}

func UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		mdPeerName := metadata.Pairs(PeerNameFlag, os.Getenv("ENV_SERVER_NAME"))
		if md, ok := metadata.FromOutgoingContext(ctx); ok {
			mdPeerName = metadata.Join(mdPeerName, md)
		}

		ctx = metadata.NewOutgoingContext(ctx, mdPeerName)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
