package discord

import (
	"fmt"
	"time"

	"thunderatz.org/thor/core/types"
	"thunderatz.org/thor/pkg/dclient"
)

var pingCmd = &dclient.Command{
	Name:        "ping",
	Category:    "Geral",
	Description: "Vê a latência.",
	Usage:       "ping",

	Enabled:   true,
	GuildOnly: false,
	Aliases:   []string{"ping"},
	PermLevel: "User",

	Run: func(c *dclient.Context) {
		msg, _ := c.Session.ChannelMessageSend(c.Message.ChannelID, "Ping?")

		t1, _ := msg.Timestamp.Parse()
		t2, _ := c.Message.Timestamp.Parse()
		dt1 := (t1.UnixNano() - t2.UnixNano()) / 1000000

		replyCh := make(types.CoreReplyCh)
		msgCh <- types.CoreMsg{
			Type:  types.PingMsg,
			Reply: replyCh,
		}

		reply := <-replyCh

		t3 := time.Now()

		if reply.Success {
			c.Logger.Info().Msg("Done")
		} else {
			c.Logger.Info().Msg("Error")
		}

		dt2 := (t3.UnixNano() - t2.UnixNano()) / 1000000

		c.Session.ChannelMessageEdit(c.Message.ChannelID, msg.ID, fmt.Sprintf("Pong!\nServer: %dms\nTotal: %dms", dt1, dt2))
	},
}
