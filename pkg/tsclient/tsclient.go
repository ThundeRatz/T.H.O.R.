// package tsclient (ThunderServer CLient) interacts with Digital Ocean API
// to provide access to the ThundeRatz's servers
package tsclient

import (
	"errors"

	"github.com/digitalocean/godo"
	"github.com/rs/zerolog"
)

// Client is the main server client
type Client struct {
	c *godo.Client

	logger zerolog.Logger
}


var (
	ErrThunderLeagueDropletAlreadyPresent = errors.New("ThunderLeague droplet already present")
)

func NewClient(token string, logger *zerolog.Logger) *Client {
	return &Client{
		c: godo.NewFromToken(token),
		logger: logger.With().Str("pkg", "tsclient").Logger(),
	}
}
