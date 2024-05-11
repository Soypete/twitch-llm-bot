package langchain

import (
	"fmt"
	"time"

	"github.com/Soypete/twitch-llm-bot/database"
	"github.com/google/uuid"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type Client struct {
	llm              *openai.LLM
	db               database.Postgres
	ChatHistory      []llms.MessageContent
	LastChatID       uuid.UUID
	CurrentStartTime time.Time
	Duration         time.Duration
}

func Setup(db database.Postgres) (*Client, error) {
	opts := []openai.Option{
		openai.WithBaseURL("http://127.0.0.1:8080"),
	}
	llm, err := openai.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenAI LLM: %w", err)
	}
	return &Client{
		llm: llm,
		db:  db,
	}, nil
}
