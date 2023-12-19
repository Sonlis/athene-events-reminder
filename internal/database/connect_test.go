package database

import (
	"testing"
)

func TestConnect(t *testing.T) {
	config, err := InitTestConfig()
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
	config, err := InitTestConfig()
	if err != nil {
		t.Errorf("Failed to initialize test config: %v", err)
	}
	connString := config.ConnString()
	want := "postgres://test_user:test_pass@localhost:5447/test"
	if connString != want {
		t.Errorf("ConnString() = %v, want %v", connString, want)
	}
}
