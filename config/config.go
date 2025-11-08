package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type OAuthConfig struct {
	Name string `yaml:"name"`
	Url  string `yaml:"url"`

	SkipDiscovery bool `yaml:"skip_discovery"`

	Endpoints struct {
		AuthURL     string `yaml:"auth"`
		TokenURL    string `yaml:"token"`
		UserInfoURL string `yaml:"user_info"`
	} `yaml:"endpoints"`

	ClientId     string   `yaml:"client_id"`
	ClientSecret string   `yaml:"client_secret"`
	Scopes       []string `yaml:"scopes"`
}

type SystemTokenConfig struct {
	Name  string `yaml:"name"`
	Token string `yaml:"token"`
}

type Config struct {
	ConfigPath string `env:"CONFIG_PATH" env-default:"config.yml"`

	CookieName     string   `yaml:"cookie_name" env:"COOKIE_NAME" env-default:"katok"`
	CookieDomain   string   `yaml:"cookie_domain" env:"COOKIE_DOMAIN" env-default:"localhost"`
	AllowedOrigins []string `yaml:"allowed_origins" env:"ALLOWED_ORIGINS" env-default:"http://localhost:14414"`

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

	Token struct {
		// DEPRECATED: Use PrivateKeyPath instead for RSA signing
		KeyPhrase string `yaml:"key_phrase" env:"KEY_Phrase"`

		// Path to RSA private key file for RS512 signing
		PrivateKeyPath string `yaml:"private_key_path" env:"PRIVATE_KEY_PATH" env-default:"private_key.pem"`

		// Path to RSA public key file for RS512 verification
		PublicKeyPath string `yaml:"public_key_path" env:"PUBLIC_KEY_PATH" env-default:"public_key.pem"`

		// Time to live for the refresh token (in seconds). The refresh token is a long-lived cookie and bound to the core domain of the auth service.
		// It'll be used to generate short-lived auth-tokens that can be used across your service.
		// Default: 864000 (10 days)
		RefreshTokenTTL uint `yaml:"refresh_token_ttl" env:"REFRESH_TOKEN_TTL" env-default:"864000"`

		// minimum age of a refresh token before a new one is sent back with any request.
		//
		// Default: 86400 (one day)
		MinAgeForRefresh uint `yaml:"refresh_token_min_age" env:"REFRESH_TOKEN_TTL" env-default:"86400"`

		// Time to live for the auth token (in seconds). Should be very small (1-2 minutes is good).
		//
		// Default: 90 (1.5 minutes)
		AuthTokenTTL uint `yaml:"auth_token_ttl" env:"AUTH_TOKEN_TTL" env-default:"90"`

		Issuer string `yaml:"issuer" env:"JWT_ISSUER" env-default:"kingdom-auth"`

		DefaultAudience string `yaml:"default_audience" env:"JWT_DEFAULT_AUDIENCE" env-default:"default-audience"`
	} `yaml:"token"`

	MainService struct {
		Port      int    `yaml:"port" env:"MAIN_PORT" env-default:"14414"`
		PublicUrl string `yaml:"public_url" env:"MAIN_PUBLIC_URL" env-default:"http://localhost:14414"`
	} `yaml:"main_service"`

	SystemService struct {
		Port   int                 `yaml:"port" env:"SYSTEM_PORT" env-default:"14415"`
		Tokens []SystemTokenConfig `yaml:"tokens"`
	} `yaml:"system_service"`
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
