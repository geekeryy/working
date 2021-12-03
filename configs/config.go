package configs

import (
	"context"
	"encoding/json"
	"log"

	"github.com/comeonjy/go-kit/pkg/xconfig"
	"github.com/comeonjy/go-kit/pkg/xconfig/apollo"
	"github.com/comeonjy/go-kit/pkg/xenv"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewConfig)

// Interface 对外暴露接口（用于功能扩展）
type Interface interface {
	Get() Config
	xconfig.ReloadConfigInterface
}

// 内部配置类载体
type config struct {
	xconfig.IConfig
}

// Get 获取配置
func (c *config) Get() Config {
	return c.LoadValue().(Config)
}

func (c *config) ReloadConfig() error {
	return c.ReLoad()
}

func storeHandler(data []byte) interface{} {
	conf := Config{}
	if err := json.Unmarshal(data, &conf); err != nil {
		log.Println(err)
		return nil
	}
	return conf
}

// NewConfig 获取配置
func NewConfig(ctx context.Context) Interface {
	cfg := &config{
		xconfig.New(ctx,
			apollo.NewSource(xenv.GetEnv(xenv.ApolloUrl), xenv.GetEnv(xenv.ApolloAppID), xenv.GetApolloCluster("default"), xenv.GetApolloSecret(), xenv.GetApolloNamespace("grpc"),xenv.GetApolloNamespace("common")),
			storeHandler),
	}
	if err := cfg.ReLoad(); err != nil {
		panic(err)
	}
	return cfg
}
