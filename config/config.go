package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/tiny911/gobase/utils"
	"gopkg.in/yaml.v2"
)

const (
	defaultConfigFileName = "config.yaml"
	defaultConfigFilePath = "./"
	defaultEnvPrefix      = "ENV"
	defaultHostIp         = "127.0.0.1"
	defaultServerVer      = "0.0.0.0"
	defaultHostName       = "localhost"
)

//这里的envCfg结构体，会在Parse时候，被导入环境变量.
//这样每个模块通过环境变量就可以拿到关心的信息了.
// TODO: add others
var envCfg = struct {
	Server_Name          string `yaml:"_Server_Name"`
	Server_Port          string `yaml:"_Server_Port"`
	Server_Ver           string `yaml:"_Server_Ver"`
	Server_Time          string `yaml:"_Server_Time"`
	Environment          string `yaml:"_Environment"`
	Log_Level            string `yaml:"_Log_Level"`
	Log_Supervisor       string `yaml:"_Log_Supervisor"`
	Host_Ip              string `yaml:"_Host_Ip"`
	Host_Idc             string `yaml:"_Host_Idc"`
	Host_Name            string `yaml:"_Host_Name"`
	Naming_Switch        string `yaml:"_Naming_Switch"`
	Naming_IDC           string `yaml:"_Naming_IDC"`
	Naming_Addr          string `yaml:"_Naming_Addr"`
	Trace_Switch         string `yaml:"_Trace_Switch"`
	Trace_Collect_Type   string `yaml:"_Trace_Collect_Type"`
	Trace_Collect_Addr   string `yaml:"_Trace_Collect_Addr"`
	Trace_Sample_Rate    string `yaml:"_Trace_Sample_Rate"`
	Prometheus_Switch    string `yaml:"_Prometheus_Switch"`
	Prometheus_Histogram string `yaml:"_Prometheus_Histogram"`
}{}

//配置文件默认是在当前路径下的config.yaml
//当然，为了方便，我们可以设置ENV_CONFIG_FILE环境变量来加载一个配置文件
func Parse(cfg interface{}) {
	var (
		envConfigFile  string
		configFilePath string = defaultConfigFilePath
		configFileName string = defaultConfigFileName
	)

	if envConfigFile = os.Getenv("ENV_CONFIG_FILE"); envConfigFile != "" {
		index := strings.LastIndexByte(envConfigFile, '/')
		configFilePath = envConfigFile[:index+1]
		configFileName = envConfigFile[index+1:]
	}

	configure := NewConfigure(configFilePath, configFileName)
	configure.LoadData()
	configure.ParseCfg(cfg)

	{
		print(cfg)
	}

	//将解析之后的envCfg导出到环境变量
	{
		configure.ParseCfg(&envCfg)
		export(envCfg)
	}
	return
}

func print(cfg interface{}) {
	t := reflect.TypeOf(cfg).Elem()
	v := reflect.ValueOf(cfg).Elem()

	for i := 0; i < v.NumField(); i++ {
		logPrint(t.Field(i).Name, v.Field(i))
	}
}

type Configure struct {
	filePath string
	fileName string
	content  []byte
}

func NewConfigure(filePath, fileName string) *Configure {
	return &Configure{
		filePath: filePath,
		fileName: fileName,
	}
}

func (this *Configure) LoadData() {
	var (
		file = this.filePath + this.fileName
		err  error
		data []byte
	)

	data, err = ioutil.ReadFile(file)
	if err != nil {
		log.Panicf("[config] read file[%s] failed, errmsg:[%s].\n", file, err)
	}

	this.content = data
	return
}

func (this *Configure) ParseCfg(cfg interface{}) {
	err := yaml.Unmarshal(this.content, cfg)
	if err != nil {
		log.Panicf("[config] parse failed, errmsg:[%s].\n", err)
	}
}

func export(cfg interface{}) {
	var (
		envKey, envVal string
		err            error
	)

	exportHostInfo()
	exportServerInfo()

	t := reflect.TypeOf(cfg)
	v := reflect.ValueOf(cfg)

	for i := 0; i < t.NumField(); i++ {
		envKey = fmt.Sprintf("%s_%s", defaultEnvPrefix, strings.ToUpper(t.Field(i).Name))
		envVal = v.Field(i).String()

		if envVal != "" {
			err = os.Setenv(envKey, envVal)
			if err != nil {
				log.Panicf("[config] export failed, errmsg:[%s].\n", err)
			}
		}
	}

	{ //打印env
		for _, env := range os.Environ() {
			if strings.HasPrefix(env, defaultEnvPrefix) {
				envSlice := strings.Split(env, "=")
				logPrint(envSlice[0], envSlice[1])
			}
		}

	}

}

func exportHostInfo() {
	hostIp, err := utils.GetIP("eth0")
	if err != nil {
		hostIp = defaultHostIp
		log.Printf("[config] get eth0 failed, err:%s.", err)
	}

	hostName, err := os.Hostname()
	if err != nil {
		hostName = defaultHostName
		log.Printf("[config] get hostname failed, err:%s.", err)
	}

	//TODO: add idc
	//...

	os.Setenv("ENV_HOST_NAME", hostName)
	os.Setenv("ENV_HOST_IP", hostIp)
}

func exportServerInfo() {
	var ver string

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Printf("[config] get server dir failed, err:%s.", err)
	}

	symlink, err := filepath.EvalSymlinks(strings.Replace(dir, "\\", "/", -1))
	if err != nil {
		log.Printf("[config] get server symlink failed, err:%s.", err)
	}

	symlinks := strings.Split(symlink, "/")
	projects := strings.Split(symlinks[len(symlinks)-1], "-")
	if len(projects) < 2 {
		ver = defaultServerVer
	} else {
		ver = projects[len(projects)-2]
	}

	os.Setenv("ENV_SERVER_VER", ver)
	os.Setenv("ENV_SERVER_TIME", time.Now().Format("2006-01-02 15:04:05"))
}

func logPrint(k, v interface{}) {
	log.Printf("%-25q : %v", k, v)
}
