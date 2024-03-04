# Telegram GPT

## Overview

This project is a Telegram bot written in Golang that provides users with access to OpenAI API capabilities and enabling users to interact with artificial intelligence directly from their chat.

## Key Features

- **Whitelist Access**: Access to the bot is restricted to a list of approved User IDs and chats to ensure security and control over who can use the bot.
- **Preset System**: Users can choose the bot's behavior using hashtags to activate predefined settings and scenarios.
- **Imagine Command**: Users can also generate images via DALL·E 3
- **Chat Support**: It is convenient to use in chats. For example, when someone asked a question. Just reply to the message mentioning the bot, it will do all the work for you!

## Getting Started

### Prerequisites

To run the bot, you must have the following software installed:
- Docker
- Docker Compose

### Installation and Launch

Follow these steps to install and launch the bot:

1. Clone the repository:
   ```sh
   git clone https://github.com/isKONSTANTIN/Telegram-GPT
   cd Telegram-GPT
   ```

2. Navigate to the production directory:
   ```sh
   cd production
   ```

3. Run the start script:
   ```sh
   ./start.sh
   ```

   > Note: On the first launch, the bot will generate the necessary configuration files.

4. Stop the bot:
   ```sh
   docker-compose stop
   ```

5. Edit the configuration file production/configs/main.json to add your API keys for OpenAI and Telegram, also set your user id for admin access

6. After configuring, launch the bot as in step 3

### Configuration

The main.json file contains the following key settings:

- telegram -> telegramToken: Your Telegram bot token.
- telegram -> mainAdminId: Your user id (you can use [this](https://t.me/userinfobot) to find out your id)
- openAI -> token: Your OpenAI API key.
- openAI -> model: ChatGPT model (full list [here](https://github.com/sashabaranov/go-openai/blob/master/completion.go))
- openAI -> defaultPreset: Default prompt for behavior

## Support

If you have any questions or issues related to the use of the bot, please create an issue in the project repository on GitHub or [ask me directly](https://t.me/isKONSTANTIN) (English or Russian)