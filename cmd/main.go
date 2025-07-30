package main

import (
	"github.com/brandao07/vzcount/internal/bot"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	botToken := os.Getenv("DISCORD_BOT_TOKEN")

	err = bot.Run(botToken)
	if err != nil {
		log.Fatal(err)
	}
}
