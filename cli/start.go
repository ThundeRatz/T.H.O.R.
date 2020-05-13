package cli

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"thunderatz.org/thor/pkg/discord"
	"thunderatz.org/thor/pkg/discord/commands"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Run T.H.O.R.",
	Run: func(cmd *cobra.Command, args []string) {
		client := discord.New(viper.GetString("discord.token"), &logger)

		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

		commands.AddAllCommands(client)

		client.Start()
		defer client.Stop()

		<-sc
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
