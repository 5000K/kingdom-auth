package service

import (
	"testing"

	"github.com/5000K/kingdom-auth/config"
	"github.com/5000K/kingdom-auth/db"
	"github.com/golang-jwt/jwt/v5"
)

func TestRS512TokenCreationAndVerification(t *testing.T) {
	// Create a test config
	cfg := &config.Config{}
	cfg.Token.PrivateKeyPath = "../private_key.pem"
	cfg.Token.PublicKeyPath = "../public_key.pem"
	cfg.Token.RefreshTokenTTL = 86400
	cfg.Token.AuthTokenTTL = 90
	cfg.Token.Issuer = "test-issuer"
	cfg.Token.DefaultAudience = "test-audience"

	// Create service
	svc, err := NewService(cfg, nil)
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}

	// Create a test user
	user := &db.User{}
	user.ID = 1
	user.PublicData = `{"username":"testuser"}`

	// Test refresh token creation
	refreshToken, err := svc.createRefreshTokenFor(user)
	if err != nil {
		t.Fatalf("Failed to create refresh token: %v", err)
	}

	// Verify refresh token
	claims, err := svc.readRefreshToken(refreshToken)
	if err != nil {
		t.Fatalf("Failed to read refresh token: %v", err)
	}

	if claims["sub"] != "1" {
		t.Errorf("Expected subject '1', got '%v'", claims["sub"])
	}

	if claims["iss"] != "test-issuer" {
		t.Errorf("Expected issuer 'test-issuer', got '%v'", claims["iss"])
	}

	// Verify the token uses RS512
	token, _ := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return svc.publicKey, nil
	})

	if token.Method.Alg() != "RS512" {
		t.Errorf("Expected RS512 algorithm, got %s", token.Method.Alg())
	}

	// Test auth token creation
	authToken, _, err := svc.createAuthTokenFor(user)
	if err != nil {
		t.Fatalf("Failed to create auth token: %v", err)
	}

	// Verify auth token
	authClaims, err := svc.readAuthToken(authToken)
	if err != nil {
		t.Fatalf("Failed to read auth token: %v", err)
	}

	if authClaims["sub"] != "1" {
		t.Errorf("Expected subject '1', got '%v'", authClaims["sub"])
	}

	if authClaims["aud"] != "test-audience" {
		t.Errorf("Expected audience 'test-audience', got '%v'", authClaims["aud"])
	}

	// Verify the auth token also uses RS512
	authTokenParsed, _ := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
		return svc.publicKey, nil
	})

	if authTokenParsed.Method.Alg() != "RS512" {
		t.Errorf("Expected RS512 algorithm for auth token, got %s", authTokenParsed.Method.Alg())
	}
}
