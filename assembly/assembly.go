package assembly

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

type (
	// IInterceptor 定义拦截器接口
	IInterceptor interface {
		Setup()
		Unload()
		With(i *Assembly)
	}

	// Assembly 拦截器集合
	Assembly struct {
		usi []grpc.UnaryServerInterceptor
		ssi []grpc.StreamServerInterceptor
		uci []grpc.UnaryClientInterceptor
		sci []grpc.StreamClientInterceptor
		in  []IInterceptor
	}
)

func newAssembly() *Assembly {
	return &Assembly{
		usi: make([]grpc.UnaryServerInterceptor, 0),
		ssi: make([]grpc.StreamServerInterceptor, 0),
		uci: make([]grpc.UnaryClientInterceptor, 0),
		sci: make([]grpc.StreamClientInterceptor, 0),
		in:  make([]IInterceptor, 0),
	}
}

// WithUnaryServer 加载unary服务端拦截器
func (a *Assembly) WithUnaryServer() grpc.ServerOption {
	return grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(a.usi...))
}

// WithStreamServer 加载stream服务端拦截器
func (a *Assembly) WithStreamServer() grpc.ServerOption {
	return grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(a.ssi...))
}

// WithUnaryClient 加载unary客户端拦截器
func (a *Assembly) WithUnaryClient() grpc.DialOption {
	return grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(a.uci...))
}

// WithStreamClient 加载stream客户端拦截器
func (a *Assembly) WithStreamClient() grpc.DialOption {
	return grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(a.sci...))
}

// AddUnaryServerInterceptor 向assembly中加入unary服务端拦截器
func (a *Assembly) AddUnaryServerInterceptor(interceptor grpc.UnaryServerInterceptor) {
	a.usi = append(a.usi, interceptor)
}

// AddStreamServerInterceptor 向assembly中加入stream服务端拦截器
func (a *Assembly) AddStreamServerInterceptor(interceptor grpc.StreamServerInterceptor) {
	a.ssi = append(a.ssi, interceptor)
}

// AddUnaryClientInterceptor 向assembly中加入unary客户端拦截器
func (a *Assembly) AddUnaryClientInterceptor(interceptor grpc.UnaryClientInterceptor) {
	a.uci = append(a.uci, interceptor)
}

// AddStreamClientInterceptor 向assembly中加入stream客户端拦截器
func (a *Assembly) AddStreamClientInterceptor(interceptor grpc.StreamClientInterceptor) {
	a.sci = append(a.sci, interceptor)
}

// Unload 卸载assembly中所有的拦截器
func (a *Assembly) Unload() {
	for _, interceptor := range a.in {
		interceptor.Unload()
	}
}

// Setup 安装assembly中所有的拦截器
func Setup(interceptors ...IInterceptor) *Assembly {
	var (
		assemble = newAssembly()
	)

	for _, interceptor := range interceptors {
		interceptor.With(assemble)
		interceptor.Setup()
		assemble.in = append(assemble.in, interceptor)
	}

	return assemble
}
