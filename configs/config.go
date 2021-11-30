package configs

import (
	"context"

	"github.com/comeonjy/go-kit/pkg/xconfig"
	"github.com/comeonjy/go-kit/pkg/xconfig/apollo"
	"github.com/comeonjy/go-kit/pkg/xenv"
	"github.com/comeonjy/working/pkg/consts"
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

// ReloadConfig 实现重载配置 xconfig.ReloadConfigInterface
func (c *config) ReloadConfig() error {
	if err := c.Load(); err != nil {
		return err
	}
	var tempConf Config
	if err := c.Scan(&tempConf); err != nil {
		return err
	}
	if err := tempConf.Validate(); err != nil {
		return err
	}
	c.StoreValue(tempConf)
	return nil
}

// NewConfig 获取配置
func NewConfig(ctx context.Context) Interface {
	cfg := &config{
		xconfig.New(
			xconfig.WithContext(ctx),
			xconfig.WithSource(apollo.NewSource(xenv.GetEnv(consts.ApolloUrl), consts.ApolloAppID, consts.ApolloCluster, consts.ApolloNamespace, xenv.GetEnv(consts.ApolloSecret))),
		),
	}
	if err := cfg.ReloadConfig(); err != nil {
		panic(err)
	}
	return cfg
}
