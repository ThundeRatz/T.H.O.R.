// Functions to sendd various embedded messages

package discord

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/google/go-github/v29/github"
)

func getBaseEmbed() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color:     0xe800ff,
		Timestamp: time.Now().Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "T.H.O.R. | ThundeRatz",
			IconURL: "https://static.thunderatz.org/ThorJoinha.png",
		},
	}
}

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
		&discordgo.MessageEmbedField{
			Name:  "Description",
			Value: body,
		},
	}

	if err := ds.client.SendEmbed(ds.AlertChannel, embed); err != nil {
		ds.logger.Error().Err(err).Msg("Failed to send Issue Alert")
	}
}

func (ds *Service) SendGitHubPRAlert() {}

func (ds *Service) SendYouTubeAlert() {}
