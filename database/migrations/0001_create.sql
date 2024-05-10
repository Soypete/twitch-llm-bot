-- +goose Up
CREATE TABLE twitch_chat (
		id serial PRIMARY KEY,
		username text,
		message text,
		isCommand BOOLEAN,
		created_at timestamptz DEFAULT NOW()
		);

-- +goose Down
DROP TABLE IF EXISTS twitch_chat;

