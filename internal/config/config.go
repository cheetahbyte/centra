package config

import (
	"sync"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	// Core
	Port        string `env:"PORT" envDefault:"3000"`
	ContentRoot string `env:"CONTENT_ROOT" envDefault:"/content"`
	APIKey      string `env:"CENTRA_API_KEY"`

	// Git & Keys
	GitRepo       string `env:"GITHUB_REPO_URL"`
	KeysDir       string `env:"KEYS_DIR" envDefault:"/keys"`
	PrivateKey    string `env:"SSH_PRIVATE_KEY"`
	PublicKey     string `env:"SSH_PUBLIC_KEY"`
	WebhookSecret string `env:"WEBHOOK_SECRET"`

	// CORS (The library handles slices automatically via comma separation)
	AllowedOrigins []string `env:"CORS_ALLOWED_ORIGINS" envDefault:"*"`
	AllowedMethods []string `env:"CORS_ALLOWED_METHODS" envDefault:"GET,HEAD,OPTIONS"`
	AllowedHeaders []string `env:"CORS_ALLOWED_HEADERS" envDefault:"*"`
	ExposedHeaders []string `env:"CORS_EXPOSED_HEADERS" envDefault:"Cache-Control,Content-Language,Content-Type,Expires,Last-Modified"`
	MaxAge         int      `env:"CORS_MAX_AGE" envDefault:"360"`
	Credentials    bool     `env:"CORS_ALLOW_CREDENTIALS"`

	// Logging & Limits
	LogLevel      string `env:"LOG_LEVEL" envDefault:"INFO"`
	LogStructured bool   `env:"LOG_STRUC"`
	RateQuota     int    `env:"RATELIMIT_QUOTA" envDefault:"100"`

	// Features
	CacheBinaries bool `env:"CACHE_BINARIES"`
}

var (
	cfg  Config
	once sync.Once
)

// Load reads config once. Call this at app startup.
func Load() (*Config, error) {
	var err error
	once.Do(func() {
		// This parses environment variables into the struct
		err = env.Parse(&cfg)
	})
	return &cfg, err
}

func Get() Config {
	_, _ = Load()

	return cfg
}
