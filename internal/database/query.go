package database

import (
	"context"
	"time"
)

func (db *DB) CreateReminder(ctx context.Context, chatId int64, eventId int, reminderTime time.Time) error {
	_, err := db.Exec(ctx, "insert into reminder(chat_id, event_id, reminder_time) values($1, $2, $3)", chatId, eventId, reminderTime)
	return err
}

func (db *DB) RemoveReminder(ctx context.Context, chatId int64, eventId int) error {
	_, err := db.Exec(ctx, "delete from reminder where chat_id = $1 and event_id = $2", chatId, eventId)
	return err
}

func (db *DB) CheckReminder(ctx context.Context, chatId int64, eventId int) (bool, error) {
	var reminderExists bool
	err := db.QueryRow(ctx, "select exists(select 1 from reminder where chat_id = $1 and event_id = $2)", chatId, eventId).Scan(&reminderExists)
	return reminderExists, err
}
