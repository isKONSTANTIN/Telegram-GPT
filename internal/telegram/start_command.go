package telegram

import "gopkg.in/telebot.v3"

func (b *GPTBot) startCommand(c telebot.Context) error {
	err := c.Send("Welcome to the Telegram bot designed to facilitate effective interaction with the artificial intelligence of ChatGPT. " +
		"This bot is created to respond to your inquiries and fulfill requests in a conversational mode.\n\n" +
		"Should you find any errors in your last message or wish to amend it, rest assured. " +
		"The bot will automatically recognize the edit and prepare an appropriate updated response.\n\n" +
		"🏷️ Additionally, for your convenience, the bot incorporates a system of hashtags that enable you to activate various dialogue modes, such as receiving brief responses.\n\n" +
		"🔍 You can find more detailed information about commands and hashtag usage in the /help section.\n\n" +
		"Start a dialogue with the bot, and you'll quickly appreciate all its capabilities.")

	if err != nil {
		return err
	}

	return nil
}
