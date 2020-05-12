package cli

import (
	"fmt"

	"github.com/manifoldco/promptui"
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
	}

	gitListCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List all organization repositories",

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("All Repositories")

			for i, v := range gitClient.ListRepos() {
				fmt.Println(i+1, v)
			}
		},
	}

	gitCloneCmd = &cobra.Command{
		Use:   "clone [repository | all | interactive]",
		Short: "Clone a repository",

		Args: cobra.ExactArgs(1),

		Run: func(cmd *cobra.Command, args []string) {
			verbose, _ := rootCmd.PersistentFlags().GetBool("verbose")

			if args[0] != "all" && args[0] != "interactive" {
				err := gitClient.CloneRepo(args[0], verbose)

				if err != nil {
					logger.Error().
						Err(err).
						Msg("Failed to clone repository")
				}

				return
			}

			repoList := gitClient.ListRepos()

			if args[0] == "all" {
				prompt := promptui.Select{
					Label: fmt.Sprintf("Are you sure you want to clone %d repositories?", len(repoList)),
					Items: []string{"Yes", "No"},
				}

				_, result, err := prompt.Run()

				if err != nil {
					result = "No"
				}

				if result == "Yes" {
					for _, v := range repoList {
						logger.Info().
							Str("repo", v).
							Msg("Cloning repository")

						err := gitClient.CloneRepo(v, verbose)

						if err != nil {
							logger.Error().
								Str("repo", v).
								Err(err).
								Msg("Failed to clone repository")
						}
					}
				}
			}

			if args[0] == "interactive" {
				for _, v := range repoList {
					prompt := promptui.Select{
						Label: fmt.Sprintf("Clone %s?", v),
						Items: []string{"Yes", "No", "Cancel"},
					}

					_, result, err := prompt.Run()

					if err != nil {
						result = "Cancel"
					}

					if result == "Cancel" {
						break
					}

					if result == "Yes" {
						logger.Info().
							Str("repo", v).
							Msg("Cloning repository")

						err := gitClient.CloneRepo(v, verbose)

						if err != nil {
							logger.Error().
								Str("repo", v).
								Err(err).
								Msg("Failed to clone repository")
						}
					}
				}
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

	gitCmd.AddCommand(gitListCmd)
	gitCmd.AddCommand(gitCloneCmd)

	rootCmd.AddCommand(gitCmd)
}
