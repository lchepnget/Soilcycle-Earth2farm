package auth

import (
	"errors"
	"strings"
)

type LoginRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

func (r *LoginRequest) Validate() error {
	r.Identifier = strings.TrimSpace(r.Identifier)
	if r.Identifier == "" || r.Password == "" {
		return errors.New("email or phone and password are required")
	}
	if len(r.Identifier) > 255 {
		return errors.New("identifier is too long")
	}
	return nil
}
