package db

import (
	"log/slog"

	"github.com/5000K/kingdom-auth/config"
	"github.com/5000K/kingdom-auth/core"
	"gorm.io/gorm"
)

// Driver wrapper for a db
type Driver struct {
	db  *gorm.DB
	log *slog.Logger
}

func NewDriver(config *config.Config) (*Driver, error) {
	db, err := connect(config)

	if err != nil {
		return nil, err
	}

	driver := &Driver{
		db:  db,
		log: slog.With("source", "db.Driver"),
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

	err = d.db.AutoMigrate(&UserData{})
	if err != nil {
		return err
	}

	err = d.db.AutoMigrate(&Authentication{})
	if err != nil {
		return err
	}

	return nil
}
