package bot

import (
	"context"
	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Handler(ctx context.Context, b *tgbot.Bot, update *models.Update) {

	b.SendMessage(ctx, &tgbot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Hi"})

}

func Start(ctx context.Context, b *tgbot.Bot, update *models.Update) {

	b.SendMessage(ctx, &tgbot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Hi"})

}
