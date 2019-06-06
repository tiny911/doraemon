package prometheus

import (
	"net/http"
	"os"

	"time"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tiny911/doraemon/assembly"
	"google.golang.org/grpc"
)

const (
	enable  string = "on"
	disable string = "off"
)

var _ assembly.IInterceptor = &Interceptor{}

type Interceptor struct {
	assembly         *assembly.Assembly
	hasRegistHandler bool
}

func (this *Interceptor) With(assembly *assembly.Assembly) {
	this.assembly = assembly
}

func (this *Interceptor) Setup() {
	if os.Getenv("ENV_PROMETHEUS_SWITCH") != enable { // prometheus关闭状态
		return
	}

	{
		this.setupUSI()
		this.setupSSI()
		this.setupUCI()
		this.setupSCI()
	}

	if os.Getenv("ENV_PROMETHEUS_HISTOGRAM") == enable {
		grpc_prometheus.EnableHandlingTimeHistogram()
	}
}

func (this *Interceptor) Unload() {
	if os.Getenv("ENV_PROMETHEUS_SWITCH") != enable { // prometheus关闭状态
		return
	}
}

func (this *Interceptor) setupUSI() {
	this.assembly.AddUnaryServerInterceptor(grpc_prometheus.UnaryServerInterceptor)
	if !this.hasRegistHandler {
		http.Handle("/metrics", promhttp.Handler())
		this.hasRegistHandler = true
	}
}

func (this *Interceptor) setupSSI() {
	this.assembly.AddStreamServerInterceptor(grpc_prometheus.StreamServerInterceptor)
	if !this.hasRegistHandler {
		http.Handle("/metrics", promhttp.Handler())
		this.hasRegistHandler = true
	}
}

func (this *Interceptor) setupUCI() {
	this.assembly.AddUnaryClientInterceptor(grpc_prometheus.UnaryClientInterceptor)
}

func (this *Interceptor) setupSCI() {
	this.assembly.AddStreamClientInterceptor(grpc_prometheus.StreamClientInterceptor)
}

func New() *Interceptor {
	return &Interceptor{}
}

func (this *Interceptor) WithServer(server **grpc.Server) *Interceptor {
	go func() {
		for *server == nil {
			time.Sleep(100 * time.Millisecond)
		}
		grpc_prometheus.Register(*server)
	}()
	return this
}
