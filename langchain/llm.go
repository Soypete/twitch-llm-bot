package langchain

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/tmc/langchaingo/llms"
)

var stopWords = []string{"twitch", "stream", "SoyPeteTech", "bot", "assistant", "silent", "software"}

// PromptWithoutChat will generate a prompt without any chat history.
func (c Client) PromptWithoutChat(ctx context.Context) (string, error) {
	content, err := llms.GenerateFromSinglePrompt(ctx,
		c.llm,
		"The SoyPeteTech twitch channel has been unusually silent lately. Please generate a creative and kind chat message to help spark a converastion about software, golang, programming, linux, or food.",
		llms.WithTemperature(0.8),
		llms.WithMaxLength(50),
		llms.WithStopWords(stopWords),
	)
	if err != nil {
		return "", fmt.Errorf("failed to get llm response: %w", err)
	}
	return content, nil
}

// clearMessageHistory will reset the prompt after a message is sent. This includes any information used to connect the chat history to the embeddings in the DB.
func (c Client) clearMessageHistory() {
	// Clear the message history
	c.LastChatID = uuid.New()
	prompt := fmt.Sprintf("This is the previous %d seconds of twitch that. In a respectful, helpful, and kind way please respond to the various users. Most of the viewers are Software engineers and are interested in software engineering.", c.Duration)
	// TODO: This might be a good place to add the context about what I am streaming about.
	c.ChatHistory = []llms.MessageContent{llms.TextParts(llms.ChatMessageTypeSystem, "system", prompt)}
}

// PromptWithChat will generate a prompt with chat history. After teh message is generated, history will be cleared and the
// bot responses will be stored in the database.
func (c Client) PromptWithChat(ctx context.Context, interval time.Duration) (string, error) {
	log.Println("Generating bot response")
	resp, err := c.llm.GenerateContent(ctx, c.ChatHistory,
		llms.WithCandidateCount(1),
		llms.WithMaxLength(50),
		llms.WithTemperature(0.8),
		llms.WithStopWords(stopWords),
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	c.clearMessageHistory()
	err = c.db.InsertAllResponses(ctx, c.LastChatID, resp)
	if err != nil {
		log.Println(err)
	}
	// TODO: should we return the first choice?
	return resp.Choices[0].Content, nil
}
