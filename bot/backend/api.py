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
    base_url = os.environ.get('BACKEND_URL')
    if not base_url:
        raise RuntimeError('BACKEND_URL environment variable not set')
    if not base_url.endswith('/'):
        base_url += '/'

    @staticmethod
    def generate_epitaph_yandex(name: str, date_birth: str, date_death: str, **kwargs) -> str:
        promt = f'Сгенирируй эпитафию для данного человека:\nФИО: {name}Дата рождения: {date_birth}Дата смерти: {date_death}\nскажи добрые слова про него'
        

if __name__ == '__main__':
    pass
