package twitchirc

import (
	"context"
	"log"
	"strings"
	"time"

	database "github.com/Soypete/twitch-llm-bot/database"
	v2 "github.com/gempir/go-twitch-irc/v2"
)

func cleanMessage(msg v2.PrivateMessage) database.TwitchMessage {
	chat := database.TwitchMessage{
		Username: msg.User.DisplayName,
		Text:     msg.Message,
		// TODO: add an embedding for the message
		Time: time.Now(),
	}

	if strings.HasPrefix(msg.Message, "!") {
		chat.IsCommand = true
	}

	if strings.Contains(msg.User.DisplayName, "RestreamBot") {
		words := strings.Split(msg.Message, "]")
		chat.Username = strings.Replace(words[0], "Youtube:", "", 1) // sets username to the first word after the video source.
		chat.Text = strings.Join(words[1:], " ")                     // create a clean message without the video source.
	}
	return chat
}

func (irc *IRC) HandleChat(ctx context.Context, msg v2.PrivateMessage) {
	chat := cleanMessage(msg)
	if msg.User.DisplayName == "Nightbot" {
		return
	}
	// TODO: respond to pedro commands
	if err := irc.db.InsertMessage(ctx, chat); err != nil {
		log.Println("Failed to insert message into database")
	}
}
