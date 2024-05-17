package twitchirc

import (
	"reflect"
	"testing"

	database "github.com/Soypete/twitch-llm-bot/database"
	v2 "github.com/gempir/go-twitch-irc/v2"
)

func Test_cleanMessage(t *testing.T) {
	tests := []struct {
		name string
		msg  v2.PrivateMessage
		want database.TwitchMessage
	}{
		{
			name: "Restream+Youtube",
			msg: v2.PrivateMessage{
				User: v2.User{
					DisplayName: "[RestreamBot]",
				},
				Message: "[YouTube: IMJONEZZ] Yeah, conversation history, is that what you're wondering about?",
			},
			want: database.TwitchMessage{
				Username: "IMJONEZZ",
				Text:     "Yeah, conversation history, is that what you're wondering about?",
			},
		},
		{
			name: "Restream+Youtube",
			msg: v2.PrivateMessage{
				User: v2.User{
					DisplayName: "[RestreamBot]",
				},
				Message: "[YouTube: MD Habib] Also i have to make a chatbot",
			},
			want: database.TwitchMessage{
				Username: "MD Habib",
				Text:     "Also i have to make a chatbot",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cleanMessage(tt.msg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cleanMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
