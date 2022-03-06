// Package k8s @Description  TODO
// @Author  	 jiangyang
// @Created  	 2021/12/4 1:58 下午
package k8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func NewClient(kubeConfig string) (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error
	config, err = rest.InClusterConfig()
	if err != nil {
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
	}
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}
