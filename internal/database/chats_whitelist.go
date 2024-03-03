package database

type ChatsWhitelistRepo struct {
	executor Executor
}

type Chat struct {
	ChatId      int64  `db:"chat_id"`
	State       bool   `db:"state"`
	Description string `db:"description"`
}

func NewChatsWhitelistRepo(executor Executor) *ChatsWhitelistRepo {
	return &ChatsWhitelistRepo{executor: executor}
}

func (r *ChatsWhitelistRepo) AddChat(chatId int64, description string) error {
	_, err := r.executor.Exec(
		"INSERT INTO chats_whitelist (chat_id, state, description) VALUES ($1, $2, $3)",
		chatId, true, description)

	return err
}

func (r *ChatsWhitelistRepo) RemoveChat(chatId int64) error {
	_, err := r.executor.Exec(
		"DELETE FROM chats_whitelist where chat_id = $1",
		chatId)

	return err
}

func (r *ChatsWhitelistRepo) SetChatState(chatId int64, enabled bool) error {
	_, err := r.executor.Exec(
		"UPDATE chats_whitelist SET state = $2 WHERE chat_id = $1",
		chatId, enabled)

	return err
}

func (r *ChatsWhitelistRepo) SetChatDescription(chatId int64, description string) error {
	_, err := r.executor.Exec(
		"UPDATE chats_whitelist SET description = $2 WHERE chat_id = $1",
		chatId, description)

	return err
}

func (r *ChatsWhitelistRepo) GetChat(chatId int64) (*Chat, error) {
	var chat Chat

	err := r.executor.Get(&chat,
		"SELECT * FROM chats_whitelist WHERE chat_id = $1",
		chatId)

	if err != nil {
		return nil, err
	}

	return &chat, nil
}

func (r *ChatsWhitelistRepo) GetList() ([]Chat, error) {
	var chats []Chat

	err := r.executor.Select(&chats,
		"SELECT * FROM chats_whitelist ORDER BY chat_id")

	if err != nil {
		return nil, err
	}

	return chats, nil
}
