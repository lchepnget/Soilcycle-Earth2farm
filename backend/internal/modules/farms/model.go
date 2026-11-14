package farms

import (
	"time"

	"github.com/google/uuid"
)

type Farm struct {
	ID                uuid.UUID `json:"id"`
	UserID            uuid.UUID `json:"user_id"`
	Name              string    `json:"name"`
	County            string    `json:"county"`
	SubCounty         *string   `json:"sub_county,omitempty"`
	Latitude          float64   `json:"latitude"`
	Longitude         float64   `json:"longitude"`
	AreaHectares      *float64  `json:"area_hectares,omitempty"`
	WaterAvailability *string   `json:"water_availability,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
}
