package db

import (
	"log/slog"
	"time"

	"github.com/5000K/kingdom-auth/config"
	"github.com/5000K/kingdom-auth/core"
	"gorm.io/gorm"
)

// Driver wrapper for a db
type Driver struct {
	db  *gorm.DB
	log *slog.Logger
	cfg *config.Config
}

func NewDriver(config *config.Config) (*Driver, error) {
	db, err := connect(config)

	if err != nil {
		return nil, err
	}

	driver := &Driver{
		db:  db,
		log: slog.With("source", "db.Driver"),
		cfg: config,
	}

	if config.Db.RunMigrations {
		err := driver.migrate()

		if err != nil {
			driver.log.Error("Failed to run migrations", "error", err)

			return nil, core.ErrFailedMigration
		}
	} else {
		driver.log.Warn("Running migrations was disabled by your configuration - this means you need to take care of schema-migrations yourself when updating kingdom-auth.")
	}

	return driver, nil
}

func (d *Driver) migrate() error {
	err := d.db.AutoMigrate(&User{})
	if err != nil {
		return err
	}

	err = d.db.AutoMigrate(&Authentication{})
	if err != nil {
		return err
	}

	return nil
}

func (d *Driver) CreateUser() (*User, error) {
	user := User{
		Authentications: make([]Authentication, 0),
		PrivateData:     "{}",
		PublicData:      "{}",
		LastLogin:       time.UnixMilli(0),
	}

	_ = user.SetPublicUserdata(UserData{
		"aud": d.cfg.Token.DefaultAudience,
	})

	return &user, d.db.Create(&user).Error
}

func (d *Driver) UpdateUser(user *User) error {
	return d.db.Save(user).Error
}

func (d *Driver) GetUser(id uint32) (*User, error) {
	user := User{}
	return &user, d.db.Preload("Authentications").First(&user, id).Error
}

func (d *Driver) TryGetAuthentication(provider string, subject string) (*Authentication, error) {
	auth := Authentication{}
	return &auth, d.db.First(&auth, "provider = ? AND subject = ?", provider, subject).Error
}

func (d *Driver) GetUserFor(auth *Authentication) (*User, error) {
	user := User{}
	return &user, d.db.Preload("Authentications").First(&user, "id = ?", auth.UserID).Error
}
func (d *Driver) CreateAuthenticationFor(user *User) (*Authentication, error) {
	auth := Authentication{
		UserID: user.ID,
	}

	err := d.db.Create(&auth).Error
	if err != nil {
		return nil, err
	}

	user.Authentications = append(user.Authentications, auth)

	return &auth, d.db.Save(user).Error
}

func (d *Driver) UpdateAuthentication(auth *Authentication) error {
	return d.db.Save(auth).Error
}
