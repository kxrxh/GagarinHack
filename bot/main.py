from enum import Enum
import os
import re
from telebot.async_telebot import AsyncTeleBot
import telebot

from storage.api import KeyValueStorage
from utils.utils import is_date_after
from utils.validator import Validator

KEY = os.getenv("TELEGRAM_API_KEY")
if not KEY:
    raise RuntimeError("No API key found")
bot = AsyncTeleBot(KEY)


class State(str, Enum):
    START = "start"
    NAME = "name"
    BIRTH_DATE = "birth_date"
    DEATH_DATE = "death_date"
    EPITAPHIA = "epitaphia"


def create_mode_keyboard():
    web_app_link = telebot.types.WebAppInfo("https://cataas.com/cat")

    keyboard = telebot.types.InlineKeyboardMarkup()
    web_app_link = telebot.types.InlineKeyboardButton(
        "📱 Режим веб-приложение", web_app=web_app_link)
    start_button = telebot.types.InlineKeyboardButton(
        "🚩 Начать заполнение", callback_data="start")
    keyboard.add(web_app_link)
    keyboard.add(start_button)
    return keyboard


@bot.message_handler(commands=['start', 'help'])
async def start(message):
    KeyValueStorage.set(message.chat.id, State.START.value)
    await bot.send_message(message.chat.id, 'Привет! Я - MemoryCode Бот, ваш помощник в заполнении страницы памяти. Моя задача - сделать процесс заполнения страницы памяти легким и приятным для вас. \
Вам нужно будет предоставлять мне информацию по одному запросу за раз, и я помогу вам создать замечательные тексты для всех полей формы. \
Готовы начать этот важный процесс вместе?\n\nТак же я могу работать в режиме веб-страницы. Для этого просто нажмите кнопку ниже 🙂',
                           reply_markup=create_mode_keyboard())


@bot.message_handler(func=lambda message: not KeyValueStorage.get(str(message.chat.id)).value)
async def handle_message(message):
    await bot.send_message(
        message.chat.id, "Выберите команду '/start' или '/mode', чтобы начать!")


@bot.message_handler(func=lambda message: message.text == "Режим чат-бота" and KeyValueStorage.get(str(message.chat.id)).value == State.START.value)
async def go_to_bot_mode(message):
    KeyValueStorage.set(message.chat.id, State.BOT.value)
    await bot.reply_to(message, 'Вы выбрали режим чат-бота. Если вы хотите сменить режим выберите команду: \'/mode\'', reply_markup=None)


@bot.message_handler(func=lambda message: KeyValueStorage.get(str(message.chat.id)).value == State.START.value)
async def handle_random_message(message):
    keyboard = telebot.types.InlineKeyboardMarkup()
    keyboard.add(telebot.types.InlineKeyboardButton(
        text="Начать заполнение", callback_data="start"))
    await bot.reply_to(message, "Если вы хотите начать заполнение данных, нажмите кнопку ниже", reply_markup=keyboard)


@bot.callback_query_handler(func=lambda call: call.data == "start")
async def start_callback(call):
    KeyValueStorage.set(call.message.chat.id, State.NAME.value)
    await bot.send_message(call.message.chat.id, "Отлично, давайте же приступим к заполнению страницы памяти! Пожалуйста, укажите полное имя человека (ФИО):")


@bot.message_handler(content_types=['text'], func=lambda message: KeyValueStorage.get(str(message.chat.id)).value == State.NAME.value)
async def handle_name(message):
    # delete repeated whitespaces
    cleaned_string = re.sub(r'\s+', ' ', message.text)
    is_valid_name = Validator.validate_name(cleaned_string)
    if not is_valid_name:
        await bot.send_message(message.chat.id, "Извините, введенное вами имя некорректно. Пожалуйста, убедитесь, что вы ввели ваше полное имя на кириллице без использования цифр или специальных символов.")
        return
    KeyValueStorage.set(message.chat.id, State.BIRTH_DATE.value)
    KeyValueStorage.set(f"{message.chat.id}.name", cleaned_string)
    await bot.send_message(message.chat.id, f"Вы успешно установили поле 'ФИО'. Пожалуйста, укажите когда родился {cleaned_string} в формате 'ДД.ММ.ГГГГ' (например, 31.12.1989):")


@bot.message_handler(content_types=['text'], func=lambda message: KeyValueStorage.get(str(message.chat.id)).value == State.BIRTH_DATE.value)
async def handle_birth_date(message):
    date = message.text.replace('-', '.').replace(' ', '.').replace('/', '.')
    is_date_valid = Validator.validate_date(date)
    if not is_date_valid:
        await bot.send_message(message.chat.id, "Дата рождения введена неверно. Пожалуйста, введите дату рождения в формате 'ДД.ММ.ГГГГ' (например, 31.12.1990).")
        return
    KeyValueStorage.set(message.chat.id, State.DEATH_DATE.value)
    KeyValueStorage.set(f"{message.chat.id}.birth_date", date)
    await bot.send_message(message.chat.id, "Вы успешно установили поле 'Дата рождения'. Пожалуйста, укажите дату смерти в формате 'ДД.ММ.ГГГГ' (например, 31.12.2020):")


@bot.message_handler(content_types=['text'], func=lambda message: KeyValueStorage.get(str(message.chat.id)).value == State.DEATH_DATE.value)
async def handle_death_date(message):
    date = message.text.replace('-', '.').replace(' ', '.').replace('/', '.')
    
    is_date_valid = Validator.validate_date(date)
    if not is_date_valid:
        await bot.send_message(message.chat.id, "Дата смерти введена неверно. Пожалуйста, введите дату смерти в формате 'ДД.ММ.ГГГГ' (например, 31.12.2020).")

    birth_date = KeyValueStorage.get(f"{message.chat.id}.birth_date").value

    if not birth_date:
        keyboard = telebot.types.InlineKeyboardMarkup()
        keyboard.add(telebot.types.InlineKeyboardButton(
            text="Начать заполнение", callback_data="start"))
        KeyValueStorage.set(message.chat.id, State.START.value)
        await bot.send_message(message.chat.id, "Произошла ошибка на стороне сервера 😿. Попробуйте начать заполнение ещё раз...", reply_markup=keyboard)
        return
    
    if not is_date_after(date, birth_date):
        await bot.send_message(message.chat.id, "Дата смерти не может быть раньше даты рождения. Пожалуйста, введите правильную дату.", reply_markup=None)
        return
    
    KeyValueStorage.set(message.chat.id, State.EPITAPHIA.value)
    KeyValueStorage.set(f"{message.chat.id}.death_date", message.text)
    await bot.send_message(message.chat.id, "Дата успешно установлена. Теперь переходим к генерации эпитафии. Пожалуйста, подождите немного, пока я создам для вас уникальный текст, который вы сможете редактировать, если захотите.")


if __name__ == "__main__":
    import asyncio
    asyncio.run(bot.polling())
