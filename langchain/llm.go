package langchain

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	database "github.com/Soypete/twitch-llm-bot/database"
	"github.com/tmc/langchaingo/llms"
)

const pedroPrompt = "Your name is Pedro_el_asistente. You are a chat bot that helps out in SoyPeteTech's twitch chat. You are allowed to use links, code, or emotes to express fun messages about software. Helpful links are always appreciated, such as SoyPeteTech's github https://github.com/Soypete, youtube https://www.youtube.com/channel/UCEkM7JXVQIdvz7Z7gG53lqw, or linktree https://linktr.ee/soypete_tech."

func (c Client) PromptWithoutChat(ctx context.Context) (string, error) {
	content, err := llms.GenerateFromSinglePrompt(ctx,
		c.llm,
		pedroPrompt,
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
	messageHistory := []llms.MessageContent{llms.TextParts(llms.ChatMessageTypeSystem, pedroPrompt),
		llms.TextParts(llms.ChatMessageTypeSystem, "Here is the twitch chat history for you to respond to:")}
	// get message history from database
	messages, err := c.db.QueryMessageHistory(interval)
	if err != nil {
		log.Fatal(err)
	}
	if len(messages) == 0 {
		return nil, fmt.Errorf("no messages found")
	}
	for _, message := range messages {
		// Experiment using just the text and no username
		prompt := message.Text
		// prompt := fmt.Sprintf("%s: %s", message.Username, message.Text)
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
		llms.WithPresencePenalty(1.0), // 2 is the largest penalty for using a work that has already been used
		llms.WithStopWords([]string{"twitch", "stream", "SoyPeteTech", "bot", "assistant", "silent", "software"}))
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
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
	messageHistory := []llms.MessageContent{llms.TextParts(llms.ChatMessageTypeSystem, pedroPrompt),
		llms.TextParts(llms.ChatMessageTypeSystem, "is the message for you to respond to: "+msg.Text)}
	resp, err := c.llm.GenerateContent(ctx, messageHistory,
		llms.WithCandidateCount(1),
		llms.WithMaxLength(500),
		llms.WithTemperature(0.7),
		llms.WithPresencePenalty(1.0), // 2 is the largest penalty for using a work that has already been used
		llms.WithStopWords([]string{"twitch", "stream", "SoyPeteTech", "bot", "assistant", "silent", "software"}))
	if err != nil {
		return "", fmt.Errorf("failed to get llm response: %w", err)
	}
	return cleanResponse(resp.Choices[0].Content), nil
}

// GenerateTimer is a response from the LLM model from the list of helpful links and reminders
func (c Client) GenerateTimer(ctx context.Context, jsonbody string) (string, error) {
	messageHistory := []llms.MessageContent{llms.TextParts(llms.ChatMessageTypeSystem, pedroPrompt),
		llms.TextParts(llms.ChatMessageTypeSystem, "here are some helpful links and reminders in a json format that you should use to help when generating your message. Pick one and them prompt chat to take action on it."+jsonbody)}
	resp, err := c.llm.GenerateContent(ctx, messageHistory,
		llms.WithCandidateCount(1),
		llms.WithMaxLength(500),
		llms.WithTemperature(0.7),
		llms.WithPresencePenalty(1.0), // 2 is the largest penalty for using a work that has already been used
		llms.WithStopWords([]string{"twitch", "stream", "SoyPeteTech", "bot", "assistant", "silent", "software"}))
	if err != nil {
		return "", fmt.Errorf("failed to get llm response: %w", err)
	}
	return cleanResponse(resp.Choices[0].Content), nil
}
