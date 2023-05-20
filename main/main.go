package main

import (
	"flag"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ccr1m3/osmc/main/db"
	"github.com/ccr1m3/osmc/main/discord"
	"github.com/ccr1m3/osmc/main/discord/slashcommands"
	"github.com/ccr1m3/osmc/main/env"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func init() {
	flag.Parse()
}

func main() {
	// load up .env and set log level
	err := godotenv.Load(".env")
	if err != nil {
		log.Warning("error loading .env file: " + err.Error())
	}
	env.Init()
	if env.LogLevel == env.DEBUG {
		log.SetLevel(log.DebugLevel)
		log.SetReportCaller(true)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	// check for prometheus authorization tokens and start services, when available
	envpromauth := env.Prometheus.Authorization
	envpromrt := env.Prometheus.Refreshtoken
	if strings.Compare(envpromauth, "") == 0 || strings.Compare(envpromrt, "") == 0 {
		log.Warning("no omega strikers authorization given, shutting down")
	} else {
		db.Init()
		discord.Init()
		slashcommands.Init()
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		log.Info("initialization done")
		<-stop
		log.Info("gracefully shutting down.")
		discord.Stop()
	}
}
