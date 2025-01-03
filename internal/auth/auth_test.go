package auth

import (
	"testing"
)

func TestHashAndUnHash(t *testing.T) {
	password := "testingpassword"
	hashedPass, err := HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password")
	}
	err = CheckPasswordHash(password, hashedPass)
	if err != nil {
		t.Fatalf("failed to check password")
	}
}
