package twitchirc

import (
	"context"
	"log"

	v2 "github.com/gempir/go-twitch-irc/v2"
)

func (irc *IRC) HandleChat(ctx context.Context, msg v2.PrivateMessage) {
	if err := irc.db.InsertMessage(ctx, msg); err != nil {
		log.Println("Failed to insert message into database")
	}
}
