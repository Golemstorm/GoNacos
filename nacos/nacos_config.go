package nacos

import "github.com/nacos-group/nacos-sdk-go/v2/common/constant"

func NewServerConfig() constant.ServerConfig {
	return constant.ServerConfig{
		Scheme:      "",
		ContextPath: "",
		IpAddr:      "",
		Port:        0,
		GrpcPort:    0,
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
