package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Sonlis/athene-events-reminder/internal/database"
	"github.com/Sonlis/athene-events-reminder/internal/event"
	"github.com/Sonlis/athene-events-reminder/internal/handle"
	"github.com/Sonlis/athene-events-reminder/internal/notifier"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sethvargo/go-envconfig"
)

func main() {
	var err error
	tgToken := os.Getenv("TG_TOKEN")
	bot, err := tgbotapi.NewBotAPI(tgToken)
	if err != nil {
		// Abort if something is wrong
		log.Panic("Error logging to telegram: ", err)
	}

	// Set this to true to log all interactions with telegram servers
	bot.Debug = false

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Create a new cancellable background context. Calling `cancel()` leads to the cancellation of the context
	ctx := context.Background()

	// `updates` is a golang channel which receives telegram updates
	updates := bot.GetUpdatesChan(u)

	var config database.Config
	if err := envconfig.Process(context.Background(), &config); err != nil {
		log.Panic("Error processing database config: ", err)
	}

	db_conn, err := database.Connect(&config)
	if err != nil {
		log.Panic("Error connecting to database: ", err)
	}
	log.Println("Connected to database")
	defer db_conn.Close(context.Background())

	if err != nil {
		log.Panic("Error connecting to database: ", err)
	}

	ilmo, err := event.Init()
	if err != nil {
		log.Panic("Error getting ilmo's API URL: ", err)
	}

	// Pass cancellable context to goroutine
	go receiveUpdates(ctx, updates, bot, db_conn, &ilmo)

	for {
		notifier.Notify(bot, db_conn, &ilmo)

		wait()
	}

}

func wait() {
	now := time.Now()

	secondsUntilNextMinute := 60 - now.Second()

	waitDuration := time.Second * time.Duration(secondsUntilNextMinute)

	time.Sleep(waitDuration)
}

func receiveUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel, bot *tgbotapi.BotAPI, db_conn *database.DB, ilmo *event.Ilmo) {
	// `for {` means the loop is infinite until we manually stop it
	for {
		select {
		// stop looping if ctx is cancelled
		case <-ctx.Done():
			return
		// receive update from channel and then handle it
		case update := <-updates:
			handle.HandleUpdate(update, bot, db_conn, ilmo)
		}
	}
}
