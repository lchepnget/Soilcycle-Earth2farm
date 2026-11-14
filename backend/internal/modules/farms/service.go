package farms

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrInvalidFarmInput = errors.New("invalid farm input")
	ErrInvalidFarmID    = errors.New("invalid farm ID")
	ErrInvalidFarmerID  = errors.New("invalid farmer ID")
)

type FarmService struct {
	repo *FarmRepository
}

func NewFarmService(repo *FarmRepository) *FarmService {
	return &FarmService{repo: repo}
}

type FarmListResult struct {
	Data       []Farm `json:"data"`
	Total      int64  `json:"total"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	TotalPages int64  `json:"total_pages"`
}

func (s *FarmService) Create(
	ctx context.Context,
	req CreateFarmRequest,
) (*Farm, error) {

	if err := req.Validate(); err != nil {
		return nil, errors.Join(ErrInvalidFarmInput, err)
	}

	return s.repo.Create(ctx, req)
}

func (s *FarmService) GetByID(
	ctx context.Context,
	id string,
) (*Farm, error) {

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrInvalidFarmID
	}

	return s.repo.GetByID(ctx, parsedID)
}

func (s *FarmService) List(
	ctx context.Context,
	userID string,
	page, limit int,
) (*FarmListResult, error) {

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	var ownerID *uuid.UUID

	if userID != "" {
		id, err := uuid.Parse(userID)
		if err != nil {
			return nil, ErrInvalidFarmerID
		}
		ownerID = &id
	}

	offset := (page - 1) * limit

	farms, total, err := s.repo.List(
		ctx,
		ownerID,
		limit,
		offset,
	)

	if err != nil {
		return nil, err
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	return &FarmListResult{
		Data:       farms,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *FarmService) Update(
	ctx context.Context,
	id, userID string,
	req UpdateFarmRequest,
) (*Farm, error) {

	farmID, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrInvalidFarmID
	}

	ownerID, err := uuid.Parse(userID)
	if err != nil {
		return nil, ErrInvalidFarmerID
	}
	if err := req.Validate(); err != nil {
		return nil, errors.Join(ErrInvalidFarmInput, err)
	}

	return s.repo.Update(ctx, farmID, ownerID, req)
}

func (s *FarmService) Delete(
	ctx context.Context,
	id, userID string,
) error {

	farmID, err := uuid.Parse(id)
	if err != nil {
		return ErrInvalidFarmID
	}

	ownerID, err := uuid.Parse(userID)
	if err != nil {
		return ErrInvalidFarmerID
	}

	return s.repo.Delete(ctx, farmID, ownerID)
}
