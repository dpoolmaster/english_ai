package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Environment string

const (
	EnvDevelopment Environment = "development"
	EnvTest        Environment = "test"
	EnvProduction  Environment = "production"
)

const (
	defaultHTTPAddr        = ":8080"
	defaultMetricsAddr     = "127.0.0.1:9090"
	defaultShutdownTimeout = 10 * time.Second
)

// Config contains the process configuration. Secret values are populated only
// from environment variables and are never read from the YAML file.
type Config struct {
	Env             Environment
	BaseURL         string
	HTTPAddr        string
	MetricsAddr     string
	ShutdownTimeout time.Duration
	LogLevel        slog.Level
	DatabaseURL     string
}

// Loader reads configuration from an optional YAML file and environment
// variables. It is a struct so tests can replace the process dependencies.
type Loader struct {
	LookupEnv func(string) (string, bool)
	ReadFile  func(string) ([]byte, error)
}

func Load() (Config, error) {
	return Loader{LookupEnv: os.LookupEnv, ReadFile: os.ReadFile}.Load()
}

func (l Loader) Load() (Config, error) {
	cfg := Config{
		HTTPAddr:        defaultHTTPAddr,
		MetricsAddr:     defaultMetricsAddr,
		ShutdownTimeout: defaultShutdownTimeout,
		LogLevel:        slog.LevelInfo,
	}

	if path, ok := l.lookup("APP_CONFIG_FILE"); ok && strings.TrimSpace(path) != "" {
		if err := l.applyFile(path, &cfg); err != nil {
			return Config{}, err
		}
	}

	if err := l.applyEnv(&cfg); err != nil {
		return Config{}, err
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (l Loader) lookup(key string) (string, bool) {
	if l.LookupEnv == nil {
		return "", false
	}
	return l.LookupEnv(key)
}

type fileConfig struct {
	Env             string `yaml:"env"`
	BaseURL         string `yaml:"base_url"`
	HTTPAddr        string `yaml:"http_addr"`
	MetricsAddr     string `yaml:"metrics_addr"`
	ShutdownTimeout string `yaml:"shutdown_timeout"`
	LogLevel        string `yaml:"log_level"`
}

var secretFileKeys = map[string]struct{}{
	"database_url":               {},
	"media_sign_key":             {},
	"smtp_password":              {},
	"llm_primary_key":            {},
	"llm_fallback_key":           {},
	"llm_strong_key":             {},
	"yc_speech_key":              {},
	"yk_secret":                  {},
	"yandex_oauth_client_secret": {},
	"vk_oauth_client_secret":     {},
	"smartcaptcha_server_key":    {},
	"telegram_bot_token":         {},
}

func (l Loader) applyFile(path string, cfg *Config) error {
	if l.ReadFile == nil {
		return errors.New("APP_CONFIG_FILE cannot be read")
	}

	data, err := l.ReadFile(path)
	if err != nil {
		return errors.New("APP_CONFIG_FILE cannot be read")
	}

	var keys map[string]yaml.Node
	if err := yaml.Unmarshal(data, &keys); err != nil {
		return errors.New("APP_CONFIG_FILE contains invalid YAML")
	}
	for key := range keys {
		if _, secret := secretFileKeys[strings.ToLower(key)]; secret {
			return fmt.Errorf("APP_CONFIG_FILE must not contain secret %q; use environment variables", key)
		}
	}

	var file fileConfig
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err := decoder.Decode(&file); err != nil && !errors.Is(err, io.EOF) {
		return errors.New("APP_CONFIG_FILE contains unsupported or invalid fields")
	}

	if file.Env != "" {
		cfg.Env = Environment(strings.TrimSpace(file.Env))
	}
	if file.BaseURL != "" {
		cfg.BaseURL = strings.TrimSpace(file.BaseURL)
	}
	if file.HTTPAddr != "" {
		cfg.HTTPAddr = strings.TrimSpace(file.HTTPAddr)
	}
	if file.MetricsAddr != "" {
		cfg.MetricsAddr = strings.TrimSpace(file.MetricsAddr)
	}
	if file.ShutdownTimeout != "" {
		timeout, err := time.ParseDuration(strings.TrimSpace(file.ShutdownTimeout))
		if err != nil {
			return errors.New("APP_CONFIG_FILE shutdown_timeout must be a duration such as 10s")
		}
		cfg.ShutdownTimeout = timeout
	}
	if file.LogLevel != "" {
		level, err := ParseLogLevel(file.LogLevel)
		if err != nil {
			return errors.New("APP_CONFIG_FILE log_level must be one of debug, info, warn, error")
		}
		cfg.LogLevel = level
	}

	return nil
}

func (l Loader) applyEnv(cfg *Config) error {
	if value, ok := l.lookup("APP_ENV"); ok {
		cfg.Env = Environment(strings.TrimSpace(value))
	}
	if value, ok := l.lookup("APP_BASE_URL"); ok {
		cfg.BaseURL = strings.TrimSpace(value)
	}
	if value, ok := l.lookup("HTTP_ADDR"); ok {
		cfg.HTTPAddr = strings.TrimSpace(value)
	}
	if value, ok := l.lookup("METRICS_ADDR"); ok {
		cfg.MetricsAddr = strings.TrimSpace(value)
	}
	if value, ok := l.lookup("SHUTDOWN_TIMEOUT"); ok {
		timeout, err := time.ParseDuration(strings.TrimSpace(value))
		if err != nil {
			return errors.New("SHUTDOWN_TIMEOUT must be a duration such as 10s")
		}
		cfg.ShutdownTimeout = timeout
	}
	if value, ok := l.lookup("LOG_LEVEL"); ok {
		level, err := ParseLogLevel(value)
		if err != nil {
			return errors.New("LOG_LEVEL must be one of debug, info, warn, error")
		}
		cfg.LogLevel = level
	}
	if value, ok := l.lookup("DATABASE_URL"); ok {
		cfg.DatabaseURL = strings.TrimSpace(value)
	}

	return nil
}

func (c Config) Validate() error {
	var errs []error

	switch c.Env {
	case EnvDevelopment, EnvTest, EnvProduction:
	default:
		errs = append(errs, errors.New("APP_ENV must be one of development, test, production"))
	}

	if c.BaseURL == "" {
		errs = append(errs, errors.New("APP_BASE_URL is required"))
	} else if parsed, err := url.Parse(c.BaseURL); err != nil || parsed.Scheme == "" || parsed.Host == "" ||
		(parsed.Scheme != "http" && parsed.Scheme != "https") {
		errs = append(errs, errors.New("APP_BASE_URL must be an absolute http(s) URL"))
	}

	if c.HTTPAddr == "" {
		errs = append(errs, errors.New("HTTP_ADDR is required"))
	}
	if c.MetricsAddr == "" {
		errs = append(errs, errors.New("METRICS_ADDR is required"))
	}
	if c.ShutdownTimeout <= 0 {
		errs = append(errs, errors.New("SHUTDOWN_TIMEOUT must be greater than zero"))
	}

	if c.DatabaseURL != "" && !validDatabaseURL(c.DatabaseURL) {
		errs = append(errs, errors.New("DATABASE_URL must be a valid postgres URL"))
	}

	return errors.Join(errs...)
}

func validDatabaseURL(raw string) bool {
	parsed, err := url.Parse(raw)
	if err != nil {
		return false
	}
	if parsed.Scheme != "postgres" && parsed.Scheme != "postgresql" {
		return false
	}
	return parsed.Host != "" && parsed.Path != ""
}

func ParseLogLevel(raw string) (slog.Level, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn", "warning":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return 0, errors.New("unknown log level")
	}
}
