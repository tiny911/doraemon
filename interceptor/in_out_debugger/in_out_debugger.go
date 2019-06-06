package in_out_debugger

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)

		extCtx := grpc_zap.Extract(ctx)
		extCtx.Debug("debug for logging request and response",
			zap.Reflect("req", req),
			zap.Reflect("rsp", resp),
			zap.Reflect("err", err),
		)

		return resp, err
	}
}
