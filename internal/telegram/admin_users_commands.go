package telegram

import (
	"fmt"
	"gopkg.in/telebot.v3"
	"strconv"
	"strings"
)

func sendAddUserUsage(c telebot.Context) error {
	return c.Send("Usage: /add_user <id> <description>")
}

func sendDeleteUserUsage(c telebot.Context) error {
	return c.Send("Usage: /del_user <id>")
}

func sendSwitchUserUsage(c telebot.Context) error {
	return c.Send("Usage: /switch_user <id>")
}

func (b *GPTBot) addUserCommand(c telebot.Context) error {
	args := c.Args()

	if len(args) < 2 {
		return sendAddUserUsage(c)
	}

	var id int64
	var err error

	if id, err = strconv.ParseInt(args[0], 10, 64); err != nil {
		return sendAddUserUsage(c)
	}

	description := strings.Join(args[1:], " ")

	err = b.usersWhitelistRepo.AddUser(id, description)

	if err != nil {
		fmt.Println(err)

		return c.Send("Fail to add user")
	}

	return c.Send("User added")
}

func (b *GPTBot) deleteUserCommand(c telebot.Context) error {
	args := c.Args()

	if len(args) != 1 {
		return sendDeleteUserUsage(c)
	}

	var id int64
	var err error

	if id, err = strconv.ParseInt(args[0], 10, 64); err != nil {
		return sendDeleteUserUsage(c)
	}

	err = b.usersWhitelistRepo.RemoveUser(id)

	if err != nil {
		fmt.Println(err)

		return c.Send("Fail to delete user")
	}

	return c.Send("User deleted")
}

func (b *GPTBot) listUserCommand(c telebot.Context) error {
	users, err := b.usersWhitelistRepo.GetList()

	if err != nil {
		fmt.Println(err)

		return c.Send("Fail to get user list")
	}

	var result string
	var state string

	if len(users) == 0 {
		return c.Send("List is empty")
	}

	for _, user := range users {
		if user.State {
			state = ""
		} else {
			state = " (disabled)"
		}

		result += strconv.FormatInt(user.UserId, 10) + ": " + user.Description + state + "\n"
	}

	return c.Send(result)
}

func (b *GPTBot) switchUserCommand(c telebot.Context) error {
	args := c.Args()

	if len(args) != 1 {
		return sendSwitchUserUsage(c)
	}

	var id int64
	var err error

	if id, err = strconv.ParseInt(args[0], 10, 64); err != nil {
		return sendSwitchUserUsage(c)
	}

	user, err := b.usersWhitelistRepo.GetUser(id)

	if err != nil {
		return c.Send("User not found")
	}

	err = b.usersWhitelistRepo.SetUserState(id, !user.State)

	if err != nil {
		fmt.Println(err)

		return c.Send("Fail to switch user")
	}

	return c.Send("User switched")
}
