package crops

import (
	"strings"
	"testing"
)

func TestCreateCropRequestValidate(t *testing.T) {
	tests := []struct {
		name     string
		request  CreateCropRequest
		wantErr  string
		wantName string
	}{
		{
			name:     "valid crop name",
			request:  CreateCropRequest{Name: "Maize"},
			wantName: "Maize",
		},
		{
			name:    "empty crop name",
			request: CreateCropRequest{Name: ""},
			wantErr: "crop name is required",
		},
		{
			name:    "whitespace-only crop name",
			request: CreateCropRequest{Name: "   "},
			wantErr: "crop name is required",
		},
		{
			name:    "crop name exceeds 100 characters",
			request: CreateCropRequest{Name: string(make([]byte, 101))},
			wantErr: "crop name cannot exceed 100 characters",
		},
		{
			name:     "crop name is trimmed",
			request:  CreateCropRequest{Name: "  Beans  "},
			wantName: "Beans",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()

			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("Validate() expected error %q, got nil", tt.wantErr)
				}

				if err.Error() != tt.wantErr {
					t.Errorf("Validate() error = %q, want %q", err.Error(), tt.wantErr)
				}

				return
			}

			if err != nil {
				t.Fatalf("Validate() unexpected error: %v", err)
			}

			if tt.request.Name != tt.wantName {
				t.Errorf("Name = %q, want %q", tt.request.Name, tt.wantName)
			}
		})
	}
}
func TestUpdateCropRequestValidate(t *testing.T) {
	tests := []struct {
		name     string
		request  UpdateCropRequest
		wantErr  string
		wantName string
	}{
		{
			name:     "valid crop name",
			request:  UpdateCropRequest{Name: "Maize"},
			wantName: "Maize",
		},
		{
			name:    "empty crop name",
			request: UpdateCropRequest{Name: ""},
			wantErr: "crop name is required",
		},
		{
			name:    "whitespace-only crop name",
			request: UpdateCropRequest{Name: "   "},
			wantErr: "crop name is required",
		},
		{
			name:    "crop name exceeds 100 characters",
			request: UpdateCropRequest{Name: strings.Repeat("a", 101)},
			wantErr: "crop name cannot exceed 100 characters",
		},
		{
			name:     "crop name is trimmed",
			request:  UpdateCropRequest{Name: "  Beans  "},
			wantName: "Beans",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()

			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("Validate() expected error %q, got nil", tt.wantErr)
				}

				if err.Error() != tt.wantErr {
					t.Errorf("Validate() error = %q, want %q", err.Error(), tt.wantErr)
				}

				return
			}

			if err != nil {
				t.Fatalf("Validate() unexpected error: %v", err)
			}

			if tt.request.Name != tt.wantName {
				t.Errorf("Name = %q, want %q", tt.request.Name, tt.wantName)
			}
		})
	}
}
func TestCropCharacteristicRequestValidate(t *testing.T) {
	zero := 0
	negative := -5
	positive := 90

	low := "low"
	medium := "medium"
	high := "high"
	uppercase := "HIGH"
	invalid := "extreme"

	tests := []struct {
		name    string
		request CropCharacteristicRequest
		wantErr string
	}{
		{
			name:    "optional fields can be empty",
			request: CropCharacteristicRequest{},
		},
		{
			name: "positive growth duration",
			request: CropCharacteristicRequest{
				GrowthDurationDays: &positive,
			},
		},
		{
			name: "zero growth duration",
			request: CropCharacteristicRequest{
				GrowthDurationDays: &zero,
			},
			wantErr: "growth duration must be greater than zero",
		},
		{
			name: "negative growth duration",
			request: CropCharacteristicRequest{
				GrowthDurationDays: &negative,
			},
			wantErr: "growth duration must be greater than zero",
		},
		{
			name: "valid water requirement",
			request: CropCharacteristicRequest{
				WaterRequirement: &medium,
			},
		},
		{
			name: "uppercase water requirement",
			request: CropCharacteristicRequest{
				WaterRequirement: &uppercase,
			},
		},
		{
			name: "invalid water requirement",
			request: CropCharacteristicRequest{
				WaterRequirement: &invalid,
			},
			wantErr: "water requirement must be low, medium or high",
		},
		{
			name: "valid nutrient demand",
			request: CropCharacteristicRequest{
				NutrientDemand: &high,
			},
		},
		{
			name: "invalid nutrient demand",
			request: CropCharacteristicRequest{
				NutrientDemand: &invalid,
			},
			wantErr: "nutrient demand must be low, medium or high",
		},
		{
			name: "all valid characteristic fields",
			request: CropCharacteristicRequest{
				WaterRequirement:   &low,
				GrowthDurationDays: &positive,
				NutrientDemand:     &medium,
				NitrogenFixing:     true,
				SuitableSeasons:    []string{"long_rains", "short_rains"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()

			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("Validate() expected error %q, got nil", tt.wantErr)
				}

				if err.Error() != tt.wantErr {
					t.Errorf("Validate() error = %q, want %q", err.Error(), tt.wantErr)
				}

				return
			}

			if err != nil {
				t.Errorf("Validate() unexpected error: %v", err)
			}
		})
	}
}
