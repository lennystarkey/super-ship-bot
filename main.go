package main

import (
	"log"
	"os"
	"ship/bot"
)

func main() {
	dcToken, ok := os.LookupEnv("SUPERSHIP_TOKEN")
	if !ok {
		log.Fatal("must set environment variable SUPERSHIP_TOKEN to the discord bot token")
	}
	// aiToken, ok := os.LookupEnv("AI_TOKEN")
	// if !ok {
	// 	log.Fatal("must set environment variable AI_TOKEN to the ai api token")
	// }

	bot.Token = dcToken
	// bot.AiToken = aiToken

	bot.Run()
}
