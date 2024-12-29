package main

import (
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	redisAddr := os.Getenv("REDIS_ADDR")

	if telegramToken == "" || redisAddr == "" {
		log.Fatal("Missing TELEGRAM_TOKEN or REDIS_ADDR environment variables")
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	botAPI, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Fatalf("Failed to create bot API: %v", err)
	}

	log.Printf("Authorized on account %s", botAPI.Self.UserName)

	StartBot(botAPI, redisClient)
}
