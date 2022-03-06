// Package xgrpc @Description  TODO
// @Author  	 jiangyang
// @Created  	 2022/3/6 11:00 上午
package xgrpc

import (
	"context"

	"github.com/sercand/kuberesolver/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

const schema = "kubernetes"

func init() {
	resolver.Register(kuberesolver.NewBuilder(nil, schema))
}

func DialContext(ctx context.Context, target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	var conn *grpc.ClientConn
	var err error
	opts = append(opts, grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`), grpc.WithInsecure())
	conn, err = grpc.DialContext(ctx, target, opts...)
	if err != nil {
		conn, err = grpc.DialContext(ctx, schema+":///"+target, opts...)
	}
	return conn, err
}
