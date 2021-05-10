// Package types exposes types for the services to communicate with the core
package types

import (
	"time"

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
	From  string

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
	GitHubIssueReplyMsg
	KVConfigGetMsg
	KVConfigSetMsg
	KVConfigListMsg
	TLTestMsg
)

var CoreMsgTypeDesc = []string{
	"PingMsg",
	"InfoMsg",
	"GitHubStatsMsg",
	"GitHubEventMsg",
	"GitHubIssueReplyMsg",
	"KVConfigGetMsg",
	"KVConfigSetMsg",
	"KVConfigListMsg",
	"TLTestMsg",
}

// GitHubEventArgs represents data sent by the GitHub Webhook service
type GitHubEventArgs struct {
	Issue      *github.Issue
	Repository *github.Repository
}

// GitHubIssueReplyArgs represents
type GitHubIssueReplyArgs struct {
	Lang string
}

// KVConfigGetArgs represents the key to be retrived from the core KV database
type KVConfigGetArgs struct {
	Key string
}

// KVConfigSetArgs represents the key and value to be set in the core KV database
type KVConfigSetArgs struct {
	Key   string
	Value string
}

// TLTestArgs represents what to be tested by the thunderleague server
type TLTestArgs struct {
	Commit string  // Commit hash of the first team
	Enemy  string  // Name of the second team or hash for a different commit
	Amount int
}

// InfoReply is the reply for the Info function
type InfoReply struct {
	NGoRoutines int
	UsedMemory  uint64
	Uptime      time.Duration
	Version     string
	BuildDate   string
}

// GitHubStatsReply is the reply for the GitHub Service GetStats function
type GitHubStatsReply struct {
	RepoStats gclient.RepoStats
}

// GitHubIssueReplyReply represents
type GitHubIssueReplyReply struct {
	Reply string
}

// KVConfigGetReply holds the value for the respective key
type KVConfigGetReply struct {
	Value string
}

// KVConfigListReply holds the list of keys saved on the database
type KVConfigListReply struct {
	Keys []string
}

// KVConfigListReply holds the list of keys saved on the database
type TLTestReply struct {
	StatusCh chan string
}
