package main

import (
	"log"
	"sync"

	"github.com/Soypete/twitch-llm-bot/langchain"
	twitchirc "github.com/Soypete/twitch-llm-bot/twitch"
)

func main() {
	// setup llm connection
	llm, err := langchain.Setup()
	if err != nil {
		log.Fatalln(err)
	}
	if llm == nil {
		log.Fatalln("llm is nil")
	}

	wg := sync.WaitGroup{}
	// setup twitch IRC
	irc, err := twitchirc.SetupTwitchIRC(wg, llm)
	if err != nil {
		log.Fatalln(err)
	}
	irc.Client.Say("soypetetech", "soy_un_bot connected")
}
