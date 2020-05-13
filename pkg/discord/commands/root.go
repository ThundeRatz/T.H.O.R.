package commands

import "thunderatz.org/thor/pkg/discord"

func AddAllCommands(c *discord.Client) {
	c.AddCommand(marcoCmd)
}
