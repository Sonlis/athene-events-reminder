package db

import (
	"context"
	"github.com/jackc/pgx/v5"
	"time"
)

func CreateReminder(conn *pgx.Conn, ctx context.Context, chatId int64, eventId int, reminderTime time.Time) error {
	_, err := conn.Exec(ctx, "insert into reminder(chat_id, event_id, reminder_time) values($1, $2, $3)", chatId, eventId, reminderTime)
	return err
}

func RemoveReminder(conn *pgx.Conn, ctx context.Context, chatId int64, eventId int) error {
	_, err := conn.Exec(ctx, "delete from reminder where chat_id = $1 and event_id = $2", chatId, eventId)
	return err
}

func CheckReminder(conn *pgx.Conn, ctx context.Context, chatId int64, eventId int) (bool, error) {
	var reminderExists bool
	err := conn.QueryRow(ctx, "select exists(select 1 from reminder where chat_id = $1 and event_id = $2)", chatId, eventId).Scan(&reminderExists)
	return reminderExists, err
}
