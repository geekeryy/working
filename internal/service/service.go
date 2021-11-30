package service

import (
	"context"

	"github.com/google/wire"
	"google.golang.org/grpc/metadata"

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
}

func NewWorkingService(conf configs.Interface, logger *xlog.Logger, workRepo data.WorkRepo) *WorkingService {
	return &WorkingService{
		conf:     conf,
		workRepo: workRepo,
		logger:   logger,
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