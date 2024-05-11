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

func (c Client) PromptWithoutChat() (string, error) {
	ctx := context.Background()
	content, err := llms.GenerateFromSinglePrompt(ctx,
		c.llm,
		"The SoyPeteTech twitch channel has been unusually silent lately. Please generate a creative and kind chat message to help spark a converastion about software, golang, programming, linux, or food.",
		llms.WithTemperature(0.8),
		llms.WithMaxLength(500),
		llms.WithStopWords(stopWords),
	)
	if err != nil {
		return "", fmt.Errorf("failed to get llm response: %w", err)
	}
	return content, nil
}

func GenerateUUID() uuid.UUID {
	return uuid.New()
}
func (c Client) clearMessageHistory() {
	// Clear the message history
	c.LastChatID = GenerateUUID()
	c.ChatHistory = []llms.MessageContent{}
}

func (c Client) PromptWithChat(interval time.Duration) (string, error) {
	ctx := context.Background()

	log.Println("Generating bot response")
	resp, err := c.llm.GenerateContent(ctx, c.ChatHistory,
		llms.WithCandidateCount(1),
		llms.WithMaxLength(500),
		llms.WithTemperature(0.8),
		llms.WithStopWords(stopWords),
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	c.clearMessageHistory()
	if err != nil {
		return "", fmt.Errorf("failed to clear message history: %w", err)
	}
	err = c.db.InsertResponse(resp)
	if err != nil {
		log.Println(err)
	}
	return resp.Choices[0].Content, nil
}

func (c Client) SendSingleChat(ctx context.Context, chat string) (string, error) {
	content, err := llms.GenerateFromSinglePrompt(ctx,
		c.llm,
		chat,
		llms.WithTemperature(0.8), // this is randomness
		llms.WithStopWords([]string{"Chat", "SoyPete", "SoyUnBot", "twitch", "stream"}),
		llms.WithMaxTokens(50),
	)
	if err != nil {
		return "", fmt.Errorf("failed to get llm response: %w", err)
	}
	return content, nil
}
