from dataclasses import dataclass, field
import os
from typing import List
import requests

from backend.types import GenerationRequest


@dataclass
class ApiResponse:
    success: bool = field(default_factory=bool)
    value: str = field(default_factory=str)
    message: str = field(default_factory=str)


class GptAPI:
    # Initialize the base URL for the key-value storage from the environment variable
    __BASE_URL = os.environ.get('BACKEND_URL')
    if not __BASE_URL:
        raise RuntimeError('BACKEND_URL environment variable not set')
    if not __BASE_URL.endswith('/'):
        __BASE_URL += '/'

    @staticmethod
    def generate_epitaph_yandex(req_body: GenerationRequest) -> str:
        url = GptAPI.__BASE_URL + 'completion/yandex/epitaph'
        response = requests.post(url, json={"human_info": req_body.__dict__, "type_of_story": "epitaph"})
        return response.json()['response']

    @staticmethod
    def generate_epitaph_gigachat(req_body: GenerationRequest) -> str:
        url = GptAPI.__BASE_URL + 'completion/gigachat/epitaph'
        response = requests.post(url, json={"human_info": req_body.__dict__, "type_of_story": "epitaph"})
        return response.json()['response']


    @staticmethod
    def generate_biography_yandex(req_body: GenerationRequest, type: str) -> str:
        url = GptAPI.__BASE_URL + 'completion/yandex/story'
        response = requests.post(url, json={"human_info": req_body.__dict__, "type_of_story": type})
        return response.json()['response']

    @staticmethod
    def generate_biography_gigachat(req_body: GenerationRequest, type: str, prev: str) -> str:
        url = GptAPI.__BASE_URL + 'completion/gigachat/biography'
        response = requests.post(url, json={"human_info": req_body.__dict__, "type_of_story": type, "previous": prev})
        return response.json()['response']
    
    @staticmethod
    def generate_questions(req_body: GenerationRequest) -> List[str]:
        url = GptAPI.__BASE_URL + 'completion/gigachat/questions'
        response = requests.post(url, json={"human_info": req_body.__dict__})
        return response.json()['response']


if __name__ == '__main__':
    print(GptAPI.generate_epitaph_yandex('Иван', '12.12.1982', '12.12.1983'))
