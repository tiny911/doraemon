package in_out_debugger

import (
	"github.com/tiny911/doraemon/assembly"
	"github.com/tiny911/doraemon/interceptor/in_out_debugger"
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
	this.assembly.AddUnaryServerInterceptor(in_out_debugger.UnaryServerInterceptor())
}

func (this *Interceptor) setupSSI() {
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
