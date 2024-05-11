-- +goose Up

CREATE TABLE IF NOT EXISTS bot_response (
		id serial PRIMARY KEY,
		model_name text,
		response text,
		stop_reason text,
		was_successful BOOLEAN, 
		created_at timestamptz DEFAULT NOW()
		);

-- +goose Down
DROP TABLE IF EXISTS bot_response;
