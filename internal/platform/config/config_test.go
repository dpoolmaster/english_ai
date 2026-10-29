package config

import (
	"errors"
	"log/slog"
	"os"
	"strings"
	"testing"
	"time"
)

type fakeLoader struct {
	env   map[string]string
	files map[string][]byte
}

func (f fakeLoader) loader() Loader {
	return Loader{
		LookupEnv: func(key string) (string, bool) {
			value, ok := f.env[key]
			return value, ok
		},
		ReadFile: func(path string) ([]byte, error) {
			value, ok := f.files[path]
			if !ok {
				return nil, os.ErrNotExist
			}
			return value, nil
		},
	}
}

func validEnv() map[string]string {
	return map[string]string{
		"APP_ENV":      "development",
		"APP_BASE_URL": "http://localhost:8080",
	}
}

func TestLoadDefaults(t *testing.T) {
	cfg, err := (fakeLoader{env: validEnv()}).loader().Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.HTTPAddr != ":8080" {
		t.Errorf("HTTPAddr = %q, want :8080", cfg.HTTPAddr)
	}
	if cfg.MetricsAddr != "127.0.0.1:9090" {
		t.Errorf("MetricsAddr = %q, want 127.0.0.1:9090", cfg.MetricsAddr)
	}
	if cfg.ShutdownTimeout != 10*time.Second {
		t.Errorf("ShutdownTimeout = %v, want 10s", cfg.ShutdownTimeout)
	}
	if cfg.LogLevel != slog.LevelInfo {
		t.Errorf("LogLevel = %v, want info", cfg.LogLevel)
	}
}

func TestLoadMissingRequiredValues(t *testing.T) {
	_, err := (fakeLoader{env: map[string]string{}}).loader().Load()
	if err == nil {
		t.Fatal("Load() error = nil, want error")
	}
	for _, field := range []string{"APP_ENV", "APP_BASE_URL"} {
		if !strings.Contains(err.Error(), field) {
			t.Errorf("error %q does not mention %s", err, field)
		}
	}
}

func TestLoadInvalidValues(t *testing.T) {
	tests := []struct {
		name  string
		env   map[string]string
		field string
	}{
		{
			name:  "unknown environment",
			env:   map[string]string{"APP_ENV": "staging", "APP_BASE_URL": "http://localhost"},
			field: "APP_ENV",
		},
		{
			name:  "relative base url",
			env:   map[string]string{"APP_ENV": "development", "APP_BASE_URL": "/relative"},
			field: "APP_BASE_URL",
		},
		{
			name:  "unsupported base url scheme",
			env:   map[string]string{"APP_ENV": "development", "APP_BASE_URL": "ftp://example.com"},
			field: "APP_BASE_URL",
		},
		{
			name:  "invalid shutdown timeout",
			env:   map[string]string{"APP_ENV": "development", "APP_BASE_URL": "http://localhost", "SHUTDOWN_TIMEOUT": "soon"},
			field: "SHUTDOWN_TIMEOUT",
		},
		{
			name:  "zero shutdown timeout",
			env:   map[string]string{"APP_ENV": "development", "APP_BASE_URL": "http://localhost", "SHUTDOWN_TIMEOUT": "0s"},
			field: "SHUTDOWN_TIMEOUT",
		},
		{
			name:  "invalid log level",
			env:   map[string]string{"APP_ENV": "development", "APP_BASE_URL": "http://localhost", "LOG_LEVEL": "verbose"},
			field: "LOG_LEVEL",
		},
		{
			name:  "empty http addr",
			env:   map[string]string{"APP_ENV": "development", "APP_BASE_URL": "http://localhost", "HTTP_ADDR": ""},
			field: "HTTP_ADDR",
		},
		{
			name:  "empty metrics addr",
			env:   map[string]string{"APP_ENV": "development", "APP_BASE_URL": "http://localhost", "METRICS_ADDR": ""},
			field: "METRICS_ADDR",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := (fakeLoader{env: test.env}).loader().Load()
			if err == nil {
				t.Fatal("Load() error = nil, want error")
			}
			if !strings.Contains(err.Error(), test.field) {
				t.Errorf("error %q does not mention %s", err, test.field)
			}
		})
	}
}

func TestLoadDatabaseSecretNeverAppearsInError(t *testing.T) {
	const secret = "super-secret-value"
	env := validEnv()
	env["DATABASE_URL"] = "postgres://user:" + secret

	_, err := (fakeLoader{env: env}).loader().Load()
	if err == nil {
		t.Fatal("Load() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "DATABASE_URL") {
		t.Errorf("error %q does not mention DATABASE_URL", err)
	}
	if strings.Contains(err.Error(), secret) {
		t.Fatalf("error %q leaks secret value", err)
	}
}

func TestLoadValidDatabaseURL(t *testing.T) {
	env := validEnv()
	env["DATABASE_URL"] = "postgres://user:password@localhost:5432/english?sslmode=disable"

	cfg, err := (fakeLoader{env: env}).loader().Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.DatabaseURL != env["DATABASE_URL"] {
		t.Errorf("DatabaseURL = %q, want configured value", cfg.DatabaseURL)
	}
}

func TestFileOverridesDefaultsAndEnvOverridesFile(t *testing.T) {
	files := map[string][]byte{
		"/etc/app.yaml": []byte(`
env: production
base_url: https://example.com
http_addr: ":8081"
metrics_addr: "127.0.0.1:9091"
shutdown_timeout: 5s
log_level: warn
`),
	}
	env := map[string]string{"APP_CONFIG_FILE": "/etc/app.yaml", "LOG_LEVEL": "error"}

	cfg, err := (fakeLoader{env: env, files: files}).loader().Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Env != EnvProduction {
		t.Errorf("Env = %q, want production", cfg.Env)
	}
	if cfg.HTTPAddr != ":8081" {
		t.Errorf("HTTPAddr = %q, want :8081", cfg.HTTPAddr)
	}
	if cfg.ShutdownTimeout != 5*time.Second {
		t.Errorf("ShutdownTimeout = %v, want 5s", cfg.ShutdownTimeout)
	}
	if cfg.LogLevel != slog.LevelError {
		t.Errorf("LogLevel = %v, want error from env", cfg.LogLevel)
	}
}

func TestFileRejectsSecretKeys(t *testing.T) {
	files := map[string][]byte{
		"/etc/app.yaml": []byte(`
env: development
base_url: http://localhost
database_url: postgres://user:super-secret@localhost/db
`),
	}

	_, err := (fakeLoader{env: map[string]string{"APP_CONFIG_FILE": "/etc/app.yaml"}, files: files}).loader().Load()
	if err == nil {
		t.Fatal("Load() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "database_url") {
		t.Errorf("error %q does not name the secret field", err)
	}
	if strings.Contains(err.Error(), "super-secret") {
		t.Fatalf("error %q leaks secret value", err)
	}
}

func TestFileInvalidYAMLDoesNotLeakValues(t *testing.T) {
	files := map[string][]byte{
		"/etc/app.yaml": []byte("env: [unterminated\nsecret-value"),
	}

	_, err := (fakeLoader{env: map[string]string{"APP_CONFIG_FILE": "/etc/app.yaml"}, files: files}).loader().Load()
	if err == nil {
		t.Fatal("Load() error = nil, want error")
	}
	if strings.Contains(err.Error(), "secret-value") {
		t.Fatalf("error %q leaks file content", err)
	}
}

func TestFileRejectsUnknownFields(t *testing.T) {
	files := map[string][]byte{
		"/etc/app.yaml": []byte("env: development\nbase_url: http://localhost\nunknown_field: value\n"),
	}

	_, err := (fakeLoader{env: map[string]string{"APP_CONFIG_FILE": "/etc/app.yaml"}, files: files}).loader().Load()
	if err == nil {
		t.Fatal("Load() error = nil, want error")
	}
}

func TestEmptyFileFallsBackToEnv(t *testing.T) {
	files := map[string][]byte{"/etc/app.yaml": []byte("")}
	env := validEnv()
	env["APP_CONFIG_FILE"] = "/etc/app.yaml"

	cfg, err := (fakeLoader{env: env, files: files}).loader().Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Env != EnvDevelopment {
		t.Errorf("Env = %q, want development", cfg.Env)
	}
}

func TestFileReadErrorIsSafe(t *testing.T) {
	_, err := (fakeLoader{env: map[string]string{"APP_CONFIG_FILE": "/missing.yaml"}}).loader().Load()
	if err == nil {
		t.Fatal("Load() error = nil, want error")
	}
	if errors.Is(err, os.ErrNotExist) {
		t.Fatalf("error %v exposes underlying filesystem error", err)
	}
}
