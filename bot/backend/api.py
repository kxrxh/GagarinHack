from dataclasses import dataclass, field
import logging
import os
import requests


@dataclass
class ApiResponse:
    success: bool = field(default_factory=bool)
    value: str = field(default_factory=str)
    message: str = field(default_factory=str)


class BackendApi:
    # Initialize the base URL for the key-value storage from the environment variable
    __BASE_URL = os.environ.get('BACKEND_URL')
    if not __BASE_URL:
        raise RuntimeError('BACKEND_URL environment variable not set')
    if not __BASE_URL.endswith('/'):
        __BASE_URL += '/'

    @staticmethod
    def generate_epitaph_yandex(name: str, date_birth: str, date_death: str, **kwargs) -> str:
        promt = f'Сгенирируй эпитафию для данного человека:\nФИО: {name}Дата рождения: {date_birth}Дата смерти: {date_death}'
        response = requests.post(f"http://127.0.0.1:3033/api/v1/completion/yandex", json={"request_message": promt})
        if response.status_code == 200:
            return response.json()['result']['alternatives'][0]['message']['text']
        else:
            return str(response.status_code)

if __name__ == '__main__':
    print(BackendApi.generate_epitaph_yandex('Иван', '12.12.1982', '12.12.1983'))

