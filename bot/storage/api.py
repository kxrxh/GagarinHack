import os
from typing import Tuple
import requests


class KeyValueStorage:
    # Initialize the base URL for the key-value storage from the environment variable
    base_url = os.environ.get('STORAGE_URL')
    if not base_url:
        raise RuntimeError('STORAGE_URL environment variable not set')
    if not base_url.endswith('/'):
        base_url += '/'

    @staticmethod
    def get(key: str) -> Tuple[bool, str]:
        """
        Retrieve the value associated with the given key from the key-value storage.

        Args:
            key (str): The key to retrieve the value for.

        Returns:
            Tuple[bool, str]: A tuple indicating whether the operation was successful and the retrieved value (if successful).
        """
        # Send a GET request to retrieve the value associated with the key
        response = requests.get(KeyValueStorage.base_url + key)
        if response.status_code == 200:
            data = response.json()
            return (True, data.get('value', None))
        else:
            return (False, f"Error: {response.status_code}: {response.text}")

    @staticmethod
    def set(key: str, value: str) -> Tuple[bool, str]:
        """
        Store the given key-value pair in the key-value storage.

        Args:
            key (str): The key to store.
            value (str): The value associated with the key.

        Returns:
            Tuple[bool, str]: A tuple indicating whether the operation was successful and a message describing the result.
        """
        # Send a POST request to store the key-value pair
        response = requests.post(KeyValueStorage.base_url, json={
                                 "key": key, "value": value})
        if response.status_code == 200:
            return (True, "Key-value pair stored successfully")
        else:
            return (False, f"Error: {response.status_code}: {response.text}")

    @staticmethod
    def delete(key: str) -> Tuple[bool, str]:
        """
        Delete the key-value pair associated with the given key from the key-value storage.

        Args:
            key (str): The key to delete.

        Returns:
            Tuple[bool, str]: A tuple indicating whether the operation was successful and a message describing the result.
        """
        # Send a DELETE request to delete the key-value pair
        response = requests.delete(KeyValueStorage.base_url + key)
        if response.status_code == 200:
            return (True, "Key deleted successfully")
        else:
            return (False, f"Error: {response.status_code}: {response.text}")


if __name__ == '__main__':
    # Example usage of the KeyValueStorage class
    # Retrieving the value associated with the key 'foo'
    print(KeyValueStorage.get('foo'))

    # Storing the key-value pair ('foo', 'bar')
    print(KeyValueStorage.set('foo', 'bar'))

    # Retrieving the value associated with the key 'foo' after storing it
    print(KeyValueStorage.get('foo'))

    # Deleting the key-value pair associated with the key 'foo'
    print(KeyValueStorage.delete('foo'))

    # Retrieving the value associated with the key 'foo' after deleting it
    print(KeyValueStorage.get('foo'))
