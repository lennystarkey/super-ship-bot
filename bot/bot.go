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
	GuildID = "" // test guild id, if not specified bot will register cmds globally
	Token   = ""
	// AiToken = ""

	commands = []discordgo.ApplicationCommand{
		{
			Name:        "ship",
			Description: "ship 2 users",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user1",
					Description: "first user to ship",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user2",
					Description: "second user to ship",
					Required:    true,
				},
			},
		},
	}
)

func Run() {
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatalf("couldn't create session: %v", err)
	}

	dg.AddHandler(newMessage) // TODO

	if err := dg.Open(); err != nil {
		log.Fatalf("couldn't open connection: %v", err)
	}
	defer dg.Close()

	// TODO app cmd

	fmt.Println("bot running!")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func newMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if strings.Contains(m.Content, "bot") {
		s.ChannelMessageSend(m.ChannelID, "Hello world!")
	}
}
