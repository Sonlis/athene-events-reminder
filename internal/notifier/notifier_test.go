package notifier

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/Sonlis/athene-events-reminder/internal/database"
	"github.com/Sonlis/athene-events-reminder/internal/event"
)

func TestGetReminders(t *testing.T) {
	ctx := context.Background()
	os.Setenv("ILMO_API_URL", "http://localhost:8080")
	i, err := event.Init()
	if err != nil {
		t.Errorf("Error while initializing event client: %v", err)
	}
	db_conf, err := database.InitTestConfig()
	if err != nil {
		t.Errorf("Error while initializing db client: %v", err)
	}

	db_conn, err := database.Connect(&db_conf)
	if err != nil {
		t.Errorf("Failed to connec to db: %v", err)
	}
	defer db_conn.Close(ctx)

	timeToInsert := time.Date(2023, 12, 8, 16, 0, 0, 0, time.UTC)

	currentTime := time.Now().Add(time.Second * 30).UTC()
	err = db_conn.CreateReminder(context.Background(), 1, 1, timeToInsert)
	if err != nil {
		t.Errorf("Failed to create reminder: %v", err)
	}

	err = db_conn.CreateReminder(context.Background(), 2, 2, currentTime)
	if err != nil {
		t.Errorf("Failed to create reminder: %v", err)
	}

	reminders, err := getReminders(ctx, db_conn, &i)
	if err != nil {
		t.Errorf("Error while getting reminders: %v", err)
	}
	want := Reminder{2, 2, currentTime}
	if len(reminders) != 1 {
		t.Errorf("Expected 1 reminder, got %v", len(reminders))
	}
	if !reminders[0].ReminderTime.Truncate(time.Second).Equal(want.ReminderTime.Truncate(time.Second)) {
		t.Errorf("Expected %v, got %v", want, reminders[0])
	}
}

func TestBuildEventsDescription(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/1" && r.URL.Path != "/2" {
			t.Errorf("Expected to request '/1' or '/2', got: %s", r.URL.Path)
		} else {
			switch r.URL.Path {
			case "/1":
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`
                {
                    "id": 1,
                    "title": "Test event 1"
                }
                `))
			case "/2":
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`
                {
                    "id": 2,
                    "title": "Test event 2"
                }
                `))
			}
		}
	}))
	defer server.Close()

	os.Setenv("ILMO_API_URL", server.URL)
	os.Setenv("ILMO_WEB_URL", "works_on_my_machine")

	reminders := []Reminder{
		{
			ChatId:  1,
			EventId: 1,
		},
		{
			ChatId:  2,
			EventId: 2,
		},
	}

	i, err := event.Init()
	if err != nil {
		t.Errorf("Error while initializing event client: %v", err)
	}

	want := map[int]string{
		1: "Registration for Test event 1 starts in 5 minutes at works_on_my_machine/1",
		2: "Registration for Test event 2 starts in 5 minutes at works_on_my_machine/2",
	}
	got, err := buildEventsDescription(reminders, &i)
	if err != nil {
		t.Errorf("Error building test event descriptiona: %v", err)
	}
	for position, description := range want {
		if description != got[position] {
			t.Errorf("Error: want %v, got %v", description, got[position])
		}
	}

}
