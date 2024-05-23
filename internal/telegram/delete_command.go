package telegram

import (
	"github.com/jackc/pgx/v5/pgtype"
	"gopkg.in/telebot.v3"
)

func (b *GPTBot) deleteCommand(c telebot.Context) error {
	replyTo := c.Message().ReplyTo
	chatId := c.Message().Chat.ID

	if replyTo == nil {
		_ = c.Reply("Reply to any message to remove all context with it")

		return nil
	}

	dbMessage, _ := b.messagesRepo.GetMessage(int64(replyTo.ID), chatId)

	if dbMessage == nil {
		_ = c.Reply("Context not found")

		return nil
	}

	var contextUUID *pgtype.UUID
	contextUUID = &dbMessage.ContextUUID

	messages, err := b.messagesRepo.GetMessages(*contextUUID)

	if err != nil || len(messages) == 0 {
		_ = c.Reply("Messages not found")

		return nil
	}

	err = b.messagesRepo.RemoveMessages(*contextUUID)

	if err != nil {
		_ = c.Reply("Error to remove messages")

		return nil
	}

	for _, m := range messages {
		_ = b.bot.Delete(m)
	}

	_ = b.bot.Delete(c.Message())

	return nil
}
