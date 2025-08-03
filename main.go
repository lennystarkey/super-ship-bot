package main

import (
	"flag"
	"log"
	"ship/bot"
)

func main() {

	flag.Parse()

	if *bot.DcToken == "" {
		log.Fatal("must provide discord application token")
	}

	// bot.AiToken = aiToken

	bot.Run()
}
