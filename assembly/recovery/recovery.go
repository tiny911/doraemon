package recovery

import (
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/tiny911/doraemon/assembly"
)

var _ assembly.IInterceptor = &Interceptor{}

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
	this.assembly.AddUnaryServerInterceptor(grpc_recovery.UnaryServerInterceptor())
}

func (this *Interceptor) setupSSI() {
	this.assembly.AddStreamServerInterceptor(grpc_recovery.StreamServerInterceptor())
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
