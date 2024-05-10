package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	database "github.com/Soypete/twitch-llm-bot/database"
	"github.com/Soypete/twitch-llm-bot/langchain"
	twitchirc "github.com/Soypete/twitch-llm-bot/twitch"
)

func main() {
	// setup postgres connection
	db, err := database.NewPostgres()
	if err != nil {
		log.Fatalln(err)
	}
	// setup llm connection
	llm, err := langchain.Setup(db)
	if err != nil {
		log.Fatalln(err)
	}
	if llm == nil {
		log.Fatalln("llm is nil")
	}

	// TODO: audit waitgroup
	wg := sync.WaitGroup{}
	// setup twitch IRC
	irc, err := twitchirc.SetupTwitchIRC(wg, llm, db)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("starting twitch IRC connection")
	// long running function
	err = irc.ConnectIRC()
	if err != nil {
		panic(err)
	}
	go func() {
		log.Println("Starting prompt loop")
		for {
			timeout := 1 * time.Minute
			time.Sleep(timeout)
			log.Println("Getting prompt")
			prompt, err := llm.PromptWithChat(timeout)
			// weird error handling
			switch {
			case err == nil:
				irc.Client.Say("soypetetech", prompt)
			case strings.Contains(err.Error(), "no messages found"):
				log.Println("No messages found, generating prompt without chat")
				prompt, err = llm.PromptWithoutChat()
				if err != nil {
					log.Println(err)
				}
				irc.Client.Say("soypetetech", prompt)
			default:
				log.Println(err)
			}
		}
	}()

	// TODO: why is this not in a goroutine?
	err = irc.Client.Connect()
	if err != nil {
		panic(fmt.Errorf("failed to connect to twitch IRC: %w", err))
	}
}

func Shutdown(ctx context.Context, wg *sync.WaitGroup) {
	ctx.Done()
	log.Println("Shutting down")

}
