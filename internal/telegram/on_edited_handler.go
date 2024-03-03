package telegram

import (
	"TelegramGPT/internal/database"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sashabaranov/go-openai"
	"gopkg.in/telebot.v3"
	"slices"
)

func (b *GPTBot) onEdit(c telebot.Context) error {
	var contextUUID *pgtype.UUID
	messageId := int64(c.Message().ID)
	chatId := c.Message().Chat.ID

	message, err := b.messagesRepo.GetMessage(messageId, chatId)

	if err != nil {
		return nil
	}

	contextUUID = &message.ContextUUID

	preset, newUserText := b.cutCustomPresetIfExist(b.removeUsernameFromMessage(c.Text()), chatId)

	err = b.messagesRepo.EditMessage(messageId, chatId, openai.ChatMessageRoleUser, newUserText)

	if err != nil {
		return err
	}

	if preset != nil {
		_ = b.messagesRepo.EditMessage(messageId, chatId, openai.ChatMessageRoleSystem, preset.Text)
	}

	allMessages, err := b.messagesRepo.GetMessages(*contextUUID)

	if err != nil {
		return err
	}

	allMessagesLastIndex := len(allMessages) - 1

	messageContextIndex := slices.IndexFunc(allMessages, func(m database.Message) bool {
		return m.MessageId == messageId && m.Role == openai.ChatMessageRoleUser
	})

	if messageContextIndex != allMessagesLastIndex-1 {
		return nil
	}

	_ = c.Notify(telebot.Typing)

	newText, err := b.generator.Continue(allMessages)

	if err != nil {
		return err
	}

	editedMessage, err := b.bot.Edit(allMessages[allMessagesLastIndex], newText)

	if err != nil {
		return err
	}

	err = b.messagesRepo.EditMessage(int64(editedMessage.ID), editedMessage.Chat.ID, openai.ChatMessageRoleAssistant, newText)

	return err
}
