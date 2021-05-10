package github

import (
	"errors"
	"strconv"
	"strings"

	"github.com/abadojack/whatlanggo"
	"go.thunderatz.org/thor/core/types"
	"go.thunderatz.org/thor/pkg/gclient"
)

// SendDefaultVSSIssueMessage send the default message to newly created issues
func (ghs *Service) SendDefaultVSSIssueMessage(repoName, body string, issueNumber int) {
	client, err := gclient.NewInstallationClient(ghs.AppID, ghs.InstallationID, ghs.PEMFile, &ghs.logger)

	if err != nil {
		ghs.logger.Error().Err(err).Msg("NewInstallationClient")
		return
	}

	lang := whatlanggo.DetectLang(body).Iso6391()

	replyCh := make(types.CoreReplyCh)

	msgCh <- types.CoreMsg{
		Type:  types.GitHubIssueReplyMsg,
		Reply: replyCh,
		From:  serviceId,
		Args: types.GitHubIssueReplyArgs{
			Lang: lang,
		},
	}

	reply := <-replyCh
	replyText := reply.Reply.(*types.GitHubIssueReplyReply).Reply

	client.IssueComment("ThundeRatz", repoName, replyText, issueNumber)
}

func (ghs *Service) ProcessIssueComment(repoName, body string, issueNumber int, isPR bool) {
	splitBody := strings.Fields(body)

	if len(splitBody) < 2 {
		ghs.logger.Warn().Msg("Need at least prefix and command, ignoring...")
		return
	}

	prefix := splitBody[0]
	cmd := splitBody[1]
	args := splitBody[2:]

	if prefix != "!thor" {
		return
	}

	ghs.logger.Info().Str("cmd", cmd).Msg("New issue command")

	// Maybe improve this commands interface to use something like the discord one
	// so it's easier to add new commands
	switch cmd {
	case "test":
		if repoName != "ThunderLeague" && repoName != "test" {
			ghs.logger.Warn().Msg("Received test command outside ThunderLeague repository, ignoring...")
			return
		}

		client, err := gclient.NewInstallationClient(ghs.AppID, ghs.InstallationID, ghs.PEMFile, &ghs.logger)

		if err != nil {
			ghs.logger.Error().Err(err).Msg("NewInstallationClient")
			return
		}

		commentMsg := "Pedido de teste recebido!"

		commentID, err := client.IssueComment("ThundeRatz", repoName, commentMsg, issueNumber)

		if err != nil {
			ghs.logger.Error().Err(err).Msg("Failed to send Issue Comment")
			return
		}

		testsAmount, enemy, err := parseTestArgs(args)

		if err != nil {
			ghs.logger.Error().Err(err).Msg("Failed to parse test arguments")
			commentMsg += "\nFormato inválido, uso: !thor test [quantidade] [nome do time inimigo ou hash do commit]"
			err = client.IssueCommentEdit("ThundeRatz", repoName, commentMsg, commentID)

			if err != nil {
				ghs.logger.Error().Err(err).Msg("Failed to edit Issue Comment")
			}

			return
		}

		commit, err := client.GetPRLastCommit("ThundeRatz", repoName, issueNumber)

		if err != nil {
			ghs.logger.Error().Err(err).Msg("Failed to get commit hash")
			commentMsg += "\nErro: Failed to get commit hash"
			err = client.IssueCommentEdit("ThundeRatz", repoName, commentMsg, commentID)

			if err != nil {
				ghs.logger.Error().Err(err).Msg("Failed to edit Issue Comment")
			}

			return
		}

		ghs.logger.Debug().Str("Commit", commit).Str("Enemy", enemy).Int("Amount", testsAmount).Msg("Sending test request")
		replyCh := make(types.CoreReplyCh)

		msgCh <- types.CoreMsg{
			Type:  types.TLTestMsg,
			Reply: replyCh,
			From:  serviceId,
			Args: types.TLTestArgs{
				Commit: commit,
				Amount: testsAmount,
				Enemy:  enemy,
			},
		}

		reply := <-replyCh

		if !reply.Success {
			ghs.logger.Error().Err(err).Msg("Core failed to start ThunderLeague testing")
			commentMsg += "\nNão consegui iniciar os testes :("
			err = client.IssueCommentEdit("ThundeRatz", repoName, commentMsg, commentID)

			if err != nil {
				ghs.logger.Error().Err(err).Msg("Failed to edit Issue Comment")
			}
		}

		statusCh := reply.Reply.(*types.TLTestReply).StatusCh

		for s := range statusCh {
			commentMsg += "\n> " + s
			err = client.IssueCommentEdit("ThundeRatz", repoName, commentMsg, commentID)

			if err != nil {
				ghs.logger.Error().Err(err).Msg("Failed to edit Issue Comment")
			}
		}
	}
}

func parseTestArgs(args []string) (int, string, error) {
	if len(args) != 2 {
		return -1, "", errors.New("invalid number of arguments")
	}

	testsAmount, err := strconv.Atoi(args[0])

	if err != nil {
		return -1, "", err
	}

	return testsAmount, args[1], nil
}
