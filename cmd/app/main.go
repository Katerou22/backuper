package main

import (
	"backuper/internal/db"
	"backuper/internal/source"
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"
)

func main() {

	//Connect to db
	database, err := db.NewSqlite("./test.db")
	if err != nil {
		log.Fatalf("failed to start db: %v", err)
	}
	//
	sourceHandler, err := source.NewHandler(database)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	proxyURL, _ := url.Parse("http://127.0.0.1:10809")

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
		Timeout: 10 * time.Second,
	}

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
		bot.WithHTTPClient(time.Second, client),
	}

	b, err := bot.New("8056205347:AAH7H0E2K7sUdWHfNZa67Khf_j3-UL0JIGg", opts...)
	if err != nil {
		panic(err)
	}

	b.Start(ctx)

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, startHandler)

	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//err = sourceHandler.CreateFromDsn("postgres://postgres:123123@localhost:5432/postgres")
	//if err != nil {
	//	panic(err)
	//}
	//db.Create()
	//Iterate to all connectins get from config or db
	// In the i
	fmt.Println("Hello World")

	//http.ListenAndServe(":8000", b.WebhookHandler())

}

func startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      "Hello, *" + bot.EscapeMarkdown(update.Message.From.FirstName) + "*",
		ParseMode: models.ParseModeMarkdown,
	})
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Say /start",
	})
}
