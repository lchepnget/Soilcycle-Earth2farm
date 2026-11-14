package earthdata

import (
	"context"
	"fmt"
)

type Service struct {
	client *POWERClient
}

func NewService(client *POWERClient) *Service {
	return &Service{
		client: client,
	}
}

// GetClimateData retrieves climate observations from NASA POWER.
func (s *Service) GetClimateData(
	ctx context.Context,
	req EarthdataRequest,
) (*EarthdataResponse, error) {

	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	observations, err := s.client.FetchDaily(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("fetch climate data: %w", err)
	}

	response := &EarthdataResponse{
		FarmID: req.FarmID,
		Location: FarmLocation{
			Latitude:  req.Latitude,
			Longitude: req.Longitude,
		},
		Source:       nasaPOWERSource,
		Observations: observations,
		TotalRecords: len(observations),
	}

	return response, nil
}
