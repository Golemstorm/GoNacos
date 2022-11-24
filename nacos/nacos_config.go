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

func NewClientConfigForACM() constant.ClientConfig {
	return constant.ClientConfig{
		Endpoint:    "acm.aliyun.com:8080",
		NamespaceId: "e525eafa-f7d7-XXXX-XXXX-XXXXXXXX",
		RegionId:    "cn-shanghai",
		AccessKey:   "LTAI4G8KxxxxxxxxxxxxxbwZLBr",
		SecretKey:   "n5jTL9YxxxxxxxxxxxxaxmPLZV9",
		OpenKMS:     true,
		TimeoutMs:   5000,
		LogLevel:    "debug",
	}
}
