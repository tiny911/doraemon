package peername

import (
	"github.com/tiny911/doraemon/assembly"
	"github.com/tiny911/doraemon/interceptor/peername"
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
	this.assembly.AddUnaryServerInterceptor(peername.UnaryServerInterceptor())
}

func (this *Interceptor) setupSSI() {
	//nothing to do
}

func (this *Interceptor) setupUCI() {
	this.assembly.AddUnaryClientInterceptor(peername.UnaryClientInterceptor())
}

func (this *Interceptor) setupSCI() {
	//nothing to do
}

func New() *Interceptor {
	return &Interceptor{}
}
