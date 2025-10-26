package main

import (
	"fmt"
	"log"

	kingdomauth "github.com/5000K/kingdom-auth"
)

func main() {
	// Initialize the Kingdom Auth client
	client, err := kingdomauth.NewClient(
		"https://auth.example.com",
		"your-service-secret",
		"./public_key.pem",
	)
	if err != nil {
		log.Fatalf("Failed to create Kingdom Auth client: %v", err)
	}

	// Example: Validate a JWT token
	token := "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9..." // Your JWT token here

	claims, err := client.ValidateToken(token)
	if err != nil {
		log.Fatalf("Token validation failed: %v", err)
	}

	// Access claims (these basic ones are always set IF the token is valid AND issued by Kingdom Auth)
	fmt.Printf("User ID (sub): %v\n", claims["sub"])
	fmt.Printf("Issuer (iss): %v\n", claims["iss"])
	fmt.Printf("Audience (aud): %v\n", claims["aud"])
	fmt.Printf("Expiration (exp): %v\n", claims["exp"])

	// only present if set (duh)
	if publicData, ok := claims["public-data"]; ok {
		fmt.Printf("Public Data: %v\n", publicData)
	}
}
