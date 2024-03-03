package main

import (
	"TelegramGPT/internal/app"
)

func main() {
	mainApp := app.App{}
	err := mainApp.Start()

	if err != nil {
		panic(err)
	}
}
