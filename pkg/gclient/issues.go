package gclient

import (
	"context"

	"github.com/google/go-github/v29/github"
)

// IssueComment comments on an issue
func (gc *Client) IssueComment(orgName, repo, text string, issueNumber int) {
	gc.c.Issues.CreateComment(context.Background(), orgName, repo, issueNumber, &github.IssueComment{
		Body: &text,
	})
}
