package crops

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
)

type CropService struct {
	repo *CropRepository
}

func NewCropService(repo *CropRepository) *CropService {
	return &CropService{repo: repo}
}

type CropListResult struct {
	Data       []Crop `json:"data"`
	Total      int64  `json:"total"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	TotalPages int64  `json:"total_pages"`
}

func (s *CropService) Create(ctx context.Context, req CreateCropRequest) (*Crop, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	return s.repo.Create(ctx, req)
}

func (s *CropService) GetByID(ctx context.Context, id string) (*Crop, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid crop ID")
	}

	return s.repo.GetByID(ctx, parsedID)
}

func (s *CropService) List(
	ctx context.Context,
	search, category string,
	page, limit int,
) (*CropListResult, error) {

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	search = strings.TrimSpace(search)
	category = strings.TrimSpace(category)

	offset := (page - 1) * limit

	crops, total, err := s.repo.List(
		ctx,
		search,
		category,
		limit,
		offset,
	)

	if err != nil {
		return nil, err
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	return &CropListResult{
		Data:       crops,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *CropService) Update(
	ctx context.Context,
	id string,
	req UpdateCropRequest,
) (*Crop, error) {

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid crop ID")
	}

	if err := req.Validate(); err != nil {
		return nil, err
	}

	return s.repo.Update(ctx, parsedID, req)
}

func (s *CropService) Delete(ctx context.Context, id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid crop ID")
	}

	return s.repo.Delete(ctx, parsedID)
}

func (s *CropService) GetCharacteristics(
	ctx context.Context,
	id string,
) (*CropCharacteristic, error) {

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid crop ID")
	}

	return s.repo.GetCharacteristics(ctx, parsedID)
}

func (s *CropService) UpsertCharacteristics(
	ctx context.Context,
	id string,
	req CropCharacteristicRequest,
) (*CropCharacteristic, error) {

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid crop ID")
	}

	if err := req.Validate(); err != nil {
		return nil, err
	}

	return s.repo.UpsertCharacteristics(ctx, parsedID, req)
}
