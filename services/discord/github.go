package discord

import (
	"fmt"
	"sort"
	"strings"

	"go.thunderatz.org/thor/core/types"
	"go.thunderatz.org/thor/pkg/dclient"
	"go.thunderatz.org/thor/pkg/gclient"
)

func getOrderedStats(repoStats *gclient.RepoStats) string {
	ordered := func(sl map[string]int) []string {
		type kv struct {
			k string
			v int
		}

		var ss []kv
		for k, v := range sl {
			ss = append(ss, kv{k, v})
		}

		sort.Slice(ss, func(i, j int) bool {
			return ss[i].v > ss[j].v
		})

		ranked := make([]string, len(sl))

		for i, kv := range ss {
			ranked[i] = kv.k
		}

		return ranked
	}

	var out strings.Builder
	out.WriteString("== Top 10 Additions, Deletions e Commits ==\n")
	out.WriteString("```\nAdditions:\n")

	count := 0
	for _, p := range ordered(repoStats.Adds) {
		if repoStats.Adds[p] == 0 {
			continue
		}

		if count++; count == 10 {
			break
		}

		fmt.Fprintf(&out, "%s: %d\n", p, repoStats.Adds[p])
	}

	out.WriteString("\nDeletions:\n")

	count = 0
	for _, p := range ordered(repoStats.Dels) {
		if repoStats.Dels[p] == 0 {
			continue
		}

		if count++; count == 10 {
			break
		}

		fmt.Fprintf(&out, "%s: %d\n", p, repoStats.Dels[p])
	}

	out.WriteString("\nCommits:\n")

	count = 0
	for _, p := range ordered(repoStats.Commits) {
		if repoStats.Commits[p] == 0 {
			continue
		}

		if count++; count == 10 {
			break
		}

		fmt.Fprintf(&out, "%s: %d\n", p, repoStats.Commits[p])
	}

	out.WriteString("```")

	return out.String()
}

var githubCmd = &dclient.Command{
	Name:        "github",
	Category:    "Tools",
	Description: "Comandos relacionados ao GitHub da equipe",
	Usage:       "github [subcommand]",

	Enabled:   true,
	GuildOnly: true,
	Aliases:   []string{"gh"},
	PermLevel: "User",

	Run: func(c *dclient.Context) {
		if len(c.Args) < 1 {
			c.Session.ChannelMessageSend(c.Message.ChannelID, "Comandos disponíveis: `stats`")
			return
		}

		replyCh := make(types.CoreReplyCh)

		switch c.Args[0] {
		case "stats":
			msg, _ := c.Session.ChannelMessageSend(c.Message.ChannelID, "Buscando, isso pode demorar um pouco...")

			msgCh <- types.CoreMsg{
				Type:  types.GitHubStatsMsg,
				Reply: replyCh,
				From:  serviceId,
			}

			reply := <-replyCh
			stats := reply.Reply.(*types.GitHubStatsReply).RepoStats

			c.Session.ChannelMessageEdit(c.Message.ChannelID, msg.ID, getOrderedStats(&stats))

		default:
			c.Session.ChannelMessageSend(c.Message.ChannelID, "Não conheço esse comando")
		}
	},
}
