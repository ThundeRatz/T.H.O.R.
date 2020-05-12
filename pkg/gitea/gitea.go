package gitea

import (
	"os"

	"code.gitea.io/sdk/gitea"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

// Client holds connection data
type Client struct {
	g     *gitea.Client
	org   *gitea.Organization
	user  *gitea.User
	token string
	url   string
}

// New creates a new Client for given organization
func New(url, token, orgName string) (*Client, error) {
	var err error
	client := &Client{
		url:   url,
		token: token,
	}

	client.g = gitea.NewClient(url, token)

	if client.org, err = client.g.GetOrg(orgName); err != nil {
		return nil, err
	}

	if client.user, err = client.g.GetMyUserInfo(); err != nil {
		return nil, err
	}

	return client, nil
}

// ListRepos return a list of the available repositories in the Client
// organization
func (c *Client) ListRepos() []string {
	repos, _ := c.g.ListOrgRepos(c.org.UserName)

	repoNames := make([]string, len(repos))
	for i, v := range repos {
		repoNames[i] = v.Name
	}

	return repoNames
}

// CloneRepo clones a repository in the current path
func (c *Client) CloneRepo(repoName string, verbose bool) error {
	repo, err := c.g.GetRepo(c.org.UserName, repoName)

	if err != nil {
		return err
	}

	opts := &git.CloneOptions{
		URL: c.url + repo.FullName,
		Auth: &http.BasicAuth{
			Username: c.user.UserName,
			Password: c.token,
		},
	}

	if verbose {
		opts.Progress = os.Stdout
	}

	_, err = git.PlainClone(repoName, false, opts)

	return err
}
