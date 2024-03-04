package telegram

import (
	"TelegramGPT/internal/config"
	"TelegramGPT/internal/database"
	"TelegramGPT/internal/gpt"
	"gopkg.in/telebot.v3"
	"time"
)

type GPTBot struct {
	settings  *telebot.Settings
	configs   *config.TelegramConfig
	bot       *telebot.Bot
	generator *gpt.Generator

	messagesRepo       *database.MessagesRepo
	usersWhitelistRepo *database.UsersWhitelistRepo
	chatsWhitelistRepo *database.ChatsWhitelistRepo
	aiPresetsRepo      *database.AiPresetsRepo
}

func InitBot(
	configs *config.TelegramConfig,
	messagesRepo *database.MessagesRepo,
	usersWhitelistRepo *database.UsersWhitelistRepo,
	chatsWhitelistRepo *database.ChatsWhitelistRepo,
	aiPresetsRepo *database.AiPresetsRepo,
	generator *gpt.Generator) (*GPTBot, error) {
	pref := telebot.Settings{
		Token:  configs.TelegramToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		return nil, err
	}

	gptBot := GPTBot{
		settings:           &pref,
		configs:            configs,
		bot:                bot,
		messagesRepo:       messagesRepo,
		usersWhitelistRepo: usersWhitelistRepo,
		chatsWhitelistRepo: chatsWhitelistRepo,
		aiPresetsRepo:      aiPresetsRepo,
		generator:          generator,
	}

	gptBot.init()

	return &gptBot, nil
}

func (b *GPTBot) init() {
	b.bot.Use(b.autoResponseMiddleware)

	b.bot.Handle("/chat_info", b.chatInfo)

	userGroup := b.bot.Group()

	userGroup.Use(b.whitelistMiddleware)
	userGroup.Handle(telebot.OnText, b.onText)
	userGroup.Handle(telebot.OnEdited, b.onEdit)

	userGroup.Handle("/start", b.startCommand)

	userGroup.Handle("/add_preset", b.addPresetCommand)
	userGroup.Handle("/del_preset", b.deletePresetCommand)
	userGroup.Handle("/edit_preset", b.editPresetCommand)
	userGroup.Handle("/list_presets", b.listPresetCommand)

	userGroup.Handle("/imagine", b.imagineCommand)

	userGroup.Handle("/help", func(c telebot.Context) error {
		return c.Send("/add_preset <tag> <preset> - add new preset\n" +
			"/del_preset <tag> - delete preset\n" +
			"/edit_preset <tag> <new preset text> - edit preset\n" +
			"/list_presets - presets list\n" +
			"/imagine <resolution 0-4> <prompt> - imagine image",
		)
	})

	adminGroup := b.bot.Group()
	adminGroup.Use(b.adminMiddleware)

	adminGroup.Handle("/add_user", b.addUserCommand)
	adminGroup.Handle("/del_user", b.deleteUserCommand)
	adminGroup.Handle("/list_users", b.listUserCommand)
	adminGroup.Handle("/switch_user", b.switchUserCommand)

	adminGroup.Handle("/add_chat", b.addChatCommand)
	adminGroup.Handle("/del_chat", b.deleteChatCommand)
	adminGroup.Handle("/list_chats", b.listChatCommand)
	adminGroup.Handle("/switch_chat", b.switchChatCommand)

	adminGroup.Handle("/help", func(c telebot.Context) error {
		return c.Send(
			"/add_user <id> <description> - add user to whitelist\n" +
				"/del_user <id> - delete user from whitelist\n" +
				"/list_users - show users in whitelist\n" +
				"/switch_user <id> - disable or enable access to bot\n" +
				"\n" +
				"/add_chat <id> <description> - add chat to whitelist\n" +
				"/del_chat <id> - delete chat from whitelist\n" +
				"/list_chats - show chats in whitelist\n" +
				"/switch_chat <id> - disable or enable access to bot\n",
		)
	})
}

func (b *GPTBot) Start() {
	b.bot.Start()
}
