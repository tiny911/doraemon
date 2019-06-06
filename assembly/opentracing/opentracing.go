package opentracing

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	opentracing "github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	"github.com/tiny911/doraemon/assembly"
)

const (
	enable  string = "on"
	disable string = "off"
)

var _ assembly.IInterceptor = &Interceptor{}

var (
	defaultTraceCollectType = "http"
	defaultTraceCollectAddr = "127.0.0.1:9411"
	defaultTraceSampleRate  = "1.0" // <1.0 不采样  >=1.0全采样
)

type Interceptor struct {
	svrName   string
	svrAddr   string
	tracer    opentracing.Tracer
	collector zipkin.Collector
	assembly  *assembly.Assembly
}

func (this *Interceptor) With(assembly *assembly.Assembly) {
	this.assembly = assembly
}

func (this *Interceptor) Setup() {
	if os.Getenv("ENV_TRACE_SWITCH") != enable { // trace关闭状态
		return
	}

	this.setupUSI()
	this.setupSSI()
	this.setupUCI()
	this.setupSCI()
}

func (this *Interceptor) Unload() {
	if os.Getenv("ENV_TRACE_SWITCH") != enable { // trace关闭状态
		return
	}

	if this.collector != nil {
		this.collector.Close()
	}
}

func (this *Interceptor) setupUSI() {
	this.assembly.AddUnaryServerInterceptor(grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(this.tracer)))
}

func (this *Interceptor) setupSSI() {
	this.assembly.AddStreamServerInterceptor(grpc_opentracing.StreamServerInterceptor(grpc_opentracing.WithTracer(this.tracer)))
}

func (this *Interceptor) setupUCI() {
	this.assembly.AddUnaryClientInterceptor(grpc_opentracing.UnaryClientInterceptor(grpc_opentracing.WithTracer(this.tracer)))
}

func (this *Interceptor) setupSCI() {
	this.assembly.AddStreamClientInterceptor(grpc_opentracing.StreamClientInterceptor(grpc_opentracing.WithTracer(this.tracer)))
}

func New(svrName, svrAddr string) *Interceptor {
	var (
		err       error
		tracer    opentracing.Tracer
		collector zipkin.Collector
	)

	if tracer, collector, err = newTracer(svrName, svrAddr); err != nil {
		tracer = opentracing.NoopTracer{}
		log.Printf("[trace] init failed, err:%s.\n", err)
	}

	return &Interceptor{
		svrName:   svrName,
		svrAddr:   svrAddr,
		tracer:    tracer,
		collector: collector,
	}
}

func newTracer(svrName, svrAddr string) (opentracing.Tracer, zipkin.Collector, error) {
	var (
		traceCollectType = defaultTraceCollectType
		traceCollectAddr = defaultTraceCollectAddr
		traceSampleRate  = defaultTraceSampleRate
		collector        zipkin.Collector
		tracer           opentracing.Tracer
		err              error
	)

	if collectType := os.Getenv("ENV_TRACE_COLLECT_TYPE"); collectType != "" {
		traceCollectType = collectType
	}

	if collectAddr := os.Getenv("ENV_TRACE_COLLECT_ADDR"); collectAddr != "" {
		traceCollectAddr = collectAddr
	}

	if sampleRate := os.Getenv("ENV_TRACE_SAMPLE_RATE"); sampleRate != "" {
		traceSampleRate = sampleRate
	}

	switch traceCollectType {
	case "http":
		collector, err = zipkin.NewHTTPCollector(fmt.Sprintf("http://%s/api/v1/spans", traceCollectAddr))
	case "kafka":
		collector, err = zipkin.NewKafkaCollector(strings.Split(traceCollectAddr, ","))
	default:
		log.Panicf("[trace] collector type[%s] is illegal.\n", traceCollectType)
	}

	if err != nil {
		return nil, nil, err
	}

	rate, _ := strconv.ParseFloat(traceSampleRate, 32)
	tracer, err = zipkin.NewTracer(
		zipkin.NewRecorder(collector, false, svrAddr, svrName),
		zipkin.ClientServerSameSpan(true),
		zipkin.TraceID128Bit(true),
		zipkin.WithSampler(zipkin.NewCountingSampler(rate)),
	)

	return tracer, collector, err
}
