package db

import (
	"log/slog"

	"github.com/5000K/kingdom-auth/config"
	"github.com/5000K/kingdom-auth/core"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func connect(config *config.Config) (*gorm.DB, error) {
	switch config.Db.Type {
	case "mysql":
		slog.Debug("connecting to mysql")
		return gorm.Open(mysql.Open(config.Db.DSN), &gorm.Config{})

	case "postgres":
		slog.Debug("connecting to postgres")
		return gorm.Open(postgres.Open(config.Db.DSN), &gorm.Config{})

	case "sqlite":
		slog.Debug("opening sqlite")
		return gorm.Open(sqlite.Open(config.Db.DSN), &gorm.Config{})

	default:
		slog.Error("Unknown database type", "db-type", config.Db.Type)
		return nil, core.ErrUnknownDbDriver
	}
}
