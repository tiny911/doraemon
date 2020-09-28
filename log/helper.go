package log

import (
	"fmt"
	"log"
	"os"
	"sync"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/tiny911/gobase/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logHelper struct {
	svrName   string        //服务的名称
	logPath   string        //日志路径名
	logName   string        //日志文件名
	logFile   string        //日志文件
	logSpvor  string        //supervisor日志文件
	logPrefix string        //日志前缀
	logKit    *zap.Logger   //zap日志对象
	config    zap.Config    //zap日志配置
	logLevel  zapcore.Level //zap日志级别
	sync.Mutex
}

const enableLogSpvor = "on"

func NewLogHelper(svr, prefix string) *logHelper {
	return &logHelper{
		svrName:   svr,
		logPrefix: prefix,
	}
}

func (this *logHelper) SetLogLevel(level string) {
	switch level {
	case "debug":
		this.logLevel = zapcore.DebugLevel
	case "info":
		this.logLevel = zapcore.InfoLevel
	case "warn":
		this.logLevel = zapcore.WarnLevel
	case "error":
		this.logLevel = zapcore.ErrorLevel
	default:
		log.Panicf("[logger] level[%s] is illegal.", level)
	}
}

func (this *logHelper) GetLogKit() *zap.Logger {
	return this.logKit
}

func (this *logHelper) Execute() {
	if exists := this.islogFileExist(); !exists {
		this.createLogFile()
	}

	this.loadKit()
}

func (this *logHelper) Cancel() {
	if Kit != nil {
		Kit.Sync()
	}
}

func (this *logHelper) LogLevelSwitch() {
	log.Printf("[log] Level switch from:%v.", this.logLevel)

	this.Lock() //锁过程,保证kit加载不竞争
	defer this.Unlock()

	if this.logLevel == zapcore.DebugLevel {
		this.logLevel = zapcore.InfoLevel
	} else {
		this.logLevel = zapcore.DebugLevel
	}

	this.config.Level.SetLevel(this.logLevel)

	log.Printf("[log] Level switch to:%v.", this.logLevel)
}

func (this *logHelper) loadKit() {
	cfg := zap.NewProductionConfig()

	cfg.OutputPaths = []string{this.logFile}
	cfg.DisableStacktrace = true
	cfg.Level = zap.NewAtomicLevelAt(this.logLevel)
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	{ // 如果设置了ENV_LOG_SUPERVISOR环境变量为on，则输出相应日志
		envLogSupervisor := os.Getenv("ENV_LOG_SUPERVISOR")
		if envLogSupervisor == enableLogSpvor {
			cfg.OutputPaths = append(cfg.OutputPaths, this.logSpvor)
		}
	}

	logger, err := cfg.Build()
	if err != nil {
		log.Panicf("[logger] config build failed, err:%s.\n", err)
	}

	logger = logger.WithOptions(zap.AddCallerSkip(1))
	grpc_zap.ReplaceGrpcLogger(logger)
	this.logKit = logger
	this.config = cfg
}

func (this *logHelper) islogFileExist() bool {
	var (
		err    error
		exists bool
	)

	logPath := fmt.Sprintf("%s/", this.logPrefix)
	logName := fmt.Sprintf("%s.log", this.svrName)
	logFile := logPath + logName

	exists, err = utils.PathExists(logFile)
	if err != nil {
		log.Panicf("[logger] logFile[%s] exists, err:%s.\n", logFile, err)
	}

	this.logPath = logPath
	this.logName = logName
	this.logFile = logFile
	this.logSpvor = fmt.Sprintf("%s/%s", logPath, "supervisor.log")

	return exists
}

func (this *logHelper) createLogFile() {
	err := os.MkdirAll(this.logPath, 0755)
	if err != nil {
		log.Panicf("[logger] create logPath[%s] failed, err:%s.\n", this.logPath, err)
	}
	_, err = os.Create(this.logFile)
	if err != nil {
		log.Panicf("[logger] create logFile[%s] failed, err:%s.\n", this.logFile, err)
	}
}
