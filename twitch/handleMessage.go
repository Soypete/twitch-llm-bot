package twitchirc

import (
	"fmt"
	"log"

	"github.com/tmc/langchaingo/llms"
)

// TODO: this should be in an different package
func (irc *IRC) HandleChat() {
	// TODO: close if channel is closed
	for {
		log.Println("Handling chat messages")
		msg := <-irc.msgQueue
		irc.llm.ChatHistory = append(
			irc.llm.ChatHistory,
			llms.TextParts(llms.ChatMessageTypeHuman,
				msg.User.DisplayName, msg.Message))
		chat := fmt.Sprintf("%s: %s", msg.User.DisplayName, msg.Message)
		go irc.db.AppendChatHistory(irc.llm.LastChatID, chat, irc.llm.CurrentStartTime, irc.llm.Duration)
	}
}
