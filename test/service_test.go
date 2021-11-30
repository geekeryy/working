package test

import (
	"context"
	"log"
	"sync/atomic"
	"testing"

	"github.com/comeonjy/go-kit/grpc/reloadconfig"
	"google.golang.org/grpc"
)

type A struct {
	name string
}

type Interface interface {
	Store(val interface{})
	Load() (val interface{})
}
type inner struct {
	value atomic.Value
}

func (i *inner) Store(val interface{}) {
	i.value.Store(val)
}

func (i *inner) Load() (val interface{}) {
	return i.value.Load()
}

func TestService_Ping(t *testing.T) {
	a:=A{name: "a"}
	i:=inner{}
	i.Store(a)
	a2 := i.Load().(A)
	log.Println(a2.name)
}


func TestReloadConfig(t *testing.T) {
	dial, err := grpc.Dial("localhost:8081", grpc.WithInsecure())
	if err != nil {
		t.Error(err)
		return
	}
	client := reloadconfig.NewReloadConfigClient(dial)
	_, err = client.ReloadConfig(context.TODO(), &reloadconfig.Empty{})
	if err != nil {
		t.Error(err)
	}
}
