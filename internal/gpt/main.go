package gpt

import (
	"TelegramGPT/internal/config"
	"TelegramGPT/internal/database"
	"context"
	"github.com/sashabaranov/go-openai"
)

type Generator struct {
	messagesRepo *database.MessagesRepo
	configs      *config.OpenAIConfig
	client       *openai.Client
}

func NewGenerator(messagesRepo *database.MessagesRepo, configs *config.OpenAIConfig) *Generator {
	client := openai.NewClient(configs.Token)

	return &Generator{
		messagesRepo: messagesRepo,
		configs:      configs,
		client:       client,
	}
}

func (g *Generator) DefaultPresetText() string {
	return g.configs.DefaultPreset
}

func (g *Generator) Continue(messages []database.Message) (string, error) {
	mappedMessages := mapMessages(messages)

	resp, err := g.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    g.configs.Model,
			Messages: mappedMessages,
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func mapMessages(messages []database.Message) []openai.ChatCompletionMessage {
	result := make([]openai.ChatCompletionMessage, len(messages))

	for i, message := range messages {
		switch message.Type {
		case "text":
			{
				result[i] = openai.ChatCompletionMessage{
					Role:    message.Role,
					Content: message.Text,
				}
			}
		case "image":
			{
				result[i] = openai.ChatCompletionMessage{
					Role: message.Role,
					MultiContent: []openai.ChatMessagePart{
						{
							Type: "image_url",
							Text: "",
							ImageURL: &openai.ChatMessageImageURL{
								URL:    message.Text,
								Detail: "",
							},
						},
					},
				}
			}

		}
	}

	return result
}

func (g *Generator) CreateImage(prompt string, horizontal bool) (string, error) {
	var resolutionChoice = openai.CreateImageSize1024x1792
	if horizontal {
		resolutionChoice = openai.CreateImageSize1792x1024
	}

	respUrl, err := g.client.CreateImage(context.Background(),
		openai.ImageRequest{
			Model:          openai.CreateImageModelDallE3,
			Prompt:         prompt,
			Size:           resolutionChoice,
			ResponseFormat: openai.CreateImageResponseFormatURL,
			N:              1,
		})

	if err != nil {
		return "", err
	}

	return respUrl.Data[0].URL, nil
}
