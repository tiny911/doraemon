package localtracing

import (
	"reflect"

	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	TraceFlag   = "TraceId"
	TraceHeader = "doraemon-trace-id"
)

func UnaryServerInterceptor(opts ...Option) grpc.UnaryServerInterceptor {
	o := evaluateOptions(opts)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var (
			extCtxTags        = grpc_ctxtags.Extract(ctx)
			extCtxVals        = extCtxTags.Values()
			traceId    string = ""
		)

		if openTraceId, exists := extCtxVals[grpc_opentracing.TagTraceId]; !exists {
			{
				{ //1.先从请求的metadata中取出traceId,如果没有则自生成一个
					md, ok := metadata.FromIncomingContext(ctx)
					if ok && len(md[TraceHeader]) != 0 {
						traceId = md[TraceHeader][0]
					} else {
						traceId = o.idGenerateFunc()
					}

				}

				//2.将traceId传入ctx中，日志中用的traceId，就来源于此
				extCtxTags.Set(grpc_opentracing.TagTraceId, traceId)

				{ //3.将traceId的metadata传入下游
					mdTraceHeader := metadata.Pairs(TraceHeader, traceId)
					if md, ok := metadata.FromOutgoingContext(ctx); ok {
						mdTraceHeader = metadata.Join(mdTraceHeader, md)
					}
					ctx = metadata.NewOutgoingContext(ctx, mdTraceHeader)
				}
			}
		} else {
			if id, ok := openTraceId.(string); ok {
				traceId = id
			}
		}

		resp, err := handler(ctx, req)

		var (
			emptyValue = reflect.Value{}
			respValue  = reflect.ValueOf(resp)
		)

		if respValue != emptyValue {
			flag := respValue.Elem().FieldByName(TraceFlag)
			if flag.CanSet() {
				flag.SetString(traceId)
			}
		}

		//将traceId写回响应的metadata中
		grpc.SendHeader(ctx, metadata.Pairs(TraceHeader, traceId))
		return resp, err
	}
}
