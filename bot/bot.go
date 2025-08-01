package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

var (
	GuildID = "" // test guild id, if not specified bot will register cmds globally
	Token   = ""
	// AiToken = ""

	command = discordgo.ApplicationCommand{
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
	}

	commandHandler = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		options := i.ApplicationCommandData().Options
		user1 := options[0].UserValue(s)
		user2 := options[1].UserValue(s)
		msg := fmt.Sprintf("<@%v> 💘 <@%v>", user1.ID, user2.ID)

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msg,
			},
		})
	}
)

func Run() {
	s, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatalf("couldn't create session: %v", err)
	}

	s.AddHandler(commandHandler)

	if err := s.Open(); err != nil {
		log.Fatalf("couldn't open connection: %v", err)
	}
	defer s.Close()

	_, err = s.ApplicationCommandCreate(s.State.User.ID, GuildID, &command)
	if err != nil {
		log.Fatalf("couldn't create '%v' command: %v", command.Name, err)
	}

	fmt.Println("bot running!")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
