package discord

import (
	"fmt"

	"thunderatz.org/thor/core/types"
	"thunderatz.org/thor/pkg/dclient"
)

var infoCmd = &dclient.Command{
	Name:        "info",
	Category:    "Core",
	Description: "Vê informações de runtime do bot",
	Usage:       "info",

	Enabled:   true,
	GuildOnly: false,
	Aliases:   []string{"status"},
	PermLevel: "User",

	Run: func(c *dclient.Context) {
		replyCh := make(types.CoreReplyCh)
		msgCh <- types.CoreMsg{
			Type:  types.InfoMsg,
			Reply: replyCh,
		}

		reply := <-replyCh
		ans := reply.Reply.(*types.InfoReply)

		var msg = ""
		if reply.Success == true {
			msg = fmt.Sprintf("Got: %d", ans.NGoRoutines)
		} else {
			msg = "Error"
		}

		c.Session.ChannelMessageSend(c.Message.ChannelID, msg)
	},
}
