package db

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Authentications []Authentication
	PublicData      UserData
	PrivateData     UserData
	LastLogin       time.Time
}
