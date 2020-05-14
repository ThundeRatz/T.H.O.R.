// Package discord is the discord module for T.H.O.R.
package discord

import (
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"thunderatz.org/thor/pkg/dclient"
)

var (
	client = dclient.New()
	logger zerolog.Logger
)

// Init initializes discord service
func Init(token string, _logger *zerolog.Logger) {
	logger = _logger.With().Str("serv", "discord").Logger()

	err := client.Init(token, &logger)

	if err != nil {
		return
	}
}

// SendMessage sends a message to the specified alert channel
func SendMessage(content string) {
	alertChannel := viper.GetString("discord.alert_channel_id")

	err := client.SendMessage(alertChannel, content)

	if err != nil {
		logger.Error().
			Err(err).
			Str("channel", alertChannel).
			Msg("Couldn't send message to channel")
	}
}

// Start starts discord service
func Start() {
}
