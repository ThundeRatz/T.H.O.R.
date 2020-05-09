package cli

import (
	homedir "github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	logger  zerolog.Logger

	rootCmd = &cobra.Command{
		Use:   "thor <command> <subcommand>",
		Short: "T.H.O.R. Command Line Interface.",
		Long: `[T.H.O.R | ThundeRatz Holistic Operational Robot]
This is a Command Line Interface application to interface with various ThundeRatz stuff.`,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.thor)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbosity")
}

func initConfig() {
	if Version == "DEV" {
		logger = zerolog.New(zerolog.NewConsoleWriter()).
			With().
			Timestamp().
			Caller().
			Logger().
			Level(zerolog.DebugLevel)
	} else {
		logger = zerolog.New(zerolog.NewConsoleWriter()).
			With().
			Timestamp().
			Logger()

		if v, _ := rootCmd.PersistentFlags().GetBool("verbose"); v {
			logger = logger.Level(zerolog.InfoLevel)
		} else {
			logger = logger.Level(zerolog.WarnLevel)
		}
	}

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		viper.SetConfigType("toml")
	} else {
		home, err := homedir.Dir()
		if err != nil {
			logger.Fatal().Msg("Couldn't detect your home directory")
		}

		viper.SetConfigName(".thor")
		viper.SetConfigType("toml")
		viper.AddConfigPath(home)
	}

	if err := viper.ReadInConfig(); err != nil {
		logger.Warn().Err(err).Msg("Couldn't read config file")
	}

	logger.Debug().Msg("Initialized")
}
