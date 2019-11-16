package main

import (
	"os"

	"github.com/haydenmcfarland/discord_chan/bot"
	"github.com/haydenmcfarland/discord_chan/config"
	log "github.com/sirupsen/logrus"
)

func main() {
	err := config.ReadConfig()
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	bot.Start()

	<-make(chan struct{})
	return
}
