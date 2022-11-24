package nacos_config

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

var configClient = &config_client.ConfigClient{}

func InitConfigClient(clientConfig *constant.ClientConfig, serverConfigs []constant.ServerConfig) {
	nConfigClient, err := clients.NewConfigClient(vo.NacosClientParam{
		ClientConfig:  clientConfig,
		ServerConfigs: serverConfigs,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	configClient = nConfigClient.(*config_client.ConfigClient)
	return
}

func GetConfig(dataId, group string) (configStr string, err error) {
	return configClient.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
}

func PublishConfig(dataId, group, content string) (success bool, err error) {
	return configClient.PublishConfig(vo.ConfigParam{
		DataId:  dataId,
		Group:   group,
		Content: content,
	})
}

func ListenConfig(dataId, group string, OnChange func(namespace, group, dataId, data string)) error {
	return configClient.ListenConfig(vo.ConfigParam{
		DataId:   dataId,
		Group:    group,
		OnChange: OnChange,
	})
}

func CancelListenConfig(dataId, group string) error {
	return configClient.CancelListenConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
}

func SearchConfigBlur(dataId, group string, pageNo, pageSize int, tag, appName string) (*model.ConfigPage, error) {
	return configClient.SearchConfig(vo.SearchConfigParam{
		Search:   "blur",
		DataId:   dataId,
		Group:    group,
		Tag:      tag,
		AppName:  appName,
		PageNo:   pageNo,
		PageSize: pageSize,
	})
}
