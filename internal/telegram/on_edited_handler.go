package telegram

import (
	"TelegramGPT/internal/database"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sashabaranov/go-openai"
	"gopkg.in/telebot.v3"
	"slices"
	"time"
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

	typingStatus := true

	go func() {
		for typingStatus {
			err := c.Notify(telebot.Typing)
			if err != nil {
				fmt.Println(err.Error())
			}
			time.Sleep(3 * time.Second)
		}
	}()

	newText, err := b.generator.Continue(allMessages)
	typingStatus = false

	if err != nil {
		return err
	}

	lastMessageIsFile := len(allMessages[allMessagesLastIndex].Text) > 4096
	var editedMessage *telebot.Message

	var content interface{}

	if lastMessageIsFile {
		content = prepareDocument(newText)
	} else {
		content = newText
	}

	editedMessage, err = b.bot.Edit(allMessages[allMessagesLastIndex], content, telebot.ModeMarkdown)

	if err != nil {
		return err
	}

	err = b.messagesRepo.EditMessage(int64(editedMessage.ID), editedMessage.Chat.ID, openai.ChatMessageRoleAssistant, newText)

	return err
}
