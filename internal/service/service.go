package service

import (
	"context"

	"github.com/google/wire"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	reloadconfig "github.com/comeonjy/go-kit/grpc/reloadconfig"
	"github.com/comeonjy/go-kit/pkg/xlog"
	v1 "github.com/comeonjy/working/api/v1"
	"github.com/comeonjy/working/configs"
	"github.com/comeonjy/working/internal/data"
)

var ProviderSet = wire.NewSet(NewWorkingService)

type WorkingService struct {
	v1.UnimplementedWorkingServer
	conf     configs.Interface
	logger   *xlog.Logger
	workRepo data.WorkRepo
	rcAccountSvc reloadconfig.ReloadConfigClient
}


func NewWorkingService(conf configs.Interface, logger *xlog.Logger, workRepo data.WorkRepo) *WorkingService {
	accountDial, err := grpc.Dial("account-grpc.jiangyang.me:80", grpc.WithInsecure())
	if err != nil {
		return nil
	}


	return &WorkingService{
		conf:     conf,
		workRepo: workRepo,
		logger:   logger,
		rcAccountSvc: reloadconfig.NewReloadConfigClient(accountDial),
	}
}

func (svc *WorkingService) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	if mdIn, ok := metadata.FromIncomingContext(ctx); ok {
		mdIn.Get("")
	}
	return ctx, nil
}

func (svc *WorkingService) Ping(ctx context.Context, in *v1.Empty) (*v1.Result, error) {
	return &v1.Result{
		Code:    200,
		Message: svc.conf.Get().Mode,
	}, nil
}