package handle

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Sonlis/athene-events-reminder/internal/event"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
	"reflect"
)

func TestHandleListEvents(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			t.Errorf("Expected to request '/' got: %s", r.URL.Path)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`
            [
                {
                    "id": 1,
                    "title": "Test event 1"
                },
                {
                    "id": 2,
                    "title": "Test event 2"
                },
                {
                    "Wrong": ""Field,
                    "Testing": "F*** it we ball"
                }
            ]
            `))
		}
	}))
	defer server.Close()

	os.Setenv("ILMO_API_URL", server.URL)
	i, err := event.Init()
	if err != nil {
		t.Errorf("Error while initializing event client: %v", err)
	}
	msg, err := handleListEvents(1, &i)
	if err != nil {
		t.Errorf("Error testing the listing of events: %v", err)
	}
	want := tgbotapi.MessageConfig{
		ParseMode: tgbotapi.ModeHTML,
	}
	want.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Test event 1", "eventInfo 1"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Test event 2", "eventInfo 2"),
		),
	)
	if reflect.DeepEqual(want.ReplyMarkup, msg.ReplyMarkup) {
		t.Errorf("Expected ReplyMarkup: %v, got: %v", want.ReplyMarkup, msg.ReplyMarkup)
	}
	if msg.ParseMode != want.ParseMode {
		t.Errorf("Expected ParseMode: %v, got: %v", want.ParseMode, msg.ParseMode)
	}
}
