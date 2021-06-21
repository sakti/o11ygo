package service

import (
	"context"

	"github.com/sakti/o11ygo/api"
	"github.com/sakti/o11ygo/pkg/client"
)

func NewMultiplyService(adderPort int) (*MultiplyService, error) {
	adderClient, err := client.NewAdderClient(adderPort)
	if err != nil {
		return nil, err
	}
	return &MultiplyService{
		adderClient: adderClient,
	}, nil
}

type MultiplyService struct {
	api.UnimplementedMultiplysvcServer
	adderClient api.AddersvcClient
}

func (m *MultiplyService) Multiply(ctx context.Context, in *api.MultiplyRequest) (*api.MultiplyReply, error) {
	var sum int64
	for i := 0; i < int(in.B); i++ {
		response, err := m.adderClient.Add(ctx, &api.AddRequest{A: in.A, B: sum})
		if err != nil {
			return nil, err
		}
		sum = response.Result
	}
	return &api.MultiplyReply{Result: sum}, nil
}
