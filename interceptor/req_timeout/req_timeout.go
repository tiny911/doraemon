package req_timeout

import (
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// UnaryClientInterceptor 处理超时
func UnaryClientInterceptor(opts ...Option) grpc.UnaryClientInterceptor {
	o := evaluateOptions(opts)
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		timeoutCtx, _ := context.WithTimeout(ctx, time.Duration(o.timeout)*time.Millisecond)
		return invoker(timeoutCtx, method, req, reply, cc, opts...)
	}
}
