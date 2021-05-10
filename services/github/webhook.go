package github

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/go-github/v29/github"
)

// WebhookAPIRequest represents the payload received by GitHub
type WebhookAPIRequest struct {
	// Common
	Action       string              `json:"action"`
	User         github.User         `json:"sender"`
	Repo         github.Repository   `json:"repository"`
	Org          github.Organization `json:"organization"`
	Installation struct {
		ID int64 `json:"id"`
	} `json:"installation"`

	// Specific
	Issue        github.Issue        `json:"issue"`
	IssueComment github.IssueComment `json:"comment"`
	PullRequest  github.PullRequest  `json:"pull_request"`
}

type contextKey int

const (
	payloadKey contextKey = iota
)

func (ghs *Service) validatePayload(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p, err := github.ValidatePayload(r, []byte(ghs.WebhookSecret))

		if err != nil {
			ghs.logger.Error().Err(err).Send()
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		pd, err := getWebhookAPIRequest(p)
		if err != nil {
			ghs.logger.Error().Err(err).Msg("Error unmarshalling json payload")
			return
		}

		ctx := context.WithValue(r.Context(), payloadKey, pd)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// helper function to convert json payload to struct
func getWebhookAPIRequest(body []byte) (*WebhookAPIRequest, error) {
	var wh = new(WebhookAPIRequest)
	err := json.Unmarshal(body, &wh)
	if err != nil {
		return nil, err
	}
	return wh, nil
}
