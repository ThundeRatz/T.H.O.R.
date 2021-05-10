package github

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"

	"go.thunderatz.org/thor/core/types"
)

var (
	msgCh types.CoreMsgCh
)

const (
	serviceId = "github"
)

// Service represents the GitHub service.
type Service struct {
	AppID          int64
	InstallationID int64
	WebhookSecret  string
	PEMFile        string

	logger zerolog.Logger
}

// Init initializes a GitHub service and adds its endpoint to the mux
func (ghs *Service) Init(_logger *zerolog.Logger, r *mux.Router, _ch types.CoreMsgCh) {
	ghs.logger = _logger.With().Str("serv", serviceId).Logger()
	msgCh = _ch

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Use(ghs.validatePayload)

	r.HandleFunc("/", ghs.process).Methods("POST")
}

func (ghs *Service) process(w http.ResponseWriter, r *http.Request) {
	requestID := r.Header.Get("X-GitHub-Delivery")

	payload := r.Context().Value(payloadKey).(*WebhookAPIRequest)

	switch githubEvent := r.Header.Get("X-GitHub-Event"); githubEvent {
	case "issues":
		if payload.Action == "opened" {
			if msgCh != nil {
				msgCh <- types.CoreMsg{
					Type:  types.GitHubEventMsg,
					Reply: nil,
					From:  serviceId,
					Args: types.GitHubEventArgs{
						Issue:      &payload.Issue,
						Repository: &payload.Repo,
					},
				}
			}

			issue := payload.Issue
			ghs.logger.Info().Int("issue", issue.GetNumber()).Str("repository", payload.Repo.GetName()).Msg("New Issue created")
		}

	case "issue_comment":
		if payload.Action == "created" {
			issue := payload.Issue
			ic := payload.IssueComment

			isPR := issue.GetPullRequestLinks() != nil

			var msg string
			if isPR {
				msg = "New Pull Request Comment"
			} else {
				msg = "New Issue Comment"
			}

			ghs.logger.Info().Int("pr", issue.GetNumber()).Str("repository", payload.Repo.GetName()).Msg(msg)

			go ghs.ProcessIssueComment(payload.Repo.GetName(), ic.GetBody(), issue.GetNumber(), isPR)
		}

	default:
		ghs.logger.Error().Str("Request ID", requestID).Str("event type", githubEvent).Msg("No Handler")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}
