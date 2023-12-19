package handle

import (
	"context"
	"github.com/Sonlis/athene-events-reminder/internal/database"
	"github.com/Sonlis/athene-events-reminder/internal/event"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestHandleEventInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/1" {
			t.Errorf("Expected to request '/1' got: %s", r.URL.Path)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`
            {
                "id": 1,
                "title": "Test event 1",
                "registrationStartDate": "2021-06-03T10:00:00Z"
            }
            `))
		}
	}))
	ctx := context.Background()
	os.Setenv("ILMO_API_URL", server.URL)
	i, err := event.Init()
	if err != nil {
		t.Errorf("Error while initializing event client: %v", err)
	}
	db_conf, err := database.InitTestConfig()
	if err != nil {
		t.Errorf("Error while initializing db client: %v", err)
	}

	db, err := database.Connect(&db_conf)
	if err != nil {
		t.Errorf("Failed to connec to db: %v", err)
	}
	defer db.Close(ctx)
	reminderTime, err := time.Parse("02.01.2006", "03.06.2021")
	query := tgbotapi.CallbackQuery{
		Data: "eventInfo 1",
		Message: &tgbotapi.Message{
			MessageID: 2,
			Chat: &tgbotapi.Chat{
				ID: 3,
			},
		},
	}
	err = db.CreateReminder(ctx, query.Message.Chat.ID, 1, reminderTime)
	if err != nil {
		t.Errorf("Failed to create reminder: %v", err)
	}

	text, _, err := handleEventInfo(ctx, &i, query, db)
	if err != nil {
		t.Errorf("Error testing the listing of events: %v", err)
	}

	want := tgbotapi.MessageConfig{
		ParseMode: tgbotapi.ModeHTML,
		Text:      "<b>Test event 1</b>\nRegistration date: Jun  3 10:00:00\nReminder is set",
	}
	if text != want.Text {
		t.Errorf("Expected text: %v, got: %v", want.Text, text)
	}

	err = db.RemoveReminder(ctx, query.Message.Chat.ID, 1)
	if err != nil {
		t.Errorf("Failed to remove reminder: %v", err)
	}
	text, _, err = handleEventInfo(ctx, &i, query, db)
	if err != nil {
		t.Errorf("Error testing the listing of events: %v", err)
	}
	want = tgbotapi.MessageConfig{
		ParseMode: tgbotapi.ModeHTML,
		Text:      "<b>Test event 1</b>\nRegistration date: Jun  3 10:00:00\nReminder is not set",
	}
	if text != want.Text {
		t.Errorf("Expected text: %v, got: %v", want.Text, text)
	}
}
