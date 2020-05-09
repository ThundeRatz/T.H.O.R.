package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Build Information, overriden on build
var (
	Version   = "DEV"
	BuildDate = ""
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version and exit",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("[T.H.O.R. - ThundeRatz Holistic Operational Robot]")
		fmt.Printf("Version %s - BuildDate %s\n", Version, BuildDate)
	},
}

func init() {
	rootCmd.Version = fmt.Sprintf("%s (%s)", Version, BuildDate)

	rootCmd.AddCommand(versionCmd)
}
