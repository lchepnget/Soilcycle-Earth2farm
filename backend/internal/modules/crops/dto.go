package crops

import (
	"errors"
	"strings"
)

type CreateCropRequest struct {
	Name           string  `json:"name"`
	ScientificName *string `json:"scientific_name"`
	CropFamily     *string `json:"crop_family"`
	Category       *string `json:"category"`
	Description    *string `json:"description"`
}

type UpdateCropRequest struct {
	Name           string  `json:"name"`
	ScientificName *string `json:"scientific_name"`
	CropFamily     *string `json:"crop_family"`
	Category       *string `json:"category"`
	Description    *string `json:"description"`
}

type CropCharacteristicRequest struct {
	WaterRequirement   *string  `json:"water_requirement"`
	GrowthDurationDays *int     `json:"growth_duration_days"`
	NitrogenFixing     bool     `json:"nitrogen_fixing"`
	NutrientDemand     *string  `json:"nutrient_demand"`
	SuitableSeasons    []string `json:"suitable_seasons"`
}

func (r *CreateCropRequest) Validate() error {
	r.Name = strings.TrimSpace(r.Name)

	if r.Name == "" {
		return errors.New("crop name is required")
	}

	if len(r.Name) > 100 {
		return errors.New("crop name cannot exceed 100 characters")
	}

	return nil
}

func (r *UpdateCropRequest) Validate() error {
	r.Name = strings.TrimSpace(r.Name)

	if r.Name == "" {
		return errors.New("crop name is required")
	}

	if len(r.Name) > 100 {
		return errors.New("crop name cannot exceed 100 characters")
	}

	return nil
}

func (r *CropCharacteristicRequest) Validate() error {
	if r.GrowthDurationDays != nil && *r.GrowthDurationDays <= 0 {
		return errors.New("growth duration must be greater than zero")
	}

	if r.WaterRequirement != nil {
		value := strings.ToLower(*r.WaterRequirement)
		switch value {
		case "low", "medium", "high":
		default:
			return errors.New("water requirement must be low, medium or high")
		}
	}

	if r.NutrientDemand != nil {
		value := strings.ToLower(*r.NutrientDemand)
		switch value {
		case "low", "medium", "high":
		default:
			return errors.New("nutrient demand must be low, medium or high")
		}
	}

	return nil
}
