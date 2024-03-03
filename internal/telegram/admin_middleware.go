package telegram

import "gopkg.in/telebot.v3"

func (b *GPTBot) adminMiddleware(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		userId := c.Sender().ID

		if userId == b.configs.MainAdminId {
			return next(c)
		}

		return nil
	}
}
