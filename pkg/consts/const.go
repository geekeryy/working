package consts

import (
	"github.com/comeonjy/go-kit/pkg/xenv"
)

var EnvMap = map[string]string{
	xenv.AppName:     "working",
	xenv.AppVersion:  "v1.0",
	xenv.ApolloAppID: "working",
	xenv.ApolloUrl:   "http://apollo.dev.jiangyang.me",
	"kube_config":    "/Users/jiangyang/.kube/k8s",
	"images_repo":    "ccr.ccs.tencentyun.com/comeonjy",
}
