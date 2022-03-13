package test

import (
	"context"
	"log"
	"sync/atomic"
	"testing"
	"time"

	"github.com/comeonjy/go-kit/grpc/reloadconfig"
	"github.com/comeonjy/working/pkg/xgrpc"
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
	a := A{name: "a"}
	i := inner{}
	i.Store(a)
	a2 := i.Load().(A)
	log.Println(a2.name)
}

func TestReloadConfig(t *testing.T) {
	dial, err := xgrpc.DialContext(context.Background(), "account.default:8081")
	if err != nil {
		t.Error(err)
		return
	}
	client := reloadconfig.NewReloadConfigClient(dial)
	t.Run("demo1", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			_, err = client.ReloadConfig(context.TODO(), &reloadconfig.Empty{})
			if err != nil {
				t.Error(err)
			}
			time.Sleep(time.Second)
		}

	})

}
