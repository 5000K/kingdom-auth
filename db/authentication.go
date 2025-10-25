package db

import (
	"gorm.io/gorm"
)

type Authentication struct {
	gorm.Model

	UserID uint

	Provider       string
	ProviderUserId string

	Email    *string
	Username *string
}
