package client

import (
	"fmt"

	"github.com/sakti/o11ygo/api"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

func NewAuthClient(port int) (api.AuthsvcClient, error) {
	conn, err := grpc.Dial(fmt.Sprintf(":%d", port), grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		return nil, err
	}
	return api.NewAuthsvcClient(conn), nil
}

func NewAdderClient(port int) (api.AddersvcClient, error) {
	conn, err := grpc.Dial(fmt.Sprintf(":%d", port), grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		return nil, err
	}
	return api.NewAddersvcClient(conn), nil
}

func NewMultiplyClient(port int) (api.MultiplysvcClient, error) {
	conn, err := grpc.Dial(fmt.Sprintf(":%d", port), grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		return nil, err
	}
	return api.NewMultiplysvcClient(conn), nil
}
