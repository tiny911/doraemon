package ctx_tags

import (
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/tiny911/doraemon/assembly"
)

type Interceptor struct {
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
	//nothing to do
}

func (this *Interceptor) setupUSI() {
	this.assembly.AddUnaryServerInterceptor(grpc_ctxtags.UnaryServerInterceptor())
}

func (this *Interceptor) setupSSI() {
	this.assembly.AddStreamServerInterceptor(grpc_ctxtags.StreamServerInterceptor())
	//nothind to do
}

func (this *Interceptor) setupUCI() {
	//nothing to do
}

func (this *Interceptor) setupSCI() {
	//nothing to do
}

func New() *Interceptor {
	return &Interceptor{}
}
