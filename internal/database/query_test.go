package database

import (
	"context"
	"testing"
	"time"
)

func TestCreateReminder(t *testing.T) {
	config, err := InitTestConfig()
	if err != nil {
		t.Errorf("Failed to initialize test config: %v", err)
	}
	db, err := Connect(&config)
	if err != nil {
		t.Errorf("Failed to connec to db: %v", err)
	}
	timeToInsert := time.Date(2023, 12, 8, 16, 0, 0, 0, time.UTC)

	err = db.CreateReminder(context.Background(), 1, 1, timeToInsert)
	if err != nil {
		t.Errorf("Failed to create reminder: %v", err)
	}
	reminderExists, err := db.CheckReminder(context.Background(), 1, 1)
	if err != nil {
		t.Errorf("Failed to check if reminder exists: %v", err)
	}
	if !reminderExists {
		t.Errorf("Reminder does not exist")
	}
}

func TestRemoveReminder(t *testing.T) {
	config, err := InitTestConfig()
	if err != nil {
		t.Errorf("Failed to initialize test config: %v", err)
	}
	db, err := Connect(&config)
	if err != nil {
		t.Errorf("Failed to connec to db: %v", err)
	}
	timeToInsert := time.Date(2023, 12, 8, 16, 0, 0, 0, time.UTC)

	err = db.CreateReminder(context.Background(), 1, 1, timeToInsert)
	if err != nil {
		t.Errorf("Failed to create reminder: %v", err)
	}

	err = db.RemoveReminder(context.Background(), 1, 1)
	if err != nil {
		t.Errorf("Failed to remove reminder: %v", err)
	}
	reminderExists, err := db.CheckReminder(context.Background(), 1, 1)
	if err != nil {
		t.Errorf("Failed to check if reminder exists: %v", err)
	}
	if reminderExists {
		t.Errorf("Reminder exists")
	}
}

func TestCheckReminder(t *testing.T) {

	config, err := InitTestConfig()
	if err != nil {
		t.Errorf("Failed to initialize test config: %v", err)
	}
	db, err := Connect(&config)
	if err != nil {
		t.Errorf("Failed to connec to db: %v", err)
	}
	timeToInsert := time.Date(2023, 12, 8, 16, 0, 0, 0, time.UTC)

	err = db.CreateReminder(context.Background(), 1, 1, timeToInsert)
	if err != nil {
		t.Errorf("Failed to create reminder: %v", err)
	}

	reminderExists, err := db.CheckReminder(context.Background(), 1, 1)
	if err != nil {
		t.Errorf("Failed to check if reminder exists: %v", err)
	}
	if !reminderExists {
		t.Errorf("Reminder does not exist, but it should. Failed to create or check if the reminder exists")
	}

	err = db.RemoveReminder(context.Background(), 1, 1)
	if err != nil {
		t.Errorf("Failed to remove reminder: %v", err)
	}
	reminderExists, err = db.CheckReminder(context.Background(), 1, 1)
	if err != nil {
		t.Errorf("Failed to check if reminder exists: %v", err)
	}
	if reminderExists {
		t.Errorf("Reminder exists, but it should not. Failed to remove or check if the reminder exists")
	}
}
