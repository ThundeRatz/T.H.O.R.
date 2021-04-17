package discord

import (
	"go.thunderatz.org/thor/pkg/dclient"
)

var modCmd = &dclient.Command{
	Name:        "ping",
	Category:    "Geral",
	Description: "Vê a latência.",
	Usage:       "ping",

	Enabled:   true,
	GuildOnly: false,
	Aliases:   []string{"ping"},
	PermLevel: "Gestão",

	Run: func(c *dclient.Context) {

	},
}
