package farmers

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidFarmerInput = errors.New("invalid farmer input")

type FarmerService struct {
	repo *FarmerRepository
}

func NewFarmerService(repo *FarmerRepository) *FarmerService {
	return &FarmerService{repo: repo}
}

type FarmerListResult struct {
	Data       []Farmer `json:"data"`
	Total      int64    `json:"total"`
	Page       int      `json:"page"`
	Limit      int      `json:"limit"`
	TotalPages int64    `json:"total_pages"`
}

func (s *FarmerService) Create(
	ctx context.Context,
	req CreateFarmerRequest,
) (*Farmer, error) {

	if err := req.Validate(); err != nil {
		return nil, errors.Join(ErrInvalidFarmerInput, err)
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return nil, err
	}

	return s.repo.Create(ctx, req, string(hash))
}

func (s *FarmerService) GetByID(
	ctx context.Context,
	id string,
) (*Farmer, error) {

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid farmer ID")
	}

	return s.repo.GetByID(ctx, parsedID)
}

func (s *FarmerService) List(
	ctx context.Context,
	page, limit int,
) (*FarmerListResult, error) {

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

	farmers, total, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	return &FarmerListResult{
		Data:       farmers,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *FarmerService) Update(
	ctx context.Context,
	id string,
	req UpdateFarmerRequest,
) (*Farmer, error) {

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid farmer ID")
	}

	if err := req.Validate(); err != nil {
		return nil, errors.Join(ErrInvalidFarmerInput, err)
	}

	return s.repo.Update(ctx, parsedID, req)
}
