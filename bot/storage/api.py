import os
from typing import Tuple
import requests


class KeyValueStorage:
    base_url = os.environ['STORAGE_URL']
    if not base_url:
        raise RuntimeError('STORAGE_URL environment variable not set')
    if base_url[-1] != '/':
        base_url += '/'

    @staticmethod
    def get(key: str) -> Tuple[bool, str]:
        # Send a GET request to retrieve the value associated with the key
        response = requests.get(KeyValueStorage.base_url + key)
        if response.status_code == 200:
            data = response.json()
            return (True, data.get('value', None))
        else:
            return (False, f"{response.status_code}: {response.text}")

    @staticmethod
    def set(key: str, value: str) -> Tuple[bool, str]:
        response = requests.post(KeyValueStorage.base_url, json={
                                 "key": key, "val": value})
        if response.status_code == 200:
            return (True, None)
        else:
            return (False, f"{response.status_code}: {response.text}")

    @staticmethod
    def delete(key: str) -> Tuple[bool, str]:
        # Send a DELETE request to delete the key-value pair
        response = requests.delete(KeyValueStorage.base_url + key)
        if response.status_code == 200:
            return (True, None)
        else:
            return (False, f"{response.status_code}: {response.text}")


if __name__ == '__main__':
    print(KeyValueStorage.get('foo'))
    print(KeyValueStorage.set('foo', 'bar'))
    print(KeyValueStorage.get('foo'))
    print(KeyValueStorage.delete('foo'))
    print(KeyValueStorage.get('foo'))
