package github

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"

	"thunderatz.org/thor/core/types"
)

// Service represents the GitHub service.
type Service struct {
	AppID          int64
	InstallationID int64
	WebhookSecret  string
	PEMFile        string

	logger zerolog.Logger
}

var (
	msgCh types.CoreMsgCh
)

// Init initializes a GitHub service and adds its endpoint to the mux
func (ghs *Service) Init(_logger *zerolog.Logger, r *mux.Router, _ch types.CoreMsgCh) {
	ghs.logger = _logger.With().Str("serv", "github").Logger()
	msgCh = _ch

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		return
	})

	r.Use(ghs.validatePayload)

	r.HandleFunc("/", ghs.process).Methods("POST")
}

func (ghs *Service) process(w http.ResponseWriter, r *http.Request) {
	requestID := r.Header.Get("X-GitHub-Delivery")

	payload := r.Context().Value(payloadKey).(*WebhookAPIRequest)

	switch githubEvent := r.Header.Get("X-GitHub-Event"); githubEvent {
	case "issues":
		msgCh <- types.CoreMsg{
			Type:  types.GitHubEventMsg,
			Reply: nil,
			Args: types.GitHubEventArgs{
				Issue:      &payload.Issue,
				Repository: &payload.Repo,
			},
		}

		if payload.Action == "opened" {
			issue := payload.Issue
			ghs.logger.Info().Int("issue", issue.GetNumber()).Str("repository", payload.Repo.GetName()).Msg("New Issue created")
		}

	case "issue_comment":
		msgCh <- types.CoreMsg{
			Type:  types.GitHubEventMsg,
			Reply: nil,
			Args: types.GitHubEventArgs{
				Issue:      &payload.Issue,
				Repository: &payload.Repo,
			},
		}

		if payload.Action == "created" {
			comment := strings.TrimSpace(*payload.IssueComment.Body)
			ghs.logger.Info().Str("comment", comment).Msg("Received issue comment")
		}

	default:
		ghs.logger.Error().Str("Request ID", requestID).Str("event type", githubEvent).Msg("No Handler")
		http.Error(w, http.StatusText(404), 404)
	}
}
