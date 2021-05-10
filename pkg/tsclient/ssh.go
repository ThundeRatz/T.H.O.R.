package tsclient

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/digitalocean/godo"
	"golang.org/x/crypto/ssh"
)

func (tsc *Client) GetSSHClient(d *godo.Droplet) (*ssh.Client, error) {
	updateDroplet, _, err := tsc.c.Droplets.Get(context.Background(), d.ID)

	if err != nil {
		tsc.logger.Error().Err(err).Msg("unable to update droplet info")
		return nil, err
	}

	d = updateDroplet

	addr, err := d.PublicIPv4()
	if err != nil || addr == "" {
		tsc.logger.Error().Err(err).Msg("unable to find public address")
		return nil, err
	}

	tsc.logger.Info().Str("addr", addr).Msg("Found Address")

	key, err := ioutil.ReadFile(filepath.Join(os.Getenv("HOME"), ".ssh", "thor"))
	if err != nil {
		tsc.logger.Error().Err(err).Msg("unable to read private key")
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		tsc.logger.Error().Err(err).Msg("unable to parse private key")
		return nil, err
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", addr), &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})

	if err != nil {
		tsc.logger.Error().Err(err).Msg("unable to connect")
		return nil, err
	}

	return conn, nil
}
