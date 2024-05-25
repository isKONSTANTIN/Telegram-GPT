package database

import (
	"github.com/jackc/pgx/v5/pgtype"
	"strconv"
)

type MessagesRepo struct {
	executor Executor
}

type Message struct {
	Id          int64       `db:"id"`
	Text        string      `db:"text"`
	Role        string      `db:"role"`
	MessageId   int64       `db:"message_id"`
	ChatId      int64       `db:"chat_id"`
	ContextUUID pgtype.UUID `db:"context_uuid"`
	Type        string      `db:"type"`
}

func (m Message) MessageSig() (string, int64) {
	return strconv.FormatInt(m.MessageId, 10), m.ChatId
}

func NewMessagesRepo(executor Executor) *MessagesRepo {
	return &MessagesRepo{executor: executor}
}

func (r *MessagesRepo) AddMessage(text string, role string, messageId int64, chatId int64, context pgtype.UUID) error {
	_, err := r.executor.Exec(
		"INSERT INTO messages (text, role, message_id, chat_id, context_uuid) VALUES ($1, $2, $3, $4, $5)",
		text, role, messageId, chatId, context)

	return err
}

func (r *MessagesRepo) AddImageURL(url string, role string, messageId int64, chatId int64, context pgtype.UUID) error {
	_, err := r.executor.Exec(
		"INSERT INTO messages (text, role, message_id, chat_id, context_uuid, type) VALUES ($1, $2, $3, $4, $5, 'image')",
		url, role, messageId, chatId, context)

	return err
}

func (r *MessagesRepo) EditMessage(messageId int64, chatId int64, role string, text string) error {
	_, err := r.executor.Exec(
		"UPDATE messages SET text = $4 WHERE message_id = $1 AND chat_id = $2 AND role = $3 and type = 'text'",
		messageId, chatId, role, text)

	return err
}

func (r *MessagesRepo) GetMessage(messageId int64, chatId int64) (*Message, error) {
	var message Message

	err := r.executor.Get(&message,
		"SELECT * FROM messages WHERE message_id = $1 and chat_id = $2",
		messageId, chatId)

	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (r *MessagesRepo) GetMessages(contextUUID pgtype.UUID) ([]Message, error) {
	var messages []Message

	err := r.executor.Select(&messages,
		"SELECT * FROM messages WHERE context_uuid = $1 ORDER BY id",
		contextUUID)

	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *MessagesRepo) RemoveMessages(contextUUID pgtype.UUID) error {
	_, err := r.executor.Exec(
		"DELETE FROM messages WHERE context_uuid = $1",
		contextUUID)

	return err
}
