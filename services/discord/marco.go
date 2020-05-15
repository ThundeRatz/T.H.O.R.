package discord

import (
	"fmt"

	"thunderatz.org/thor/pkg/dclient"
)

var marcoCmd = &dclient.Command{
	Name:        "marco",
	Category:    "Geral",
	Description: "Vê a latência.",
	Usage:       "marco",

	Enabled:   true,
	GuildOnly: true,
	Aliases:   []string{"ping"},
	PermLevel: "User",

	Run: func(c *dclient.Context) {
		q := "Marco?"
		a := "Polo!"

		if c.Command == "ping" {
			q = "Ping?"
			a = "Pong!"
		}

		msg, _ := c.Session.ChannelMessageSend(c.Message.ChannelID, q)

		t1, _ := msg.Timestamp.Parse()
		t2, _ := c.Message.Timestamp.Parse()

		c.Session.ChannelMessageEdit(c.Message.ChannelID, msg.ID, fmt.Sprintf("%s %dms", a, (t1.UnixNano()-t2.UnixNano())/1000000))
	},
}
