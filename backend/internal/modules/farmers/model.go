package farmers

import (
	"time"

	"github.com/google/uuid"
)

type Farmer struct {
	ID                uuid.UUID `json:"id"`
	FullName          string    `json:"full_name"`
	Email             *string   `json:"email,omitempty"`
	Phone             *string   `json:"phone,omitempty"`
	Role              string    `json:"role"`
	PreferredLanguage string    `json:"preferred_language"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
