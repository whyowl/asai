package telegram

import (
	"context"
	"log"
	"time"

	"asai/internal/core"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Run(ctx context.Context, a *core.Agent, token string) {

	opts := []bot.Option{
		bot.WithDefaultHandler(func(ctx context.Context, b *bot.Bot, update *models.Update) {

			chn := make(chan bool)

			go func(chn chan bool, id int64) {
				for {
					select {
					case <-chn:
						return
					default:
						_, err := b.SendChatAction(ctx, &bot.SendChatActionParams{ChatID: id, Action: models.ChatActionTyping})
						if err != nil {
							log.Println(err)
						}
						time.Sleep(5 * time.Second)
					}
				}
			}(chn, update.Message.Chat.ID)
			defer func(chn chan bool) {
				chn <- true
			}(chn)

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
