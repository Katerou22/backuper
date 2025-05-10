package main

import (
	"backuper/internal/bot"
	"backuper/internal/db"
	"backuper/internal/source"
	"backuper/pkg/env"
	"context"
	"fmt"
	tgbot "github.com/go-telegram/bot"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"
)

func main() {

	env.LoadEnv(".env")

	//Connect to db
	database, err := db.NewSqlite(env.Get("SQLITE_PATH", "./test.db"))
	if err != nil {
		log.Fatalf("failed to start db: %v", err)
	}
	//
	sourceHandler, err := source.NewHandler(database)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	proxyURL, _ := url.Parse(os.Getenv("PROXY"))

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
		Timeout: 10 * time.Second,
	}

	opts := []tgbot.Option{
		tgbot.WithDefaultHandler(bot.Handler),
		tgbot.WithHTTPClient(time.Second, client),
	}

	fmt.Println(os.Getenv("TELEGRAM_TOKEN"))

	b, err := tgbot.New(os.Getenv("TELEGRAM_TOKEN"), opts...)
	if err != nil {
		panic(err)
	}

	b.Start(ctx)

	b.RegisterHandler(tgbot.HandlerTypeMessageText, "/start", tgbot.MatchTypeExact, bot.Start)

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
