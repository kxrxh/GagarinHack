from telebot.types import InlineKeyboardMarkup, InlineKeyboardButton


def regenerate_keyboard(accept_callback: str, repeat_callback: str, edit_callback: str) -> InlineKeyboardMarkup:
    keyboard = InlineKeyboardMarkup()
    keyboard.add(InlineKeyboardButton(
        text="✅ Принять", callback_data=accept_callback), InlineKeyboardButton(
        text="🔁 Сгенирировать занова", callback_data=repeat_callback))
    keyboard.add(InlineKeyboardButton(
        text="✏️ Редактировать", callback_data=edit_callback))
    return keyboard

# Function to handle the main inline keyboard


def generate_keyboard():
    keyboard = InlineKeyboardMarkup(row_width=2)
    next_button = InlineKeyboardButton("След", callback_data="next_question")
    prev_button = InlineKeyboardButton("Пред", callback_data="prev_question")
    finish_button = InlineKeyboardButton(
        "Закончить", callback_data="finish_questions")
    keyboard.add(prev_button, next_button)
    keyboard.add(finish_button)
    return keyboard
