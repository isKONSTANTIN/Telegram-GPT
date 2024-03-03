package telegram

import (
	"gopkg.in/telebot.v3"
)

func (b *GPTBot) whitelistMiddleware(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		userId := c.Sender().ID

		if userId == b.configs.MainAdminId {
			return next(c)
		}

		switch c.Chat().Type {
		case telebot.ChatPrivate:
			userId := c.Sender().ID

			u, _ := b.usersWhitelistRepo.GetUser(userId)

			if u != nil && u.State {
				return next(c)
			}
		case telebot.ChatGroup, telebot.ChatSuperGroup:
			chatId := c.Chat().ID

			chatRecord, _ := b.chatsWhitelistRepo.GetChat(chatId)

			if chatRecord != nil && chatRecord.State {
				return next(c)
			}
		}

		_ = c.Send("Not in whitelist")

		return nil
	}
}
