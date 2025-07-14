package main

import (
	"backuper/internal/bot"
	"backuper/internal/db"
	"backuper/internal/schedule"
	"backuper/internal/source"
	"backuper/pkg/env"
	"bytes"
	"context"
	"fmt"
	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"
)

var cronTime = "0 4 * *"
var channelId string

func main() {

	//env.LoadEnv(".env")

	cronTime = env.Get("CRON", cronTime)

	channelId = env.Get("CHANNEL_ID", "")

	//Connect to db
	database, err := db.NewSqlite(env.Get("SQLITE_PATH", "./test.db"))
	if err != nil {
		log.Fatalf("failed to start db: %v", err)
	}

	sourceHandler, err := source.NewHandler(database)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	scheduler := schedule.NewScheduler()

	proxyURL, _ := url.Parse(env.Get("PROXY", ""))

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
		Timeout: 10 * time.Second,
	}

	botHandler := &bot.NewBotHandler{SourceHandler: sourceHandler}

	opts := []tgbot.Option{
		tgbot.WithDefaultHandler(botHandler.Handler),
		tgbot.WithHTTPClient(time.Second, client),
	}

	b, err := tgbot.New(env.Get("TELEGRAM_TOKEN", ""), opts...)
	if err != nil {
		fmt.Println(err)
	}

	b.RegisterHandler(tgbot.HandlerTypeMessageText, "/start", tgbot.MatchTypeExact, botHandler.Start)
	b.RegisterHandler(tgbot.HandlerTypeMessageText, "/add", tgbot.MatchTypePrefix, botHandler.Add)
	b.RegisterHandler(tgbot.HandlerTypeMessageText, "/list", tgbot.MatchTypeExact, botHandler.List)
	b.RegisterHandler(tgbot.HandlerTypeMessageText, "/backup", tgbot.MatchTypePrefix, botHandler.Backup)
	b.RegisterHandler(tgbot.HandlerTypeMessageText, "/bAll", tgbot.MatchTypePrefix, botHandler.BackupAll)

	scheduler.AddTask(BackupAll(ctx, b, sourceHandler), cronTime)

	scheduler.Run()

	b.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: channelId,
		Text:   "Bot started with cron: " + cronTime,
	})

	b.Start(ctx)

}

func BackupAll(ctx context.Context, b *tgbot.Bot, sourceHandler *source.Handler) func() {

	return func() {
		sources := sourceHandler.List()

		for _, src := range sources {

			backup, err := sourceHandler.Backup(src.ID)

			if err != nil {
				b.SendMessage(ctx, &tgbot.SendMessageParams{
					ChatID: channelId,
					Text:   "Failed to backup: " + err.Error(),
				})
				continue

			}
			fileContent, _ := os.ReadFile(backup)

			params := &tgbot.SendDocumentParams{
				ChatID:   channelId,
				Document: &models.InputFileUpload{Filename: backup[len("backup/"):], Data: bytes.NewReader(fileContent)},
				Caption:  "#" + src.Title,
			}

			b.SendDocument(ctx, params)

		}
	}
}
