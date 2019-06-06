package logger

import (
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/tiny911/doraemon/assembly"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

var _ assembly.IInterceptor = &Interceptor{}

type Interceptor struct {
	log      *zap.Logger
	assembly *assembly.Assembly
}

func (this *Interceptor) With(assembly *assembly.Assembly) {
	this.assembly = assembly
}

func (this *Interceptor) Setup() {
	this.setupUSI()
	this.setupSSI()
	this.setupUCI()
	this.setupSCI()
}

func (this *Interceptor) Unload() {
	if this.log != nil {
		this.log.Sync()
	}
}

func call(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
	return true
}

func (this *Interceptor) setupUSI() {
	this.assembly.AddUnaryServerInterceptor(grpc_zap.UnaryServerInterceptor(this.log))
	//this.assembly.AddUnaryServerInterceptor(grpc_zap.PayloadUnaryServerInterceptor(this.log, call))
}

func (this *Interceptor) setupSSI() {
	this.assembly.AddStreamServerInterceptor(grpc_zap.StreamServerInterceptor(this.log))
}

func (this *Interceptor) setupUCI() {
	this.assembly.AddUnaryClientInterceptor(grpc_zap.UnaryClientInterceptor(this.log))
}

func (this *Interceptor) setupSCI() {
	this.assembly.AddStreamClientInterceptor(grpc_zap.StreamClientInterceptor(this.log))
}

func New(log *zap.Logger) *Interceptor {
	return &Interceptor{log: log}
}
