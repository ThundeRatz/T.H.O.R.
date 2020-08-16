// Package gclient wraps google/go-github for use with T.H.O.R.
package gclient

import (
	"context"
	"net/http"
	"time"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v29/github"
	"github.com/rs/zerolog"
)

// Client holds connection data
type Client struct {
	c *github.Client

	logger zerolog.Logger
}

// RepoStats holds repository statistics, it's a map of username:<int> for additions
// deletions and commits
type RepoStats struct {
	Adds    map[string]int
	Dels    map[string]int
	Commits map[string]int
}

// NewInstallationClient creates a new GitHub client for installations
func NewInstallationClient(appID, instID int64, PEMFile string, logger *zerolog.Logger) (*Client, error) {
	tr := http.DefaultTransport

	itr, err := ghinstallation.NewKeyFromFile(
		tr,
		appID,
		instID,
		PEMFile,
	)

	if err != nil {
		return nil, err
	}

	c := github.NewClient(&http.Client{Transport: itr})

	return &Client{
		c:      c,
		logger: logger.With().Str("pkg", "gclient").Logger(),
	}, nil
}

// NewAppClient creates a new GitHub client for apps
func NewAppClient(appID int64, PEMFile string, logger *zerolog.Logger) (*Client, error) {
	tr := http.DefaultTransport

	itr, err := ghinstallation.NewAppsTransportKeyFromFile(
		tr,
		appID,
		PEMFile,
	)

	if err != nil {
		return nil, err
	}

	c := github.NewClient(&http.Client{Transport: itr})

	return &Client{
		c:      c,
		logger: logger.With().Str("pkg", "gclient").Logger(),
	}, nil
}

// GetRepositories returns a slice of all repository names for the organization
func (gc *Client) GetRepositories(orgName string) []string {
	repos, _, _ := gc.c.Repositories.ListByOrg(context.Background(), orgName, &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	})

	repoNames := make([]string, len(repos))

	for i, v := range repos {
		repoNames[i] = v.GetName()
	}

	return repoNames
}

// GetMembers returns a slice of all members users for the organization
func (gc *Client) GetMembers(orgName string) []string {
	members, _, _ := gc.c.Organizations.ListMembers(context.Background(), orgName, &github.ListMembersOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	})

	memberNames := make([]string, len(members))

	for i, v := range members {
		memberNames[i] = v.GetLogin()
	}

	return memberNames
}

// GetStats returns contributor statistics for all repositories in the organization
func (gc *Client) GetStats(orgName string) RepoStats {
	repos := gc.GetRepositories(orgName)
	members := gc.GetMembers(orgName)

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
		gc.logger.Info().
			Str("repo", r).
			Msg("Processing Repository")

		stats, resp, err := gc.c.Repositories.ListContributorsStats(context.Background(), orgName, r)

		for resp.StatusCode == 202 {
			time.Sleep(500 * time.Millisecond)
			stats, resp, err = gc.c.Repositories.ListContributorsStats(context.Background(), orgName, r)
		}

		if err != nil {
			gc.logger.Error().
				Err(err).
				Send()

			continue
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
