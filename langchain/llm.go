package langchain

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/tmc/langchaingo/llms"
)

func (c Client) PromptWithoutChat(ctx context.Context) (string, error) {
	content, err := llms.GenerateFromSinglePrompt(ctx,
		c.llm,
		"SoyPeteTech's twitch chat has been unusually silent lately. Send a kind chat message to help spark a converastion about software, golang, programming, linux, or food.",
		llms.WithTemperature(0.8),
		llms.WithMaxLength(500),
		llms.WithStopWords([]string{"twitch, SoyPeteTech, bot, assistant, silent, stream, software"}),
	)
	if err != nil {
		return "", fmt.Errorf("failed to get llm response: %w", err)
	}
	return content, nil
}

func (c Client) GetMessageHistory(interval time.Duration) ([]llms.MessageContent, error) {
	messageHistory := []llms.MessageContent{llms.TextParts(llms.ChatMessageTypeSystem, "The following messages are from SoyPeteTech's twitch chat. Please respond to the chat messages to help spark a conversation. You can talk about software, golang, programming, linux, or large language models.")}
	// get message history from database
	messages, err := c.db.QueryMessageHistory(interval)
	if err != nil {
		log.Fatal(err)
	}
	if len(messages) == 0 {
		return nil, fmt.Errorf("no messages found")
	}
	for _, message := range messages {
		prompt := fmt.Sprintf("%s: %s", message.Username, message.Message)
		messageHistory = append(messageHistory, llms.TextParts(llms.ChatMessageTypeHuman, prompt))
	}
	return messageHistory, nil
}

func (c Client) PromptWithChat(ctx context.Context, interval time.Duration) (string, error) {
	log.Println("Getting message history")
	messageHistory, err := c.GetMessageHistory(interval)
	if err != nil {
		return "", fmt.Errorf("failed to get message history: %w", err)
	}
	log.Println("Generating bot response")
	resp, err := c.llm.GenerateContent(ctx, messageHistory,
		llms.WithCandidateCount(1),
		llms.WithMaxLength(500),
		llms.WithTemperature(0.7),
		llms.WithPresencePenalty(0.1), // 2 is the largest penalty for using a work that has already been used
		llms.WithStopWords([]string{"twitch", "stream", "SoyPeteTech", "bot", "assistant", "silent", "software"}))
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	err = c.db.InsertResponse(ctx, resp)
	if err != nil {
		log.Println(err)
	}
	return resp.Choices[0].Content, nil
}
