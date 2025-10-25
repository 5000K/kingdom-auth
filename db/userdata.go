package db

import "gorm.io/gorm"

type UserData struct {
	gorm.Model

	UserID uint
	Data   string
}
