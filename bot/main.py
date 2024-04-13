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

    EDUCATION = "education"
    WORK = "work"
    PLACE_OF_BIRTH = "place_of_birth"
    PLACE_OF_DEATH = "place_of_death"
    KIDS = "kids"
    CITIZENSHIP = "citizenship"
    AWARDS = "awards"

    # Save result
    FINISHED = "finished"


def create_mode_keyboard():
    web_app_link = telebot.types.WebAppInfo(
        "https://themixadev.github.io/GagarinHackView/")

    keyboard = telebot.types.InlineKeyboardMarkup()
    web_app_link = telebot.types.InlineKeyboardButton(
        "📱 Начать в режиме веб-приложение", web_app=web_app_link)
    start_button = telebot.types.InlineKeyboardButton(
        "🚩 Начать в простом режиме", callback_data="start")
    keyboard.add(web_app_link)
    keyboard.add(start_button)
    return keyboard


@bot.message_handler(commands=['start', 'help'])
async def start(message):
    KeyValueStorage.set(message.chat.id, State.START.value)
    await bot.send_message(message.chat.id, 'Привет! Я - MemoryCode Бот, ваш помощник в заполнении страницы памяти. Моя задача - сделать процесс заполнения страницы памяти легким и приятным для вас. \
Вам нужно будет предоставлять мне информацию по одному запросу за раз, и я помогу вам создать замечательные тексты для всех полей формы. \
Готовы начать этот важный процесс вместе?\n\nТак же я могу работать в режиме веб-страницы (полная версия). Для этого просто нажмите кнопку ниже 🙂',
                           reply_markup=create_mode_keyboard())


@bot.message_handler(func=lambda message: not KeyValueStorage.get(str(message.chat.id)).value)
async def handle_message(message):
    await bot.send_message(
        message.chat.id, "Выберите команду '/start' или '/help', чтобы начать!")


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
    keyboard.add(
        # telebot.types.InlineKeyboardButton(text="❓ Далее", callback_data=f"questions_{message.id}"),
        telebot.types.InlineKeyboardButton(text="⏩ Генерация", callback_data=f"generate_{message.chat.id}"))
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
    if call.data.split('_')[1] == State.EPITAPH.value:
        KeyValueStorage.set(chat_id, State.EDUCATION.value)
        await bot.send_message(chat_id=chat_id, text="Теперь мне нужна информация об образовании:")


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
        KeyValueStorage.set(str(message.from_user.id), State.EDUCATION.value)
        await bot.send_message(chat_id=message.chat.id, text="Теперь мне нужна информация об образовании:")

    elif KeyValueStorage.get(str(message.from_user.id)).value == State.EDIT_BIOGRAPHY.value:
        KeyValueStorage.set(f"{message.chat.id}.biography", message.text)
        await bot.send_message(chat_id=message.chat.id, text="Текст был успешно изменен!")


@bot.message_handler(content_types=['text'], func=lambda message: KeyValueStorage.get(str(message.from_user.id)).value == State.EDUCATION.value)
async def handle_education(message):
    # Store education information
    KeyValueStorage.set(f"{message.chat.id}.education", message.text)
    await bot.reply_to(message, "Отлично, я сохранил информацию об образовании!")

    # Go to place of birth mode
    KeyValueStorage.set(message.chat.id, State.PLACE_OF_BIRTH.value)
    await bot.send_message(message.chat.id, "Теперь, пожалуйста, укажите место рождения:")


@bot.message_handler(content_types=['text'], func=lambda message: KeyValueStorage.get(str(message.from_user.id)).value == State.PLACE_OF_BIRTH.value)
async def handle_place_of_birth(message):
    # Store place of birth
    KeyValueStorage.set(f"{message.chat.id}.place_of_birth", message.text)
    await bot.reply_to(message, "Место рождения сохранено!")

    # Go to place of death mode
    KeyValueStorage.set(message.chat.id, State.PLACE_OF_DEATH.value)
    await bot.send_message(message.chat.id, "Теперь, пожалуйста, укажите место смерти:")


@bot.message_handler(content_types=['text'], func=lambda message: KeyValueStorage.get(str(message.from_user.id)).value == State.PLACE_OF_DEATH.value)
async def handle_place_of_death(message):
    # Store place of death
    KeyValueStorage.set(f"{message.chat.id}.place_of_death", message.text)
    await bot.reply_to(message, "Место смерти сохранено!")

    # Go to kids mode
    KeyValueStorage.set(message.chat.id, State.KIDS.value)
    await bot.send_message(message.chat.id, "Теперь, пожалуйста, укажите, были ли у этого человека дети:")


@bot.message_handler(content_types=['text'], func=lambda message: KeyValueStorage.get(str(message.from_user.id)).value == State.KIDS.value)
async def handle_kids(message):
    # Store kids information
    KeyValueStorage.set(f"{message.chat.id}.kids", message.text)
    await bot.reply_to(message, "Информация о детях сохранена!")

    # Go to citizenship mode
    KeyValueStorage.set(message.chat.id, State.CITIZENSHIP.value)
    await bot.send_message(message.chat.id, "Теперь, пожалуйста, укажите гражданство этого человека:")


@bot.message_handler(content_types=['text'], func=lambda message: KeyValueStorage.get(str(message.from_user.id)).value == State.CITIZENSHIP.value)
async def handle_citizenship(message):
    # Store citizenship information
    KeyValueStorage.set(f"{message.chat.id}.citizenship", message.text)
    await bot.reply_to(message, "Информация о гражданстве сохранена!")

    # Go to awards mode
    KeyValueStorage.set(message.chat.id, State.AWARDS.value)
    await bot.send_message(message.chat.id, "Наконец, пожалуйста, укажите, были ли у этого человека какие-либо награды или достижения:")


@bot.message_handler(content_types=['text'], func=lambda message: KeyValueStorage.get(str(message.from_user.id)).value == State.AWARDS.value)
async def handle_awards(message):
    # Store awards information
    KeyValueStorage.set(f"{message.chat.id}.awards", message.text)
    await bot.reply_to(message, "Информация о наградах и достижениях сохранена!")

    KeyValueStorage.set(f"{message.chat.id}", State.WORK.value)
    await bot.send_message(message.chat.id, "Теперь, пожалуйста, укажите, чем этот человек занимался (кем работал):")


@bot.message_handler(content_types=['text'], func=lambda message: KeyValueStorage.get(str(message.from_user.id)).value == State.WORK.value)
async def handle_work(message):
    KeyValueStorage.set(f"{message.chat.id}.work", message.text)
    await bot.reply_to(message, "Информация о работе сохранена!")

    KeyValueStorage.set(f"{message.chat.id}", State.BIOGRAPHY.value)
    await bot.send_message(message.chat.id, "Поздравляю, вы успешно заполнили все необходимые поля! Теперь я попытаюсь написать биографию...")
    name = KeyValueStorage.get(f"{message.chat.id}.name").value
    sex = KeyValueStorage.get(f"{message.chat.id}.sex").value
    education = KeyValueStorage.get(f"{message.chat.id}.education").value
    place_of_birth = KeyValueStorage.get(
        f"{message.chat.id}.place_of_birth").value
    place_of_death = KeyValueStorage.get(
        f"{message.chat.id}.place_of_death").value
    birth_date = KeyValueStorage.get(f"{message.chat.id}.birth_date").value
    death_date = KeyValueStorage.get(f"{message.chat.id}.death_date").value
    kids = KeyValueStorage.get(f"{message.chat.id}.kids").value
    awards = KeyValueStorage.get(f"{message.chat.id}.awards").value
    citizenship = KeyValueStorage.get(f"{message.chat.id}.citizenship").value

    person = GenerationRequest(name, sex, birth_date, death_date,
                               {"Кем работал?": message.text, "Какое образование получил человек?": education, "Где родился?": place_of_birth, "Где умер?": place_of_death, "Были ли у этого человека дети?": kids, "Гражданство?": citizenship, "Были ли у этого человека какие-либо награды или достижения?": awards})
    bio1 = GptAPI.generate_biography_gigachat(person, "youth", "Пустая биография")
    await bot.send_message(message.chat.id, f"Молодось:\n{bio1}")
    bio2 = GptAPI.generate_biography_gigachat(person, "middle_age", bio1)
    await bot.send_message(message.chat.id, f"Средние года:\n{bio2}")
    bio3 = GptAPI.generate_biography_gigachat(person, "old_age", bio2)
    await bot.send_message(message.chat.id, f"Последние года:\n{bio3}")
    


if __name__ == "__main__":
    import asyncio
    asyncio.run(bot.polling())
