package config

import (
	"encoding/json"
	"github.com/sashabaranov/go-openai"
)

type Config struct {
	Telegram TelegramConfig `json:"telegram"`
	OpenAI   OpenAIConfig   `json:"openAI"`
}

type TelegramConfig struct {
	TelegramToken string `json:"telegramToken"`
	MainAdminId   int64  `json:"mainAdminId"`
}

type OpenAIConfig struct {
	Token         string `json:"token"`
	Model         string `json:"model"`
	DefaultPreset string `json:"defaultPreset"`
}

func CreateDefault() *Config {
	return &Config{
		Telegram: TelegramConfig{
			TelegramToken: "XXX",
			MainAdminId:   -1,
		},
		OpenAI: OpenAIConfig{
			Token:         "YYY",
			Model:         openai.GPT4Turbo1106,
			DefaultPreset: "You are the best artificial intelligence that helps a person answer his questions",
		},
	}
}

func (c *Config) String() string {
	marshal, err := json.Marshal(&c)
	if err != nil {
		return "{}"
	}

	return string(marshal)
}
