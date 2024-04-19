package langchain

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms"
)

func (c Client) SendChat(ctx context.Context, chat string) (string, error) {
	// llmContent:= llms.TextParts(chat)
	// content, err := c.llm.GenerateContent(ctx, llmContent)
	content, err := llms.GenerateFromSinglePrompt(ctx,
		c.llm,
		chat,
		llms.WithTemperature(0.8),
		llms.WithStopWords([]string{"Chat"}),
	)
	if err != nil {
		return "", fmt.Errorf("failed to get llm response: %w", err)
	}
	return content, nil
}
