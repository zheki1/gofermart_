package auth

import (
	"testing"
)

func TestGenerateAndParseToken(t *testing.T) {
	userID := 42

	token, err := GenerateToken(userID)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	parsedID, err := ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken failed: %v", err)
	}

	if parsedID != userID {
		t.Fatalf("expected userID %d, got %d", userID, parsedID)
	}
}

func TestParseToken_InvalidToken(t *testing.T) {
	_, err := ParseToken("invalid.token.string")
	if err == nil {
		t.Fatal("expected error for invalid token, got nil")
	}
}

func TestHashAndCheckPassword(t *testing.T) {
	password := "supersecret"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if err := CheckPassword(hash, password); err != nil {
		t.Fatalf("CheckPassword failed: %v", err)
	}
}

func TestCheckPassword_WrongPassword(t *testing.T) {
	password := "correct"
	wrong := "wrong"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if err := CheckPassword(hash, wrong); err == nil {
		t.Fatal("expected error for wrong password, got nil")
	}
}
