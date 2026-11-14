package crops

import (
	"context"
	"testing"
)

func TestCropServiceCreateRejectsEmptyName(t *testing.T) {
	service := NewCropService(nil)

	req := CreateCropRequest{
		Name: "",
	}

	got, err := service.Create(context.Background(), req)

	if err == nil {
		t.Fatal("Create() expected an error for an empty crop name, got nil")
	}

	if got != nil {
		t.Errorf("Create() returned crop %v, want nil", got)
	}

	if err.Error() != "crop name is required" {
		t.Errorf("Create() error = %q, want %q",
			err.Error(), "crop name is required")
	}
}

func TestCropServiceGetByIDRejectsInvalidID(t *testing.T) {
	service := NewCropService(nil)

	got, err := service.GetByID(context.Background(), "invalid-id")

	if err == nil {
		t.Fatal("GetByID() expected an error for an invalid ID, got nil")
	}

	if got != nil {
		t.Errorf("GetByID() returned crop %v, want nil", got)
	}

	if err.Error() != "invalid crop ID" {
		t.Errorf("GetByID() error = %q, want %q",
			err.Error(), "invalid crop ID")
	}
}

func TestCropServiceUpdateRejectsInvalidID(t *testing.T) {
	service := NewCropService(nil)

	req := UpdateCropRequest{Name: "Maize"}

	got, err := service.Update(context.Background(), "invalid-id", req)

	if err == nil {
		t.Fatal("Update() expected an error for an invalid ID, got nil")
	}

	if got != nil {
		t.Errorf("Update() returned crop %v, want nil", got)
	}

	if err.Error() != "invalid crop ID" {
		t.Errorf("Update() error = %q, want %q",
			err.Error(), "invalid crop ID")
	}
}

func TestCropServiceUpdateRejectsEmptyName(t *testing.T) {
	service := NewCropService(nil)

	req := UpdateCropRequest{Name: ""}

	got, err := service.Update(
		context.Background(),
		"550e8400-e29b-41d4-a716-446655440000",
		req,
	)

	if err == nil {
		t.Fatal("Update() expected an error for an empty crop name, got nil")
	}

	if got != nil {
		t.Errorf("Update() returned crop %v, want nil", got)
	}

	if err.Error() != "crop name is required" {
		t.Errorf("Update() error = %q, want %q",
			err.Error(), "crop name is required")
	}
}

func TestCropServiceDeleteRejectsInvalidID(t *testing.T) {
	service := NewCropService(nil)

	err := service.Delete(context.Background(), "invalid-id")

	if err == nil {
		t.Fatal("Delete() expected an error for an invalid ID, got nil")
	}

	if err.Error() != "invalid crop ID" {
		t.Errorf("Delete() error = %q, want %q",
			err.Error(), "invalid crop ID")
	}
}

func TestCropServiceGetCharacteristicsRejectsInvalidID(t *testing.T) {
	service := NewCropService(nil)

	got, err := service.GetCharacteristics(
		context.Background(),
		"invalid-id",
	)

	if err == nil {
		t.Fatal("GetCharacteristics() expected an error for an invalid ID, got nil")
	}

	if got != nil {
		t.Errorf("GetCharacteristics() returned %v, want nil", got)
	}

	if err.Error() != "invalid crop ID" {
		t.Errorf("GetCharacteristics() error = %q, want %q",
			err.Error(), "invalid crop ID")
	}
}

func TestCropServiceUpsertCharacteristicsRejectsInvalidID(t *testing.T) {
	service := NewCropService(nil)

	got, err := service.UpsertCharacteristics(
		context.Background(),
		"invalid-id",
		CropCharacteristicRequest{},
	)

	if err == nil {
		t.Fatal("UpsertCharacteristics() expected an error for an invalid ID, got nil")
	}

	if got != nil {
		t.Errorf("UpsertCharacteristics() returned %v, want nil", got)
	}

	if err.Error() != "invalid crop ID" {
		t.Errorf("UpsertCharacteristics() error = %q, want %q",
			err.Error(), "invalid crop ID")
	}
}
func TestCropServiceUpsertCharacteristicsRejectsInvalidGrowthDuration(t *testing.T) {
	service := NewCropService(nil)

	zero := 0

	req := CropCharacteristicRequest{
		GrowthDurationDays: &zero,
	}

	got, err := service.UpsertCharacteristics(
		context.Background(),
		"550e8400-e29b-41d4-a716-446655440000",
		req,
	)

	if err == nil {
		t.Fatal("UpsertCharacteristics() expected an error for zero growth duration, got nil")
	}

	if got != nil {
		t.Errorf("UpsertCharacteristics() returned %v, want nil", got)
	}

	if err.Error() != "growth duration must be greater than zero" {
		t.Errorf("UpsertCharacteristics() error = %q, want %q",
			err.Error(), "growth duration must be greater than zero")
	}
}
