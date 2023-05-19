package main

import (
	"flag"
	"os"
	"os/signal"
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
