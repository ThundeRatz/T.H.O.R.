package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"thunderatz.org/thor/core"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Run T.H.O.R.",
	Run: func(cmd *cobra.Command, args []string) {
		core.Start()
	},
}

func init() {
	startCmd.PersistentFlags().String("dtoken", "", "Discord Token")

	viper.BindPFlag("discord.token", startCmd.PersistentFlags().Lookup("dtoken"))

	rootCmd.AddCommand(startCmd)
}
