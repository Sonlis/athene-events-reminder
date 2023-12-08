package handle

import (
	"context"
	"strings"
	"time"

	"github.com/Sonlis/athene-events-notifier/internal/event"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	pgx "github.com/jackc/pgx/v5"
)

// HandleUpdate handles incoming updates from Telegram.
// This is the entrypoint for all incoming messages.
func HandleUpdate(update tgbotapi.Update, bot *tgbotapi.BotAPI, db_conn *pgx.Conn, ilmo *event.Ilmo) {
	switch {
	// Handle messages
	case update.Message != nil:
		msg, err := handleCommand(update.Message, ilmo)
		if err != nil {
			panic(err)
		}
		bot.Send(msg)
		break

	// Handle button clicks
	case update.CallbackQuery != nil:

		query := update.CallbackQuery
		callbackCfg := tgbotapi.NewCallback(query.ID, "")
		bot.Send(callbackCfg)

		msg, err := handleButton(query, db_conn, ilmo)
		if err != nil {
			panic(err)
		}
		bot.Send(msg)
		break
	}
}

// When we get a command, we react accordingly.
func handleCommand(message *tgbotapi.Message, ilmo *event.Ilmo) (tgbotapi.MessageConfig, error) {
	var err error
	command := message.Text
	chatId := message.Chat.ID

	switch command {

	case "/list":
		msg, err := handleListEvents(chatId, ilmo)
		return msg, err
	}
	return tgbotapi.MessageConfig{}, err
}

// When we get a button click, we react accordingly.
func handleButton(query *tgbotapi.CallbackQuery, db_conn *pgx.Conn, ilmo *event.Ilmo) (tgbotapi.EditMessageTextConfig, error) {
	var text string
	var err error
	var markup tgbotapi.InlineKeyboardMarkup
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	buttonType := strings.Split(query.Data, " ")[0]
	message := query.Message

	switch buttonType {
	case "eventInfo":
		text, markup, err = handleEventInfo(ctx, ilmo, *query, db_conn)
	case "Back":
		text, markup, err = handleBackButton(ilmo, *query)
	case "setReminder":
		text, markup, err = handleSetReminder(ctx, ilmo, db_conn, *query, *message)
	case "removeReminder":
		text, markup, err = handleRemoveReminder(ctx, ilmo, db_conn, *query, *message)
	}
	if err != nil {
		return tgbotapi.EditMessageTextConfig{}, err
	}
	return formatEditMessage(message, text, markup), nil
}
