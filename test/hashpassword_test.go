package auth_test

import (
	"testing"

	"github.com/Tran-Duc-Hoa/chirphy/internal/auth"
)

func TestHashPassword(t *testing.T) {
	password := "securepassword123"
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error: %v", err)
	}

	if len(hashedPassword) == 0 {
		t.Error("HashPassword() returned an empty hash")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "securepassword123"
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error: %v", err)
	}

	// Test with correct password
	err = auth.CheckPasswordHash(password, hashedPassword)
	if err != nil {
		t.Errorf("CheckPasswordHash() failed for correct password: %v", err)
	}

	// Test with incorrect password
	incorrectPassword := "wrongpassword"
	err = auth.CheckPasswordHash(incorrectPassword, hashedPassword)
	if err == nil {
		t.Error("CheckPasswordHash() did not fail for incorrect password")
	}
}

func TestHashPasswordConsistency(t *testing.T) {
	password := "securepassword123"
	hashedPassword1, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error: %v", err)
	}

	hashedPassword2, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error: %v", err)
	}

	if hashedPassword1 == hashedPassword2 {
		t.Error("HashPassword() returned the same hash for the same password, but it should be different")
	}
}