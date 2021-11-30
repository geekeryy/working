package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/comeonjy/go-kit/pkg/xenv"
	"github.com/comeonjy/go-kit/pkg/xsync"
	"github.com/comeonjy/working/configs"
	"google.golang.org/grpc"
)

type App struct {
	ctx  context.Context
	grpc *grpc.Server
	http *http.Server
	conf configs.Interface
}

func newApp( ctx context.Context,grpc *grpc.Server, http *http.Server, conf configs.Interface) *App {
	return &App{
		grpc: grpc,
		http: http,
		conf: conf,
		ctx:  ctx,
	}
}

func (app *App) Run(cancel context.CancelFunc) error {
	var g xsync.Group
	g.Go(func(ctx context.Context) error {
		return app.runGrpc()
	})
	g.Go(func(ctx context.Context) error {
		return app.runHttp()
	})
	if xenv.IsDebug(app.conf.Get().Mode) {
		g.Go(func(ctx context.Context) error {
			return app.runPprof()
		})
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGKILL)
	for {
		select {
		case sig, _ := <-quit:
			log.Println("signal stop ...")
			app.grpc.GracefulStop()
			log.Println("grpc closed")
			_ = app.http.Shutdown(app.ctx)
			log.Println("http closed")
			cancel()
			return errors.New(fmt.Sprintf("%v", sig))
		}
	}
}

func (app *App) runHttp() error {
	log.Printf("http run success in %s \n", app.http.Addr)
	return app.http.ListenAndServe()
}

func (app *App) runGrpc() error {
	listen, err := net.Listen("tcp", app.conf.Get().GrpcAddr)
	if err != nil {
		return err
	}

	log.Printf("grpc run success in %s \n", app.conf.Get().GrpcAddr)
	return app.grpc.Serve(listen)
}

func (app *App) runPprof() error {
	s := http.Server{
		Addr:    app.conf.Get().PprofAddr,
		Handler: http.DefaultServeMux,
	}
	log.Printf("pprof run success in %s \n", app.conf.Get().PprofAddr)
	return s.ListenAndServe()
}
