package telegram

import (
	"fmt"
	"gopkg.in/telebot.v3"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func sendImagineUsage(c telebot.Context) error {
	return c.Send("Usage: /imagine <<horizontal/vertical> <prompt> | <horizontal/vertical> [prompt] with reply to message>")
}

func (b *GPTBot) imagineCommand(c telebot.Context) error {
	args := c.Args()

	validReply := c.Message().IsReply() && c.Message().ReplyTo.Text != ""

	if !(len(args) >= 2 || len(args) == 1 && validReply) {
		return sendImagineUsage(c)
	}

	horizontal := args[0] == "horizontal"
	prompt := strings.Join(args[1:], " ")

	if validReply {
		prompt = c.Message().ReplyTo.Text + " " + prompt
	}

	_ = c.Notify(telebot.UploadingPhoto)

	image, err := b.generator.CreateImage(prompt, horizontal)

	if err != nil {
		_ = c.Send("Fail to imagine image")
		return err
	}

	go func() {
		path, err := downloadToTempFile(image)

		if err != nil {
			_ = c.Send("Fail to upload image")

			fmt.Println(err.Error())
			return
		}

		_ = c.Notify(telebot.UploadingPhoto)

		_, err = b.bot.Send(c.Recipient(), &telebot.Photo{
			File:    telebot.FromDisk(path),
			Caption: "[Result](" + image + ")",
		}, telebot.ModeMarkdown)

		if err != nil {
			_ = c.Send("Fail to upload image")

			fmt.Println(err.Error())
			return
		}

		_ = os.Remove(path)

	}()

	return nil
}

func downloadToTempFile(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	path := filepath.Join("/tmp/gptbot")

	_ = os.Mkdir(path, 0770)

	file, err := os.CreateTemp(path, "dalle_image*.png")
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return "", err
	}

	return file.Name(), nil
}
