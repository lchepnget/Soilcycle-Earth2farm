package soil

import (
	"time"

	"github.com/google/uuid"
)

type SoilData struct {
	ID              uuid.UUID  `json:"id"`
	FarmID          uuid.UUID  `json:"farm_id"`
	PH              *float64   `json:"ph,omitempty"`
	OrganicCarbon   *float64   `json:"organic_carbon,omitempty"`
	Nitrogen        *float64   `json:"nitrogen,omitempty"`
	Phosphorus      *float64   `json:"phosphorus,omitempty"`
	Potassium       *float64   `json:"potassium,omitempty"`
	SoilTexture     *string    `json:"soil_texture,omitempty"`
	DataSource      string     `json:"data_source"`
	MeasurementDate *time.Time `json:"measurement_date,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}
