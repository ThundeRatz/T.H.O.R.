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

		var key string

		// Discord service can get and set any config key and is not restricted to its namespace
		if msg.From == "discord" {
			key = args.Key
		} else {
			key = msg.From + "/" + args.Key
		}

		value, err := GetConfig(key)

		msg.Reply <- types.CoreReply{
			Success: err == nil,
			Reply: &types.KVConfigGetReply{
				Value: value,
			},
		}

	case types.KVConfigListMsg:
		keys, err := GetKeyList()

		msg.Reply <- types.CoreReply{
			Success: err == nil,
			Reply: &types.KVConfigListReply{
				Keys: keys,
			},
		}

	case types.KVConfigSetMsg:
		args := msg.Args.(types.KVConfigSetArgs)

		var key string

		// Discord service can get and set any config key and is not restricted to its namespace
		if msg.From == "discord" {
			key = args.Key
		} else {
			key = msg.From + "/" + args.Key
		}

		err := SetConfig(key, args.Value)

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

		value, err := GetConfig("github/issue-reply/" + args.Lang)

		if err != nil {
			// Defaults to english if not found
			value, err = GetConfig("github/issue-reply/en")
		}

		msg.Reply <- types.CoreReply{
			Success: err == nil,
			Reply: &types.GitHubIssueReplyReply{
				Reply: value,
			},
		}

	case types.TLTestMsg:
		args := msg.Args.(types.TLTestArgs)

		scOut := make(chan string)
		scIn := make(chan string)
		err := ThunderServerService.StartThunderLeagueDroplet()

		msg.Reply <- types.CoreReply{
			Success: err == nil,
			Reply: &types.TLTestReply{
				StatusCh: scOut,
			},
		}

		if err != nil {
			close(scOut)
			close(scIn)
			return
		}

		go func() {
			err = ThunderServerService.RunThunderLeagueTest(args.Commit, args.Enemy, args.Amount, scIn)
			close(scIn)

			if err != nil {
				logger.Error().Err(err).Msg("Error on RunThunderLeagueTest")
				return
			}
		}()

		for msg := range scIn {
			logger.Debug().Msg(msg)
			scOut <- msg
		}

		scOut <- "Test finished, results NYI"
		close(scOut)

		err = ThunderServerService.StopThunderLeagueDroplet()

		if err != nil {
			logger.Error().Err(err).Msg("Failed to stop thunderleague droplet")
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
