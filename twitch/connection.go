package twitchirc

import (
	"fmt"
	"sync"

	v2 "github.com/gempir/go-twitch-irc/v2"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

const peteTwitchChannel = "soypetetech"

// IRC Connection to the twitch IRC server.
type IRC struct {
	wg       sync.WaitGroup
	Client   *v2.Client
	tok      *oauth2.Token
	msgQueue chan string
}

// SetupTwitchIRC sets up the IRC, configures oauth, and inits connection functions.
func SetupTwitchIRC(wg sync.WaitGroup) (*IRC, error) {
	irc := &IRC{
		wg:       wg,
		msgQueue: make(chan string),
	}
	err := irc.AuthTwitch()
	if err != nil {
		return nil, errors.Wrap(err, "failed to authenticate with twitch")
	}

	return irc, nil
}

// connectIRC gets the auth and connects to the twitch IRC server for channel.
func (irc *IRC) connectIRC() error {
	c := v2.NewClient(peteTwitchChannel, "oauth:"+irc.tok.AccessToken)
	c.Join(peteTwitchChannel)
	c.OnConnect(func() { c.Say(peteTwitchChannel, "grpc twitch bot connected") })
	c.OnPrivateMessage(func(msg v2.PrivateMessage) {
		fmt.Println(msg.Message)
		fmt.Println(irc.msgQueue)
		for m := range irc.msgQueue {
			c.Say(peteTwitchChannel, m)
		}
	})
	err := c.Connect()
	if err != nil {
		return errors.Wrap(err, "failed to connect over IRC")
	}
	irc.Client = c
	return nil
}
