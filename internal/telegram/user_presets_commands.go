package telegram

import (
	"gopkg.in/telebot.v3"
	"strings"
)

func sendAddPresetUsage(c telebot.Context) error {
	return c.Send("/add_preset <tag> <preset>")
}

func sendRemovePresetUsage(c telebot.Context) error {
	return c.Send("/del_preset <tag>")
}

func sendEditPresetUsage(c telebot.Context) error {
	return c.Send("/edit_preset <tag> <new preset text>")
}

func (b *GPTBot) addPresetCommand(c telebot.Context) error {
	args := c.Args()

	if len(args) < 2 {
		return sendAddPresetUsage(c)
	}

	preset, _ := b.aiPresetsRepo.GetPresetByTag(args[0], c.Chat().ID)

	if preset != nil {
		return c.Send("Preset with that tag already exists")
	}

	err := b.aiPresetsRepo.AddPreset(c.Chat().ID, strings.Join(args[1:], " "), args[0])

	if err != nil {
		_ = c.Send("Fail to add preset")

		return err
	}

	return c.Send("Preset added")
}

func (b *GPTBot) deletePresetCommand(c telebot.Context) error {
	args := c.Args()

	if len(args) != 1 {
		return sendRemovePresetUsage(c)
	}

	preset, err := b.aiPresetsRepo.GetPresetByTag(args[0], c.Chat().ID)

	if err != nil {
		return c.Send("Preset not found")
	}

	err = b.aiPresetsRepo.RemovePreset(preset.Id)

	if err != nil {
		_ = c.Send("Fail to remove preset")

		return err
	}

	return c.Send("Preset removed")
}

func (b *GPTBot) editPresetCommand(c telebot.Context) error {
	args := c.Args()

	if len(args) < 2 {
		return sendEditPresetUsage(c)
	}

	preset, err := b.aiPresetsRepo.GetPresetByTag(args[0], c.Chat().ID)

	if err != nil {
		return c.Send("Preset not found")
	}

	err = b.aiPresetsRepo.EditPreset(preset.Id, strings.Join(args[1:], " "), preset.Tag)

	if err != nil {
		_ = c.Send("Fail to edit preset")

		return err
	}

	return c.Send("Preset edited")
}

func (b *GPTBot) listPresetCommand(c telebot.Context) error {
	presets, err := b.aiPresetsRepo.GetChatPresets(c.Chat().ID)

	if err != nil {
		_ = c.Send("Fail to fetch presets")

		return err
	}

	if len(presets) == 0 {
		return c.Send("List is empty")
	}

	var result string

	for _, preset := range presets {
		result += "#" + preset.Tag + ": " + preset.Text + "\n"
	}

	return c.Send(result)
}
