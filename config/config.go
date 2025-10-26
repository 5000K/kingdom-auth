package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type OAuthConfig struct {
	Name         string   `yaml:"name"`
	Url          string   `yaml:"url"`
	ClientId     string   `yaml:"client_id"`
	ClientSecret string   `yaml:"client_secret"`
	Scopes       []string `yaml:"scopes"`
}

type Config struct {
	ConfigPath string `env:"CONFIG_PATH" env-default:"config.yml"`

	KeyPhrase string `yaml:"key_phrase" env:"KEY_Phrase"`

	CookieName   string `yaml:"cookie_name" env:"COOKIE_NAME" env-default:"ka_token"`
	CookieDomain string `yaml:"cookie_domain" env:"COOKIE_DOMAIN" env-default:"localhost"`

	Db struct {
		// supported: "sqlite", "mysql", "postgres"
		// default: sqlite
		Type string `yaml:"type" env:"DB_TYPE" env-default:"sqlite"`

		// DSN=Database Source Name, see https://github.com/go-sql-driver/mysql#dsn-data-source-name
		// default: kingdom-auth.db (useful default for the default of sqlite)
		DSN string `yaml:"dsn" env:"DB_DSN" env-default:"kingdom-auth.db"`

		// Guaranteed to be non-destructive by gorm
		//
		// Is allowed to be deactivated to give users more control,
		// but might (and over time WILL) lead to the need to manually migrate after an update of kingdom-auth
		RunMigrations bool `yaml:"run_migrations" env:"DB_RUN_MIGRATIONS" env-default:"true"`
	} `yaml:"db"`

	OAuthProviders []OAuthConfig `yaml:"providers"`

	MainService struct {
		Port      int    `yaml:"port" env-default:"14414"`
		PublicUrl string `yaml:"public_url" env-default:"http://localhost:14414"`
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
