package langchain

import (
	"context"
	"fmt"
	"strings"

	database "github.com/Soypete/twitch-llm-bot/database"
	"github.com/tmc/langchaingo/llms"
)

const pedroPrompt = "Your name is Pedro_el_asistente. You are a chat bot that helps out in SoyPeteTech's twitch chat. You are allowed to use links, code, or emotes to express fun messages about software. Helpful links are always appreciated, such as SoyPeteTech's github https://github.com/Soypete, youtube https://www.youtube.com/channel/UCEkM7JXVQIdvz7Z7gG53lqw, or linktree https://linktr.ee/soypete_tech."

func (c Client) callLLM(ctx context.Context, injection []string) (string, error) {
	messageHistory := []llms.MessageContent{llms.TextParts(llms.ChatMessageTypeSystem, pedroPrompt),
		llms.TextParts(llms.ChatMessageTypeHuman, strings.Join(injection, " "))}

	resp, err := c.llm.GenerateContent(ctx, messageHistory,
		llms.WithCandidateCount(1),
		llms.WithMaxLength(100),
		llms.WithTemperature(0.7),
		llms.WithPresencePenalty(1.0), // 2 is the largest penalty for using a work that has already been used
		llms.WithStopWords([]string{"twitch", "stream", "SoyPeteTech", "bot", "assistant", "silent", "software"}))
	if err != nil {
		return "", fmt.Errorf("failed to get llm response: %w", err)
	}
	err = c.db.InsertResponse(ctx, resp)
	if err != nil {
		return cleanResponse(resp.Choices[0].Content), fmt.Errorf("failed to write to db: %w", (err))
	}
	return cleanResponse(resp.Choices[0].Content), nil
}

// cleanResponse removes any newlines from the response
func cleanResponse(resp string) string {
	// remove any newlines
	resp = strings.ReplaceAll(resp, "\n", " ")
	resp = strings.ReplaceAll(resp, "<|im_start|>user", " ")
	return strings.TrimSpace(resp)
}

// SingleMessageResponse is a response from the LLM model to a single message, but to work it needs to have context of chat history
func (c Client) SingleMessageResponse(ctx context.Context, msg database.TwitchMessage) (string, error) {
	prompt, err := c.callLLM(ctx, []string{fmt.Sprintf("%s: %s", msg.Username, msg.Text)})
	if err != nil {
		return "", err
	}
	return prompt, nil
}

// GenerateTimer is a response from the LLM model from the list of helpful links and reminders
func (c Client) GenerateTimer(ctx context.Context, jsonbody string) (string, error) {
	prompt, err := c.callLLM(ctx, []string{fmt.Sprintf("here are some helpful links and reminders in a json format that you should use to help when generating your message. Pick one and them prompt chat to take action on it." + jsonbody)})
	if err != nil {
		return "", err
	}
	return prompt, nil
}
