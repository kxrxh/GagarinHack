from datetime import datetime


def is_date_after(date_str1, date_str2):
    try:
        date1 = datetime.strptime(date_str1, '%d.%m.%Y')
        date2 = datetime.strptime(date_str2, '%d.%m.%Y')
        return date1 >= date2
    except ValueError:
        return False
