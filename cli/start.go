package cli

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"thunderatz.org/thor/services/discord"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Run T.H.O.R.",
	Run: func(cmd *cobra.Command, args []string) {
		discordToken := viper.GetString("discord.token")

		if discordToken == "" {
			logger.Fatal().Msg("Discord token can't be empty")
		}

		discord.Init(discordToken, &logger)

		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

		<-sc
	},
}

func init() {
	startCmd.PersistentFlags().String("dtoken", "", "Discord Token")

	viper.BindPFlag("discord.token", startCmd.PersistentFlags().Lookup("dtoken"))

	rootCmd.AddCommand(startCmd)
}
