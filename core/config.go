package core

import (
	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

// Build Information, overriden on build
var (
	Version   = "DEV"
	BuildDate = ""
)

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
			Logger().
			Level(zerolog.InfoLevel)
	}

	home, err := homedir.Dir()

	if err != nil {
		logger.Fatal().Msg("Couldn't detect your home directory")
	} else {
		viper.SetConfigName(".thor")
		viper.SetConfigType("toml")
		viper.AddConfigPath(home)
	}

	if err := viper.ReadInConfig(); err != nil {
		logger.Fatal().Err(err).Msg("Couldn't read config file")
	}

	logger.Debug().Msg("Initialized")
}
