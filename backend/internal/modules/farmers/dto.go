package farmers

import (
	"errors"
	"strings"
)

type CreateFarmerRequest struct {
	FullName          string  `json:"full_name"`
	Email             *string `json:"email"`
	Phone             *string `json:"phone"`
	Password          string  `json:"password"`
	PreferredLanguage string  `json:"preferred_language"`
}

type UpdateFarmerRequest struct {
	FullName          string  `json:"full_name"`
	Email             *string `json:"email"`
	Phone             *string `json:"phone"`
	PreferredLanguage string  `json:"preferred_language"`
}

func (r *CreateFarmerRequest) Validate() error {
	r.FullName = strings.TrimSpace(r.FullName)

	if r.FullName == "" {
		return errors.New("full name is required")
	}

	if len(r.FullName) > 150 {
		return errors.New("full name cannot exceed 150 characters")
	}

	if (r.Email == nil || strings.TrimSpace(*r.Email) == "") &&
		(r.Phone == nil || strings.TrimSpace(*r.Phone) == "") {
		return errors.New("either email or phone is required")
	}

	if len(r.Password) < 8 {
		return errors.New("password must contain at least 8 characters")
	}

	if r.Email != nil {
		email := strings.TrimSpace(*r.Email)
		if email != "" && !strings.Contains(email, "@") {
			return errors.New("invalid email address")
		}
		*r.Email = email
	}

	if r.Phone != nil {
		*r.Phone = strings.TrimSpace(*r.Phone)
	}

	if r.PreferredLanguage == "" {
		r.PreferredLanguage = "en"
	}

	return nil
}

func (r *UpdateFarmerRequest) Validate() error {
	r.FullName = strings.TrimSpace(r.FullName)

	if r.FullName == "" {
		return errors.New("full name is required")
	}

	if len(r.FullName) > 150 {
		return errors.New("full name cannot exceed 150 characters")
	}

	if (r.Email == nil || strings.TrimSpace(*r.Email) == "") &&
		(r.Phone == nil || strings.TrimSpace(*r.Phone) == "") {
		return errors.New("either email or phone is required")
	}

	if r.Email != nil {
		email := strings.TrimSpace(*r.Email)
		if email != "" && !strings.Contains(email, "@") {
			return errors.New("invalid email address")
		}
		*r.Email = email
	}

	if r.Phone != nil {
		*r.Phone = strings.TrimSpace(*r.Phone)
	}

	if r.PreferredLanguage == "" {
		r.PreferredLanguage = "en"
	}

	return nil
}
