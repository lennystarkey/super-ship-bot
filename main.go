package main

import (
	"log"
	"os"
	"ship/bot"
)

func main() {
	// temporarily parse the token from cmd line, for testing purposes
	dcToken := ""
	if len(os.Args) > 1 {
		dcToken = os.Args[1]
	} else {
		log.Fatal("token not included as argument")
	}
	// dcToken, ok := os.LookupEnv("SUPERSHIP_TOKEN")
	// if !ok {
	// 	log.Fatal("must set environment variable SUPERSHIP_TOKEN to the discord bot token")
	// }

	// aiToken, ok := os.LookupEnv("AI_TOKEN")
	// if !ok {
	// 	log.Fatal("must set environment variable AI_TOKEN to the ai api token")
	// }

	bot.Token = dcToken
	// bot.AiToken = aiToken

	bot.Run()
}
