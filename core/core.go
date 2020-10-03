package core

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"

	"thunderatz.org/thor/core/types"
	"thunderatz.org/thor/services/discord"
	"thunderatz.org/thor/services/github"
)

var logger zerolog.Logger

// ThorCore is
type ThorCore struct{}

// Services
var (
	DiscordService *discord.Service
	GitHubService  *github.Service
)

var (
	MsgCh types.CoreMsgCh
	root  *mux.Router
)

func init() {
	root = mux.NewRouter()

	root.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Debug().Str("url", r.URL.String()).Msg("New Request")
			next.ServeHTTP(w, r)
		})
	})
}

// Initialize logger
func initLogger() {
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

	err := DiscordService.Init(&logger, MsgCh)

	if err != nil {
		logger.Fatal().Err(err).Msg("Couldn't start discord service")
	}
}

func initGitHubService() {
	ghRouter := root.PathPrefix("/gh").Subrouter()

	GitHubService = &github.Service{
		AppID:          viper.GetInt64("github.app_id"),
		InstallationID: viper.GetInt64("github.installation_id"),
		PEMFile:        viper.GetString("github.pem_file"),
		WebhookSecret:  viper.GetString("github.webhook_secret"),
	}

	GitHubService.Init(&logger, ghRouter, MsgCh)
}

func initServices() {
	initDiscordService()
	initGitHubService()
}

func initAPI() *http.Server {
	srv := &http.Server{
		Addr:         "127.0.0.1:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      root,
	}

	go func() {
		logger.Info().Msg("Starting API server")

		if err := http.ListenAndServe(":8080", root); err != nil {
			logger.Error().Err(err).Send()
		}
	}()

	return srv
}

func processForever() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	for {
		select {
		case msg := <-MsgCh:
			logger.Info().Int("type", int(msg.Type)).Msg("Received Message")
			go ProcessMsg(msg)

		case <-sc:
			return
		}
	}
}

// Start starts the core T.H.O.R. process
func Start() {
	initConfig()

	MsgCh = make(types.CoreMsgCh, 10)

	initLogger()

	initServices()
	defer DiscordService.Stop()

	server := initAPI()

	processForever()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	server.Shutdown(ctx)

	logger.Info().Msg("Shutting down")
	os.Exit(0)
}
