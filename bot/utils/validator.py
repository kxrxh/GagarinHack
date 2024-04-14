import re


class Validator:
    __cyrillic_pattern = re.compile(r'^[а-яА-ЯёЁ\-]+$')
    __date_pattern = re.compile(
        r'^(?:(?:31(\/|-|\.)(?:0?[13578]|1[02]))\1|(?:(?:29|30)(\/|-|\.)(?:0?[13-9]|1[0-2])\2))(?:(?:1[6-9]|[2-9]\d)?\d{2})$|^(?:29(\/|-|\.)0?2\3(?:(?:(?:1[6-9]|[2-9]\d)?(?:0[48]|[2468][048]|[13579][26])|(?:(?:16|[2468][048]|[3579][26])00))))$|^(?:0?[1-9]|1\d|2[0-8])(\/|-|\.)(?:(?:0?[1-9])|(?:1[0-2]))\4(?:(?:1[6-9]|[2-9]\d)?\d{2})$')

    @staticmethod
    def validate_name(name: str) -> bool:
        parts = name.split(' ')
        if len(parts) not in [2, 3]:
            return False

        for part in parts:
            if not Validator.__cyrillic_pattern.match(part):
                return False
            if part.startswith('-') or part.endswith('-'):
                return False

        return True

    @staticmethod
    def validate_date(date: str) -> bool:
        return Validator.__date_pattern.match(date) is not None
