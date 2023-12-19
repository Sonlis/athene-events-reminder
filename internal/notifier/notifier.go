package notifier

import (
	"fmt"
	"log"
	"time"

	"context"

	"github.com/Sonlis/athene-events-reminder/internal/database"
	"github.com/Sonlis/athene-events-reminder/internal/event"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Reminder struct {
	ChatId       int64
	EventId      int
	ReminderTime time.Time
}

func Notify(bot *tgbotapi.BotAPI, db *database.DB, i *event.Ilmo) {
	ctx := context.Background()

	reminders, err := getReminders(ctx, db, i)
	if err != nil {
		log.Printf("Error getting reminders: %v", err)
	}

	eventsDescription, err := buildEventsDescription(reminders, i)
	if err != nil {
		log.Printf("Error building events description: %v", err)
	}

	for _, reminder := range reminders {
		msg := eventsDescription[reminder.EventId]
		msgConfig := tgbotapi.NewMessage(reminder.ChatId, msg)

		log.Printf("Sending reminder to %v", reminder.ChatId)

		_, err := bot.Send(msgConfig)
		if err != nil {
			log.Printf("Error sending message to chat %d: %v", reminder.ChatId, err)
		}
	}
}

func buildEventsDescription(reminders []Reminder, i *event.Ilmo) (map[int]string, error) {
	eventsDescription := make(map[int]string)
	for _, reminder := range reminders {
		if _, ok := eventsDescription[reminder.EventId]; !ok {
			eventIdString := fmt.Sprint(reminder.EventId)
			description, err := i.BuildEventDescription(eventIdString)
			if err != nil {
				log.Printf("Error building event description for event %d: %v", reminder.EventId, err)
			}
			eventsDescription[reminder.EventId] = description
		}
	}
	return eventsDescription, nil
}

func getReminders(ctx context.Context, db *database.DB, i *event.Ilmo) ([]Reminder, error) {
	// Checks for reminders that are set to be sent in less than a minute.
	rows, err := db.Query(ctx, "SELECT chat_id, event_id, reminder_time FROM reminder WHERE reminder_time - localtimestamp BETWEEN interval '0 minute' AND interval '1 minute'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reminders := make([]Reminder, 0)
	for rows.Next() {
		var reminder Reminder
		if err := rows.Scan(&reminder.ChatId, &reminder.EventId, &reminder.ReminderTime); err != nil {
			return nil, err
		}
		reminders = append(reminders, reminder)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return reminders, nil
}
