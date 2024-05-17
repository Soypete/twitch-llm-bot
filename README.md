# go-template-repo

this is a template for starting go projects

[![Actions Status](https://github.com/soypete/{}/workflows/build/badge.svg)](https://github.com/soypete/{}/actions/workflows/go.yml)
[![wakatime](https://wakatime.com/badge/user/953eeb5a-d347-44af-9d8b-a5b8a918cecf/project/018ef728-5089-4148-b326-592f7a744f7e.svg)](https://wakatime.com/badge/user/953eeb5a-d347-44af-9d8b-a5b8a918cecf/project/018ef728-5089-4148-b326-592f7a744f7e)

## To Use

install [lama.cpp]() and run there server on `127.0.0.1` and port `8080`

```bash
source .secrets
go run main.go
```

your secrets should contain

`bash export OPENAI_API_KEY export TWITCH_CLIENT_ID export TWITCH_CLIENT_SECRET`

## Notes:

So far the longest that the bot has taken to respond is 5 minutes, so we need to account for that in the tmeout the api call.

## TODO

* change bot name
* git bot moderator permissions
* add more tokens to llm in llama cpp
* batch twitch chat to set via the langchain [GenerateContent](https://github.com/tmc/langchaingo/blob/3a36972919a83b119825de4ea6216e175ae20cb3/examples/openai-chat-example/openai_chat_example.go#L25C19-L25C34)
* Add embeddings
* connect Supabase for managing vector embeddings
* add config for managing the bot [channel commands, prompts, links, stream title etc]
* move to windows machine
