package bot

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"backuper/internal/source"
	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type NewBotHandler struct {
	SourceHandler *source.Handler
}

func (h *NewBotHandler) Handler(ctx context.Context, b *tgbot.Bot, update *models.Update) {

	if update.ChannelPost != nil {
		if update.ChannelPost.Text == "/myID" {
			b.SendMessage(ctx, &tgbot.SendMessageParams{ChatID: update.ChannelPost.Chat.ID, Text: strconv.FormatInt(update.ChannelPost.Chat.ID, 10)})

		}

	}

	if update.Message != nil {

		if update.Message.Text == "/myID" {
			b.SendMessage(ctx, &tgbot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: strconv.FormatInt(update.Message.Chat.ID, 10)})

		} else {
			b.SendMessage(ctx, &tgbot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Use list of commands please."})

		}

	}

	// b.SendMessage(ctx, &tgbot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "TEST"})

}

func (h *NewBotHandler) MyID(ctx context.Context, b *tgbot.Bot, update *models.Update) {

	fmt.Println(update.ChannelPost)

	if update.ChannelPost != nil {
		b.SendMessage(ctx, &tgbot.SendMessageParams{ChatID: update.ChannelPost.Chat.ID, Text: strconv.FormatInt(update.ChannelPost.Chat.ID, 10)})

	}

	if update.Message != nil {
		b.SendMessage(ctx, &tgbot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: strconv.FormatInt(update.Message.Chat.ID, 10)})

	}

	// b.SendMessage(ctx, &tgbot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "TEST"})

}

func (h *NewBotHandler) Start(ctx context.Context, b *tgbot.Bot, update *models.Update) {

	b.SendMessage(ctx, &tgbot.SendMessageParams{ChatID: update.Message.Chat.ID, Text: "Hi"})

}
func (h *NewBotHandler) List(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	list := h.SourceHandler.List()

	if len(list) == 0 {
		b.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "No sources found.",
		})
		return
	}

	var sb strings.Builder
	for _, src := range list {
		sb.WriteString(fmt.Sprintf("ID: %d\nTitle: %s\n\nLink: %s\n\n", src.ID, src.Title, src.Link))
	}

	b.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   sb.String(),
	})
}
func (h *NewBotHandler) Add(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	if update.Message.Text == "/add" {
		b.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Please provide a dsn.",
		})

		return
	}

	parts := strings.SplitN(update.Message.Text[len("/add "):], "|", 2)
	if len(parts) != 2 {
		// Handle invalid input
		b.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Invalid format. Use: Title|DSN",
		})
		return
	}

	title := strings.TrimSpace(parts[0])
	dsn := strings.TrimSpace(parts[1])

	err := h.SourceHandler.CreateFromDsn(title, dsn)
	if err != nil {
		b.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Failed to add source: " + err.Error(),
		})
		return
	}

	b.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Source added successfully!",
	})

}

func (h *NewBotHandler) BackupAll(ctx context.Context, b *tgbot.Bot, update *models.Update) {

	sources := h.SourceHandler.List()

	for _, src := range sources {

		backup, err := h.SourceHandler.Backup(src.ID)

		if err != nil {
			b.SendMessage(ctx, &tgbot.SendMessageParams{
				ChatID: "64872989",
				Text:   "Failed to backup: " + err.Error(),
			})
			continue

		}
		fileContent, _ := os.ReadFile(backup)

		params := &tgbot.SendDocumentParams{
			ChatID:   "64872989",
			Document: &models.InputFileUpload{Filename: backup[len("backup/"):], Data: bytes.NewReader(fileContent)},
			Caption:  "#" + src.Title,
		}

		b.SendDocument(ctx, params)

	}

}

func (h *NewBotHandler) Backup(ctx context.Context, b *tgbot.Bot, update *models.Update) {

	if update.Message.Text == "/backup" {
		b.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Please provide a id for backup like: /backup ID",
		})

		return
	}

	msg, err := b.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Please wait...",
	})

	id := update.Message.Text[len("/backup "):]

	idUint64, _ := strconv.ParseUint(id, 10, 64)

	backup, err := h.SourceHandler.Backup(uint(idUint64))
	// defer os.Remove(backup)

	if err != nil {
		b.SendMessage(ctx, &tgbot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Failed to backup: " + err.Error(),
		})
		return
	}

	fileContent, _ := os.ReadFile(backup)

	src := h.SourceHandler.Find(uint(idUint64))

	params := &tgbot.SendDocumentParams{
		ChatID:   update.Message.Chat.ID,
		Document: &models.InputFileUpload{Filename: backup[len("backup/"):], Data: bytes.NewReader(fileContent)},
		Caption:  "#" + src.Title,
	}

	b.SendDocument(ctx, params)

	b.DeleteMessage(ctx, &tgbot.DeleteMessageParams{MessageID: msg.ID, ChatID: update.Message.Chat.ID})

}
