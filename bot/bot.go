package bot

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	DcToken      = flag.String("token", "", "discord application token")
	GuildID      = flag.String("guildID", "", "(optional) guild id for testing")
	MessageCount = flag.Int("messageCount", 10, "minimum messages per user")

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
		u1data, u2data, err := getHistory(s, i.ChannelID, *user1, *user2)
		if err != nil {
			log.Fatalf("error getting history: %v", err)
		}
		msg += "\n" + u1data + "\n" + u2data

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: msg,
			},
		})
	}
)

func Run() {
	s, err := discordgo.New("Bot " + *DcToken)
	if err != nil {
		log.Fatalf("couldn't create session: %v", err)
	}

	s.AddHandler(commandHandler)

	if err := s.Open(); err != nil {
		log.Fatalf("couldn't open connection: %v", err)
	}
	defer s.Close()

	_, err = s.ApplicationCommandCreate(s.State.User.ID, *GuildID, &command)
	if err != nil {
		log.Fatalf("couldn't create '%v' command: %v", command.Name, err)
	}

	fmt.Println("bot running!")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

// gets the last *MessageCount messages for two users in the server, and returns them as a string separated by \n.
func getHistory(s *discordgo.Session, cid string, user1, user2 discordgo.User) (string, string, error) {
	history, err := s.ChannelMessages(cid, 100, "", "", "")
	if err != nil {
		return "", "", err
	}
	uid1 := user1.ID
	uid2 := user2.ID
	u1data := []string{}
	u2data := []string{}
	for len(u1data) < *MessageCount && len(u2data) < *MessageCount {
		for _, msg := range history {
			if msg.Author.ID == uid1 && len(u1data) < *MessageCount {
				u1data = append(u1data, msg.Content)
			}
			if msg.Author.ID == uid2 && len(u2data) < *MessageCount {
				u2data = append(u2data, msg.Content)
			}
			if len(u1data) == *MessageCount && len(u2data) == *MessageCount {
				break // we did it!
			}
		}
		// we couldn't find enough (*MessageCount) messages by both users in the last 100 messages
		// in this channel, so we'll have to scan another batch.
		first := history[0].ID
		history, err = s.ChannelMessages(cid, 100, first, "", "")
		if err != nil {
			if len(u1data) < *MessageCount || len(u2data) < *MessageCount {
				fmt.Println("one or both of the users hasn't sent enough messages in this channel, proceeding anyway")
				break
			} else {
				return "", "", err
			}
		}
	}

	return strings.Join(u1data, "\n"), strings.Join(u2data, "\n"), nil
}
