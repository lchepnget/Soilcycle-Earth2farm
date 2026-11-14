package soil

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var ErrInvalidSoilInput = errors.New("invalid soil input")

type SoilService struct {
	repo *SoilRepository
}

func NewSoilService(repo *SoilRepository) *SoilService {
	return &SoilService{repo: repo}
}

type SoilListResult struct {
	Data       []SoilData `json:"data"`
	Total      int64      `json:"total"`
	Page       int        `json:"page"`
	Limit      int        `json:"limit"`
	TotalPages int64      `json:"total_pages"`
}

// Create a soil record.
func (s *SoilService) Create(
	ctx context.Context,
	req CreateSoilRequest,
) (*SoilData, error) {
	if err := req.Validate(); err != nil {
		return nil, errors.Join(ErrInvalidSoilInput, err)
	}

	return s.repo.Create(ctx, req)
}

// Retrieve a soil record by ID.
func (s *SoilService) GetByID(
	ctx context.Context,
	id string,
) (*SoilData, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid soil record ID")
	}

	return s.repo.GetByID(ctx, parsedID)
}

// List all soil records.
func (s *SoilService) List(
	ctx context.Context,
	page, limit int,
) (*SoilListResult, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	records, total, err := s.repo.List(
		ctx,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	return &SoilListResult{
		Data:       records,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// List soil records for a specific farm.
func (s *SoilService) ListByFarm(
	ctx context.Context,
	farmID string,
	page, limit int,
) (*SoilListResult, error) {
	parsedFarmID, err := uuid.Parse(farmID)
	if err != nil {
		return nil, errors.New("invalid farm ID")
	}

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	records, total, err := s.repo.ListByFarm(
		ctx,
		parsedFarmID,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	return &SoilListResult{
		Data:       records,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// Update a soil record.
func (s *SoilService) Update(
	ctx context.Context,
	id string,
	req UpdateSoilRequest,
) (*SoilData, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid soil record ID")
	}

	if err := req.Validate(); err != nil {
		return nil, errors.Join(ErrInvalidSoilInput, err)
	}

	return s.repo.Update(ctx, parsedID, req)
}

// Delete a soil record.
func (s *SoilService) Delete(
	ctx context.Context,
	id string,
) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid soil record ID")
	}

	return s.repo.Delete(ctx, parsedID)
}
