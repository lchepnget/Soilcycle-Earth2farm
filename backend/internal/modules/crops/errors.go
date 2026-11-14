package crops

import "errors"

// ValidationError represents invalid input supplied by an API client.
type ValidationError struct {
	Message string
}

// Error returns the validation error message.
func (e *ValidationError) Error() string {
	return e.Message
}

// NewValidationError creates a new validation error.
func NewValidationError(message string) error {
	return &ValidationError{Message: message}
}

var (
	ErrCropNotFound  = errors.New("crop not found")
	ErrDuplicateCrop = errors.New("crop already exists")
)
