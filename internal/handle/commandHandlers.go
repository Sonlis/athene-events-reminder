package handle

import (
	"github.com/Sonlis/athene-events-reminder/internal/event"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func handleListEvents(chatId int64, ilmo *event.Ilmo) (tgbotapi.MessageConfig, error) {
	events, err := ilmo.GetEvents()
	if err != nil {
		return tgbotapi.MessageConfig{}, err
	}
	eventsMarkup := buildEventsMarkup(events)
	msg := tgbotapi.NewMessage(chatId, "List of incoming events")
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = eventsMarkup
	return msg, nil
}
