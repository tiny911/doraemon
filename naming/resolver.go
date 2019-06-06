package naming

import (
	consul "github.com/hashicorp/consul/api"
	"google.golang.org/grpc/naming"
)

// ConsulResolver is the implementaion of grpc.naming.Resolver
type ConsulResolver struct {
	ServiceName string //service name
	DataCenter  string
}

// NewResolver return ConsulResolver with service name
func NewResolver(name, env, datacenter string) *ConsulResolver {
	return &ConsulResolver{
		ServiceName: serviceName(name, env),
		DataCenter:  datacenter,
	}
}

// Resolve to resolve the service from consul, target is the dial address of consul
func (cr *ConsulResolver) Resolve(target string) (naming.Watcher, error) {
	// generate consul client, return if error
	conf := &consul.Config{
		Scheme:     "http",
		Address:    target,
		Datacenter: cr.DataCenter,
	}
	client, err := consul.NewClient(conf)
	if err != nil {
		return nil, err
	}

	// return ConsulWatcher
	watcher := &ConsulWatcher{
		cr: cr,
		cc: client,
	}
	return watcher, nil
}
