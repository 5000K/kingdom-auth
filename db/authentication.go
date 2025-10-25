package db

import (
	"time"

	"gorm.io/gorm"
)

type Authentication struct {
	gorm.Model

	UserID uint

	Provider       string
	ProviderUserId string

	Email    string
	Username string

	// these will only be stored, if set up in the provider configuration.
	// they are only accessible via the system facing backend.

	// YOU ARE RESPONSIBLE FOR STORING YOUR USERS DATA SECURELY.

	AccessToken        *string
	RefreshToken       *string
	RefreshTokenExpiry *time.Time
	TokenType          *string
}
