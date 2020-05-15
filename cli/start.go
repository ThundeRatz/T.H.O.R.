package cli

import (
	"net/rpc"

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

var testCmd = &cobra.Command{
	Use: "test [msg]",

	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		conn, err := rpc.Dial("unix", viper.GetString("core.socket"))

		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to send")
		}
		defer conn.Close()

		ans := core.SendReply{}

		conn.Call("ThorCore.SendDiscordAlert", core.SendArgs{Msg: args[0]}, &ans)

		if ans.Success {
			logger.Info().Msg("Done")
		} else {
			logger.Info().Msg("Error")
		}
	},
}

func init() {
	startCmd.PersistentFlags().String("dtoken", "", "Discord Token")

	viper.BindPFlag("discord.token", startCmd.PersistentFlags().Lookup("dtoken"))

	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(testCmd)
}
