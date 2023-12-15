package handle

import (
	"fmt"
	"github.com/Sonlis/athene-events-notifier/internal/event"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"time"
)

var backButton = "Back"

// buildEventsMarkup builds the markup structure returned when a single event is requested.
func buildEventInfoMarkup(event event.Event, reminderSet bool) tgbotapi.InlineKeyboardMarkup {
	// Set the reminde to be 5 minutes before the event registration starts.
	// Also remove 2 hours to the time, because the event is in UTC+2, while the db time is in UTC.
	reminderTime := event.RegistrationStartDate.Add(time.Duration(-5) * time.Minute).Add(time.Duration(-2) * time.Hour).Format("2006-01-02T15:04:05Z07:00")
	if reminderSet {
		return tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Remove reminder", "removeReminder "+strconv.Itoa(event.ID)),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(backButton, backButton+" events"),
			),
		)
	} else {
		return tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Set a reminder", "setReminder "+strconv.Itoa(event.ID)+" "+reminderTime),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(backButton, backButton+" events"),
			),
		)
	}
}

// buildEventsMarkup builds the markup structure returned when all events are requested.
func buildEventsMarkup(events []event.Event) tgbotapi.InlineKeyboardMarkup {
	eventRows := [][]tgbotapi.InlineKeyboardButton{}
	for _, event := range events {
		eventRow := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(event.Title, "eventInfo "+strconv.Itoa(event.ID)),
		)
		eventRows = append(eventRows, eventRow)
	}
	return tgbotapi.NewInlineKeyboardMarkup(eventRows...)
}

// buildReminderMarkup builds the markup structure returned when a reminder is created or deleted.
// It gives the eventId in the answer, so it will go back to the event info when clicked.
func buildReminderMarkup(eventId int) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(backButton, backButton+" event "+strconv.Itoa(eventId)),
		),
	)
}

// formatEventInfo returns a string containing the event's info, formatted in markdown for the user.
func formatEventInfo(event event.Event, reminderSet bool) string {
	var reminder string
	formattedDate := event.RegistrationStartDate.Format(time.Stamp)
	if reminderSet {
		reminder = "set"
	} else {
		reminder = "not set"
	}
	return fmt.Sprintf("<b>%s</b>\nRegistration date: %s\nReminder is %s", event.Title, formattedDate, reminder)
}

func formatReminderTime(reminderTime string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05Z07:00", reminderTime)
}

// formatEditMessage returns a tgbotapi.EditMessageTextConfig with the given text and markup.
// This is used to edit a message instead of sending a new one. This is useful to edit the message
// the user is interacting with with buttons, instead of sending messages over and over.
func formatEditMessage(message *tgbotapi.Message, text string, markup tgbotapi.InlineKeyboardMarkup) tgbotapi.EditMessageTextConfig {
	msg := tgbotapi.NewEditMessageTextAndMarkup(message.Chat.ID, message.MessageID, text, markup)
	msg.ParseMode = tgbotapi.ModeHTML
	return msg
}
