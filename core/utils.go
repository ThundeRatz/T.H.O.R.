package core

import (
	"runtime"

	"thunderatz.org/thor/core/types"
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

	case types.GitHubStatsMsg:
		reply := &types.GitHubStatsReply{
			RepoStats: GitHubService.GetStats(),
		}

		msg.Reply <- types.CoreReply{
			Success: true,
			Reply:   reply,
		}

	case types.GitHubEventMsg:
		logger.Debug().Msg("Received GitHubEvent request")
		args := msg.Args.(types.GitHubEventArgs)

		DiscordService.SendGitHubIssueAlert(args.Issue, args.Repository.GetName())
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

	reply.NGoRoutines = runtime.NumGoroutine()

	return nil
}
