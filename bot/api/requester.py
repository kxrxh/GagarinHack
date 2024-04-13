import requests
from enum import Enum
from typing import Union, Dict, Any


class RequestType(Enum):
    GET = 1
    POST = 2


class Request:
    def __init__(self, url: str) -> None:
        self.url = url

    def make_request(self, request_type: RequestType, data: Union[None, Dict[str, Any]] = None) -> Union[None, Dict[str, Any]]:
        match request_type:
            case RequestType.GET:
                response = requests.get(self.url)
            case RequestType.POST:
                response = requests.post(self.url, json=data)
            case _:
                raise ValueError(
                    "Invalid request type. Please use GET or POST.")
        if response.status_code == 200:
            return response.json()
        else:
            print(f"Request failed with status code: {response.status_code}")
            return None


if __name__ == "__main__":
    url = "https://jsonplaceholder.typicode.com/todos/1"
    api_request = Request(url)
    response_data = api_request.make_request(RequestType.GET)
    print("GET Request Response:", response_data)
