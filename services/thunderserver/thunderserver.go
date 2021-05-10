package thunderserver

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/digitalocean/godo"
	"github.com/rs/zerolog"
	"go.thunderatz.org/thor/core/types"
	"go.thunderatz.org/thor/pkg/tsclient"
	"golang.org/x/crypto/ssh"
)

var (
	msgCh types.CoreMsgCh
)

const (
	serviceId = "thunderserver"
)

// Service represents the Thunderserver service.
type Service struct {
	Token string

	client         *tsclient.Client
	createdDroplet *godo.Droplet
	logger         zerolog.Logger
}

func (tss *Service) Init(_logger *zerolog.Logger, _ch types.CoreMsgCh) {
	tss.logger = _logger.With().Str("serv", serviceId).Logger()
	tss.client = tsclient.NewClient(tss.Token, &tss.logger)
	msgCh = _ch
}

func (tss *Service) StartThunderLeagueDroplet() error {
	sshKeyID, err := getIntConfigValueFromCore("tleague/sshkeyid")

	if err != nil {
		tss.logger.Error().Err(err).Msg("Couldn't find tleague/sshkeyid")
		return err
	}

	snapshotID, err := getIntConfigValueFromCore("tleague/snapshotid")

	if err != nil {
		tss.logger.Error().Err(err).Msg("Couldn't find tleague/snapshotid")
		return err
	}

	tss.createdDroplet, err = tss.client.NewThunderLeagueDroplet(snapshotID, sshKeyID)

	if err != nil {
		fmt.Printf("start new thunderleague droplet: %v\n", err)
		tss.logger.Error().Err(err).Msg("Couldn't start new thunderleague droplet")
		return err
	}

	return nil
}

func (tss *Service) RunThunderLeagueTest(team1, team2 string, amount int, statusCh chan string) error {
	if tss.createdDroplet == nil {
		statusCh <- "Error: can't find created droplet"
		tss.logger.Error().Msg("Tried to Run server without a created droplet")
		return errors.New("can't find created droplet")
	}

	// Try to connect a few times
	retries := 10
	var client *ssh.Client = nil
	var err error

	for retries > 0 && client == nil {
		tss.logger.Info().Int("try", 11-retries).Msg("Trying to connect")
		client, err = tss.client.GetSSHClient(tss.createdDroplet)

		if err != nil {
			retries--
			time.Sleep(10 * time.Second)
		}
	}

	if client == nil {
		statusCh <- "Error: couldn't connect to droplet"
		tss.logger.Error().Msg("can't connect to droplet")
		return errors.New("can't connect to droplet")
	}

	statusCh <- "Status: Connected to droplet"
	s, err := client.NewSession()

	if err != nil {
		statusCh <- "Error: couldn't start session in droplet"
		tss.logger.Error().Err(err).Msg("can't start session in droplet")
		return err
	}

	statusCh <- "Status: Started session"
	stdOutPipe, err := s.StdoutPipe()

	if err != nil {
		tss.logger.Error().Msg("stdoutpipe")
		return err
	}

	go func() {
		err = s.Start(fmt.Sprintf("echo test && ./run.sh %s %s %d", team1, team2, amount))

		if err != nil {
			tss.logger.Error().Err(err).Msg("run")
			return
		}
	}()

	scanner := bufio.NewScanner(stdOutPipe)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		statusCh <- scanner.Text()
	}

	return nil
}

func (tss *Service) StopThunderLeagueDroplet() error {
	return tss.client.RemoveThunderLeagueDroplet()
}

func getIntConfigValueFromCore(key string) (int, error) {
	replyCh := make(types.CoreReplyCh)

	msgCh <- types.CoreMsg{
		Type:  types.KVConfigGetMsg,
		Reply: replyCh,
		From:  serviceId,
		Args: types.KVConfigGetArgs{
			Key: key,
		},
	}

	reply := <-replyCh

	if !reply.Success {
		return -1, errors.New("couldn't find key")
	}

	value := reply.Reply.(*types.KVConfigGetReply).Value
	return strconv.Atoi(value)
}
