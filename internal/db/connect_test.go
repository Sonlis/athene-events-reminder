package db

import (
	"context"
	"github.com/sethvargo/go-envconfig"
	"os"
	"testing"
)

func initTestConfig() (Config, error) {
	var config Config
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5447")
	os.Setenv("DB_USER", "test_user")
	os.Setenv("DB_PASSWORD", "test_pass")
	os.Setenv("DB_NAME", "test")
	err := envconfig.Process(context.Background(), &config)
	return config, err
}

func TestConnect(t *testing.T) {
	config, err := initTestConfig()
	if err != nil {
		t.Errorf("Failed to initialize test config: %v", err)
	}
	db, err := Connect(&config)
	if err != nil {
		t.Errorf("Failed to initialize db config: %v", err)
	}
	if db == nil {
		t.Errorf("Connecting to the db returned nil")
	}
}

func TestConnString(t *testing.T) {
	config, err := initTestConfig()
	if err != nil {
		t.Errorf("Failed to initialize test config: %v", err)
	}
	connString := config.ConnString()
	want := "postgres://test_user:test_pass@localhost:5447/test"
	if connString != want {
		t.Errorf("ConnString() = %v, want %v", connString, want)
	}
}
