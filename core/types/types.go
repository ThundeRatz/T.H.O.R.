// Package types exposes types for the services to communicate with the core
package types

import (
	"github.com/google/go-github/v29/github"
	"go.thunderatz.org/thor/pkg/gclient"
)

// Msg and Reply Channels
type (
	CoreMsgCh   chan CoreMsg
	CoreReplyCh chan CoreReply
)

// CoreMsg represents a Msg received by the core via a CoreMsgCh
type CoreMsg struct {
	Type  CoreMsgType
	Reply CoreReplyCh

	// Args must be one of the *Args related to the CoreMsg Type
	Args interface{}
}

// CoreReply represents a Msg sent by the core via a CoreReplyCh
type CoreReply struct {
	Success bool

	// Reply must be one of the *Reply related to the CoreMsg Type
	Reply interface{}
}

// CoreMsgType represents a message type that the core can process
type CoreMsgType int

// CoreMsgType constants
const (
	PingMsg CoreMsgType = iota
	InfoMsg
	GitHubStatsMsg
	GitHubEventMsg
)

// GitHubEventArgs represents data sent by the GitHub Webhook service
type GitHubEventArgs struct {
	Issue      *github.Issue
	Repository *github.Repository
}

// InfoReply is the reply for the Info function
type InfoReply struct {
	NGoRoutines int
}

// GitHubStatsReply is the reply for the GitHub Service GetStats function function
type GitHubStatsReply struct {
	RepoStats gclient.RepoStats
}
