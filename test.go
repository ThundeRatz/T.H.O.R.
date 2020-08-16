package main

import (
	"fmt"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"thunderatz.org/thor/services/github"
)

var GitHubService = github.Service{
	AppID:          75012,
	InstallationID: 11072487,
	PEMFile:        "t-h-o-r-gh.2020-07-30.private-key.pem",
	WebhookSecret:  "S2G6w!9PUHK*%cx3",
}

func main() {
	logger := zerolog.New(zerolog.NewConsoleWriter()).
		With().
		Timestamp().
		Caller().
		Logger().
		Level(zerolog.DebugLevel)

	GitHubService.Init(&logger, mux.NewRouter())

	fmt.Println("Starting")
	repoStats := GitHubService.GetStats()
	fmt.Println(repoStats)
}
