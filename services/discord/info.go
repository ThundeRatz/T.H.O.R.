package discord

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"github.com/hako/durafmt"
	"go.thunderatz.org/thor/core/types"
	"go.thunderatz.org/thor/pkg/dclient"
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
			From:  serviceId,
		}

		reply := <-replyCh
		ans := reply.Reply.(*types.InfoReply)

		embed := getBaseEmbed()
		embed.Title = "T.H.O.R."
		embed.Description = "ThundeRatz Holistic Operational Robot"

		var info strings.Builder
		fmt.Fprintf(&info, "Used Mem: %s\n", humanize.IBytes(ans.UsedMemory))
		fmt.Fprintf(&info, "GoRoutines: %d", ans.NGoRoutines)

		embed.Fields = []*discordgo.MessageEmbedField{
			{
				Name:  "Info",
				Value: info.String(),
			},
			{
				Name:  "Uptime",
				Value: durafmt.Parse(ans.Uptime.Round(time.Minute)).String(),
			},
			{
				Name:   "Version",
				Value:  ans.Version,
				Inline: true,
			},
			{
				Name:   "Build Date",
				Value:  ans.BuildDate,
				Inline: true,
			},
		}

		if reply.Success == true {
			if _, err := c.Session.ChannelMessageSendEmbed(c.Message.ChannelID, embed); err != nil {
				c.Logger.Error().Err(err).Msg("Failed to send embed")
				c.Session.ChannelMessageSend(c.Message.ChannelID, err.Error())
			}
		} else {
			c.Session.ChannelMessageSend(c.Message.ChannelID, "Error")
		}

	},
}
