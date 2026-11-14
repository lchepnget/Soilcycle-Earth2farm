package crops

import (
	"time"

	"github.com/google/uuid"
)

type Crop struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	ScientificName *string   `json:"scientific_name,omitempty"`
	CropFamily     *string   `json:"crop_family,omitempty"`
	Category       *string   `json:"category,omitempty"`
	Description    *string   `json:"description,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

type CropCharacteristic struct {
	ID                 uuid.UUID `json:"id"`
	CropID             uuid.UUID `json:"crop_id"`
	WaterRequirement   *string   `json:"water_requirement,omitempty"`
	GrowthDurationDays *int      `json:"growth_duration_days,omitempty"`
	NitrogenFixing     bool      `json:"nitrogen_fixing"`
	NutrientDemand     *string   `json:"nutrient_demand,omitempty"`
	SuitableSeasons    []string  `json:"suitable_seasons"`
	CreatedAt          time.Time `json:"created_at"`
}
