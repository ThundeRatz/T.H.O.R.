package core

import (
	"runtime"
	"time"

	"go.thunderatz.org/thor/core/types"
)

func ProcessMsg(msg types.CoreMsg) {
	switch msg.Type {
	case types.InfoMsg:
		reply := &types.InfoReply{}
		err := Info(reply)

		msg.Reply <- types.CoreReply{
			Success: err == nil,
			Reply:   reply,
		}

	case types.PingMsg:
		err := Ping()

		msg.Reply <- types.CoreReply{
			Success: err == nil,
			Reply:   nil,
		}

	case types.KVConfigGetMsg:
		args := msg.Args.(types.KVConfigGetArgs)

		value, err := GetConfig(msg.From + "/" + args.Key)

		msg.Reply <- types.CoreReply{
			Success: err == nil,
			Reply: &types.KVConfigGetReply{
				Value: value,
			},
		}

	case types.KVConfigSetMsg:
		args := msg.Args.(types.KVConfigSetArgs)

		err := SetConfig(msg.From+"/"+args.Key, args.Value)

		msg.Reply <- types.CoreReply{
			Success: err == nil,
			Reply:   nil,
		}

	case types.GitHubStatsMsg:
		reply := &types.GitHubStatsReply{
			RepoStats: GitHubService.GetStats(),
		}

		msg.Reply <- types.CoreReply{
			Success: true,
			Reply:   reply,
		}

	case types.GitHubEventMsg:
		args := msg.Args.(types.GitHubEventArgs)
		issue := args.Issue
		repo := args.Repository

		DiscordService.SendGitHubIssueAlert(issue, repo.GetName())

		// For now, will only comment back on travesim repo
		if repo.GetName() == "travesim" || repo.GetName() == "test" {
			GitHubService.SendDefaultVSSIssueMessage(repo.GetName(), issue.GetBody(), issue.GetNumber())
		}

	case types.GitHubIssueReplyMsg:
		args := msg.Args.(types.GitHubIssueReplyArgs)

		// This config is set by the discord service, so we get it from its config pool
		value, err := GetConfig("discord/issue-reply/" + args.Lang)

		if err != nil {
			// Defaults to english if not found
			value, err = GetConfig("discord/issue-reply/en")
		}

		msg.Reply <- types.CoreReply{
			Success: err == nil,
			Reply: &types.GitHubIssueReplyReply{
				Reply: value,
			},
		}
	}

}

// Ping does nothing, can be used to check if core is up and response time
func Ping() error {
	logger.Debug().Msg("Received ping request")
	return nil
}

// Info returns information about the current running thor core and services
func Info(reply *types.InfoReply) error {
	logger.Debug().Msg("Received info request")

	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)

	reply.NGoRoutines = runtime.NumGoroutine()
	reply.UsedMemory = m.Alloc

	reply.Uptime = time.Since(StartTime)
	reply.Version = Version
	reply.BuildDate = BuildDate

	return nil
}
