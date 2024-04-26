package langchain

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms"
)

func (c Client) SendChat(ctx context.Context, chat string) (string, error) {
	prompt := fmt.Sprintf("This was the last message from twitch chat: %s\n please respond.", chat)
	content, err := llms.GenerateFromSinglePrompt(ctx,
		c.llm,
		prompt,
		llms.WithTemperature(0.8), // this is randomness
		llms.WithStopWords([]string{"Chat", "SoyPete", "SoyUnBot", "twitch", "stream"}),
		llms.WithMaxTokens(50),
	)
	if err != nil {
		return "", fmt.Errorf("failed to get llm response: %w", err)
	}
	return content, nil
}
