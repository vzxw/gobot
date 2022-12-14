package config

import (
	"github.com/joho/godotenv"
	zeroLog "github.com/rs/zerolog/log"
)

type readerStruct struct {
	envs map[string]string
}

type List struct {
	SlackSigningSecret string
	TelegramAuthToken  string
	SlackEventsPath    string
	SlackEventsPort    uint64
}

func Read(envFile ...string) List {
	envs, err := godotenv.Read(envFile...)
	if err != nil {
		zeroLog.Fatal().Err(err)
	}

	reader := &readerStruct{
		envs: envs,
	}

	return List{
		SlackSigningSecret: reader.getString("SLACK_SIGNING_SECRET"),
		TelegramAuthToken:  reader.getString("TELEGRAM_BOT_AUTH_TOKEN"),
		SlackEventsPort:    3000,
		SlackEventsPath:    "/slack/events",
	}
}

func (r *readerStruct) getString(key string) string {
	result, ok := r.envs[key]
	if !ok {
		zeroLog.Fatal().Msgf("Undefined <%v> env variable", key)
	}

	return result
}
