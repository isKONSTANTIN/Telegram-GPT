package app

import (
	"TelegramGPT/internal/config"
	"TelegramGPT/internal/database"
	"TelegramGPT/internal/gpt"
	"TelegramGPT/internal/telegram"
	"fmt"
	"github.com/jmoiron/sqlx"
	"os"
)

type App struct {
	configs *config.Config
	db      *sqlx.DB

	messagesRepo       *database.MessagesRepo
	usersWhitelistRepo *database.UsersWhitelistRepo
	chatsWhitelistRepo *database.ChatsWhitelistRepo
	aiPresetsRepo      *database.AiPresetsRepo

	bot       *telegram.GPTBot
	generator *gpt.Generator
}

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func (a *App) Start() error {
	var err error

	fmt.Println("Loading configs...")

	a.configs, err = config.LoadConfig("./configs/main.json")

	if err != nil {
		return err
	}

	fmt.Println("Connect to database...")

	a.db, err = database.Connect(
		getEnv("DB_USER", "tgpt"),
		getEnv("DB_PASSWORD", "change_me"),
		getEnv("DB_HOST", "postgres:5432"),
		getEnv("DB_TABLE", "tgpt"),
	)

	if err != nil {
		return err
	}

	fmt.Println("Init database repositories...")

	a.messagesRepo = database.NewMessagesRepo(a.db)
	a.usersWhitelistRepo = database.NewUsersWhitelistRepo(a.db)
	a.chatsWhitelistRepo = database.NewChatsWhitelistRepo(a.db)
	a.aiPresetsRepo = database.NewAiPresetsRepo(a.db)

	a.generator = gpt.NewGenerator(a.messagesRepo, &a.configs.OpenAI)

	fmt.Println("Init bot...")

	a.bot, err = telegram.InitBot(&a.configs.Telegram, a.messagesRepo, a.usersWhitelistRepo, a.chatsWhitelistRepo, a.aiPresetsRepo, a.generator)

	if err != nil {
		return err
	}

	fmt.Println("All done. Listening events...")

	a.bot.Start()

	return nil
}
