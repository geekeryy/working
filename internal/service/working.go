// Package service @Description  TODO
// @Author  	 jiangyang
// @Created  	 2021/11/30 10:24 下午
package service

import (
	"context"
	"fmt"

	v1 "github.com/comeonjy/working/api/v1"
)

func (svc *WorkingService) K8S(ctx context.Context, in *v1.Empty) (*v1.Result, error) {
	image := fmt.Sprintf("%s/%s:%s", svc.conf.Get().ImagesRepo, "form-online-web", "v0.0.1")
	err := svc.restartDeploy("form-online-web", image)
	if err != nil {
		return nil, err
	}

	return &v1.Result{
		Code:    0,
		Message: "",
		Data:    nil,
	}, nil
}
