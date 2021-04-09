package github

import (
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
