// Package discord is the discord module for T.H.O.R.
package discord

import (
	"github.com/rs/zerolog"
	"thunderatz.org/thor/pkg/dclient"
)

var (
	client = dclient.New()
	logger *zerolog.Logger
)

// Init initializes discord service
func Init(token string, _logger *zerolog.Logger) {
	logger = _logger

	client.Init(token, logger)
}

// Start starts discord service
func Start() {
	client.Start()
}
