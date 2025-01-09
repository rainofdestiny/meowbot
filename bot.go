package main

import (
	"context"
	"fmt"
	"strings"
  "time"

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
		response := fmt.Sprintf("%s стал котенком уже %d %s!", name, count, getDeclension(int(count)))

		sentMsg, _ := bot.Send(tgbotapi.NewMessage(msg.Chat.ID, response))

		go func() {
			time.Sleep(1 * time.Hour)
			deleteConfig := tgbotapi.DeleteMessageConfig{
				ChatID:    sentMsg.Chat.ID,
				MessageID: sentMsg.MessageID,
			}
			bot.Request(deleteConfig)
		}()

		return
  }
}

func isMeowMessage(text string) bool {
	meowMessages := []string{"мяу", "мур", "meow", "мяуу", "мяу мяу", "purr", "мурр", "мур мур"}
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

func getDeclension(count int) string {
	switch {
	case count%10 == 1 && count%100 != 11:
		return "раз"
	case count%10 >= 2 && count%10 <= 4 && (count%100 < 10 || count%100 >= 20):
		return "раза"
	default:
		return "раз"
	}
}

