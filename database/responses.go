package database

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/llms"
)

func (p Postgres) InsertResponse(ctx context.Context, resp *llms.ContentResponse) error {
	var isUsed bool
	for i, choice := range resp.Choices {
		if i == 0 {
			isUsed = true
		}
		query := "INSERT INTO bot_response (model_name, response, stop_reason, was_successful) VALUES ($1, $2, $3, $4)"
		_, err := p.connections.ExecContext(ctx, query, p.modelName, choice.Content, choice.StopReason, isUsed)
		if err != nil {
			return fmt.Errorf("error upserting response: %w", err)
		}
	}
	return nil
}

func (p Postgres) InsertAllResponses(ctx context.Context, chatId int, responses *llms.ContentResponse) error {
	text := []string{}
	stopReason := []string{}
	for _, choice := range responses.Choices {
		text = append(text, choice.Content)
		stopReason = append(stopReason, choice.StopReason)
	}
	query := "INSERT INTO bot_response (model_name, responses, stop_reasons, twitch_chat_id) VALUES ($1, $2, $3, $4)"
	_, err := p.connections.ExecContext(ctx, query, p.modelName, text, stopReason, chatId)
	if err != nil {
		return fmt.Errorf("error upserting response: %w", err)
	}
	return nil
}
