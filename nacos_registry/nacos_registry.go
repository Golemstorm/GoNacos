package nacos_registry

import (
	"context"
	"errors"
	mnet "github.com/micro/go-micro/v2/util/net"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"net"
	"strconv"
	"time"

	"github.com/micro/go-micro/v2/config/cmd"
	"github.com/micro/go-micro/v2/registry"
)

type nacosRegistry struct {
	namingClient naming_client.INamingClient
	iClient      config_client.IConfigClient
	opts         registry.Options
	clientConfig *constant.ClientConfig
}

func getNodeIpPort(s *registry.Service) (host string, port int, err error) {
	if len(s.Nodes) == 0 {
		return "", 0, errors.New("you must deregister at least one node")
	}
	node := s.Nodes[0]
	host, pt, err := net.SplitHostPort(node.Address)
	if err != nil {
		return "", 0, err
	}
	port, err = strconv.Atoi(pt)
	if err != nil {
		return "", 0, err
	}
	return
}
func (n *nacosRegistry) Init(option ...registry.Option) error {
	return configure(n, option...)
}

func (n *nacosRegistry) Options() registry.Options {
	return n.opts
}

func (n *nacosRegistry) Register(service *registry.Service, option ...registry.RegisterOption) error {
	var options registry.RegisterOptions
	for _, o := range option {
		o(&options)
	}
	withContext := false
	param := vo.RegisterInstanceParam{}
	if options.Context != nil {
		if p, ok := options.Context.Value("register_instance_param").(vo.RegisterInstanceParam); ok {
			param = p
			withContext = ok
		}
	}
	if !withContext {
		host, port, err := getNodeIpPort(service)
		if err != nil {
			return err
		}
		service.Nodes[0].Metadata["version"] = service.Version
		param.Ip = host
		param.Port = uint64(port)
		param.Metadata = service.Nodes[0].Metadata
		param.ServiceName = service.Name
		param.Enable = true
		param.Healthy = true
		param.Weight = 1.0
		param.Ephemeral = true
	}
	_, err := n.namingClient.RegisterInstance(param)
	return err
}

func (n *nacosRegistry) Deregister(service *registry.Service, option ...registry.DeregisterOption) error {
	var options registry.DeregisterOptions
	for _, o := range option {
		o(&options)
	}
	withContext := false
	param := vo.DeregisterInstanceParam{}
	if options.Context != nil {
		if p, ok := options.Context.Value("deregister_instance_param").(vo.DeregisterInstanceParam); ok {
			param = p
			withContext = ok
		}
	}
	if !withContext {
		host, port, err := getNodeIpPort(service)
		if err != nil {
			return err
		}
		param.Ip = host
		param.Port = uint64(port)
		param.ServiceName = service.Name
	}

	_, err := n.namingClient.DeregisterInstance(param)
	return err
}

func (n *nacosRegistry) GetService(name string, option ...registry.GetOption) ([]*registry.Service, error) {
	var options registry.GetOptions
	for _, o := range option {
		o(&options)
	}
	withContext := false
	param := vo.GetServiceParam{}
	if options.Context != nil {
		if p, ok := options.Context.Value("select_instances_param").(vo.GetServiceParam); ok {
			param = p
			withContext = ok
		}
	}
	if !withContext {
		param.ServiceName = name
	}
	service, err := n.namingClient.GetService(param)
	if err != nil {
		return nil, err
	}
	services := make([]*registry.Service, 0)
	for _, v := range service.Hosts {
		nodes := make([]*registry.Node, 0)
		nodes = append(nodes, &registry.Node{
			Id:       v.InstanceId,
			Address:  mnet.HostPort(v.Ip, v.Port),
			Metadata: v.Metadata,
		})
		s := registry.Service{
			Name:     v.ServiceName,
			Version:  v.Metadata["version"],
			Metadata: v.Metadata,
			Nodes:    nodes,
		}
		services = append(services, &s)
	}

	return services, nil
}

func (n *nacosRegistry) ListServices(option ...registry.ListOption) ([]*registry.Service, error) {
	var options registry.ListOptions
	for _, o := range option {
		o(&options)
	}
	withContext := false
	param := vo.GetAllServiceInfoParam{}
	if options.Context != nil {
		if p, ok := options.Context.Value("get_all_service_info_param").(vo.GetAllServiceInfoParam); ok {
			param = p
			withContext = ok
		}
	}
	if !withContext {
		services, err := n.namingClient.GetAllServicesInfo(param)
		if err != nil {
			return nil, err
		}
		param.PageNo = 1
		param.PageSize = uint32(services.Count)
	}
	services, err := n.namingClient.GetAllServicesInfo(param)
	if err != nil {
		return nil, err
	}
	var registryServices []*registry.Service
	for _, v := range services.Doms {
		registryServices = append(registryServices, &registry.Service{Name: v})
	}
	return registryServices, nil
}

func (n *nacosRegistry) Watch(option ...registry.WatchOption) (registry.Watcher, error) {
	return NewNacosWatcher(n, option...)
}

func (n *nacosRegistry) String() string {
	return "go-nacos_config"
}

func init() {
	cmd.DefaultRegistries["nacos_config"] = NewRegistry
}

func configure(c *nacosRegistry, opts ...registry.Option) error {
	// set opts
	for _, o := range opts {
		o(&c.opts)
	}
	if c.opts.Context != nil {
		if client, ok := c.opts.Context.Value("naming_client").(naming_client.INamingClient); ok {
			c.namingClient = client
			return nil
		}
	}
	serverConfigs := make([]constant.ServerConfig, 0)
	contextPath := "/nacos_config"
	// iterate the options addresses
	for _, address := range c.opts.Addrs {
		// check we have a port
		addr, port, err := net.SplitHostPort(address)
		if ae, ok := err.(*net.AddrError); ok && ae.Err == "missing port in address" {
			serverConfigs = append(serverConfigs, constant.ServerConfig{
				IpAddr:      addr,
				Port:        8848,
				ContextPath: contextPath,
			})
		} else if err == nil {
			p, err := strconv.ParseUint(port, 10, 64)
			if err != nil {
				continue
			}
			serverConfigs = append(serverConfigs, constant.ServerConfig{
				IpAddr:      addr,
				Port:        p,
				ContextPath: contextPath,
			})
		}
	}

	if c.opts.Timeout == 0 {
		c.opts.Timeout = time.Second * 1
	}
	c.clientConfig = nacosConfigTmp
	c.clientConfig.TimeoutMs = uint64(c.opts.Timeout.Milliseconds())

	client, err := clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig:  c.clientConfig,
		ServerConfigs: serverConfigs,
	})
	if err != nil {
		return err
	}

	c.namingClient = client
	return nil
}

var nacosConfigTmp = &constant.ClientConfig{}

func NewRegistry(opts ...registry.Option) registry.Registry {
	nacos := &nacosRegistry{
		opts: registry.Options{
			Context: context.Background(),
		},
		clientConfig: nacosConfigTmp,
	}
	configure(nacos, opts...)
	return nacos
}

func SetNamespaceId(namespaceId string) {
	nacosConfigTmp.NamespaceId = namespaceId
}
func SetClientConfigForAcm(endpoint, namespaceId, regionId, ak, sk string) {
	nacosConfigTmp = &constant.ClientConfig{
		Endpoint:    endpoint,
		NamespaceId: namespaceId,
		RegionId:    regionId,
		AccessKey:   ak,
		SecretKey:   sk,
		OpenKMS:     true,
		TimeoutMs:   5000,
		LogLevel:    "debug",
	}
}
func NewClientConfigForACM(endpoint, namespaceId, regionId, ak, sk string) constant.ClientConfig {
	return constant.ClientConfig{
		Endpoint:    endpoint,
		NamespaceId: namespaceId,
		RegionId:    regionId,
		AccessKey:   ak,
		SecretKey:   sk,
		OpenKMS:     true,
		TimeoutMs:   5000,
		LogLevel:    "debug",
	}
}
