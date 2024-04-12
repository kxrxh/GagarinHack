import telebot
from gpt.api import RequestCreator
from utils.const import Constant

bot = telebot.TeleBot(Constant.get_bot_key())


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

@bot.message_handler(func=lambda message: message.text == "Режим чат-бота")
def go_to_web_app(message):
    bot.send_message(
        message.chat.id, "Прекрасно! Рад, что вы решили остаться здесь! ")


@bot.message_handler(content_types=['text'])
def get_text_message(msg):
    body = {
        "modelUri": "gpt://b1gpt2d1mrgpkitoo7k6/yandexgpt-pro",
        "completionOptions": {
            "stream": False,
            "temperature": 0.25,
            "maxTokens": "2000"
        },
        "messages": [
            {
                "role": "system",
                "text": "Ты помошник MemoryCode."
            },
            {
                "role": "user",
                "text": msg.text
            }
        ]
    }
    loading = bot.send_message(msg.from_user.id, "Generating response...")
    code, res = RequestCreator.create_yandex_request(
        Constant.get_yandex_url(), body, Constant.get_ya_api_key()).send()
    bot.edit_message_text(
        f"{code}: {res['result']['alternatives'][0]['message']['text']}", message_id=loading.message_id, chat_id=msg.from_user.id)


if __name__ == "__main__":
    bot.polling(non_stop=True, interval=0)
