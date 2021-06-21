package service

import (
	"context"

	"github.com/sakti/o11ygo/api"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func NewAdderService() (*AdderService, error) {
	return &AdderService{}, nil
}

type AdderService struct {
	api.UnimplementedAddersvcServer
}

func (a *AdderService) Add(ctx context.Context, in *api.AddRequest) (*api.AddReply, error) {
	tracer := otel.GetTracerProvider()
	_, span := tracer.Tracer("").Start(ctx, "add")
	span.SetAttributes(attribute.Int64("a", in.A), attribute.Int64("b", in.B))
	defer span.End()
	return &api.AddReply{Result: in.A + in.B}, nil
}
