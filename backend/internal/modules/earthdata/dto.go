package earthdata

import (
	"errors"
	"strings"
	"time"
)

type EarthdataRequest struct {
	FarmID    string  `json:"farm_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	StartDate string  `json:"start_date"`
	EndDate   string  `json:"end_date"`
}

type FarmLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type EarthdataResponse struct {
	FarmID       string               `json:"farm_id"`
	Location     FarmLocation         `json:"location"`
	Source       string               `json:"source"`
	Observations []ClimateObservation `json:"observations"`
	TotalRecords int                  `json:"total_records"`
}

type ClimateObservation struct {
	Date            time.Time `json:"date"`
	RainfallMM      *float64  `json:"rainfall_mm,omitempty"`
	TemperatureC    *float64  `json:"temperature_c,omitempty"`
	MaxTemperatureC *float64  `json:"max_temperature_c,omitempty"`
	MinTemperatureC *float64  `json:"min_temperature_c,omitempty"`
	HumidityPercent *float64  `json:"humidity_percent,omitempty"`
	Source          string    `json:"source"`
}

func (r *EarthdataRequest) Validate() error {
	if strings.TrimSpace(r.FarmID) == "" {
		return errors.New("farm ID is required")
	}

	if r.Latitude < -90 || r.Latitude > 90 {
		return errors.New("latitude must be between -90 and 90")
	}

	if r.Longitude < -180 || r.Longitude > 180 {
		return errors.New("longitude must be between -180 and 180")
	}

	start, err := time.Parse("2006-01-02", r.StartDate)
	if err != nil {
		return errors.New("start date must use YYYY-MM-DD format")
	}

	end, err := time.Parse("2006-01-02", r.EndDate)
	if err != nil {
		return errors.New("end date must use YYYY-MM-DD format")
	}

	if end.Before(start) {
		return errors.New("end date cannot be before start date")
	}

	return nil
}
