package server

import (
	"context"
	"net/http"
	"time"

	"github.com/comeonjy/go-kit/pkg/xlog"
	"github.com/comeonjy/go-kit/pkg/xmiddleware"
	"github.com/comeonjy/working/api/v1"
	"github.com/comeonjy/working/configs"
	"github.com/comeonjy/working/internal/service"
	"github.com/comeonjy/working/pkg/consts"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func NewHttpServer(ctx context.Context, conf configs.Interface, logger *xlog.Logger ,workingService *service.WorkingService) *http.Server {
	mux := runtime.NewServeMux(runtime.WithErrorHandler(xmiddleware.HttpErrorHandler(logger)))
	server := http.Server{
		Addr:              conf.Get().HttpAddr,
		Handler:           xmiddleware.HttpUse(mux, xmiddleware.HttpLogger(consts.TraceName, logger)),
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      2 * time.Second,
	}
	Router(mux,workingService)
	if err := v1.RegisterWorkingHandlerFromEndpoint(ctx, mux, conf.Get().GrpcAddr, []grpc.DialOption{grpc.WithInsecure()}); err != nil {
		panic("RegisterSchedulerHandlerFromEndpoint" + err.Error())
	}
	return &server
}

func Router(mux *runtime.ServeMux,svc *service.WorkingService) {
	AddRouter(mux,"POST","/github-event",svc.GithubEvent)
	AddRouter(mux,"POST","/apollo-event",svc.ApolloEvent)
}

func AddRouter(mux *runtime.ServeMux, meth string, pathPattern string, h runtime.HandlerFunc) {
	if err := mux.HandlePath(meth, pathPattern, h); err != nil {
		panic(err)
	}
}
