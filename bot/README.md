# MemoryCode Bot

This is a Telegram bot that helps users create a memorial page for a deceased person. It guides the user through a series of questions to collect the necessary information, and then generates an epitaph and a biography based on the provided details.

## Features

- Collect basic information about the deceased person (name, sex, birth date, death date)
- Allow users to edit the generated epitaph and biography
- Gather additional details about the person (education, place of birth, place of death, children, citizenship, awards, work history)
- Generate a comprehensive biography in three parts (youth, middle age, old age)
- Provide a web app interface for a more immersive experience

## Prerequisites

- Python 3.7 or higher
- `python-telegram-bot` library
- `dotenv` library
- `backend.api` and `backend.types` modules
- `utils.keyboards`, `utils.utils`, and `utils.validator` modules
- `storage.api` module
- `os` and `re` modules

## Installation

### Manual

1. Clone the repository:

   ```
   git clone https://github.com/kxrxh/GagarinHack.git
   cd GagarinHack/bot
   ```

2. Install the required dependencies:

   ```
   pip install -r requirements.txt
   ```

3. Set the environment variables:

   - `STORAGE_URL`: the URL of the key-value storage API
   - `TELEGRAM_API_KEY`: the API key for your Telegram bot
   - `BACKEND_URL`: the URL of the backend API

4. Run the bot:

   ```
   python main.py
   ```
### Docker Compose
```
docker-compose up --build
```
*tip: You need to run this command from repo root folder*

*Docker compose will run kvdb and bot instances*



## Usage

1. Start the bot by sending the `/start` or `/help` command to your Telegram bot.
2. Follow the prompts to provide the necessary information about the deceased person.
3. The bot will generate an epitaph and a biography, which you can review and edit if needed.
4. The bot will then guide you through additional questions to collect more details about the person.
5. Once all the information is gathered, the bot will generate the full biography in three parts.

## Web App Interface

The bot also provides a web app interface that can be accessed by clicking the "Начать в режиме веб-приложение" button in the initial message. This interface offers a more immersive and visually appealing experience for creating the memorial page.


