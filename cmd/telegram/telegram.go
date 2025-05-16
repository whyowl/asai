package telegram

import (
	"context"
	"log"

	"asai/internal/core"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Run(ctx context.Context, a *core.Agent, token string) {

	opts := []bot.Option{
		bot.WithDefaultHandler(func(ctx context.Context, b *bot.Bot, update *models.Update) {
			msg := update.Message
			reply, err := a.HandleInput(msg.Chat.ID, msg.Text)
			if err != nil {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: msg.Chat.ID,
					Text:   "⚠️ Ошибка: " + err.Error(),
				})
				return
			}

			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: msg.Chat.ID,
				Text:   reply,
			})
		}),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("[telegram] бот запущен...")
	b.Start(ctx)
}
