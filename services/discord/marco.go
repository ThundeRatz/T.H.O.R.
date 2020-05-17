package discord

import (
	"fmt"
	"net/rpc"
	"time"

	"github.com/spf13/viper"
	"thunderatz.org/thor/core/types"
	"thunderatz.org/thor/pkg/dclient"
)

var marcoCmd = &dclient.Command{
	Name:        "marco",
	Category:    "Geral",
	Description: "Vê a latência.",
	Usage:       "marco",

	Enabled:   true,
	GuildOnly: false,
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
		dt1 := (t1.UnixNano() - t2.UnixNano()) / 1000000

		conn, err := rpc.Dial("unix", viper.GetString("core.socket"))

		if err != nil {
			// logger.Fatal().Err(err).Msg("Failed to send")
		}
		defer conn.Close()

		ans := types.PingReply{}
		conn.Call("ThorCore.Ping", types.PingArgs{}, &ans)

		t3 := time.Now()

		// if ans.Success {
		// 	logger.Info().Msg("Done")
		// } else {
		// 	logger.Info().Msg("Error")
		// }

		dt2 := (t3.UnixNano() - t2.UnixNano()) / 1000000

		c.Session.ChannelMessageEdit(c.Message.ChannelID, msg.ID, fmt.Sprintf("%s\nServer: %dms\nTotal: %dms", a, dt1, dt2))
	},
}
