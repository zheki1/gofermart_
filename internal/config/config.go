package config

import (
	"flag"
	"os"
)

// Config хранит настройки приложения.
type Config struct {
	RunAddress           string // адрес, на котором запускается сервер
	DatabaseURI          string // URI базы данных
	AccrualSystemAddress string // адрес внешней системы начислений
	JWTSecret            string // ключ подписи JWT
}

// MustLoad загружает конфигурацию приложения (основная для production).
func MustLoad() *Config {
	return MustLoadWithFlags(nil)
}

// MustLoadWithFlags загружает конфигурацию через указанный FlagSet.
// Если fs == nil, используется глобальный flag.CommandLine.
func MustLoadWithFlags(fs *flag.FlagSet) *Config {
	cfg := &Config{}

	if fs == nil {
		fs = flag.CommandLine
	}

	fs.StringVar(&cfg.RunAddress, "a", getEnv("RUN_ADDRESS", "localhost:8888"), "run address")
	fs.StringVar(&cfg.DatabaseURI, "d", getEnv("DATABASE_URI", "postgres://postgres:postgres@localhost:5432/gofermart?sslmode=disable"), "database uri")
	fs.StringVar(&cfg.AccrualSystemAddress, "r", getEnv("ACCRUAL_SYSTEM_ADDRESS", "localhost:9999"), "accrual address")
	fs.StringVar(&cfg.JWTSecret, "s", getEnv("JWT_SECRET", "super-secret-key"), "jwt secret")

	// для тестового FlagSet передаем пустой слайс аргументов
	if fs == flag.CommandLine {
		fs.Parse(os.Args[1:])
	} else {
		fs.Parse([]string{})
	}

	return cfg
}

// getEnv возвращает значение переменной окружения по ключу key.
// Если переменная не задана, возвращает значение по умолчанию def.
func getEnv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
