package telegram

import "gopkg.in/telebot.v3"

func (b *GPTBot) autoResponseMiddleware(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		if c.Callback() != nil {
			defer c.Respond()
		}
		return next(c)
	}
}
