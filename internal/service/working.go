// Package service @Description  TODO
// @Author  	 jiangyang
// @Created  	 2021/11/30 10:24 下午
package service

import (
	"context"
	"encoding/json"

	v1 "github.com/comeonjy/working/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (svc *WorkingService) K8S(ctx context.Context,in *v1.Empty) (*v1.Result, error) {

	list, err := svc.k8sClient.CoreV1().Pods("default").List(ctx, metav1.ListOptions{
		LabelSelector:        "app=account",
	})
	if err != nil {
		return nil, err
	}
	marshal, err := json.Marshal(list)
	if err != nil {
		return nil, err
	}

	return &v1.Result{
		Code:    0,
		Message: string(marshal),
		Data:    nil,
	}, nil
}