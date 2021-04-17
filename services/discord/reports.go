// Functions to sendd various embedded messages

package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/google/go-github/v29/github"
)

func (ds *Service) SendGitHubIssueAlert(issue *github.Issue, repoName string) {
	embed := getBaseEmbed()
	embed.Author = &discordgo.MessageEmbedAuthor{
		Name:    issue.User.GetLogin(),
		IconURL: issue.User.GetAvatarURL(),
		URL:     issue.User.GetHTMLURL(),
	}
	embed.Description = "**" + issue.GetTitle() + "**"
	embed.Title = fmt.Sprintf("%s | Nova Issue #%d", repoName, issue.GetNumber())
	embed.URL = issue.GetHTMLURL()

	body := issue.GetBody()

	if len(body) > 500 {
		body = body[:497] + "..."
	}

	embed.Fields = []*discordgo.MessageEmbedField{
		{
			Name:  "Description",
			Value: body,
		},
	}

	if _, err := ds.client.SendEmbed(ds.AlertChannel, embed); err != nil {
		ds.logger.Error().Err(err).Msg("Failed to send Issue Alert")
	}
}

func (ds *Service) SendGitHubPRAlert() {}

func (ds *Service) SendYouTubeAlert() {}
