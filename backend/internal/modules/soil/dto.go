package soil

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type CreateSoilRequest struct {
	FarmID          string   `json:"farm_id"`
	PH              *float64 `json:"ph"`
	OrganicCarbon   *float64 `json:"organic_carbon"`
	Nitrogen        *float64 `json:"nitrogen"`
	Phosphorus      *float64 `json:"phosphorus"`
	Potassium       *float64 `json:"potassium"`
	SoilTexture     *string  `json:"soil_texture"`
	DataSource      string   `json:"data_source"`
	MeasurementDate *string  `json:"measurement_date"`
}

type UpdateSoilRequest struct {
	PH              *float64 `json:"ph"`
	OrganicCarbon   *float64 `json:"organic_carbon"`
	Nitrogen        *float64 `json:"nitrogen"`
	Phosphorus      *float64 `json:"phosphorus"`
	Potassium       *float64 `json:"potassium"`
	SoilTexture     *string  `json:"soil_texture"`
	DataSource      string   `json:"data_source"`
	MeasurementDate *string  `json:"measurement_date"`
}

func validateSoil(
	ph, organicCarbon, nitrogen, phosphorus, potassium *float64,
	soilTexture *string,
	dataSource string,
	measurementDate *string,
) error {
	if strings.TrimSpace(dataSource) == "" {
		return errors.New("data source is required")
	}

	if len(dataSource) > 50 {
		return errors.New("data source cannot exceed 50 characters")
	}

	if ph != nil && (*ph < 0 || *ph > 14) {
		return errors.New("soil pH must be between 0 and 14")
	}

	values := map[string]*float64{
		"organic carbon": organicCarbon,
		"nitrogen":       nitrogen,
		"phosphorus":     phosphorus,
		"potassium":      potassium,
	}

	for name, value := range values {
		if value != nil && *value < 0 {
			return errors.New(name + " cannot be negative")
		}
	}

	if soilTexture != nil && len(*soilTexture) > 100 {
		return errors.New("soil texture cannot exceed 100 characters")
	}

	if measurementDate != nil {
		if _, err := time.Parse("2006-01-02", *measurementDate); err != nil {
			return errors.New("measurement date must use YYYY-MM-DD format")
		}
	}

	return nil
}

func (r *CreateSoilRequest) Validate() error {
	if _, err := uuid.Parse(r.FarmID); err != nil {
		return errors.New("valid farm ID is required")
	}

	return validateSoil(
		r.PH,
		r.OrganicCarbon,
		r.Nitrogen,
		r.Phosphorus,
		r.Potassium,
		r.SoilTexture,
		r.DataSource,
		r.MeasurementDate,
	)
}

func (r *UpdateSoilRequest) Validate() error {
	return validateSoil(
		r.PH,
		r.OrganicCarbon,
		r.Nitrogen,
		r.Phosphorus,
		r.Potassium,
		r.SoilTexture,
		r.DataSource,
		r.MeasurementDate,
	)
}
