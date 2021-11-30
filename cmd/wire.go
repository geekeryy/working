//go:build wireinject
// +build wireinject

package cmd

import (
	"context"

	"github.com/comeonjy/go-kit/pkg/xlog"
	"github.com/comeonjy/working/configs"
	"github.com/google/wire"

	"github.com/comeonjy/working/internal/data"
	"github.com/comeonjy/working/internal/server"
	"github.com/comeonjy/working/internal/service"
)

func InitApp(ctx context.Context,logger *xlog.Logger) *App {
	panic(wire.Build(
		server.ProviderSet,
		service.ProviderSet,
		newApp,
		configs.ProviderSet,
		data.ProviderSet,
	))
}
