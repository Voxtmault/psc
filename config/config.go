package config

import (
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

var loadedConfig AppConfig

type DBConfig struct {
	DBHost string `validate:"required"`
	DBPort string `validate:"required"`
	DBUser string `validate:"required"`
	DBPass string `validate:"required"`
	DBName string `validate:"required"`
}

type RedisConfig struct {
}

type LoggingConfig struct{}

type SSLConfig struct {
	CertPath string `validate:"omitempty"`
	KeyPath  string `validate:"omitempty"`
}

type AppConfig struct {
	Mode     string `validate:"required,oneof=dev prod test debug"`
	TZ       string `validate:"required"` // Asia/Jakarta, Europe/Paris, etc.
	HTTPPort string `validate:"required"`
	HTTPRoot string `validate:"required"`

	DBConfig
	SSLConfig
	RedisConfig
	LoggingConfig
}

func New(envPath string) *AppConfig {
	if err := godotenv.Load(envPath); err != nil {
		slog.Warn("failed to load provided file", "path", envPath, "error", err, "action", "proceeding with defaults / previously set env vars")
	}

	loadedConfig = AppConfig{
		Mode: GetEnv("MODE", "dev"),
		TZ:   GetEnv("TZ", "Asia/Jakarta"),
		DBConfig: DBConfig{
			DBHost: GetEnv("DB_HOST", "localhost"),
			DBPort: GetEnv("DB_PORT", "3306"),
			DBUser: GetEnv("DB_USER", "root"),
			DBPass: GetEnv("DB_PASS", "root"),
			DBName: GetEnv("DB_NAME", "psc"),
		},
	}

	return &loadedConfig
}

// Simple helper function to read an environment or return a default value.
func GetEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	if nextValue := os.Getenv(key); nextValue != "" {
		return nextValue
	}

	return defaultVal
}

// Simple helper function to read an environment variable into integer or return a default value.
func GetEnvAsInt(name string, defaultVal int) int {
	valueStr := GetEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

// Helper to read an environment variable into a bool or return default value.
func GetEnvAsBool(name string, defaultVal bool) bool {
	valStr := GetEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}

	return defaultVal
}

// Helper to read an environment variable into a slice of a specific type or return default value.
func GetEnvAsSlice[T any](name string, defaultVal []T, sep string) []T {
	valStr := GetEnv(name, "")

	if valStr == "" {
		return defaultVal
	}

	vals := strings.Split(valStr, sep)
	result := make([]T, len(vals))

	for i, v := range vals {
		switch any(result).(type) {
		case []string:
			result[i] = any(v).(T)
		case []int:
			intVal, _ := strconv.Atoi(v)
			result[i] = any(intVal).(T)
		case []bool:
			boolVal, _ := strconv.ParseBool(v)
			result[i] = any(boolVal).(T)
		default:
			return defaultVal
		}
	}

	return result
}
