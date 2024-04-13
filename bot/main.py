import os
import telebot

KEY = os.getenv("TELEGRAM_API_KEY")
if not KEY:
    raise RuntimeError("No API key found")
bot = telebot.TeleBot(KEY)


def create_ask_mode_keyboard():
    keyboard = telebot.types.ReplyKeyboardMarkup(
        one_time_keyboard=True, resize_keyboard=True)
    web_app_link = telebot.types.WebAppInfo("https://github.com/KXRXH")

    webapp_button = telebot.types.KeyboardButton(
        text="Режим веб-приложение", web_app=web_app_link)
    bot_button = telebot.types.KeyboardButton(text="Режим чат-бота")
    keyboard.add(webapp_button)
    keyboard.add(bot_button)
    return keyboard


@bot.message_handler(commands=['start'])
def start(message):
    bot.send_message(message.chat.id, 'Привет, я MemoryCode бот. В каком режиме вы хотите работать?',
                     reply_markup=create_ask_mode_keyboard())


@bot.message_handler(func=lambda message: message.text == "Режим чат-бота")
def go_to_web_app(message):
    bot.send_message(
        message.chat.id, "Прекрасно! Рад, что вы решили остаться здесь! ")


if __name__ == "__main__":
    bot.polling(non_stop=True, interval=0)
