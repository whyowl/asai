package telegram

import (
	"asai/internal/config"
	"context"
	"fmt"
	"log"
	"strconv"
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
			if config.AppConfig.Telegram.WhiteList[strconv.FormatInt(msg.Chat.ID, 10)] {
				reply, err := a.HandleInput(ctx, msg.Chat.ID, msg.Text)
				if err != nil {
					b.SendMessage(ctx, &bot.SendMessageParams{
						ChatID: msg.Chat.ID,
						Text:   "⚠️ Error: " + err.Error(),
					})
					return
				}

				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: msg.Chat.ID,
					Text:   reply,
				})
			} else {
				fmt.Println(msg.Chat.ID, " try chat with bot")
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: msg.Chat.ID,
					Text:   "Forbidden access",
				})
			}
		}),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Telegram bot started...")
	b.Start(ctx)
}
