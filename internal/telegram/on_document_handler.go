package telegram

import (
	"TelegramGPT/internal/database"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sashabaranov/go-openai"
	"gopkg.in/telebot.v3"
	"io"
	"strings"
	"time"
	"unicode/utf8"
)

func (b *GPTBot) onDocument(c telebot.Context) error {
	doc := c.Message().Document
	if doc.FileSize > 256*1024 { // 256 KiB
		return c.Reply("File is too big")
	}

	docText, err := b.getTextFromDocument(doc)

	if err != nil {
		return err
	}

	if !utf8.ValidString(docText) {
		return c.Reply("File is invalid")
	}

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

			if err != nil {
				return err
			}

			needAnswer = true
		}
	}

	err = b.messagesRepo.AddMessage(text+" (User sent document: "+docText+")", openai.ChatMessageRoleUser, int64(c.Message().ID), chatId, *contextUUID)

	if err != nil {
		return err
	}

	if len(text) > 0 {
		needAnswer = true
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

func (b *GPTBot) getTextFromDocument(doc *telebot.Document) (string, error) {
	if doc == nil {
		return "", errors.New("document is nil")
	}

	resp, err := b.bot.File(&doc.File)
	if err != nil {
		return "", err
	}

	text, err := readCloserToString(resp)

	if err != nil {
		return "", err
	}

	return text, nil
}

func readCloserToString(rc io.ReadCloser) (string, error) {
	builder := new(strings.Builder)
	_, err := io.Copy(builder, rc)
	if err != nil {
		return "", err
	}

	return builder.String(), nil
}
