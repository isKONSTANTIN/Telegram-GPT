package telegram

import (
	"TelegramGPT/internal/database"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sashabaranov/go-openai"
	"gopkg.in/telebot.v3"
	"strings"
	"time"
)

func (b *GPTBot) onPhoto(c telebot.Context) error {
	photo := c.Message().Photo
	text := c.Message().Caption

	startWithMe := strings.HasPrefix(text, "@"+b.bot.Me.Username)
	replyToMe := c.Message().IsReply() && c.Message().ReplyTo.Sender.ID == b.bot.Me.ID

	if c.Chat().Type != telebot.ChatPrivate {
		if !startWithMe && !replyToMe {
			return nil
		}

		text = b.removeUsernameFromMessage(text)
	} else if c.Message().IsForwarded() || c.Message().OriginalSenderName != "" {
		return nil
	}

	var contextUUID *pgtype.UUID
	replyTo := c.Message().ReplyTo
	chatId := c.Message().Chat.ID

	var preset *database.Preset
	var presetText string

	preset, text = b.cutCustomPresetIfExist(text, chatId)

	var needAnswer = false

	if replyTo != nil && preset == nil {
		message, _ := b.messagesRepo.GetMessage(int64(replyTo.ID), chatId)

		if message != nil {
			contextUUID = &message.ContextUUID
		}
	}

	if contextUUID == nil {
		contextUUID = &pgtype.UUID{
			Bytes: uuid.New(),
			Valid: true,
		}

		if preset != nil {
			presetText = preset.Text
		} else {
			presetText = b.generator.DefaultPresetText()
		}

		err := b.messagesRepo.AddMessage(presetText, openai.ChatMessageRoleSystem, int64(c.Message().ID), chatId, *contextUUID)

		if err != nil {
			return err
		}

		if replyTo != nil {
			err = b.messagesRepo.AddMessage("Forwarded message: "+replyTo.Text, openai.ChatMessageRoleUser, int64(replyTo.ID), chatId, *contextUUID)

			if replyTo.Photo != nil {
				f, _ := b.bot.FileByID(replyTo.Photo.FileID)

				err = b.messagesRepo.AddImageURL(
					"https://api.telegram.org/file/bot"+b.bot.Token+"/"+f.FilePath,
					openai.ChatMessageRoleUser,
					int64(replyTo.ID),
					chatId,
					*contextUUID)
			}

			if err != nil {
				return err
			}

			needAnswer = true
		}
	}

	if len(text) > 0 {
		err := b.messagesRepo.AddMessage(text, openai.ChatMessageRoleUser, int64(c.Message().ID), chatId, *contextUUID)

		if err != nil {
			return err
		}

		needAnswer = true
	}

	f, _ := b.bot.FileByID(photo.FileID)

	err := b.messagesRepo.AddImageURL(
		"https://api.telegram.org/file/bot"+b.bot.Token+"/"+f.FilePath,
		openai.ChatMessageRoleUser,
		int64(c.Message().ID),
		chatId,
		*contextUUID)

	if err != nil {
		return err
	}

	if !needAnswer {
		return nil
	}

	messages, err := b.messagesRepo.GetMessages(*contextUUID)
	if err != nil {
		return err
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

	answer, err := b.generator.Continue(messages)
	typingStatus = false

	if err != nil {
		return err
	}

	var content interface{}

	if len(answer) > 4096 {
		content = prepareDocument(answer)
	} else {
		content = answer
	}

	sentMessage, err := b.bot.Reply(c.Message(), content, telebot.ModeMarkdown)
	if err != nil {
		return err
	}

	err = b.messagesRepo.AddMessage(answer, openai.ChatMessageRoleAssistant, int64(sentMessage.ID), chatId, *contextUUID)
	if err != nil {
		return err
	}

	return nil
}
