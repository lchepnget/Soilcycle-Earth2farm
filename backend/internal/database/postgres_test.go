package database

import (
	"strings"
	"testing"

	"github.com/soilcycle-earth2farm/backend/internal/config"
)

func TestConnectReturnsErrorForInvalidPort(t *testing.T) {
	cfg := &config.Config{
		DBHost:    "localhost",
		DBPort:    "not-a-port",
		DBUser:    "soilcycle",
		DBName:    "soilcycle",
		DBSSLMode: "disable",
	}

	pool, err := Connect(cfg)

	if err == nil {
		t.Fatal("Connect() expected an error for an invalid port, got nil")
	}

	if pool != nil {
		pool.Close()
		t.Fatal("Connect() expected a nil pool when configuration is invalid")
	}

	if !strings.Contains(err.Error(), "parse database config:") {
		t.Errorf("Connect() error = %q, want an error containing %q",
			err.Error(), "parse database config:")
	}
}
func TestConnectReturnsErrorForInvalidSSLMode(t *testing.T) {
	cfg := &config.Config{
		DBHost:    "localhost",
		DBPort:    "5432",
		DBUser:    "soilcycle",
		DBName:    "soilcycle",
		DBSSLMode: "invalid-mode",
	}

	pool, err := Connect(cfg)

	if err == nil {
		t.Fatal("Connect() expected an error for an invalid SSL mode, got nil")
	}

	if pool != nil {
		pool.Close()
		t.Fatal("Connect() expected a nil pool when configuration is invalid")
	}

	if !strings.Contains(err.Error(), "parse database config:") {
		t.Errorf("Connect() error = %q, want an error containing %q",
			err.Error(), "parse database config:")
	}
}
func TestConnectReturnsErrorWhenDatabaseUnavailable(t *testing.T) {
	cfg := &config.Config{
		DBHost:    "127.0.0.1",
		DBPort:    "1",
		DBUser:    "soilcycle",
		DBName:    "soilcycle",
		DBSSLMode: "disable",
	}

	pool, err := Connect(cfg)

	if err == nil {
		t.Fatal("Connect() expected an error when the database is unavailable, got nil")
	}

	if pool != nil {
		pool.Close()
		t.Fatal("Connect() expected a nil pool when the database connection fails")
	}

	if !strings.Contains(err.Error(), "database ping failed:") {
		t.Errorf("Connect() error = %q, want an error containing %q",
			err.Error(), "database ping failed:")
	}
}
