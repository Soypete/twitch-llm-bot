package langchain

import (
	"fmt"

	"github.com/tmc/langchaingo/llms/openai"
)

type Client struct {
	llm *openai.LLM
}

// TODO: add OpenAI API key as an option
// TODO: add config for options
func Setup() (*Client, error) {
	opts := []openai.Option{
		openai.WithBaseURL("http://127.0.0.1:8080"),
	}
	llm, err := openai.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenAI LLM: %w", err)
	}
	return &Client{llm: llm}, nil
}
