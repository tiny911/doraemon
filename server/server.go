package server

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "net/http/pprof" //开启pprof

	logpkg "log"

	"github.com/tiny911/doraemon/assembly"
	"github.com/tiny911/doraemon/assembly/ctx_tags"
	"github.com/tiny911/doraemon/assembly/in_out_debugger"
	"github.com/tiny911/doraemon/assembly/localtracing"
	"github.com/tiny911/doraemon/assembly/logger"
	"github.com/tiny911/doraemon/assembly/opentracing"
	"github.com/tiny911/doraemon/assembly/peername"
	"github.com/tiny911/doraemon/assembly/prometheus"
	"github.com/tiny911/doraemon/assembly/recovery"
	"github.com/tiny911/doraemon/log"
	"github.com/tiny911/gobase/utils"

	"github.com/gorilla/handlers"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// Server 包含了grpc server, 集成assembly下不同的拦截器
type Server struct {
	name        string //服务名称
	hostIP      string //主机ip
	environment string //服务环境

	pidTag  string //进程号
	pidFile string //进程文件

	RPCSvr      *grpc.Server //rpc服务
	rpcPort     int          //rpc端口
	rpcRegister func()       //rpc注册

	HTTPSvr       *http.Server   //http服务
	httpPort      int            //http端口
	httpRegisters []HTTPRegister //http注册

	//naming       *naming.Naming     //naming

	assembly     *assembly.Assembly //assembly
	interceptors []assembly.IInterceptor

	httpGatewayOptions []runtime.ServeMuxOption

	quitChan    chan interface{}
	quitTimeout time.Duration
	//quitContext
	//cancel context.CancelFunc
}

// HTTPRegister 注册http回调类型
type HTTPRegister func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)

var defaultQuitTimeout = 3 * time.Second

// New 生产Server实例
func New() *Server {
	var (
		hostIP      = os.Getenv("ENV_HOST_IP")
		serverName  = os.Getenv("ENV_SERVER_NAME")
		environment = os.Getenv("ENV_ENVIRONMENT")
	)

	return &Server{
		name:               serverName,
		hostIP:             hostIP,
		environment:        environment,
		interceptors:       make([]assembly.IInterceptor, 0),
		httpGatewayOptions: make([]runtime.ServeMuxOption, 0),
		quitChan:           make(chan interface{}),
		quitTimeout:        defaultQuitTimeout,
	}
}

// SetHostIP 设置服务ip
func (s *Server) SetHostIP(ip string) {
	s.hostIP = ip
}

// SetRPCRegister rpc注册
func (s *Server) SetRPCRegister(register func()) {
	s.rpcRegister = register
}

// SetHTTPRegisters http注册
func (s *Server) SetHTTPRegisters(registers []HTTPRegister) {
	s.httpRegisters = registers
}

// SetPort 设置服务端口
func (s *Server) SetPort(port int) {
	s.rpcPort = port
	s.httpPort = port + 1
}

// AddInterceptor 增加拦截器
func (s *Server) AddInterceptor(interceptor assembly.IInterceptor) {
	s.interceptors = append(s.interceptors, interceptor)
}

// AddhttpOptions 增加http选型
func (s *Server) AddHttpGatewayOption(option runtime.ServeMuxOption) {
	s.httpGatewayOptions = append(s.httpGatewayOptions, option)
}

// touchPidFile 创建pid文件
func (s *Server) touchPidFile() {
	var (
		file = fmt.Sprintf("./%s.pid", s.name)
		pid  = strconv.Itoa(os.Getpid())
	)

	err := ioutil.WriteFile(file, []byte(pid), 0777)
	if err != nil {
		log.WithField(log.Fields{
			"error": err,
			"file":  file,
		}).Panic("server touch pid file failed!")
	}

	s.pidFile = file
	s.pidTag = pid

	logpkg.Printf("[pid] Process:%s file touched success.", s.pidTag)
}

// deletePidFile 删除pid文件
func (s *Server) deletePidFile() {
	if exists, _ := utils.PathExists(s.pidFile); exists {
		os.Remove(s.pidFile)
	}

	logpkg.Printf("[pid] Process:%s file deleted success.", s.pidTag)
}

// Run 启动server
func (s *Server) Run() {
	logpkg.Printf("[svr] Server Running.")

	{ //启动时候创建pid文件
		s.touchPidFile()
	}

	go s.rpcServer()  //rpc服务
	go s.httpServer() //代理http server

	// { //注册naming
	// 	s.naming = naming.New(
	// 		s.name,
	// 		s.hostIP,
	// 		s.environment,
	// 		s.rpcPort,
	// 		s.httpPort,
	// 	)
	// 	go s.naming.Regist(os.Getenv("ENV_NAMING_ADDR"), os.Getenv("ENV_NAMING_IDC"))
	// }

	<-s.quitChan
}

// Stop 停止server
func (s *Server) Stop() {
	ctx, _ := context.WithTimeout(context.Background(), s.quitTimeout)

	//s.naming.UnRegist()
	s.HTTPSvr.Shutdown(ctx)
	s.RPCSvr.GracefulStop()
	s.assembly.Unload()
	s.deletePidFile()
	log.Stop()
	close(s.quitChan)
}

func (s *Server) rpcServer() error {
	interceptors := []assembly.IInterceptor{
		ctx_tags.New(),
		logger.New(log.Kit),
		opentracing.New(s.name, fmt.Sprintf("%s:%d", s.hostIP, s.rpcPort)),
		in_out_debugger.New(),
		localtracing.New(),
		peername.New(),
		recovery.New(),
		prometheus.New().WithServer(&s.RPCSvr),
	}
	interceptors = append(interceptors, s.interceptors...)
	s.assembly = assembly.Setup(interceptors...)

	s.RPCSvr = grpc.NewServer(s.assembly.WithUnaryServer())

	s.rpcRegister()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.rpcPort))
	if err != nil {
		logpkg.Printf("[rpc] Listen failed! error:%s.", err)
	}
	return s.RPCSvr.Serve(lis)
}

func (s *Server) httpServer() error {
	var (
		err error
		ctx context.Context
	)

	ctx, _ = context.WithCancel(context.Background())
	runtimeMux := runtime.NewServeMux(s.httpGatewayOptions...)

	for _, register := range s.httpRegisters {
		err = register(ctx, runtimeMux, fmt.Sprintf("127.0.0.1:%d", s.rpcPort), []grpc.DialOption{grpc.WithInsecure()})
		if err != nil {
			logpkg.Printf("[http] Regist failed, error:%s.", err)
			return err
		}
	}

	handler := http.DefaultServeMux
	handler.Handle("/", handlers.CompressHandler(runtimeMux))

	s.HTTPSvr = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.httpPort),
		Handler: handler,
	}

	err = s.HTTPSvr.ListenAndServe()
	if err != nil {
		logpkg.Printf("[http] Listen failed! error:%s.", err)
	}
	return err
}
