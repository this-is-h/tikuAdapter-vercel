package registry

import (
	"github.com/this-is-h/tikuAdapter-vercel/api/configs"
	"github.com/this-is-h/tikuAdapter-vercel/api/internal/service"
)

// RegisterEs 注册配置文件
func RegisterEs(cfg configs.Config) service.Elasticsearch {
	client, err := service.NewElasticsearchClient(cfg.Elasticsearch.Addresses)
	if err != nil {
		return nil
	}
	return client
}
