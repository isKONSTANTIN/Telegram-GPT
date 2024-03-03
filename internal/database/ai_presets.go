package database

type AiPresetsRepo struct {
	executor Executor
}

type Preset struct {
	Id     int64  `db:"id"`
	ChatId int64  `db:"chat_id"`
	Text   string `db:"text"`
	Tag    string `db:"tag"`
}

func NewAiPresetsRepo(executor Executor) *AiPresetsRepo {
	return &AiPresetsRepo{executor: executor}
}

func (r *AiPresetsRepo) AddPreset(chatId int64, text string, tag string) error {
	_, err := r.executor.Exec(
		"INSERT INTO ai_presets (chat_id, text, tag) VALUES ($1, $2, $3)",
		chatId, text, tag)

	return err
}

func (r *AiPresetsRepo) RemovePreset(id int64) error {
	_, err := r.executor.Exec("DELETE FROM ai_presets where id = $1", id)

	return err
}

func (r *AiPresetsRepo) EditPreset(id int64, text string, tag string) error {
	_, err := r.executor.Exec(
		"UPDATE ai_presets SET text = $2, tag = $3 WHERE id = $1",
		id, text, tag)

	return err
}

func (r *AiPresetsRepo) GetPresetById(id int64) (*Preset, error) {
	var result Preset

	err := r.executor.Get(&result,
		"SELECT * FROM ai_presets WHERE id = $1",
		id)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *AiPresetsRepo) GetPresetByTag(tag string, chatId int64) (*Preset, error) {
	var result Preset

	err := r.executor.Get(&result,
		"SELECT * FROM ai_presets WHERE tag = $1 AND chat_id = $2",
		tag, chatId)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *AiPresetsRepo) GetChatPresets(chatId int64) ([]Preset, error) {
	var presets []Preset

	err := r.executor.Select(&presets,
		"SELECT * FROM ai_presets WHERE chat_id = $1 ORDER BY chat_id", chatId)

	if err != nil {
		return nil, err
	}

	return presets, nil
}
