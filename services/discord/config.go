package discord

import (
	"strings"

	"go.thunderatz.org/thor/core/types"
	"go.thunderatz.org/thor/pkg/dclient"
)

var configCmd = &dclient.Command{
	Name:        "config",
	Category:    "Core",
	Description: "Lê e escreve configurações do bot",
	Usage:       "config [get|set|list] <key> [value]",

	Enabled:   true,
	GuildOnly: false,
	PermLevel: "Admin",

	Run: func(c *dclient.Context) {
		if len(c.Args) < 1 {
			c.Session.ChannelMessageSend(c.Message.ChannelID, "Comandos disponíveis: `get <key>`, `set <key> <value>`")
			return
		}

		replyCh := make(types.CoreReplyCh)

		switch c.Args[0] {
		case "get":
			msgCh <- types.CoreMsg{
				Type:  types.KVConfigGetMsg,
				Reply: replyCh,
				From:  serviceId,
				Args: types.KVConfigGetArgs{
					Key: c.Args[1],
				},
			}

			reply := <-replyCh
			value := reply.Reply.(*types.KVConfigGetReply).Value

			if reply.Success {
				c.Session.ChannelMessageSend(c.Message.ChannelID, value)
			} else {
				c.Session.ChannelMessageSend(c.Message.ChannelID, "Chave não encontrada")
			}

		case "set":
			if len(c.Args) < 3 {
				c.Session.ChannelMessageSend(c.Message.ChannelID, "Uso: `set <key> <value>`")
				return
			}

			msgCh <- types.CoreMsg{
				Type:  types.KVConfigSetMsg,
				Reply: replyCh,
				From:  serviceId,
				Args: types.KVConfigSetArgs{
					Key:   c.Args[1],
					Value: strings.Join(c.Args[2:], " "),
				},
			}

			reply := <-replyCh

			if reply.Success {
				c.Session.ChannelMessageSend(c.Message.ChannelID, "Valor setado com sucesso")
			} else {
				c.Session.ChannelMessageSend(c.Message.ChannelID, "Ocorreu um erro")
			}

		case "list":
			msgCh <- types.CoreMsg{
				Type:  types.KVConfigListMsg,
				Reply: replyCh,
				From:  serviceId,
			}

			reply := <-replyCh
			keys := reply.Reply.(*types.KVConfigListReply).Keys

			if reply.Success {
				c.Session.ChannelMessageSend(c.Message.ChannelID, strings.Join(keys, "\n"))
			} else {
				c.Session.ChannelMessageSend(c.Message.ChannelID, "Erro listando chaves")
			}

		default:
			c.Session.ChannelMessageSend(c.Message.ChannelID, "Não conheço esse comando")
		}
	},
}
