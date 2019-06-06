package naming

import (
	"fmt"
	"os"
	"time"

	"io/ioutil"
	"os/exec"

	consul "github.com/hashicorp/consul/api"
	"github.com/tiny911/doraemon/log"
)

const (
	enable  string = "on"
	disable string = "off"
)

var (
	defaultDep   = "tiny"
	defaultCop   = "tiny"
	defaultTags  = "doraemon"
	defaultItval = 10
	defaultTTL   = 15
)

// Naming 名字服务
type Naming struct {
	serviceName string
	serviceID   string
	serviceHost string

	servicePort int
	rpcPort     int
	httpPort    int

	rawName     string
	target      string
	datacenter  string
	environment string
	interval    int
	ttl         int
	cmd         string
	client      *consul.Client
	quit        chan interface{}
}

// New 生成naming实例
func New(name, host, env string, rpcPort, httpPort int) *Naming {
	n := &Naming{
		rawName:     name,
		serviceHost: switchIP(host, env),
		servicePort: rpcPort,
		rpcPort:     rpcPort,
		httpPort:    httpPort,

		environment: env,
		interval:    defaultItval,
		ttl:         defaultTTL,
		cmd:         fmt.Sprintf("top -b -n1 -p %d | sed -n '7,8p'", os.Getpid()),
		quit:        make(chan interface{}),
	}

	n.serviceName = serviceName(name, env)
	n.setServiceID()
	return n
}

// Regist 注册名字
func (n *Naming) Regist(target, datacenter string) error {
	if os.Getenv("ENV_NAMING_SWITCH") != enable {
		return nil
	}

	n.target = target
	n.datacenter = datacenter

	var err error

	conf := &consul.Config{Scheme: "http", Datacenter: n.datacenter, Address: n.target}
	n.client, err = consul.NewClient(conf)
	if err != nil {
		log.WithField(log.Fields{
			"error":       err,
			"serviceName": n.serviceName,
			"serviceID":   n.serviceID,
			"target":      n.target,
			"datacenter":  n.datacenter,
		}).Error("naming new client failed.")
		return err
	}

	go n.updateTTL()

	if err = n.registService(); err != nil {
		return err
	}

	if err = n.registCheck(); err != nil {
		return err
	}

	log.WithField(log.Fields{
		"serviceName": n.serviceName,
		"serviceID":   n.serviceID,
		"target":      n.target,
		"datacenter":  n.datacenter,
	}).Info("naming regist success.")
	return nil
}

// UnRegist 取消注册名字
func (n *Naming) UnRegist() error {
	if os.Getenv("ENV_NAMING_SWITCH") != enable {
		return nil
	}

	err := n.client.Agent().ServiceDeregister(n.serviceID)
	if err != nil {
		log.WithField(log.Fields{
			"error":       err,
			"serviceName": n.serviceName,
			"serviceID":   n.serviceID,
			"target":      n.target,
			"datacenter":  n.datacenter,
		}).Error("naming deregister service failed.")
	} else {
		log.WithField(log.Fields{
			"serviceName": n.serviceName,
			"serviceID":   n.serviceID,
			"target":      n.target,
			"datacenter":  n.datacenter,
		}).Info("naming deregister service success.")
	}

	err = n.client.Agent().CheckDeregister(n.serviceID)
	if err != nil {
		log.WithField(log.Fields{
			"error":       err,
			"serviceName": n.serviceName,
			"serviceID":   n.serviceID,
			"target":      n.target,
			"datacenter":  n.datacenter,
		}).Error("naming deregister check failed.")
	} else {
		log.WithField(log.Fields{
			"serviceName": n.serviceName,
			"serviceID":   n.serviceID,
			"target":      n.target,
			"datacenter":  n.datacenter,
		}).Info("naming deregister check success.")

	}

	close(n.quit)
	return err
}

func serviceName(name string, env string) string {
	serviceName := name + "."
	if env != "" {
		serviceName += env + "."
	}
	serviceName += defaultDep + "."
	serviceName += defaultCop
	return serviceName
}

func (n *Naming) setServiceID() {
	n.serviceID = fmt.Sprintf("%s:%d", n.serviceHost, n.servicePort)
}

func (n *Naming) getOutput() string {
	cmd := exec.Command("/bin/sh", "-c", n.cmd)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err.Error()
	}

	if err := cmd.Start(); err != nil {
		return err.Error()
	}

	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		return err.Error()
	}

	if err := cmd.Wait(); err != nil {
		return err.Error()
	}

	cut := "-------------------------------------------------------------------------------\n"
	time := fmt.Sprintf("TIME:%s", time.Now().Format("2006-01-02 15:04:05"))
	return string(bytes) + cut + time
}

func (n *Naming) registCheck() error {
	check := &consul.AgentCheckRegistration{
		ID:        n.serviceID,
		Name:      n.serviceName,
		ServiceID: n.serviceID,
		Notes:     fmt.Sprintf("Gather service:%s process information.", n.rawName),
		AgentServiceCheck: consul.AgentServiceCheck{
			TTL:    fmt.Sprintf("%ds", n.ttl),
			Status: "passing",
		},
	}
	err := n.client.Agent().CheckRegister(check)
	if err != nil {
		log.WithField(log.Fields{
			"error":       err,
			"serviceName": n.serviceName,
			"serviceID":   n.serviceID,
			"target":      n.target,
			"datacenter":  n.datacenter,
		}).Error("naming regist check failed.")
		return err
	}
	return nil
}

func (n *Naming) registService() error {
	var err error

	tags := []string{defaultTags}
	{
		tags = append(
			tags,
			fmt.Sprintf("rpc:%d", n.rpcPort),
			fmt.Sprintf("http:%d", n.httpPort),
		)

		if os.Getenv("ENV_SERVER_VER") != "" {
			tags = append(tags, os.Getenv("ENV_SERVER_VER"))
		}

		if os.Getenv("ENV_SERVER_TIME") != "" {
			tags = append(tags, os.Getenv("ENV_SERVER_TIME"))
		}
	}

	regis := &consul.AgentServiceRegistration{
		ID:      n.serviceID,
		Name:    n.serviceName,
		Address: n.serviceHost,
		Port:    n.servicePort,
		Tags:    tags,
	}

	err = n.client.Agent().ServiceRegister(regis)
	if err != nil {
		log.WithField(log.Fields{
			"error":       err,
			"serviceName": n.serviceName,
			"serviceID":   n.serviceID,
			"target":      n.target,
			"datacenter":  n.datacenter,
		}).Error("naming regist service failed.")
		return err
	}
	return nil
}

func (n *Naming) updateTTL() {
	var (
		err    error
		ticker = time.NewTicker(time.Duration(n.interval) * time.Second)
	)

	for {
		select {
		case <-ticker.C:
			err = n.client.Agent().UpdateTTL(n.serviceID, n.getOutput(), "passing")
			if err != nil {
				log.WithField(log.Fields{
					"error":       err,
					"serviceName": n.serviceName,
					"serviceID":   n.serviceID,
					"target":      n.target,
					"datacenter":  n.datacenter,
				}).Error("naming update ttl failed.")
			}
		case <-n.quit:
			return
		}
	}
}
