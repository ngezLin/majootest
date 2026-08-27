package utils

import (
	"testing"

	"majootest/case2-go/internal/models"
)

func TestPasswordHashing(t *testing.T) {
	rawPassword := "SecurePass123!"

	hash, err := HashPassword(rawPassword)
	if err != nil {
		t.Fatalf("expected no error hashing password, got %v", err)
	}

	if !CheckPasswordHash(rawPassword, hash) {
		t.Errorf("expected password hash to match raw password")
	}

	if CheckPasswordHash("WrongPassword", hash) {
		t.Errorf("expected wrong password to fail check")
	}
}

func TestJWTGenerationAndValidation(t *testing.T) {
	secret := "test-secret-key-12345"
	userID := int64(42)
	email := "test@example.com"

	token, err := GenerateJWT(userID, email, secret, 1)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	claims, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("failed to validate token: %v", err)
	}

	if claims.UserID != userID || claims.Email != email {
		t.Errorf("expected userID %d and email %s, got %d and %s", userID, email, claims.UserID, claims.Email)
	}

	// Test invalid secret
	_, err = ValidateJWT(token, "wrong-secret")
	if err == nil {
		t.Errorf("expected error with invalid secret, got nil")
	}
}

func TestValidator(t *testing.T) {
	validReq := models.RegisterRequest{
		Name:     "Alice Smith",
		Email:    "alice@example.com",
		Password: "password123",
	}

	details, err := ValidateStruct(validReq)
	if err != nil || len(details) > 0 {
		t.Errorf("expected valid struct, got error: %v, details: %v", err, details)
	}

	invalidReq := models.RegisterRequest{
		Name:     "A",          // too short (min 2)
		Email:    "notanemail", // invalid email
		Password: "short",      // too short (min 8)
	}

	details, err = ValidateStruct(invalidReq)
	if err == nil || len(details) != 3 {
		t.Errorf("expected 3 validation errors, got %d errors", len(details))
	}
}

func TestExpiredJWT(t *testing.T) {
	secret := "test-secret"
	token, err := GenerateJWT(1, "test@example.com", secret, -1) // expired 1 hour ago
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Errorf("expected error on expired token, got nil")
	}
}
