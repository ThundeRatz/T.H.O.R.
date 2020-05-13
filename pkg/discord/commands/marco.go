package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"thunderatz.org/thor/pkg/discord"
)

var marcoCmd = &discord.Command{
	Name:        "marco",
	Category:    "Geral",
	Description: "Vê a latência.",
	Usage:       "marco",

	Enabled:   true,
	GuildOnly: false,
	Aliases:   []string{"ping"},
	PermLevel: "User",

	Run: func(s *discordgo.Session, m *discordgo.Message, c *discord.Context) {
		msg, _ := s.ChannelMessageSend(m.ChannelID, "Marco?")

		t1, _ := msg.Timestamp.Parse()
		t2, _ := m.Timestamp.Parse()

		s.ChannelMessageEdit(m.ChannelID, msg.ID, fmt.Sprintf("Polo! %dms", (t1.UnixNano()-t2.UnixNano())/1000000))
	},
}
