package rotation

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

type CreateRotationRequest struct {
	FarmID           string   `json:"farm_id"`
	FarmerPriorities []string `json:"farmer_priorities"`
	NumberOfSeasons  int      `json:"number_of_seasons"`
}

func (r *CreateRotationRequest) Validate() error {
	if r == nil {
		return errors.New("request body is required")
	}

	r.FarmID = strings.TrimSpace(r.FarmID)

	if _, err := uuid.Parse(r.FarmID); err != nil {
		return errors.New("valid farm ID is required")
	}

	if r.NumberOfSeasons < 1 || r.NumberOfSeasons > 10 {
		return errors.New("number of seasons must be between 1 and 10")
	}

	if len(r.FarmerPriorities) > 10 {
		return errors.New("maximum of 10 farmer priorities allowed")
	}

	for i, priority := range r.FarmerPriorities {
		priority = strings.TrimSpace(priority)

		if priority == "" {
			return errors.New("farmer priorities cannot contain empty values")
		}

		r.FarmerPriorities[i] = priority
	}

	return nil
}

// RotationResponse represents the result of generating a rotation plan.
type RotationResponse struct {
	Plan    *RotationPlan    `json:"plan"`
	Results []RotationResult `json:"results"`
}

// RotationPlanListResponse represents a list of saved plans.
type RotationPlanListResponse struct {
	Plans []RotationPlan `json:"plans"`
	Total int            `json:"total"`
}

// APIResponse provides a consistent API response structure.
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
