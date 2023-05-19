package main

import (
	"flag"

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
}
