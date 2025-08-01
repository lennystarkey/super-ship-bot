package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	Token = ""
	// AiToken = ""
)

func Run() {
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatalf("couldn't connect to discord: %v", err)
	}

	dg.AddHandler(newMessage)

	dg.Open()
	defer dg.Close()

	fmt.Println("bot running!")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func newMessage(dg *discordgo.Session, msg *discordgo.MessageCreate) {
	if strings.Contains(msg.Content, "bot") {
		dg.ChannelMessageSend(msg.ChannelID, "Hello world!")
	}
}
