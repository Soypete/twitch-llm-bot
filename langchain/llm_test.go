package langchain

import (
	"context"
	"reflect"
	"testing"
	"time"

	database "github.com/Soypete/twitch-llm-bot/database"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

func TestClient_PromptWithoutChat(t *testing.T) {
	type fields struct {
		llm *openai.LLM
		db  database.Postgres
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Client{
				llm: tt.fields.llm,
				db:  tt.fields.db,
			}
			got, err := c.PromptWithoutChat(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.PromptWithoutChat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Client.PromptWithoutChat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetMessageHistory(t *testing.T) {
	type fields struct {
		llm *openai.LLM
		db  database.Postgres
	}
	type args struct {
		interval time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []llms.MessageContent
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Client{
				llm: tt.fields.llm,
				db:  tt.fields.db,
			}
			got, err := c.GetMessageHistory(tt.args.interval)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetMessageHistory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetMessageHistory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_PromptWithChat(t *testing.T) {
	type fields struct {
		llm *openai.LLM
		db  database.Postgres
	}
	type args struct {
		ctx      context.Context
		interval time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Client{
				llm: tt.fields.llm,
				db:  tt.fields.db,
			}
			got, err := c.PromptWithChat(tt.args.ctx, tt.args.interval)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.PromptWithChat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Client.PromptWithChat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cleanResponse(t *testing.T) {
	tests := []struct {
		name string
		resp string
		want string
	}{
		{
			name: "Test 1",
			resp: "Hello\nWorld",
			want: "Hello World",
		},
		{
			name: "Test 2",
			resp: "<|im_start|>user \nTtocsNeb: hi",
			want: "TtocsNeb: hi",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cleanResponse(tt.resp); got != tt.want {
				t.Errorf("cleanResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}
