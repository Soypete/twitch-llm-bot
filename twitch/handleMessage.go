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
		text := strings.ReplaceAll(msg.Message, "]", ":")
		words := strings.Split(text, ":")

		chat.Username = strings.TrimSpace(words[1]) // sets username to the first word after the video source.
		chat.Text = strings.TrimSpace(words[2])     // create a clean message without the video source.
	}
	return chat
}

func needsResponseChat(msg database.TwitchMessage) bool {
	switch {
	case strings.Contains(msg.Text, "pedro"):
		return true
	case strings.Contains(msg.Text, "Pedro"):
		return true
	case strings.Contains(msg.Text, "llm"):
		return true
	case strings.Contains(msg.Text, "LLM"):
		return true
	case strings.Contains(msg.Text, "bot"):
		return true
	default:
		return false
	}
}

func (irc *IRC) HandleChat(ctx context.Context, msg v2.PrivateMessage) {
	chat := cleanMessage(msg)
	if needsResponseChat(chat) {
		resp, err := irc.llm.SingleMessageResponse(ctx, chat)
		if err != nil {
			log.Println("Failed to get response from LLM")
		}
		log.Println(resp)
	}
	// TODO: replace nitbot commands with a classifier model that prompts the LLM
	if err := irc.db.InsertMessage(ctx, chat); err != nil {
		log.Println("Failed to insert message into database")
	}
}
