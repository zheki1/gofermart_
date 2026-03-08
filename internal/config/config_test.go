package config

import (
	"flag"
	"os"
	"testing"
)

func backupEnv(keys ...string) map[string]string {
	old := make(map[string]string)
	for _, k := range keys {
		old[k] = os.Getenv(k)
	}
	return old
}

func restoreEnv(old map[string]string) {
	for k, v := range old {
		os.Setenv(k, v)
	}
}

// Тест значений по умолчанию
func TestMustLoad_Defaults(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	old := backupEnv("RUN_ADDRESS", "DATABASE_URI", "ACCRUAL_SYSTEM_ADDRESS", "JWT_SECRET")
	defer restoreEnv(old)

	os.Unsetenv("RUN_ADDRESS")
	os.Unsetenv("DATABASE_URI")
	os.Unsetenv("ACCRUAL_SYSTEM_ADDRESS")
	os.Unsetenv("JWT_SECRET")

	cfg := MustLoadWithFlags(fs)

	if cfg.RunAddress != "localhost:8888" {
		t.Errorf("expected default RunAddress, got %s", cfg.RunAddress)
	}
	if cfg.DatabaseURI != "postgres://postgres:postgres@localhost:5432/gofermart?sslmode=disable" {
		t.Errorf("expected default DatabaseURI, got %s", cfg.DatabaseURI)
	}
	if cfg.AccrualSystemAddress != "localhost:9999" {
		t.Errorf("expected default AccrualSystemAddress, got %s", cfg.AccrualSystemAddress)
	}
	if cfg.JWTSecret != "super-secret-key" {
		t.Errorf("expected default JWTSecret got %s", cfg.JWTSecret)
	}
}

// Тест переменных окружения
func TestMustLoad_Env(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	old := backupEnv("RUN_ADDRESS", "DATABASE_URI", "ACCRUAL_SYSTEM_ADDRESS")
	defer restoreEnv(old)

	os.Setenv("RUN_ADDRESS", "env:8888")
	os.Setenv("DATABASE_URI", "env:dburi")
	os.Setenv("ACCRUAL_SYSTEM_ADDRESS", "env:9999")

	cfg := MustLoadWithFlags(fs)

	if cfg.RunAddress != "env:8888" {
		t.Errorf("expected RunAddress from env, got %s", cfg.RunAddress)
	}
	if cfg.DatabaseURI != "env:dburi" {
		t.Errorf("expected DatabaseURI from env, got %s", cfg.DatabaseURI)
	}
	if cfg.AccrualSystemAddress != "env:9999" {
		t.Errorf("expected AccrualSystemAddress from env, got %s", cfg.AccrualSystemAddress)
	}
}

// Тест функции getEnv
func TestGetEnv(t *testing.T) {
	old := os.Getenv("TEST_KEY")
	defer os.Setenv("TEST_KEY", old)

	os.Unsetenv("TEST_KEY")
	if getEnv("TEST_KEY", "default") != "default" {
		t.Error("expected default value when env is not set")
	}

	os.Setenv("TEST_KEY", "value")
	if getEnv("TEST_KEY", "default") != "value" {
		t.Error("expected value from env")
	}
}
