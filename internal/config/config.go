package config

import (
	"strings"
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

	// CORS
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
	CacheBinaries   bool     `env:"CACHE_BINARIES"`
	AllowedBinaries []string `env:"ALLOWED_BINARIES" envDefault:"*"`
	AnyBinaries     bool     `env:"-"`
}

func (c *Config) Normalize() {
	// ALLOWED_BINARIES="*"
	if len(c.AllowedBinaries) == 1 && c.AllowedBinaries[0] == "*" {
		c.AnyBinaries = true
		c.AllowedBinaries = nil
		return
	}

	// normalize extensions: lowercase, ensure leading dot
	for i, ext := range c.AllowedBinaries {
		ext = strings.ToLower(strings.TrimSpace(ext))
		if ext != "" && !strings.HasPrefix(ext, ".") {
			ext = "." + ext
		}
		c.AllowedBinaries[i] = ext
	}
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
	cfg.Normalize()
	InitBinaryAllowList(&cfg)
	return &cfg, err
}

func Get() Config {
	_, _ = Load()

	return cfg
}

var BinaryAllowList map[string]bool

func InitBinaryAllowList(conf *Config) {
	BinaryAllowList = make(map[string]bool)
	for _, ext := range conf.AllowedBinaries {
		BinaryAllowList[ext] = true
	}
}
