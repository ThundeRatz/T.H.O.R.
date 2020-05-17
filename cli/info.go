package cli

import (
	"net/rpc"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"thunderatz.org/thor/core/types"
)

var infoCmd = &cobra.Command{
	Use: "info",

	Run: func(cmd *cobra.Command, args []string) {
		conn, err := rpc.Dial("unix", viper.GetString("core.socket"))

		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to connect to thor core")
		}
		defer conn.Close()

		ans := types.InfoReply{}

		conn.Call("ThorCore.Info", types.InfoArgs{}, &ans)

		if ans.Success {
			logger.Info().Int("Goroutines", ans.NGoRoutines).Msg("Done")
		} else {
			logger.Info().Msg("Error")
		}
	},
}

func init() {
	infoCmd.PersistentFlags().BoolP("full", "f", false, "get all info")

	rootCmd.AddCommand(infoCmd)
}
