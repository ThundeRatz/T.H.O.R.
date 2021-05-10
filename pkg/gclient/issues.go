package gclient

import (
	"context"
	"errors"

	"github.com/google/go-github/v29/github"
)

// IssueComment comments on an issue
func (gc *Client) IssueComment(orgName, repo, text string, issueNumber int) (int64, error) {
	ic, _, err := gc.c.Issues.CreateComment(context.Background(), orgName, repo, issueNumber, &github.IssueComment{
		Body: &text,
	})

	if err != nil {
		return -1, err
	}

	return ic.GetID(), nil
}

// IssueCommentEdit edits an issue comment
func (gc *Client) IssueCommentEdit(orgName, repo, text string, commentID int64) error {
	_, _, err := gc.c.Issues.EditComment(context.Background(), orgName, repo, commentID, &github.IssueComment{
		Body: &text,
	})

	return err
}

// GetPRLastCommit gets the last commit in a PR
func (gc *Client) GetPRLastCommit(orgName, repo string, issueNumber int) (string, error) {
	rc, _, err := gc.c.PullRequests.ListCommits(context.Background(), orgName, repo, issueNumber, &github.ListOptions{
		PerPage: 250, Page: 1,
	})

	if err != nil {
		return "", err
	}

	if len(rc) == 0 {
		return "", errors.New("didn't find any commit")
	}

	return rc[len(rc)-1].GetSHA(), nil
}
