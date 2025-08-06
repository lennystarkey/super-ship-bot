package main

import (
	"flag"
	"log"
	"ship/bot"
	"ship/compatibility"
)

func main() {
	flag.Parse()
	if len(*bot.DcToken) == 0 {
		log.Fatal("must provide discord application token")
	}
	if len(*compatibility.HfToken) == 0 {
		log.Fatal("must provide hf inference api token")
	}

	// bot.AiToken = aiToken

	bot.Run()
}
