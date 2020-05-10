package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	gitCmd = &cobra.Command{
		Use:   "git",
		Short: "Utility for ThunderGit",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("NIY")
		},
	}
)

func init() {
	rootCmd.AddCommand(gitCmd)
}
