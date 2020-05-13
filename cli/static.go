package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"thunderatz.org/thor/pkg/static"
)

// Local Flags
var (
	parent    bool
	overwrite bool
)

var staticClient *static.Conn

// Commands
var (
	staticCmd = &cobra.Command{
		Use:   "static",
		Short: "Interact with static files on the server",

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			addr := viper.GetString("server.addr")
			user := viper.GetString("server.user")
			pass := viper.GetString("server.password")
			base := viper.GetString("server.static.root")
			port := viper.GetInt("server.shh_port")

			if addr == "" {
				return fmt.Errorf("Static folder can't be empty")
			}

			if user == "" {
				return fmt.Errorf("Server user can't be empty")
			}

			if addr == "" {
				return fmt.Errorf("Server addr can't be empty")
			}

			var err error

			logger.Debug().
				Str("Addr", addr).
				Str("User", user).
				Str("Base Path", base).
				Int("SSH Port", port).
				Msg("Connecting to SFTP...")

			staticClient, err = static.New(addr, user, pass, base, port)

			if err != nil {
				logger.Error().Err(err).Msg("Couldn't connect to server")
				return nil
			}

			logger.Debug().Msg("Done")

			return nil
		},

		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if staticClient != nil {
				staticClient.Close()
				staticClient = nil
			}
		},
	}

	staticListCmd = &cobra.Command{
		Use:     "list [directory]",
		Aliases: []string{"ls"},
		Short:   "List files and directories.",

		Run: func(cmd *cobra.Command, args []string) {
			folder := ""

			if len(args) > 0 {
				folder = "/" + args[0]
			}

			logger.Debug().Msg("Getting files")

			files, err := staticClient.List(folder)

			if err != nil {
				logger.Error().Err(err).Msg("Couldn't list files")
				return
			}

			for _, f := range files {
				if f.IsDir() {
					fmt.Println(f.Name() + "/")
				} else {
					fmt.Println(f.Name())
				}
			}

			logger.Debug().Msg("Done")
		},
	}

	staticMkdirCmd = &cobra.Command{
		Use:   "mkdir <directory>",
		Short: "Creates a directory",

		Args: cobra.ExactArgs(1),

		Run: func(cmd *cobra.Command, args []string) {
			err := staticClient.Mkdir(args[0], parent)

			if err != nil {
				logger.Error().Msg("Couldn't create folder, maybe it already exists or you need to create it's parents first (-p)")
				return
			}

			logger.Info().Str("Directory", args[0]).Msg("Done!")
		},
	}

	staticGetCmd = &cobra.Command{
		Use:   "get <file>",
		Short: "Get a file",

		Args: cobra.ExactArgs(1),

		Run: func(cmd *cobra.Command, args []string) {
			err := staticClient.Get(args[0])

			if err != nil {
				logger.Error().Err(err).Str("File", args[0]).Msg("Couldn't get file")
				return
			}

			logger.Info().Str("File", args[0]).Msg("Done!")
		},
	}

	staticPutCmd = &cobra.Command{
		Use:   "put <file> [remote directory]",
		Short: "Puts a file on the static server",

		Args: cobra.RangeArgs(1, 2),

		Run: func(cmd *cobra.Command, args []string) {
			remote := "/"

			if len(args) >= 2 {
				remote = args[1]
			}

			err := staticClient.Put(args[0], remote, overwrite)

			if err != nil {
				logger.Error().Err(err).Str("File", args[0]).Str("Folder", remote).Msg("Couldn't put file")
				return
			}

			logger.Info().Str("File", args[0]).Msg("Done!")
		},
	}
)

func init() {
	staticCmd.PersistentFlags().String("root", "", "Static apps root directory")
	staticCmd.PersistentFlags().String("user", "", "Server user")
	staticCmd.PersistentFlags().String("addr", "", "ThundeRatz Server Address")
	staticCmd.PersistentFlags().Int("ssh-port", 22, "ThundeRatz Server SSH port")

	viper.BindPFlag("server.static.root", staticCmd.PersistentFlags().Lookup("root"))
	viper.BindPFlag("server.user", staticCmd.PersistentFlags().Lookup("user"))
	viper.BindPFlag("server.addr", staticCmd.PersistentFlags().Lookup("addr"))
	viper.BindPFlag("server.shh_port", staticCmd.PersistentFlags().Lookup("ssh-port"))

	viper.SetDefault("server.shh_port", 22)

	staticCmd.AddCommand(staticListCmd)

	staticMkdirCmd.Flags().BoolVarP(&parent, "parents", "p", false, "Creates all parent directories as needed")

	staticCmd.AddCommand(staticMkdirCmd)
	staticCmd.AddCommand(staticGetCmd)

	staticPutCmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrites file if it already exists")

	staticCmd.AddCommand(staticPutCmd)
	rootCmd.AddCommand(staticCmd)
}
