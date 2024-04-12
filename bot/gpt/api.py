from typing import Tuple
import requests


class RequestCreator:
    def __init__(self, url, auth_token, body, prefix="Api-Key"):
        self.url = url
        self.auth_token = auth_token
        self.body = body
        self.prefix = prefix

    def send(self) -> Tuple[int, dict]:
        headers = {
            'Authorization': f'{self.prefix} {self.auth_token}',
            'Content-Type': 'application/json'
        }

        response = requests.post(self.url, headers=headers, json=self.body)
        return (response.status_code, response.json())

    @staticmethod
    def create_yandex_request(url, body, auth_token):
        return RequestCreator(url, auth_token, body)

    @staticmethod
    def create_gigachad_request(url, body, auth_token):
        # Todo
        return None
