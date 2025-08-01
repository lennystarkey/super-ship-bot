package main

import (
	"flag"
	"ship/bot"
)

func main() {
	// temporarily parse the token from cmd line, for testing purposes
	dcToken := flag.String("token", "", "discord application token")
	guildID := flag.String("guildID", "", "guild id for testing")

	flag.Parse()

	bot.Token = *dcToken
	bot.GuildID = *guildID

	// bot.AiToken = aiToken

	bot.Run()
}
