package farmers

import (
	"context"
	"errors"
	"testing"
)

func TestFarmerServiceCreateRejectsInvalidRequest(t *testing.T) {
	service := NewFarmerService(nil)

	req := CreateFarmerRequest{
		FullName: "",
		Email:    stringPointer("ian@example.com"),
		Password: "password123",
	}

	farmer, err := service.Create(context.Background(), req)

	if err == nil {
		t.Fatal("expected validation error, got nil")
	}

	if !errors.Is(err, ErrInvalidFarmerInput) {
		t.Errorf("expected ErrInvalidFarmerInput, got %v", err)
	}

	if farmer != nil {
		t.Errorf("expected nil farmer, got %+v", farmer)
	}
}
func TestFarmerServiceGetByIDRejectsInvalidUUID(t *testing.T) {
	service := NewFarmerService(nil)

	farmer, err := service.GetByID(context.Background(), "not-a-valid-uuid")

	if err == nil {
		t.Fatal("expected error for invalid farmer ID, got nil")
	}

	if farmer != nil {
		t.Errorf("expected nil farmer, got %+v", farmer)
	}

	if err.Error() != "invalid farmer ID" {
		t.Errorf("expected invalid farmer ID error, got %q", err.Error())
	}
}
func TestFarmerServiceUpdateRejectsInvalidUUID(t *testing.T) {
	service := NewFarmerService(nil)

	req := UpdateFarmerRequest{
		FullName: "John Farmer",
		Email:    stringPointer("john@example.com"),
	}

	farmer, err := service.Update(
		context.Background(),
		"not-a-valid-uuid",
		req,
	)

	if err == nil {
		t.Fatal("expected error for invalid farmer ID, got nil")
	}

	if farmer != nil {
		t.Errorf("expected nil farmer, got %+v", farmer)
	}

	if err.Error() != "invalid farmer ID" {
		t.Errorf("expected invalid farmer ID error, got %q", err.Error())
	}
}

func TestFarmerServiceUpdateRejectsInvalidRequest(t *testing.T) {
	service := NewFarmerService(nil)

	req := UpdateFarmerRequest{
		FullName: "",
		Email:    stringPointer("john@example.com"),
	}

	farmer, err := service.Update(
		context.Background(),
		"550e8400-e29b-41d4-a716-446655440000",
		req,
	)

	if err == nil {
		t.Fatal("expected validation error, got nil")
	}

	if !errors.Is(err, ErrInvalidFarmerInput) {
		t.Errorf("expected ErrInvalidFarmerInput, got %v", err)
	}

	if farmer != nil {
		t.Errorf("expected nil farmer, got %+v", farmer)
	}
}
