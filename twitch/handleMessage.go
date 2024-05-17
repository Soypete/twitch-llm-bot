package twitchirc

import (
	"context"
	"log"
	"time"

	v2 "github.com/gempir/go-twitch-irc/v2"
	"github.com/tmc/langchaingo/llms"
)

// TODO: this should be in an different package

// HandleChat receives chat messages from IRC goroutine and appends them to the chat history.
func (irc *IRC) HandleChat(msg v2.PrivateMessage) {
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, 10*time.Second)
	irc.mu.Lock()
	defer irc.mu.Unlock()
	irc.llm.ChatHistory = append(irc.llm.ChatHistory,
		llms.TextParts(llms.ChatMessageTypeHuman, msg.User.DisplayName, msg.Message))

	if err := irc.db.InsertMessage(ctx, msg); err != nil {
		log.Printf("Error appending chat history: %v", err)
	}
}
