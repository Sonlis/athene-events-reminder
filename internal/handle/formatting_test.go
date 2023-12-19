package handle

import (
	"github.com/Sonlis/athene-events-reminder/internal/event"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"testing"
	"time"
)

func TestBuildEventInfoMarkup(t *testing.T) {
	reminderSet := true
	registrationDate, err := time.Parse("02.01.2006", "01.01.2021")
	if err != nil {
		panic(err)
	}
	event := event.Event{
		ID:                    1,
		Title:                 "This is a test event",
		RegistrationStartDate: registrationDate,
	}
	eventInfoMarkup := buildEventInfoMarkup(event, reminderSet)
	want := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
			{
				tgbotapi.NewInlineKeyboardButtonData("Remove reminder", "removeReminder 1"),
			},
			{
				tgbotapi.NewInlineKeyboardButtonData(backButton, backButton+" events"),
			},
		},
	}
	for i, row := range eventInfoMarkup.InlineKeyboard {
		for j, button := range row {
			if button.Text != want.InlineKeyboard[i][j].Text {
				t.Errorf("buildEventInfoMarkup() = %v, want %v", eventInfoMarkup, want)
			}
		}
	}
	reminderSet = false
	eventInfoMarkup = buildEventInfoMarkup(event, reminderSet)
	want = tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
			{
				tgbotapi.NewInlineKeyboardButtonData("Set a reminder", "setReminder 1 2021-01-01T00:00:00Z"),
			},
			{
				tgbotapi.NewInlineKeyboardButtonData(backButton, backButton+" events"),
			},
		},
	}
	for i, row := range eventInfoMarkup.InlineKeyboard {
		for j, button := range row {
			if button.Text != want.InlineKeyboard[i][j].Text {
				t.Errorf("buildEventInfoMarkup() = %v, want %v", eventInfoMarkup, want)
			}
		}
	}
}

func TestBuildEventMarkup(t *testing.T) {
	event := []event.Event{
		{
			ID:                    1,
			Title:                 "This is a test event",
			RegistrationStartDate: time.Now(),
		},
		{
			ID:                    2,
			Title:                 "This is another test event",
			RegistrationStartDate: time.Now().AddDate(0, 0, 1),
		},
	}
	eventMarkup := buildEventsMarkup(event)
	want := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
			{
				tgbotapi.NewInlineKeyboardButtonData("This is a test event", "1"),
			},
			{
				tgbotapi.NewInlineKeyboardButtonData("This is another test event", "2"),
			},
		},
	}
	for i, row := range eventMarkup.InlineKeyboard {
		for j, button := range row {
			if button.Text != want.InlineKeyboard[i][j].Text {
				t.Errorf("buildEventMarkup() = %v, want %v", eventMarkup, want)
			}
		}
	}
}

func TestFormatEventInfo(t *testing.T) {
	reminderSet := true
	eventDate, err := time.Parse("02.01.2006", "01.01.2021")
	if err != nil {
		t.Errorf("Error parsing date: %v", err)
	}
	event := event.Event{
		ID:                    1,
		Title:                 "This is a test event",
		RegistrationStartDate: eventDate,
	}

	want := "<b>This is a test event</b>\nRegistration date: Jan  1 00:00:00\nReminder is set"
	got := formatEventInfo(event, reminderSet)
	if got != want {
		t.Errorf("formatEventInfo() = %v, want %v", got, want)
	}

	reminderSet = false
	want = "<b>This is a test event</b>\nRegistration date: Jan  1 00:00:00\nReminder is not set"
	got = formatEventInfo(event, reminderSet)
	if got != want {
		t.Errorf("formatEventInfo() = %v, want %v", got, want)
	}
}
