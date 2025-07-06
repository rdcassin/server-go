package auth

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWT_Success(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "my-secret-key"
	expiresIn := time.Hour
	
	// Create a JWT
	tokenString, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	// Validate the JWT
	validatedID, err := ValidateJWT(tokenString, tokenSecret)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	// Check that the user ID matches
	if validatedID != userID {
		t.Errorf("Expected user ID %v, got %v", userID, validatedID)
	}
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "my-secret-key"
	expiresIn := -time.Hour // Expired 1 hour ago
	
	// Create an expired JWT
	tokenString, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Expected no error creating token, got %v", err)
	}
	
	// Try to validate the expired JWT
	_, err = ValidateJWT(tokenString, tokenSecret)
	if err == nil {
		t.Error("Expected error for expired token, but got none")
	}
}

func TestValidateJWT_WrongSecret(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "my-secret-key"
	wrongSecret := "wrong-secret-key"
	expiresIn := time.Hour
	
	// Create a JWT with one secret
	tokenString, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Expected no error creating token, got %v", err)
	}
	
	// Try to validate with a different secret
	_, err = ValidateJWT(tokenString, wrongSecret)
	if err == nil {
		t.Error("Expected error for wrong secret, but got none")
	}
}

func TestGetBearerToken_Success(t *testing.T) {
	expectedTokenString := "abc"
	authorization := fmt.Sprintf("Bearer %s", expectedTokenString)
	header := http.Header{
		"Authorization": []string{authorization},
	}

	// Extract tokenString
	tokenString, err := GetBearerToken(header)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Check that the tokenString matches
	if expectedTokenString != tokenString {
		t.Errorf("Expected tokenString %v, got %v", expectedTokenString, tokenString)
	}
}

func TestGetBearerToken_AuthorationFieldMissing(t *testing.T) {
	header := http.Header{
		"Foo": []string{"ABC"},
		"Boots": []string{"Alcoholic Booze Hound", "Greatest Bear Mage"},
	}

	// Try to get tokenString
	_, err := GetBearerToken(header)
	if err == nil {
		t.Fatalf("Expected error for authorization missing")
	}
}

func TestGetBearerToken_AuthorizationTooLong(t *testing.T) {
	expectedTokenString := "abc def"
	authorization := fmt.Sprintf("Bearer %s", expectedTokenString)
	header := http.Header{
		"Authorization": []string{authorization},
	}

	// Try to get tokenString
	_, err := GetBearerToken(header)
	if err == nil {
		t.Fatalf("Expected error for invalid authorization length")
	}
}

func TestGetBearerToken_AuthorizationTooShort(t *testing.T) {
	expectedTokenString := ""
	authorization := fmt.Sprintf("Bearer %s", expectedTokenString)
	header := http.Header{
		"Authorization": []string{authorization},
	}

	// Try to get tokenString
	_, err := GetBearerToken(header)
	if err == nil {
		t.Fatalf("Expected error for invalid authorization length")
	}
}

func TestGetBearerToken_AuthorizationWrongFormat(t *testing.T) {
	expectedTokenString := "totally great server"
	authorization := fmt.Sprintf("Bearclaw %s", expectedTokenString)
	header := http.Header{
		"Authorization": []string{authorization},
	}

	// Try to get tokenString
	_, err := GetBearerToken(header)
	if err == nil {
		t.Fatalf("Expected error for invalid format")
	}
}