package data

import (
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/comeonjy/working/configs"
)

var ProviderSet = wire.NewSet(NewData, NewWorkRepo)

type Data struct {
	Mongo *mongo.Collection
}

func NewData(cfg configs.Interface) *Data {
	return &Data{}
}
