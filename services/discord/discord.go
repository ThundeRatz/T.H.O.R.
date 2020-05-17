// Package discord is the discord module for T.H.O.R.
package discord

import (
	"github.com/rs/zerolog"
	"thunderatz.org/thor/pkg/dclient"
)

var ()

// Service is the main discord service struct
type Service struct {
	AlertChannel string
	Token        string

	client *dclient.Client
	logger zerolog.Logger
}

// Init initializes discord service
func (ds *Service) Init(_logger *zerolog.Logger) error {
	ds.logger = _logger.With().Str("serv", "discord").Logger()
	ds.client = &dclient.Client{}

	err := ds.client.Init(ds.Token, &ds.logger)

	if err != nil {
		return err
	}

	ds.client.AddCommand(marcoCmd)
	ds.client.AddCommand(infoCmd)

	return nil
}

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
