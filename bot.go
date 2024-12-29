package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-redis/redis/v8"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func StartBot(bot *tgbotapi.BotAPI, redisClient *redis.Client) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			go handleMessage(bot, redisClient, update.Message)
		}
	}
}

func handleMessage(bot *tgbotapi.BotAPI, redisClient *redis.Client, msg *tgbotapi.Message) {
	ctx := context.Background()

	if msg.ReplyToMessage != nil && isMeowMessage(msg.Text) && mirrorShield(msg) {
		key := fmt.Sprintf("%d:%d", msg.Chat.ID, msg.ReplyToMessage.From.ID)
		count, _ := redisClient.Incr(ctx, key).Result()

		name := msg.ReplyToMessage.From.FirstName
		response := fmt.Sprintf("%s был помяукан уже %d раз!", name, count)
		bot.Send(tgbotapi.NewMessage(msg.Chat.ID, response))
		return
	}
}

func isMeowMessage(text string) bool {
	meowMessages := []string{"мяу", "мур", "meow"}
	for _, m := range meowMessages {
		if strings.EqualFold(text, m) {
			return true
		}
	}
	return false
}

func mirrorShield(msg *tgbotapi.Message) bool {
  if msg.ReplyToMessage.From.ID == msg.From.ID || msg.ReplyToMessage.From.IsBot {
    return false 
  }

  return true 
}
