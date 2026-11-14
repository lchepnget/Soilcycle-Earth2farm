package rotation

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type RotationPlan struct {
	ID               uuid.UUID       `json:"id"`
	FarmID           uuid.UUID       `json:"farm_id"`
	FarmerPriorities json.RawMessage `json:"farmer_priorities"`
	InputSnapshot    json.RawMessage `json:"input_snapshot"`
	NumberOfSeasons  int             `json:"number_of_seasons"`
	CreatedAt        time.Time       `json:"created_at"`
}

type RotationResult struct {
	ID           uuid.UUID       `json:"id"`
	PlanID       uuid.UUID       `json:"plan_id"`
	PlanName     string          `json:"plan_name"`
	CropSequence json.RawMessage `json:"crop_sequence"`
	Scores       json.RawMessage `json:"scores"`
	Explanation  json.RawMessage `json:"explanation"`
	OverallScore *float64        `json:"overall_score,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
}

type CropRecommendation struct {
	CropID   uuid.UUID `json:"crop_id"`
	CropName string    `json:"crop_name"`
	Season   int       `json:"season"`
	Reason   string    `json:"reason"`
	Benefits []string  `json:"benefits"`
}

type RotationScore struct {
	SoilHealthScore float64 `json:"soil_health_score"`
	ClimateScore    float64 `json:"climate_score"`
	DiversityScore  float64 `json:"diversity_score"`
	OverallScore    float64 `json:"overall_score"`
}
