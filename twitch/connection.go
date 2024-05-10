package twitchirc

import (
	"log"
	"sync"

	"github.com/Soypete/twitch-llm-bot/database"
	"github.com/Soypete/twitch-llm-bot/langchain"
	v2 "github.com/gempir/go-twitch-irc/v2"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

const peteTwitchChannel = "soypetetech"

// IRC Connection to the twitch IRC server.
type IRC struct {
	db       database.Postgres
	wg       sync.WaitGroup
	Client   *v2.Client
	tok      *oauth2.Token
	llm      *langchain.Client
	msgQueue chan v2.PrivateMessage
}

// SetupTwitchIRC sets up the IRC, configures oauth, and inits connection functions.
func SetupTwitchIRC(wg sync.WaitGroup, llm *langchain.Client, db database.Postgres) (*IRC, error) {
	irc := &IRC{
		db:       db,
		wg:       wg,
		msgQueue: make(chan v2.PrivateMessage),
		llm:      llm,
	}
	err := irc.AuthTwitch()
	if err != nil {
		return nil, errors.Wrap(err, "failed to authenticate with twitch")
	}

	return irc, nil
}

// connectIRC gets the auth and connects to the twitch IRC server for channel.
func (irc *IRC) ConnectIRC() error {
	log.Println("Connecting to twitch IRC")
	c := v2.NewClient(peteTwitchChannel, "oauth:"+irc.tok.AccessToken)
	c.Join(peteTwitchChannel)
	c.OnConnect(func() { c.Say(peteTwitchChannel, "soy_un_bot esta lista") })
	c.OnPrivateMessage(func(msg v2.PrivateMessage) {
		irc.msgQueue <- msg
	})

	go irc.HandleChat()
	irc.Client = c
	return nil
}
