package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"thunderatz.org/thor/pkg/gitea"
)

var gitClient *gitea.Client

var (
	gitCmd = &cobra.Command{
		Use:   "git",
		Short: "Utility for ThunderGit",

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			token := viper.GetString("gitea.token")
			orgName := viper.GetString("gitea.organization")
			addr := viper.GetString("gitea.addr")

			if token == "" {
				return fmt.Errorf("gitea.token can't be empty")
			}

			var err error

			logger.Debug().
				Str("Token", token).
				Str("addr", addr).
				Str("org", orgName).
				Msg("Creating Gitea Client")

			gitClient, err = gitea.New(addr, token, orgName)

			if err != nil {
				logger.Error().Err(err).Msg("Couldn't create Gitea client")
				return nil
			}

			logger.Debug().Msg("Done")

			return nil
		},

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("All Repositories")

			for i, v := range gitClient.ListRepos() {
				fmt.Println(i+1, v)
			}
		},
	}
)

func init() {
	gitCmd.PersistentFlags().String("token", "", "GitHub Access Token")
	gitCmd.PersistentFlags().String("org", "", "GitHub Organization (default ThunderEletrica)")
	gitCmd.PersistentFlags().String("addr", "", `GitHub Organization (default "https://git.thunderatz.org/")`)

	viper.BindPFlag("gitea.token", gitCmd.PersistentFlags().Lookup("token"))
	viper.BindPFlag("gitea.organization", gitCmd.PersistentFlags().Lookup("org"))
	viper.BindPFlag("gitea.addr", gitCmd.PersistentFlags().Lookup("addr"))

	viper.SetDefault("gitea.organization", "ThunderEletrica")
	viper.SetDefault("gitea.addr", "https://git.thunderatz.org/")

	rootCmd.AddCommand(gitCmd)
}
