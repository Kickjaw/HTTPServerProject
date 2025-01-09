package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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

func TestMakeJWTAndValidateJWT(t *testing.T) {
	// Define test parameters
	tokenSecret := "test-secret-key"
	userID := uuid.New()
	expiresIn := time.Minute * 10

	// Generate a JWT
	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	assert.NoError(t, err, "error generating JWT")
	assert.NotEmpty(t, token, "generated token should not be empty")

	// Validate the JWT
	validatedUserID, err := ValidateJWT(token, tokenSecret)
	assert.NoError(t, err, "error validating JWT")
	assert.Equal(t, userID, validatedUserID, "validated user ID should match original user ID")
}

func TestExpiredJWT(t *testing.T) {
	tokenSecret := "test-secrect-key"
	userID := uuid.New()
	expiresIn := -time.Minute

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	assert.NoError(t, err, "error generating JWT")
	assert.NotEmpty(t, token, "generated token should not be empty")

	_, err = ValidateJWT(token, tokenSecret)
	assert.Error(t, err, "expired token should generate error")
	assert.Contains(t, err.Error(), "token is expired", "error should indicate token is expired")
}

func TestInvalidSecret(t *testing.T) {
	// Define test parameters
	tokenSecret := "test-secret-key"
	invalidSecret := "wrong-secret-key"
	userID := uuid.New()
	expiresIn := time.Minute * 10

	// Generate a valid JWT with the correct secret
	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	assert.NoError(t, err, "error generating JWT")
	assert.NotEmpty(t, token, "generated token should not be empty")

	// Validate the JWT with an incorrect secret
	_, err = ValidateJWT(token, invalidSecret)
	assert.Error(t, err, "token signed with the wrong secret should produce an error")
	assert.Contains(t, err.Error(), "signature is invalid", "error should indicate invalid signature")
}

func TestGetBearerToken(t *testing.T) {
	testHeader := http.Header{}
	testHeader.Set("Authorization", "Bearer testToken")

	badTestHead := http.Header{}
	badTestHead.Set("Authorization", "testToken")

	testTokenCorrect := "testToken"

	token, err := GetBearerToken(testHeader)
	assert.NoError(t, err, "error getting auth token from header")
	assert.Equal(t, testTokenCorrect, token)

	_, err = GetBearerToken(badTestHead)
	assert.Error(t, err, "improper header format should produce an error")
}
