from enum import Enum
import os
import re
from backend.api import GptAPI
from backend.types import GenerationRequest
from utils.keyboards import regenerate_keyboard
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
    # Required fields
    START = "start"
    NAME = "name"
    SEX = "sex"
    BIRTH_DATE = "birth_date"
    DEATH_DATE = "death_date"
    EDIT_EPITAPH = "edit_epitaph"
    EDIT_BIOGRAPHY = "edit_biography"
    # Optional fields
    QUESTIONS = "questions"

    # Generated fields
    EPITAPH = "epitaph"
    BIOGRAPHY = "biography"

    # Save result
    FINISHED = "finished"


def create_mode_keyboard():
    web_app_link = telebot.types.WebAppInfo(
        "https://themixadev.github.io/GagarinHackView/")

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
    # clear old data
    KeyValueStorage.delete_prefix(call.message.chat.id)
    KeyValueStorage.set(call.message.chat.id, State.SEX.value)
    keyboard = telebot.types.ReplyKeyboardMarkup(
        resize_keyboard=True, one_time_keyboard=True)
    keyboard.add(telebot.types.KeyboardButton(text="♂️ Мужской"),
                 telebot.types.KeyboardButton(text="♀️ Женский"))
    await bot.send_message(call.message.chat.id, "Отлично, давайте же приступим к заполнению страницы памяти! Пожалуйста, выберите пол человека, которого хотите внести в память.",
                           reply_markup=keyboard)


@bot.message_handler(content_types=['text'], func=lambda message:
                     KeyValueStorage.get(str(message.chat.id)).value == State.SEX.value and
                     (message.text == "♂️ Мужской" or message.text == "♀️ Женский"))
async def handle_sex(message):
    # store sex in database
    KeyValueStorage.set(f"{message.chat.id}.sex", message.text[2:])
    await bot.reply_to(message, "Хорошо, я обязательно запомню ваш ответ!")

    # go to name mode
    KeyValueStorage.set(message.chat.id, State.NAME.value)
    await bot.send_message(message.chat.id, "Теперь, пожалуйста, укажите фамилию, имя и отчество (при наличии):")


@bot.message_handler(content_types=['text'], func=lambda message: KeyValueStorage.get(str(message.chat.id)).value == State.NAME.value)
async def handle_name(message):
    # delete repeated whitespaces
    cleaned_string = re.sub(r'\s+', ' ', message.text)
    # validate name
    is_valid_name = Validator.validate_name(cleaned_string)
    if not is_valid_name:
        await bot.send_message(message.chat.id, "😔 Извините, введенное вами имя некорректно. Пожалуйста, убедитесь, что вы ввели ваше полное имя на кириллице без использования цифр или специальных символов.")
        return
    # store name
    KeyValueStorage.set(f"{message.chat.id}.name", cleaned_string)
    await bot.reply_to(message, "Замечательно, буду держать это имя у себя в голове! 🧠")

    # go to birth date mode
    KeyValueStorage.set(message.chat.id, State.BIRTH_DATE.value)
    await bot.send_message(message.chat.id, f"Теперь, пожалуйста, укажите когда родился(-ась) {cleaned_string} в формате 'ДД.ММ.ГГГГ' (например, 31.12.1989):")


@bot.message_handler(content_types=['text'], func=lambda message: KeyValueStorage.get(str(message.chat.id)).value == State.BIRTH_DATE.value)
async def handle_birth_date(message):
    # change format of date, if user input is invalid
    date = message.text.replace('-', '.').replace(' ', '.').replace('/', '.')

    # validate date
    is_date_valid = Validator.validate_date(date)
    if not is_date_valid:
        await bot.send_message(message.chat.id, "😥 Дата рождения введена неверно. Пожалуйста, введите дату рождения в формате 'ДД.ММ.ГГГГ' (например, 31.12.1990):")
        return

    # store birth_date
    KeyValueStorage.set(f"{message.chat.id}.birth_date", date)

    # go to death date mode
    KeyValueStorage.set(message.chat.id, State.DEATH_DATE.value)
    await bot.send_message(message.chat.id, "Вы успешно установили поле 'Дата рождения'. Пожалуйста, укажите дату смерти в формате 'ДД.ММ.ГГГГ' (например, 31.12.2020):")


@bot.message_handler(content_types=['text'], func=lambda message: KeyValueStorage.get(str(message.chat.id)).value == State.DEATH_DATE.value)
async def handle_death_date(message):
    date = message.text.replace('-', '.').replace(' ', '.').replace('/', '.')

    is_date_valid = Validator.validate_date(date)
    if not is_date_valid:
        await bot.send_message(message.chat.id, "😥 Дата введена неверно. Пожалуйста, введите дату в формате 'ДД.ММ.ГГГГ' (например, 31.12.2020):")

    birth_date = KeyValueStorage.get(f"{message.chat.id}.birth_date").value

    if not birth_date:
        keyboard = telebot.types.InlineKeyboardMarkup()
        keyboard.add(telebot.types.InlineKeyboardButton(
            text="Начать заполнение", callback_data="start"))
        KeyValueStorage.set(message.chat.id, State.START.value)
        await bot.send_message(message.chat.id, "😿 Произошла ошибка на стороне сервера. Попробуйте начать заполнение ещё раз...", reply_markup=keyboard)
        return

    if not is_date_after(date, birth_date):
        await bot.send_message(message.chat.id, "😥 Дата смерти не может быть раньше даты рождения. Пожалуйста, введите правильную дату.", reply_markup=None)
        return

    # store death_date
    KeyValueStorage.set(f"{message.chat.id}.death_date", date)
    await bot.reply_to(message, "Поле с датой смерти сохранено!")

    # ask youser for mode
    keyboard = telebot.types.InlineKeyboardMarkup(row_width=2)
    keyboard.add(telebot.types.InlineKeyboardButton(
        text="❓ Далее", callback_data=f"questions_{message.id}"), telebot.types.InlineKeyboardButton(
        text="⏩ Генерация", callback_data=f"generate_{message.chat.id}"))
    await bot.send_message(message.chat.id,
                           "Поздравляю! Все основные поля уже заполнены.\n\nНажмите кнопку 'Далее', чтобы ответить ещё на пару дополнительных вопросов, чтобы я смог узнать побольше о вносимом человеке.\nЛибо нажмите кнопку 'Генерация', если вы хотите сразу перейти к созданию биографии и эпитафии.",
                           reply_markup=keyboard)


@bot.callback_query_handler(func=lambda call: call.data.startswith("questions_") and KeyValueStorage.get(str(call.from_user.id)).value == State.DEATH_DATE.value)
async def questions(call: telebot.types.CallbackQuery):
    await bot.delete_message(message_id=call.message.id, chat_id=call.from_user.id)
    chat_id = call.from_user.id
    # go to question mode

    name = KeyValueStorage.get(f"{chat_id}.name").value
    sex = KeyValueStorage.get(f"{chat_id}.sex").value
    birth_date = KeyValueStorage.get(f"{chat_id}.birth_date").value
    death_date = KeyValueStorage.get(f"{chat_id}.death_date").value

    request_body = GenerationRequest(name, sex, birth_date, death_date)

    q = GptAPI.generate_questions(request_body)
    await bot.send_message(chat_id=chat_id, text="\n".join(q))
    KeyValueStorage.set(call.from_user.id, State.QUESTIONS.value)


@bot.callback_query_handler(func=lambda call: call.data.startswith("generate_") and KeyValueStorage.get(str(call.from_user.id)).value == State.DEATH_DATE.value)
async def generate(call: telebot.types.CallbackQuery):
    await bot.delete_message(message_id=call.message.id, chat_id=call.from_user.id)
    chat_id = int(call.data.split("_")[1])

    # go to epitaph mode
    KeyValueStorage.set(chat_id, State.EPITAPH.value)

    name = KeyValueStorage.get(f"{chat_id}.name").value
    sex = KeyValueStorage.get(f"{chat_id}.sex").value
    birth_date = KeyValueStorage.get(f"{chat_id}.birth_date").value
    death_date = KeyValueStorage.get(f"{chat_id}.death_date").value

    request_body = GenerationRequest(name, sex, birth_date, death_date)

    generation_msg = await bot.send_message(chat_id, f"Пожалуйста, подождите немного, пока я создам для вас уникальный текст, который вы сможете редактировать, если захотите.")
    epitah = GptAPI.generate_epitaph_gigachat(request_body)
    KeyValueStorage.set(f"{chat_id}.epitaph", epitah)
    keyboard = regenerate_keyboard(
        f"accept_{State.EPITAPH.value}", f"regenerate_{State.EPITAPH.value}", f"edit_{State.EPITAPH.value}")
    await bot.edit_message_text(text=f"Эпитафия:\n{epitah}",
                                reply_markup=keyboard, chat_id=chat_id, message_id=generation_msg.id)


@bot.callback_query_handler(func=lambda call: call.data.startswith("regenerate_"))
async def regenerate(call: telebot.types.CallbackQuery):
    obj_type = call.data.split("_")[1]
    chat_id = call.message.chat.id
    msg_id = call.message.id

    name = KeyValueStorage.get(f"{chat_id}.name").value
    sex = KeyValueStorage.get(f"{chat_id}.sex").value
    birth_date = KeyValueStorage.get(f"{chat_id}.birth_date").value
    death_date = KeyValueStorage.get(f"{chat_id}.death_date").value

    request_body = GenerationRequest(name, sex, birth_date, death_date)
    if obj_type == State.EPITAPH.value:
        epitah = GptAPI.generate_epitaph_gigachat(request_body)
        KeyValueStorage.set(f"{chat_id}.epitaph", epitah)
        keyboard = regenerate_keyboard(
            f"accept_{State.EPITAPH.value}", call.data, f"edit_{State.EPITAPH.value}")
        await bot.edit_message_text(text=f"Эпитафия:\n{epitah}", reply_markup=keyboard, chat_id=chat_id, message_id=msg_id)
    elif obj_type == State.BIOGRAPHY.value:
        biography = GptAPI.generate_biography_gigachat(request_body)
        KeyValueStorage.set(f"{chat_id}.biography", biography)
        keyboard = regenerate_keyboard(
            f"accept_{State.BIOGRAPHY.value}", call.data, f"edit_{State.BIOGRAPHY.value}")
        await bot.edit_message_text(text=f"Биография:\n{biography}", reply_markup=keyboard, chat_id=chat_id, message_id=msg_id)


@bot.callback_query_handler(func=lambda call: call.data.startswith("accept_"))
async def accept_epitaph(call: telebot.types.CallbackQuery):
    chat_id = call.message.chat.id
    await bot.edit_message_reply_markup(chat_id=chat_id, message_id=call.message.id, reply_markup=None)
    await bot.reply_to(text=f"{'Эпитафия' if call.data.split('_')[1] == State.EPITAPH.value else 'Биография'} была успешно создана!", message=call.message)


@bot.callback_query_handler(func=lambda call: call.data.startswith("edit_"))
async def edit(call: telebot.types.CallbackQuery):
    obj_type = call.data.split("_")[1]

    if obj_type == State.EPITAPH.value:
        KeyValueStorage.set(call.from_user.id, State.EDIT_EPITAPH.value)
    elif obj_type == State.BIOGRAPHY.value:
        KeyValueStorage.set(call.from_user.id, State.EDIT_BIOGRAPHY.value)
    await bot.edit_message_reply_markup(chat_id=call.message.chat.id, message_id=call.message.id, reply_markup=None)
    await bot.send_message(chat_id=call.message.chat.id, text="Пришли мне новым сообщением свою версию текста и я его сохраню!")

@bot.message_handler(content_types=['text'], func=lambda message: KeyValueStorage.get(str(message.from_user.id)).value in [State.EDIT_EPITAPH.value, State.EDIT_BIOGRAPHY.value])
async def handle_edited_text(message: telebot.types.Message):
    if KeyValueStorage.get(str(message.from_user.id)).value == State.EDIT_EPITAPH.value:
        KeyValueStorage.set(f"{message.chat.id}.epitaph", message.text)
        await bot.send_message(chat_id=message.chat.id, text="Текст был успешно изменен!")
        KeyValueStorage.set(str(message.from_user.id), State.BIOGRAPHY.value)

    elif KeyValueStorage.get(str(message.from_user.id)).value == State.EDIT_BIOGRAPHY.value:
        KeyValueStorage.set(f"{message.chat.id}.biography", message.text)
        await bot.send_message(chat_id=message.chat.id, text="Текст был успешно изменен!")
        
if __name__ == "__main__":
    import asyncio
    asyncio.run(bot.polling())
