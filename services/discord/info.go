package discord

import (
	"fmt"
	"net/rpc"

	"github.com/spf13/viper"
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
		conn, err := rpc.Dial("unix", viper.GetString("core.socket"))

		if err != nil {
			c.Session.ChannelMessageSend(c.Message.ChannelID, "Failed to connect to T.H.O.R. core")
			return
		}
		defer conn.Close()

		ans := types.InfoReply{}
		conn.Call("ThorCore.Info", types.InfoArgs{}, &ans)

		var msg = ""
		if ans.Success == true {
			msg = fmt.Sprintf("gor: %d", ans.NGoRoutines)
		} else {
			msg = "Error"
		}

		c.Session.ChannelMessageSend(c.Message.ChannelID, msg)
	},
}
