// Package github wraps google/go-github for use with T.H.O.R. CLI
// it's meant to access organization data only
package github

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v31/github"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"
)

// Client holds connection data
type Client struct {
	org *github.Organization
	c   *github.Client

	logger *zerolog.Logger
}

// RepoStats holds repository statistics, it's a map of username:<int> for additions
// deletions and commits
type RepoStats struct {
	Adds    map[string]int
	Dels    map[string]int
	Commits map[string]int
}

// New creates a new GitHub client with the provided access token
// and organization name
func New(token, orgName string, logger *zerolog.Logger) (*Client, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)

	org, _, err := client.Organizations.Get(context.Background(), orgName)

	if err != nil {
		return nil, err
	}

	return &Client{
		org:    org,
		c:      client,
		logger: logger,
	}, nil
}

// GetRepositories returns a slice of all repository names for the organization
func (gh *Client) GetRepositories() []string {
	repos, _, _ := gh.c.Repositories.ListByOrg(context.Background(), gh.org.GetLogin(), &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	})

	repoNames := make([]string, len(repos))

	for i, v := range repos {
		repoNames[i] = v.GetName()
	}

	return repoNames
}

// GetMembers returns a slice of all members users for the organization
func (gh *Client) GetMembers() []string {
	members, _, _ := gh.c.Organizations.ListMembers(context.Background(), "ThundeRatz", &github.ListMembersOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	})

	memberNames := make([]string, len(members))

	for i, v := range members {
		memberNames[i] = v.GetLogin()
	}

	return memberNames
}

// GetStats returns contributor statistics for all repositories in the organization
func (gh *Client) GetStats() RepoStats {
	repos := gh.GetRepositories()
	members := gh.GetMembers()

	repoStats := RepoStats{
		Adds:    make(map[string]int),
		Dels:    make(map[string]int),
		Commits: make(map[string]int),
	}

	for _, m := range members {
		repoStats.Adds[m] = 0
		repoStats.Dels[m] = 0
		repoStats.Commits[m] = 0
	}

	for _, r := range repos {
		if gh.logger != nil {
			gh.logger.Info().
				Str("repo", r).
				Msg("Processing Repository")
		}

		stats, resp, err := gh.c.Repositories.ListContributorsStats(context.Background(), gh.org.GetLogin(), r)

		for resp.StatusCode == 202 {
			time.Sleep(500 * time.Millisecond)
			stats, resp, err = gh.c.Repositories.ListContributorsStats(context.Background(), gh.org.GetLogin(), r)
		}

		if err != nil {
			fmt.Println(err)
		}

		for _, stat := range stats {
			username := stat.Author.GetLogin()

			for _, w := range stat.Weeks {
				if _, ok := repoStats.Adds[username]; !ok {
					continue
				}

				repoStats.Adds[username] += w.GetAdditions()
				repoStats.Dels[username] += w.GetDeletions()
				repoStats.Commits[username] += w.GetCommits()
			}
		}
	}

	return repoStats
}
