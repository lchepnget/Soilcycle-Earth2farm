package crops

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
)

func TestRespondError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{
			name:       "crop not found",
			err:        ErrCropNotFound,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "duplicate crop",
			err:        ErrDuplicateCrop,
			wantStatus: http.StatusConflict,
		},
		{
			name:       "unexpected database error",
			err:        errors.New("database unavailable"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()

			app.Get("/test", func(c fiber.Ctx) error {
				return respondError(c, tt.err)
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("app.Test() returned an error: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("HTTP status = %d, want %d",
					resp.StatusCode, tt.wantStatus)
			}
		})
	}
}
