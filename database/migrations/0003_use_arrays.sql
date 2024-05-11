-- +goose Up
CREATE TABLE IF NOT EXISTS chat_prompts (
		id uuid PRIMARY KEY,
		chats text[], -- array of twitch/youtube chat messages
		start_time timestamptz, -- time twitch chat prompt started
		end_time timestamptz, -- time twitch chat prompt ended
		created_at timestamptz DEFAULT NOW()
		);

CREATE TABLE IF NOT EXISTS bot_responses (
		id uuid PRIMARY KEY,
		model_name text,
		responses text[], -- array of messages produced by the model
		stop_reason text[],
		twith_chat_prompt_id int references twitch_chat_promts (id),		
		created_at timestamptz DEFAULT NOW()
		);

-- +goose Down
DROP TABLE IF NOT EXISTS chat_prompts;
DROP TABLE IF NOT EXISTS bot_responses;
