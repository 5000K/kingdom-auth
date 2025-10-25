package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type OAuthConfig struct {
	Name         string `yaml:"name" envDefault:""`
	AuthUrl      string `yaml:"auth_url" envDefault:""`
	TokenUrl     string `yaml:"token_url" envDefault:""`
	ClientId     string `yaml:"client_id" envDefault:""`
	ClientSecret string `yaml:"client_secret" envDefault:""`
}

type Config struct {
	ConfigPath string `env:"CONFIG_PATH" envDefault:"./config.yml"`

	Db struct {
		// supported: "sqlite", "mysql", "postgres"
		// default: sqlite
		Type string `yaml:"type" env:"DB_TYPE" envDefault:"sqlite"`

		// DSN=Database Source Name, see https://github.com/go-sql-driver/mysql#dsn-data-source-name
		// default: kingdom-auth.db (useful default for the default of sqlite)
		DSN string `yaml:"dsn" env:"DB_DSN" envDefault:"kingdom-auth.db"`

		// Guaranteed to be non-destructive by gorm
		//
		// Is allowed to be deactivated to give users more control,
		// but might (and over time WILL) lead to the need to manually migrate after an update of kingdom-auth
		RunMigrations bool `yaml:"run_migrations" env:"DB_RUN_MIGRATIONS" envDefault:"true"`
	} `yaml:"db"`

	OAuthProviders []OAuthConfig `yaml:"providers"`

	MainService struct {
		Port      int    `yaml:"port" envDefault:"14414"`
		PublicUrl string `yaml:"public_url" envDefault:"http://localhost:14414"`
	} `yaml:"main_service"`
}

func Get() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, err
	}

	err = cleanenv.ReadConfig(cfg.ConfigPath, &cfg)

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
