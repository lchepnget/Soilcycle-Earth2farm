package config

import (
	"testing"
)

func TestGetEnvReturnsEnvironmentValue(t *testing.T) {
	t.Setenv("APP_PORT", "9000")

	got := getEnv("APP_PORT", "8090")

	want := "9000"

	if got != want {
		t.Errorf("getEnv() = %q, want %q", got, want)
	}
}
func TestGetEnvReturnsFallbackWhenEmpty(t *testing.T) {
	t.Setenv("APP_PORT", "")

	got := getEnv("APP_PORT", "8090")

	want := "8090"

	if got != want {
		t.Errorf("getEnv() = %q, want %q", got, want)
	}
}
func TestLoadReturnsDefaultConfiguration(t *testing.T) {
	t.Setenv("APP_PORT", "")
	t.Setenv("DB_HOST", "")
	t.Setenv("DB_PORT", "")
	t.Setenv("DB_USER", "")
	t.Setenv("DB_PASSWORD", "")
	t.Setenv("DB_NAME", "")
	t.Setenv("DB_SSLMODE", "")

	got := Load()

	want := &Config{
		AppPort:    "8090",
		DBHost:     "localhost",
		DBPort:     "5433",
		DBUser:     "soilcycle",
		DBPassword: "",
		DBName:     "soilcycle",
		DBSSLMode:  "disable",
	}

	if *got != *want {
		t.Errorf("Load() = %+v, want %+v", *got, *want)
	}
}
func TestLoadReturnsCustomConfiguration(t *testing.T) {
	t.Setenv("APP_PORT", "9000")
	t.Setenv("DB_HOST", "192.168.1.10")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_USER", "test_user")
	t.Setenv("DB_PASSWORD", "test_password")
	t.Setenv("DB_NAME", "test_database")
	t.Setenv("DB_SSLMODE", "require")

	got := Load()

	want := &Config{
		AppPort:    "9000",
		DBHost:     "192.168.1.10",
		DBPort:     "5432",
		DBUser:     "test_user",
		DBPassword: "test_password",
		DBName:     "test_database",
		DBSSLMode:  "require",
	}

	if *got != *want {
		t.Errorf("Load() = %+v, want %+v", *got, *want)
	}
}
