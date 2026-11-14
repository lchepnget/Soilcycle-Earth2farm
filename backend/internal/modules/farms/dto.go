package farms

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

type CreateFarmRequest struct {
	UserID            string   `json:"user_id"`
	Name              string   `json:"name"`
	County            string   `json:"county"`
	SubCounty         *string  `json:"sub_county"`
	Latitude          float64  `json:"latitude"`
	Longitude         float64  `json:"longitude"`
	AreaHectares      *float64 `json:"area_hectares"`
	WaterAvailability *string  `json:"water_availability"`
}

type UpdateFarmRequest struct {
	Name              string   `json:"name"`
	County            string   `json:"county"`
	SubCounty         *string  `json:"sub_county"`
	Latitude          float64  `json:"latitude"`
	Longitude         float64  `json:"longitude"`
	AreaHectares      *float64 `json:"area_hectares"`
	WaterAvailability *string  `json:"water_availability"`
}

func validateFarm(
	name, county string,
	latitude, longitude float64,
	area *float64,
	water *string,
) error {

	if strings.TrimSpace(name) == "" {
		return errors.New("farm name is required")
	}

	if len(name) > 150 {
		return errors.New("farm name cannot exceed 150 characters")
	}

	if strings.TrimSpace(county) == "" {
		return errors.New("county is required")
	}

	if latitude < -90 || latitude > 90 {
		return errors.New("latitude must be between -90 and 90")
	}

	if longitude < -180 || longitude > 180 {
		return errors.New("longitude must be between -180 and 180")
	}

	if area != nil && *area <= 0 {
		return errors.New("farm area must be greater than zero")
	}

	if water != nil {
		switch strings.ToLower(strings.TrimSpace(*water)) {
		case "low", "medium", "high", "seasonal", "none":
		default:
			return errors.New("water availability must be low, medium, high, seasonal or none")
		}
	}

	return nil
}

func (r *CreateFarmRequest) Validate() error {
	if _, err := uuid.Parse(r.UserID); err != nil {
		return errors.New("valid farmer ID is required")
	}

	return validateFarm(
		r.Name,
		r.County,
		r.Latitude,
		r.Longitude,
		r.AreaHectares,
		r.WaterAvailability,
	)
}

func (r *UpdateFarmRequest) Validate() error {
	return validateFarm(
		r.Name,
		r.County,
		r.Latitude,
		r.Longitude,
		r.AreaHectares,
		r.WaterAvailability,
	)
}
