package service

import (
	"context"
	"log"

	"github.com/comeonjy/go-kit/pkg/xerror"
	"github.com/comeonjy/working/pkg/consts"
	"github.com/comeonjy/working/pkg/errcode"
	"github.com/comeonjy/working/pkg/k8s"
	"github.com/google/wire"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"k8s.io/client-go/kubernetes"

	reloadconfig "github.com/comeonjy/go-kit/grpc/reloadconfig"
	"github.com/comeonjy/go-kit/pkg/xlog"
	v1 "github.com/comeonjy/working/api/v1"
	"github.com/comeonjy/working/configs"
	"github.com/comeonjy/working/internal/data"
)

var ProviderSet = wire.NewSet(NewWorkingService)

type WorkingService struct {
	v1.UnimplementedWorkingServer
	conf         configs.Interface
	logger       *xlog.Logger
	workRepo     data.WorkRepo
	rcAccountSvc reloadconfig.ReloadConfigClient
	k8sClient    *kubernetes.Clientset
}

func NewWorkingService(conf configs.Interface, logger *xlog.Logger, workRepo data.WorkRepo) *WorkingService {
	accountDial, err := grpc.Dial(conf.Get().AccountGrpc, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client, err := k8s.NewClient(consts.EnvMap["kube_config"])
	if err != nil {
		panic(err)
	}

	return &WorkingService{
		conf:         conf,
		workRepo:     workRepo,
		logger:       logger,
		rcAccountSvc: reloadconfig.NewReloadConfigClient(accountDial),
		k8sClient:    client,
	}
}

func (svc *WorkingService) AuthFuncOverride(ctx context.Context, fullMethodName string) (context.Context, error) {
	if mdIn, ok := metadata.FromIncomingContext(ctx); ok {
		mdIn.Get("")
	}
	return ctx, nil
}

func (svc *WorkingService) Ping(ctx context.Context, in *v1.Empty) (*v1.Result, error) {
	accountDial, err := grpc.Dial(svc.conf.Get().AccountGrpc, grpc.WithInsecure())
	if err != nil {
		return &v1.Result{}, xerror.NewError(errcode.SystemErr, "", err.Error())
	}

	rc := reloadconfig.NewReloadConfigClient(accountDial)
	if _, err = rc.ReloadConfig(ctx, &reloadconfig.Empty{}); err != nil {
		log.Println("ReloadConfig", err.Error())
	}

	return &v1.Result{
		Code:    200,
		Message: svc.conf.Get().Mode,
	}, nil
}
