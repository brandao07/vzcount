package main

import (
	"fmt"
	"github.com/brandao07/vzcount/internal/hypixel"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	hypixelAPIKey := os.Getenv("HYPIXEL_API_KEY")

	hypeClient := hypixel.API{
		URL: "https://api.hypixel.net/v2",
		Key: hypixelAPIKey,
	}

	count, err := hypeClient.VampireZCount()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Hypixel VampireZ Count: ", count)
}
