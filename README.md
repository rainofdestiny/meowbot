# MeowBot

MeowBot is a Telegram bot that responds to "meow", tracks statistics, and provides information about user interactions.

## Key Features

- Responds to messages containing "meow", "мур", or "мяу".
- Tracks how many times a user has sent messages to the bot or other users.
- Uses Redis for data storage.

## Installation

### Requirements
- Docker and Docker Compose

### Installation Steps
1. Clone the repository:
   ```shell
   git clone https://github.com/rainofdestiny/meowbot.git
   cd meowbot
   ```

2. Create a .env file and add the following variables:
   ```shell
   TELEGRAM_TOKEN=your_telegram_bot_token
   REDIS_ADDR=redis:6379
   ```

3. Start the project using Docker Compose:
   ```shell
   docker-compose up --build
   ```

The bot will automatically connect to Telegram and start working.
