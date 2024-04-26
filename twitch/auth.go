package twitchirc

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/twitch"
)

func parseAuthCode(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		fmt.Printf("could not parse query: %v", err)
		http.Error(w, "could not parse query", http.StatusBadRequest)
	}
	code := req.FormValue("code")
	fmt.Fprint(os.Stdout, code)
}

// AuthTwitch use oauth2 protocol to retrieve oauth2 token for twitch IRC.
// _NOTE_: this has not been tested on long standing projects.
func (irc *IRC) AuthTwitch() error {
	http.HandleFunc("/oauth/redirect", parseAuthCode)
	go http.ListenAndServe("localhost:3000", nil)

	ctx := context.Background()
	conf := &oauth2.Config{
		// TODO: use const for the following.
		ClientID:     os.Getenv("TWITCH_ID"),
		ClientSecret: os.Getenv("TWITCH_SECRET"),
		Scopes:       []string{"chat:read", "chat:edit", "channel:moderate"},
		RedirectURL:  "http://localhost:3000/oauth/redirect",
		Endpoint:     twitch.Endpoint,
	}
	irc.wg.Add(1)

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	go func() {
		defer irc.wg.Done()
		url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
		fmt.Printf("Visit the URL for the auth dialog: %v\n", url)

		time.Sleep(5 * time.Second)

		var code string
		_, err := fmt.Scan(&code)
		if err != nil {
			fmt.Println(fmt.Errorf("cannot get input from standard in: %w", err))
		}

		log.Printf("code: %v", code)

		irc.tok, err = conf.Exchange(ctx, code)
		if err != nil {
			fmt.Println(fmt.Errorf("failed to get token with auth code: %w", err))
		}
		fmt.Println("token received")
		// _ = conf.Client(ctx, irc.tok)

		fmt.Println("start twitch connection")
		err = irc.connectIRC()
		if err != nil {
			fmt.Println(fmt.Errorf("failed to connect over IRC: %w", err))
		}
		fmt.Println("connection completed")
	}()
	irc.wg.Wait()
	return nil
}
