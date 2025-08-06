package bot

import (
	"bytes"
	"flag"
	"fmt"
	"image/color"
	"log"
	"math"
	"os"
	"os/signal"
	"ship/compatibility"
	"slices"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/fogleman/gg"
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

func handleShip(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// defer response immediately
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	options := i.ApplicationCommandData().Options
	user1 := options[0].UserValue(s)
	user2 := options[1].UserValue(s)
	u1data, u2data, err := getHistory(s, i.ChannelID, *user1, *user2)
	if err != nil {
		log.Fatalf("error getting history: %v", err)
	}
	result, err := compatibility.Assess(u1data, u2data)
	if err != nil {
		log.Fatalf("error getting compatibility assessment: %v", err)
	}

	percentage := result.Compatibility * 2 // progress bar is really low otherwise
	if percentage > 1 {
		fmt.Println("percentage > 1")
		percentage = 1
	}
	imgData, err := generateProgressBarImage(percentage)
	if err != nil {
		log.Printf("Error generating image: %v", err)
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "SuperShip! 🚀",
		Description: "Checking compatibility...",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "SHIP",
				Value: fmt.Sprintf("<@%v> 💘 <@%v>\n**SUCCESS**", user1.ID, user2.ID),
			},
			{
				Value: result.Story,
			},
		},
		Image: &discordgo.MessageEmbedImage{
			URL: "attachment://progress.png",
		},
		Color: 0xff0000,
	}

	_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{embed},
		Files: []*discordgo.File{
			{
				Name:   "progress.png",
				Reader: bytes.NewReader(imgData),
			},
		},
	})
	if err != nil {
		fmt.Println(err)
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

func generateProgressBarImage(percentage float64) ([]byte, error) {
	// Create a 200x50 pixel image
	dc := gg.NewContext(200, 50)
	dc.SetColor(color.RGBA{100, 100, 100, 255}) // Gray background
	dc.Clear()

	// Draw filled portion
	filledWidth := int(math.Round((percentage) * 200))

	c := color.RGBA{}
	switch {
	case percentage > 0.67:
		c = color.RGBA{0, 255, 0, 255}
	case percentage > 0.33:
		c = color.RGBA{255, 255, 0, 255}
	default:
		c = color.RGBA{255, 0, 0, 255}
	}
	dc.SetColor(c)
	dc.DrawRectangle(0, 0, float64(filledWidth), 50)
	dc.Fill()

	// Encode to PNG
	var buf bytes.Buffer
	err := dc.EncodePNG(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
