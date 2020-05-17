package core

import (
	"net"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"thunderatz.org/thor/core/types"
	"thunderatz.org/thor/services/discord"
)

var logger zerolog.Logger

// ThorCore is
type ThorCore struct{}

// Services
var (
	DiscordService *discord.Service
)

// Initialize logger
func initLogger() {
	logger = zerolog.New(zerolog.NewConsoleWriter()).
		With().
		Timestamp().
		Logger()

	if v := viper.GetInt("core.verbosity"); v >= 2 {
		logger = logger.Level(zerolog.DebugLevel)
	} else if v == 1 {
		logger = logger.Level(zerolog.InfoLevel)
	} else {
		logger = logger.Level(zerolog.WarnLevel)
	}
}

func initDiscordService() {
	DiscordService = &discord.Service{
		AlertChannel: viper.GetString("discord.alert_channel_id"),
		Token:        viper.GetString("discord.token"),
	}

	err := DiscordService.Init(&logger)

	if err != nil {
		logger.Fatal().Err(err).Msg("Couldn't start discord service")
	}
}

func initServices() {
	initDiscordService()
}

func initSocket() net.Listener {
	socket := viper.GetString("core.socket")

	if err := os.RemoveAll(socket); err != nil {
		logger.Error().Err(err).Send()
	}

	rpcServer := rpc.NewServer()
	rpcServer.Register(&ThorCore{})

	la, err := net.Listen("unix", socket)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to start rpc server")
	}

	go rpcServer.Accept(la)

	logger.Info().Msg("Core RPC server initialized")
	return la
}

// SendDiscordAlert sends an alert to the alert discord channel
func (tc *ThorCore) SendDiscordAlert(args types.SendArgs, reply *types.SendReply) error {
	logger.Debug().Str("msg", args.Msg).Msg("Received discord alert request")

	err := DiscordService.SendMessage(args.Msg)
	if err != nil {
		reply.Success = false
	} else {
		reply.Success = true
	}
	return nil
}

// Start starts the core T.H.O.R. process
func Start() {
	initLogger()

	sock := initSocket()
	defer sock.Close()

	initServices()
	defer DiscordService.Stop()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	<-sc
}
