package telegram

import (
	"fmt"
	"gopkg.in/telebot.v3"
)

func (b *GPTBot) chatInfo(c telebot.Context) error {
	err := c.Send(fmt.Sprintf("Chat id: %d\nSender id: %d", c.Chat().ID, c.Sender().ID))

	if err != nil {
		return err
	}

	return nil
}
