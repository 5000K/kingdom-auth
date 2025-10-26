package db

import (
	"gorm.io/gorm"
)

type Authentication struct {
	gorm.Model

	UserID uint

	Provider string
	Subject  string

	Email string
}
