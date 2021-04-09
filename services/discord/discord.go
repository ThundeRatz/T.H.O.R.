// Package discord is the discord module for T.H.O.R.
package discord

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"

	"go.thunderatz.org/thor/core/types"
	"go.thunderatz.org/thor/pkg/dclient"
)

var (
	msgCh     types.CoreMsgCh
	serviceId string
)

// Service is the main discord service struct
type Service struct {
	AlertChannel string
	Token        string

	client *dclient.Client
	logger zerolog.Logger
}

// Init initializes discord service
func (ds *Service) Init(_logger *zerolog.Logger, _ch types.CoreMsgCh) error {
	ds.logger = _logger.With().Str("serv", "discord").Logger()
	ds.client = &dclient.Client{}
	serviceId = "discord"
	msgCh = _ch

	err := ds.client.Init(ds.Token, &ds.logger)

	if err != nil {
		return err
	}

	ds.client.AddCommand(pingCmd)
	ds.client.AddCommand(infoCmd)
	ds.client.AddCommand(githubCmd)
	ds.client.AddCommand(configCmd)
	ds.client.AddHelpCmd()

	return nil
}

// Stop stops the service
func (ds *Service) Stop() {
	ds.client.Stop()
}

// SendMessage sends a message to the specified alert channel
func (ds *Service) SendMessage(content string) error {
	err := ds.client.SendMessage(ds.AlertChannel, content)

	if err != nil {
		ds.logger.Error().
			Err(err).
			Str("channel", ds.AlertChannel).
			Msg("Couldn't send message to channel")
	}

	return err
}

func getBaseEmbed() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color:     0xe800ff,
		Timestamp: time.Now().Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "T.H.O.R. | ThundeRatz",
			IconURL: "https://static.thunderatz.org/ThorJoinha.png",
		},
	}
}
