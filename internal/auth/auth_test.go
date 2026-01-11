package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "test-secret-key"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT() error = %v, want nil", err)
	}

	if token == "" {
		t.Error("MakeJWT() returned empty token")
	}
}

func TestValidateJWT_ValidToken(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "test-secret-key"
	expiresIn := time.Hour

	// Create a valid token
	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT() error = %v, want nil", err)
	}

	// Validate the token
	validatedUserID, err := ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Fatalf("ValidateJWT() error = %v, want nil", err)
	}

	if validatedUserID != userID {
		t.Errorf("ValidateJWT() userID = %v, want %v", validatedUserID, userID)
	}
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "test-secret-key"
	expiresIn := -time.Hour // Expired 1 hour ago

	// Create an expired token
	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT() error = %v, want nil", err)
	}

	// Try to validate the expired token
	_, err = ValidateJWT(token, tokenSecret)
	if err == nil {
		t.Error("ValidateJWT() expected error for expired token, got nil")
	}
}

func TestValidateJWT_WrongSecret(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "test-secret-key"
	wrongSecret := "wrong-secret-key"
	expiresIn := time.Hour

	// Create a token with the correct secret
	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT() error = %v, want nil", err)
	}

	// Try to validate with the wrong secret
	_, err = ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Error("ValidateJWT() expected error for wrong secret, got nil")
	}
}

func TestValidateJWT_InvalidToken(t *testing.T) {
	tokenSecret := "test-secret-key"
	invalidToken := "invalid.token.string"

	// Try to validate an invalid token
	_, err := ValidateJWT(invalidToken, tokenSecret)
	if err == nil {
		t.Error("ValidateJWT() expected error for invalid token, got nil")
	}
}

func TestMakeJWT_DifferentUserIDs(t *testing.T) {
	tokenSecret := "test-secret-key"
	expiresIn := time.Hour

	userID1 := uuid.New()
	userID2 := uuid.New()

	// Create tokens for different users
	token1, err := MakeJWT(userID1, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT() error = %v, want nil", err)
	}

	token2, err := MakeJWT(userID2, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT() error = %v, want nil", err)
	}

	// Tokens should be different
	if token1 == token2 {
		t.Error("MakeJWT() generated identical tokens for different user IDs")
	}

	// Validate both tokens
	validatedUserID1, err := ValidateJWT(token1, tokenSecret)
	if err != nil {
		t.Fatalf("ValidateJWT() error = %v, want nil", err)
	}

	validatedUserID2, err := ValidateJWT(token2, tokenSecret)
	if err != nil {
		t.Fatalf("ValidateJWT() error = %v, want nil", err)
	}

	// Ensure each token returns the correct user ID
	if validatedUserID1 != userID1 {
		t.Errorf("ValidateJWT() userID = %v, want %v", validatedUserID1, userID1)
	}

	if validatedUserID2 != userID2 {
		t.Errorf("ValidateJWT() userID = %v, want %v", validatedUserID2, userID2)
	}
}

func TestValidateJWT_ShortExpirationTime(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "test-secret-key"
	expiresIn := 100 * time.Millisecond // Very short expiration

	// Create a token with short expiration
	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT() error = %v, want nil", err)
	}

	// Immediately validate - should work
	validatedUserID, err := ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Fatalf("ValidateJWT() immediate validation error = %v, want nil", err)
	}

	if validatedUserID != userID {
		t.Errorf("ValidateJWT() userID = %v, want %v", validatedUserID, userID)
	}

	// Wait for token to expire
	time.Sleep(150 * time.Millisecond)

	// Validate again - should fail
	_, err = ValidateJWT(token, tokenSecret)
	if err == nil {
		t.Error("ValidateJWT() expected error for expired token after sleep, got nil")
	}
}
