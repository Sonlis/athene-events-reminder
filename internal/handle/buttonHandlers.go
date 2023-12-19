package handle

import (
	"context"
	"fmt"

	"strconv"
	"strings"

	"github.com/Sonlis/athene-events-notifier/internal/db"
	"github.com/Sonlis/athene-events-notifier/internal/event"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	pgx "github.com/jackc/pgx/v5"
)

func handleEventInfo(ctx context.Context, ilmo *event.Ilmo, query tgbotapi.CallbackQuery, db_conn *pgx.Conn) (string, tgbotapi.InlineKeyboardMarkup, error) {
	var text string
	var markup tgbotapi.InlineKeyboardMarkup
	message := query.Message
	eventId := strings.Split(query.Data, " ")[1]
	event, err := ilmo.GetSingleEvent(eventId)
	if err != nil {
		return text, markup, err
	}
	fmt.Println(event.RegistrationStartDate)
	reminderSet, err := db.CheckReminder(db_conn, ctx, message.Chat.ID, event.ID)
	if err != nil {
		return text, markup, err
	}
	text = formatEventInfo(event, reminderSet)
	markup = buildEventInfoMarkup(event, reminderSet)
	return text, markup, nil
}

func handleBackButton(ilmo *event.Ilmo, query tgbotapi.CallbackQuery) (string, tgbotapi.InlineKeyboardMarkup, error) {
	var text string
	var markup tgbotapi.InlineKeyboardMarkup
	events, err := ilmo.GetEvents()
	if err != nil {
		return text, markup, err
	}
	text = "List of incoming events"
	markup = buildEventsMarkup(events)
	if err != nil {
		return text, markup, err
	}
	return text, markup, nil
}

// Set a reminder for an event, then returns the same menu but with the option to remove the reminder.
func handleSetReminder(ctx context.Context, ilmo *event.Ilmo, db_conn *pgx.Conn, query tgbotapi.CallbackQuery, message tgbotapi.Message) (string, tgbotapi.InlineKeyboardMarkup, error) {
	var text string
	var markup tgbotapi.InlineKeyboardMarkup
	eventIdStr := strings.Split(query.Data, " ")[1]
	eventId, err := strconv.Atoi(eventIdStr)
	if err != nil {
		return text, markup, err
	}
	reminderTime, err := formatReminderTime(strings.Split(query.Data, " ")[2])
	if err != nil {
		return text, markup, err
	}
	err = db.CreateReminder(db_conn, ctx, message.Chat.ID, eventId, reminderTime)
	if err != nil {
		return text, markup, err
	}
	event, err := ilmo.GetSingleEvent(eventIdStr)
	if err != nil {
		return text, markup, err
	}
	text = formatEventInfo(event, true)
	markup = buildEventInfoMarkup(event, true)
	return text, markup, err
}

// Remove a reminder for an event, then returns the same menu but with the option to set the reminder.
func handleRemoveReminder(ctx context.Context, ilmo *event.Ilmo, db_conn *pgx.Conn, query tgbotapi.CallbackQuery, message tgbotapi.Message) (string, tgbotapi.InlineKeyboardMarkup, error) {
	var text string
	var markup tgbotapi.InlineKeyboardMarkup
	eventIdStr := strings.Split(query.Data, " ")[1]
	eventId, err := strconv.Atoi(eventIdStr)
	if err != nil {
		return text, markup, err
	}
	err = db.RemoveReminder(db_conn, ctx, message.Chat.ID, eventId)
	if err != nil {
		return text, markup, err
	}

	event, err := ilmo.GetSingleEvent(eventIdStr)
	if err != nil {
		return text, markup, err
	}

	text = formatEventInfo(event, false)
	markup = buildEventInfoMarkup(event, false)
	return text, markup, err
}
