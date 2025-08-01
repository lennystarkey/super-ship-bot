package compatibility

import "github.com/bwmarrin/discordgo"

type Score struct {
}

func Assess(u1, u2 discordgo.User) Score {
	return Score{}
}
