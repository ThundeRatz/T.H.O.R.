package tsclient

import (
	ctx "context"

	"github.com/digitalocean/godo"
)

// ListAllDroplets returns all droplets in the account
func (tsc *Client) ListAllDroplets() ([]godo.Droplet, error) {
	droplets, _, err := tsc.c.Droplets.List(ctx.Background(), &godo.ListOptions{})

	return droplets, err
}

// NewThunderLeagueDroplet creates a new droplet suitable for ThunderLeague testing
func (tsc *Client) NewThunderLeagueDroplet(snapshotID, sshKeyID int) (*godo.Droplet, error) {
	droplets, err := tsc.ListAllDroplets()

	if err != nil {
		return nil, err
	}

	for _, d := range droplets {
		if d.Name == "thunderleague" {
			return nil, ErrThunderLeagueDropletAlreadyPresent
		}
	}

	droplet, _, err := tsc.c.Droplets.Create(ctx.Background(), &godo.DropletCreateRequest{
		Name:   "thunderleague",
		Region: "nyc1",
		Size:   "c-32",
		Image: godo.DropletCreateImage{
			ID: snapshotID,
		},
		SSHKeys: []godo.DropletCreateSSHKey{{ID: sshKeyID}},
	})

	if err != nil {
		return nil, err
	}

	return droplet, nil
}

// RemoveThunderLeagueDroplet removes all droplets with name "thunderleague" from the server
func (tsc *Client) RemoveThunderLeagueDroplet() error {
	droplets, _, err := tsc.c.Droplets.List(ctx.Background(), &godo.ListOptions{})

	if err != nil {
		return err
	}

	for _, d := range droplets {
		if d.Name == "thunderleague" {
			_, err = tsc.c.Droplets.Delete(ctx.Background(), d.ID)

			if err != nil {
				return err
			}

			return nil
		}
	}

	return nil
}
