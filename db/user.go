package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type UserData map[string]any

type User struct {
	gorm.Model

	Authentications []Authentication
	PublicData      string
	PrivateData     string
	LastLogin       time.Time
}

func (u *User) GetPublicUserdata() (UserData, error) {
	ud := UserData{}
	return ud, ud.Scan(u.PublicData)
}

func (u *User) GetPrivateUserdata() (UserData, error) {
	ud := UserData{}
	return ud, ud.Scan(u.PrivateData)
}

func (u *User) SetPublicUserdata(ud UserData) error {
	serialized, err := ud.Value()

	if err != nil {
		return err
	}

	u.PublicData = serialized
	return nil
}

func (u *User) SetPrivateUserdata(ud UserData) error {
	serialized, err := ud.Value()

	if err != nil {
		return err
	}

	u.PrivateData = serialized
	return nil
}

func (v *UserData) Scan(value interface{}) error {
	str, ok := value.(string)

	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal SheetStorageState value:", value))
	}

	err := json.Unmarshal([]byte(str), v)

	if err != nil {
		return err
	}

	return nil
}

func (v *UserData) Value() (string, error) {
	if v == nil {
		return "{}", nil
	}

	bytes, err := json.Marshal(v)
	if err != nil {
		return "{}", err
	}
	return string(bytes), nil
}
