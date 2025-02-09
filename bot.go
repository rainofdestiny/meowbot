package main

import (
	"context"
	"fmt"
	"strings"
  "strconv"
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

	if isMeowMessage(msg.Text) {
		var key string
		var name string
		var count int64

		if msg.ReplyToMessage != nil && mirrorShield(msg) {
			key = fmt.Sprintf("%d:%d", msg.Chat.ID, msg.ReplyToMessage.From.ID)
			var err error
			count, err = redisClient.Incr(ctx, key).Result()
			if err != nil {
				return
			}
			name = msg.ReplyToMessage.From.FirstName
		} else {
			key = fmt.Sprintf("%d:%d", msg.Chat.ID, msg.From.ID)
			strCount, err := redisClient.Get(ctx, key).Result()
			if err != nil {
				return
			}
			count, err = strconv.ParseInt(strCount, 10, 64)
			if err != nil {
				return
			}
			name = msg.From.FirstName
		}

		response := fmt.Sprintf("%s стал котенком уже %d %s", name, count, getDeclension(int(count)))
		sentMsg, err := bot.Send(tgbotapi.NewMessage(msg.Chat.ID, response))
		if err != nil {
			return
		}

		go func() {
			time.Sleep(15 * time.Minute)
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
  meowMessages := []string{"мяу", "мур", "meow", "purr",}

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

