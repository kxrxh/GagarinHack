from enum import Enum
import os
from telebot.async_telebot import AsyncTeleBot
import telebot

from storage.api import KeyValueStorage

KEY = os.getenv("TELEGRAM_API_KEY")
if not KEY:
    raise RuntimeError("No API key found")
bot = AsyncTeleBot(KEY)


class State(str, Enum):
    START = "start"
    NAME = "name"
    BIRTH_DATE = "birth_date"
    BIRTH_PLACE = "birth_place"


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
    await bot.send_message(message.chat.id, 'Привет! Я - MemoryCode Бот, ваш помощник в заполнении страницы памяти. Моя задача - сделать процесс заполнения страницы памяти легким и приятным для вас. Вам нужно будет предоставлять мне информацию по одному запросу за раз, и я помогу вам создать замечательные тексты для всех полей формы. Готовы начать этот важный процесс вместе?\n\nТак же я могу работать в режиме вебстраницы. Для этог просто нажмите кнопку ниже 🙂',
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


if __name__ == "__main__":
    import asyncio
    asyncio.run(bot.polling())
