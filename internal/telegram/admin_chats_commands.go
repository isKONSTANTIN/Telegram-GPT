package telegram

import (
	"fmt"
	"gopkg.in/telebot.v3"
	"strconv"
	"strings"
)

func sendAddChatUsage(c telebot.Context) error {
	return c.Send("Usage: /add_chat <id> <description>")
}

func sendDeleteChatUsage(c telebot.Context) error {
	return c.Send("Usage: /del_chat <id>")
}

func sendSwitchChatUsage(c telebot.Context) error {
	return c.Send("Usage: /switch_chat <id>")
}

func (b *GPTBot) addChatCommand(c telebot.Context) error {
	args := c.Args()

	if len(args) < 2 {
		return sendAddChatUsage(c)
	}

	var id int64
	var err error

	if id, err = strconv.ParseInt(args[0], 10, 64); err != nil {
		return sendAddChatUsage(c)
	}

	description := strings.Join(args[1:], " ")

	err = b.chatsWhitelistRepo.AddChat(id, description)

	if err != nil {
		fmt.Println(err)

		return c.Send("Fail to add chat")
	}

	return c.Send("Chat added")
}

func (b *GPTBot) deleteChatCommand(c telebot.Context) error {
	args := c.Args()

	if len(args) != 1 {
		return sendDeleteChatUsage(c)
	}

	var id int64
	var err error

	if id, err = strconv.ParseInt(args[0], 10, 64); err != nil {
		return sendDeleteChatUsage(c)
	}

	err = b.chatsWhitelistRepo.RemoveChat(id)

	if err != nil {
		fmt.Println(err)

		return c.Send("Fail to delete chat")
	}

	return c.Send("Chat deleted")
}

func (b *GPTBot) listChatCommand(c telebot.Context) error {
	chats, err := b.chatsWhitelistRepo.GetList()

	if err != nil {
		fmt.Println(err)

		return c.Send("Fail to get chat list")
	}

	var result string
	var state string

	if len(chats) == 0 {
		return c.Send("List is empty")
	}

	for _, chat := range chats {
		if chat.State {
			state = ""
		} else {
			state = " (disabled)"
		}

		result += strconv.FormatInt(chat.ChatId, 10) + ": " + chat.Description + state + "\n"
	}

	return c.Send(result)
}

func (b *GPTBot) switchChatCommand(c telebot.Context) error {
	args := c.Args()

	if len(args) != 1 {
		return sendSwitchChatUsage(c)
	}

	var id int64
	var err error

	if id, err = strconv.ParseInt(args[0], 10, 64); err != nil {
		return sendSwitchChatUsage(c)
	}

	chat, err := b.chatsWhitelistRepo.GetChat(id)

	if err != nil {
		return c.Send("Chat not found")
	}

	err = b.usersWhitelistRepo.SetUserState(id, !chat.State)

	if err != nil {
		fmt.Println(err)

		return c.Send("Fail to switch chat")
	}

	return c.Send("Chat switched")
}
