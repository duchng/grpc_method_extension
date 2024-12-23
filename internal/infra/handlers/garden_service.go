package handlers

import (
	"context"

	gardenservicev1 "grpc_method_extension/gen/api/garden-service/v1"
)

type GardenService struct {
	gardenservicev1.UnimplementedGardenServiceServer
}

func (s *GardenService) GetFlowers(_ context.Context, _ *gardenservicev1.GetFlowersRequest) (*gardenservicev1.GetFlowersResponse, error) {
	return &gardenservicev1.GetFlowersResponse{
		Flowers: []*gardenservicev1.Flower{
			{
				Name:  "Rose",
				Color: "Red",
			},
			{
				Name:  "Violet",
				Color: "Blue",
			},
			{
				Name:  "Lily",
				Color: "White",
			},
		},
	}, nil
}

func (s *GardenService) GetMushrooms(_ context.Context, _ *gardenservicev1.GetMushroomsRequest) (*gardenservicev1.GetMushroomsResponse, error) {
	return &gardenservicev1.GetMushroomsResponse{Mushrooms: []*gardenservicev1.Mushroom{
		{
			Name: "Shiitake",
			Size: 8,
		},
		{
			Name: "Truffle",
			Size: 1,
		},
		{
			Name: "Cremini",
			Size: 6,
		},
		{
			Name: "Shimeji",
			Size: 4,
		},
	}}, nil
}

func NewGardenService() gardenservicev1.GardenServiceServer {
	return &GardenService{}
}
