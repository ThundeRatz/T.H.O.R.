package cli

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"thunderatz.org/thor/pkg/github"
)

var ghClient *github.Client

var ghCmd = &cobra.Command{
	Use:   "gh",
	Short: "GitHub information",

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		token := viper.GetString("github.token")

		if token == "" {
			return fmt.Errorf("github.token can't be empty")
		}

		var err error

		logger.Debug().
			Str("Token", token).
			Msg("Creating GitHub Client")

		ghClient, err = github.New(token, "ThundeRatz", &logger)

		if err != nil {
			logger.Error().Err(err).Msg("Couldn't create GitHub client")
			return nil
		}

		logger.Debug().Msg("Done")

		return nil
	},
}

var ghStatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Members contributions to all repositories",

	Run: func(cmd *cobra.Command, args []string) {
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

		rs := ghClient.GetStats()

		fmt.Println("Additions:")

		for _, p := range ordered(rs.Adds) {
			if rs.Adds[p] == 0 {
				continue
			}

			fmt.Printf("%s: %d\n", p, rs.Adds[p])
		}

		fmt.Println("\nDeletions:")

		for _, p := range ordered(rs.Dels) {
			if rs.Dels[p] == 0 {
				continue
			}

			fmt.Printf("%s: %d\n", p, rs.Dels[p])
		}

		fmt.Println("\nCommits:")

		for _, p := range ordered(rs.Commits) {
			if rs.Commits[p] == 0 {
				continue
			}

			fmt.Printf("%s: %d\n", p, rs.Commits[p])
		}
	},
}

func init() {
	ghCmd.PersistentFlags().String("token", "", "GitHub Access Token")

	viper.BindPFlag("github.token", ghCmd.PersistentFlags().Lookup("token"))

	ghCmd.AddCommand(ghStatsCmd)
	rootCmd.AddCommand(ghCmd)
}
