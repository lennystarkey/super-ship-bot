package bot

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"ship/bar"
	"ship/compatibility"
	"slices"
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

	commandHandler = handleShip
)

func Run() {
	s, err := discordgo.New("Bot " + *DcToken)
	if err != nil {
		log.Fatalf("couldn't create session: %v", err)
	}

	s.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages

	s.AddHandler(commandHandler)
	// s.AddHandler(handleComponentInteraction)

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

func handleShip(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// defer response immediately
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	options := i.ApplicationCommandData().Options

	user1 := options[0].UserValue(s)
	user2 := options[1].UserValue(s)
	u1data, u2data, err := getHistory(s, i.ChannelID, *user1, *user2)
	var errs []error
	if err != nil {
		errs = append(errs, fmt.Errorf("error getting history: %w", err))
	}
	result, err := compatibility.Assess(u1data, u2data)
	if err != nil {
		errs = append(errs, fmt.Errorf("error getting compatibility assessment: %w", err))
	}
	if len(errs) > 0 {
		log.Print(errs)
		sendErrorMsg(s, i)
		return
	}

	percentage := result.Compatibility * 2 // progress bar is really low otherwise
	if percentage > 1 {
		log.Print("percentage > 1")
		percentage = 1
	}

	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Ask someone... 🔮",
					Style:    discordgo.PrimaryButton,
					CustomID: "ask",
				},
			},
		},
	}

	embed := &discordgo.MessageEmbed{
		Title: "SuperShip! 🚀",
		Fields: []*discordgo.MessageEmbedField{
			{
				Value: fmt.Sprintf("<@%v> %v\n<@%v> %v\n%v", user1.ID, result.User1.Emoji, user2.ID, result.User2.Emoji, bar.Generate(percentage)),
			},
		},
		Color: rand.Intn(0xFFFFFF),
	}

	_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Embeds:     []*discordgo.MessageEmbed{embed},
		Components: components,
	})
	if err != nil {
		log.Printf("error sending follow-up message: %v", err)
	}
}

// func handleComponentInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
// 	if i.Type != discordgo.InteractionMessageComponent {
// 		return
// 	}
// 	if i.MessageComponentData().CustomID == "ask" {
// 		styles := []string{"Shakespeare", "a biblical prophet", "a toddler", "a medieval knight", "gen alpha brainrot slang"}
// 		// TODO add random gif
// 		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
// 			Type: discordgo.InteractionResponseChannelMessageWithSource,
// 			Data: &discordgo.InteractionResponseData{
// 				Content: fmt.Sprintf("*Transmitting your message to **%v**...* 🔮", styles[rand.Intn(len(styles))]),
// 			},
// 		})
// 		if err != nil {
// 			log.Printf("error responding to button interaction: %v", err)
// 			return
// 		}
// 		msg, err := compatibility.AskSomeone(result)
// 		if err != nil {
// 			log.Printf("error getting story: %v", err)
// 			sendErrorMsg(s, i)
// 			return
// 		}
// 		embed := &discordgo.MessageEmbed{
// 			Title: fmt.Sprintf("According to %v:", result.Style),
// 			Fields: []*discordgo.MessageEmbedField{
// 				{
// 					Value: msg,
// 				},
// 			},
// 			Color: 0xAA00AA,
// 		}
// 		e := ""
// 		_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
// 			Embeds:  &[]*discordgo.MessageEmbed{embed},
// 			Content: &e,
// 		})

// 	}
// }

func sendErrorMsg(s *discordgo.Session, i *discordgo.InteractionCreate) {
	_, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{
			{
				Fields: []*discordgo.MessageEmbedField{
					{
						Value: "😖 oops, we had an error. try again later.",
					},
				},
			},
		},
	})
	if err != nil {
		log.Print(err)
	}
}

// gets the last *MessageCount messages for two users in the server, and returns them as strings separated by \n.
func getHistory(s *discordgo.Session, cid string, user1, user2 discordgo.User) (string, string, error) {
	history, err := s.ChannelMessages(cid, 100, "", "", "")
	if err != nil {
		return "", "", err
	}
	uid1 := user1.ID
	uid2 := user2.ID
	u1data := []string{}
	u2data := []string{}
out:
	for len(u1data) < *MessageCount && len(u2data) < *MessageCount {
		for _, msg := range history {
			if msg.Author.ID == uid1 && len(u1data) < *MessageCount {
				u1data = append(u1data, msg.Content)
				// fmt.Println("appended to u1data")
			}
			if msg.Author.ID == uid2 && len(u2data) < *MessageCount {
				u2data = append(u2data, msg.Content)
				// fmt.Println("appended to u2data")
			}
			if len(u1data) == *MessageCount && len(u2data) == *MessageCount {
				fmt.Println("did it!")
				break out // we did it!
			}
		}
		// we couldn't find enough (*MessageCount) messages by both users in the last 100 messages
		// in this channel, so we'll have to scan another batch.
		// fmt.Println("scanning new batch...")
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
	// put the msgs in chronological order
	slices.Reverse(u1data)
	slices.Reverse(u2data)

	return preprocess(strings.Join(u1data, "\n")), preprocess(strings.Join(u2data, "\n")), nil
}

// TODO
func preprocess(text string) string {
	return text
}
