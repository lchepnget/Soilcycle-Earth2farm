package farmers

import (
	"strings"
	"testing"
)

func TestCreateFarmerRequestValidateRejectsEmptyName(t *testing.T) {
	req := CreateFarmerRequest{
		FullName: "",
		Password: "password123",
		Email:    stringPointer("farmer@example.com"),
	}

	err := req.Validate()

	if err == nil {
		t.Fatal("expected an error for an empty full name, got nil")
	}
}

func stringPointer(value string) *string {
	return &value
}
func TestCreateFarmerRequestValidateRejectsWhitespaceName(t *testing.T) {
	req := CreateFarmerRequest{
		FullName: "   ",
		Password: "password123",
		Email:    stringPointer("farmer@example.com"),
	}

	err := req.Validate()

	if err == nil {
		t.Fatal("expected an error for a whitespace-only full name, got nil")
	}
}
func TestCreateFarmerRequestValidateTrimsFullName(t *testing.T) {
	req := CreateFarmerRequest{
		FullName: "  Ian Farmer  ",
		Password: "password123",
		Email:    stringPointer("farmer@example.com"),
	}

	err := req.Validate()

	if err != nil {
		t.Fatalf("expected valid request, got error: %v", err)
	}

	if req.FullName != "Ian Farmer" {
		t.Errorf("FullName = %q, want %q", req.FullName, "Ian Farmer")
	}
}
func TestCreateFarmerRequestValidateRejectsNameOver150Characters(t *testing.T) {
	req := CreateFarmerRequest{
		FullName: strings.Repeat("A", 151),
		Password: "password123",
		Email:    stringPointer("farmer@example.com"),
	}

	err := req.Validate()

	if err == nil {
		t.Fatal("expected an error for a name longer than 150 characters, got nil")
	}
}
func TestCreateFarmerRequestValidateRequiresEmailOrPhone(t *testing.T) {
	req := CreateFarmerRequest{
		FullName: "Ian Farmer",
		Password: "password123",
	}

	err := req.Validate()

	if err == nil {
		t.Fatal("expected an error when both email and phone are missing, got nil")
	}
}
func TestCreateFarmerRequestValidateRejectsInvalidEmail(t *testing.T) {
	email := "not-an-email"

	req := CreateFarmerRequest{
		FullName: "Ian Farmer",
		Password: "password123",
		Email:    &email,
	}

	err := req.Validate()

	if err == nil {
		t.Fatal("expected an error for an invalid email, got nil")
	}
}
func TestCreateFarmerRequestValidateTrimsEmail(t *testing.T) {
	email := "  farmer@example.com  "

	req := CreateFarmerRequest{
		FullName: "Ian Farmer",
		Password: "password123",
		Email:    &email,
	}

	err := req.Validate()

	if err != nil {
		t.Fatalf("expected valid request, got error: %v", err)
	}

	if *req.Email != "farmer@example.com" {
		t.Errorf("Email = %q, want %q", *req.Email, "farmer@example.com")
	}
}
func TestCreateFarmerRequestValidateRejectsShortPassword(t *testing.T) {
	req := CreateFarmerRequest{
		FullName: "Ian Farmer",
		Password: "1234567",
		Email:    stringPointer("farmer@example.com"),
	}

	err := req.Validate()

	if err == nil {
		t.Fatal("expected an error for a password shorter than 8 characters, got nil")
	}
}
func TestCreateFarmerRequestValidateAcceptsEightCharacterPassword(t *testing.T) {
	req := CreateFarmerRequest{
		FullName: "Ian Farmer",
		Password: "12345678",
		Email:    stringPointer("farmer@example.com"),
	}

	err := req.Validate()

	if err != nil {
		t.Fatalf("expected an 8-character password to be accepted, got error: %v", err)
	}
}
func TestCreateFarmerRequestValidateAcceptsPhoneWithoutEmail(t *testing.T) {
	phone := "+254712345678"

	req := CreateFarmerRequest{
		FullName: "Ian Farmer",
		Password: "password123",
		Phone:    &phone,
	}

	err := req.Validate()

	if err != nil {
		t.Fatalf("expected phone-only registration to be accepted, got error: %v", err)
	}
}
func TestCreateFarmerRequestValidateAcceptsEmailWithoutPhone(t *testing.T) {
	email := "farmer@example.com"

	req := CreateFarmerRequest{
		FullName: "Ian Farmer",
		Password: "password123",
		Email:    &email,
	}

	err := req.Validate()

	if err != nil {
		t.Fatalf("expected email-only registration to be accepted, got error: %v", err)
	}
}
func TestCreateFarmerRequestValidateDefaultsPreferredLanguage(t *testing.T) {
	req := CreateFarmerRequest{
		FullName: "Ian Farmer",
		Email:    stringPointer("ian@example.com"),
		Password: "password123",
	}

	err := req.Validate()
	if err != nil {
		t.Fatalf("expected valid request, got error: %v", err)
	}

	if req.PreferredLanguage != "en" {
		t.Errorf("expected preferred language to default to en, got %q", req.PreferredLanguage)
	}
}
func TestCreateFarmerRequestValidateTrimsPhone(t *testing.T) {
	phone := "  +254712345678  "

	req := CreateFarmerRequest{
		FullName: "Ian Farmer",
		Phone:    &phone,
		Password: "password123",
	}

	err := req.Validate()
	if err != nil {
		t.Fatalf("expected valid request, got error: %v", err)
	}

	if *req.Phone != "+254712345678" {
		t.Errorf("expected phone to be trimmed, got %q", *req.Phone)
	}
}
func TestUpdateFarmerRequestValidateRejectsEmptyName(t *testing.T) {
	req := UpdateFarmerRequest{
		FullName: "",
		Email:    stringPointer("ian@example.com"),
	}

	err := req.Validate()

	if err == nil {
		t.Fatal("expected validation error for empty full name, got nil")
	}
}
func TestUpdateFarmerRequestValidateRejectsWhitespaceName(t *testing.T) {
	req := UpdateFarmerRequest{
		FullName: "   ",
		Email:    stringPointer("ian@example.com"),
	}

	err := req.Validate()

	if err == nil {
		t.Fatal("expected validation error for whitespace-only name, got nil")
	}
}
func TestUpdateFarmerRequestValidateTrimsFullName(t *testing.T) {
	req := UpdateFarmerRequest{
		FullName: "  Ian Farmer  ",
		Email:    stringPointer("ian@example.com"),
	}

	err := req.Validate()
	if err != nil {
		t.Fatalf("expected valid request, got error: %v", err)
	}

	if req.FullName != "Ian Farmer" {
		t.Errorf("expected full name to be trimmed, got %q", req.FullName)
	}
}
func TestUpdateFarmerRequestValidateRejectsNameOver150Characters(t *testing.T) {
	req := UpdateFarmerRequest{
		FullName: strings.Repeat("a", 151),
		Email:    stringPointer("ian@example.com"),
	}

	err := req.Validate()

	if err == nil {
		t.Fatal("expected validation error for name longer than 150 characters, got nil")
	}
}
func TestUpdateFarmerRequestValidateAcceptsNameOf150Characters(t *testing.T) {
	req := UpdateFarmerRequest{
		FullName: strings.Repeat("a", 150),
		Email:    stringPointer("ian@example.com"),
	}

	err := req.Validate()

	if err != nil {
		t.Fatalf("expected 150-character name to be accepted, got error: %v", err)
	}
}
func TestUpdateFarmerRequestValidateRequiresEmailOrPhone(t *testing.T) {
	req := UpdateFarmerRequest{
		FullName: "Ian Farmer",
	}

	err := req.Validate()

	if err == nil {
		t.Fatal("expected validation error when both email and phone are missing, got nil")
	}
}
