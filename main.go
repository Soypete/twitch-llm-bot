package main

import (
	"log"
	"sync"

	twitchirc "github.com/Soypete/twitch-llm-bot/twitch"
)

func main() {
	wg := sync.WaitGroup{}
	// setup twitch IRC
	irc, err := twitchirc.SetupTwitchIRC(wg)
	if err != nil {
		log.Fatalln(err)
	}
	irc.Client.Say("soypetetech", "soy_un_bot connected")
}
