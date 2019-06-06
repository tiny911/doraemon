package req_timeout

import (
	"github.com/tiny911/doraemon/assembly"
	"github.com/tiny911/doraemon/interceptor/req_timeout"
)

var _ assembly.IInterceptor = &Interceptor{}

// Interceptor 超时拦截器结构体
type Interceptor struct {
	timeout  int
	assembly *assembly.Assembly
}

// New 生成拦截器
func New(timeout int) *Interceptor {
	return &Interceptor{timeout: timeout}
}

// With 集成拦截器
func (i *Interceptor) With(assembly *assembly.Assembly) {
	i.assembly = assembly
}

// Setup 启动拦截器
func (i *Interceptor) Setup() {
	i.setupUSI()
	i.setupSSI()
	i.setupUCI()
	i.setupSCI()
}

// Unload 卸载拦截器
func (i *Interceptor) Unload() {
	//nothing to do
}

func (i *Interceptor) setupUSI() {
	//nothing to do
}

func (i *Interceptor) setupSSI() {
	//nothing to do
}

func (i *Interceptor) setupUCI() {
	i.assembly.AddUnaryClientInterceptor(req_timeout.UnaryClientInterceptor(req_timeout.WithTimeout(i.timeout)))
}

func (i *Interceptor) setupSCI() {
	//nothing to do
}
