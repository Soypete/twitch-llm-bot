package database

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	v2 "github.com/gempir/go-twitch-irc/v2"
)

func (p Postgres) InsertMessage(ctx context.Context, msg v2.PrivateMessage) error {
	var isCommand bool
	username := msg.User.DisplayName
	message := msg.Message
	if strings.HasPrefix(msg.Message, "!") {
		isCommand = true
	}
	if strings.Contains(msg.User.DisplayName, "RestreamBot") {
		words := strings.Split(strings.Trim(msg.Message, "[]"), " ")
		username = words[1]                    // sets username to the first word after the video source.
		message = strings.Join(words[2:], " ") // create a clean message without the video source.
	}
	query := "INSERT INTO twitch_chat (username, message, isCommand, created_at) VALUES ($1, $2, $3, $4)"
	_, err := p.connections.ExecContext(ctx, query, username, message, isCommand, msg.Time)
	if err != nil {
		log.Println("error inserting message: ", err)
		return fmt.Errorf("error inserting message: %w", err)
	}
	return nil
}

func (p Postgres) InsertChatHistory(ctx context.Context, messages []string) error {
	query := "INSERT INTO twitch_chat (chats, created_at) VALUES ($1, $2)"
	_, err := p.connections.ExecContext(ctx, query, messages, time.Now())
	if err != nil {
		return fmt.Errorf("error inserting chat history: %w", err)
	}
	return nil
}

type TwitchMessage struct {
	Username string
	Message  string
}

func (p Postgres) QueryMessageHistory(interval time.Duration) ([]TwitchMessage, error) {
	var messages []TwitchMessage
	date := time.Now().Add(-interval)
	rows, err := p.connections.Query("SELECT username, message FROM twitch_chat WHERE isCommand = false and created_at > $1 ", date)
	if err != nil {
		return nil, fmt.Errorf("error querying message history: %w", err)
	}
	for rows.Next() {
		var message TwitchMessage
		err := rows.Scan(&message.Username, &message.Message)
		if err != nil {
			return nil, fmt.Errorf("error scanning message: %w", err)
		}
		messages = append(messages, message)
	}
	return messages, nil
}
