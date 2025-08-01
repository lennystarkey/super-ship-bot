package compatibility

import (
	"github.com/bwmarrin/discordgo"
)

type Score struct {
	User1      Analysis
	User2      Analysis
	Similarity float64
}

type Analysis struct {
	Formality float64
	Sentiment float64
	Favorites []string
}

func Assess(u1, u2 discordgo.User) Score {
	return Score{}
}
