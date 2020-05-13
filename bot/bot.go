package bot

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"

	"thunderatz.org/thor/services/discord"
)

var logger *zerolog.Logger

// RunForever runs the bot
func RunForever(_logger *zerolog.Logger) {
	logger = _logger

	discordToken := viper.GetString("discord.token")

	if discordToken == "" {
		logger.Fatal().Msg("Discord token can't be empty")
	}

	discord.Init(discordToken, logger)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	discord.Start()
	<-sc
}
