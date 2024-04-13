from dataclasses import dataclass, field
import logging
import os
import requests


@dataclass
class StorageResponse:
    success: bool = field(default_factory=bool)
    value: str = field(default_factory=str)
    message: str = field(default_factory=str)


class KeyValueStorage:
    # Initialize the base URL for the key-value storage from the environment variable
    __base_url = os.environ.get('STORAGE_URL')
    if not __base_url:
        raise RuntimeError('STORAGE_URL environment variable not set')
    if not __base_url.endswith('/'):
        __base_url += '/'

    @staticmethod
    def get(key: str) -> StorageResponse:
        """
        Retrieve the value associated with the given key from the key-value storage.

        Args:
            key (str): The key to retrieve the value for.

        Returns:
            StorageResponse: A tuple indicating whether the operation was successful and the retrieved value (if successful).
        """
        # Send a GET request to retrieve the value associated with the key
        if not key:
            raise ValueError("Key cannot be empty")
        response = requests.get(KeyValueStorage.__base_url + key)
        if response.status_code == 200:
            data = response.json()
            return StorageResponse(success=True, value=data.get('value', None))
        else:
            logging.error("error: " + response.text)
            return StorageResponse(success=False, message="error: " + response.text)

    @staticmethod
    def set(key: str, value: str) -> StorageResponse:
        """
        Store the given key-value pair in the key-value storage.

        Args:
            key (str): The key to store.
            value (str): The value associated with the key.

        Returns:
            StorageResponse: A tuple indicating whether the operation was successful and a message describing the result.
        """
        # Send a POST request to store the key-value pair
        response = requests.post(KeyValueStorage.__base_url, json={
                                 "key": str(key), "value": value})
        if response.status_code == 200:
            return StorageResponse(success=True, message="Key-value pair stored successfully")
        else:
            logging.error("error: " + response.text)
            return StorageResponse(success=False, message="error: " + response.text)

    @staticmethod
    def delete(key: str) -> StorageResponse:
        """
        Delete the key-value pair associated with the given key from the key-value storage.

        Args:
            key (str): The key to delete.

        Returns:
            StorageResponse: A tuple indicating whether the operation was successful and a message describing the result.
        """
        # Send a DELETE request to delete the key-value pair
        response = requests.delete(KeyValueStorage.__base_url + key)
        if response.status_code == 200:
            return StorageResponse(success=True, message="Key-value pair deleted successfully")
        else:
            logging.error("error: " + response.text)
            return StorageResponse(success=False, message="error: " + response.text)

    @staticmethod
    def delete_prefix(prefix: str) -> StorageResponse:
        """
        Delete all key-value pairs whose keys start with the given prefix from the key-value storage.
                Args:
                prefix (str): The prefix to delete.
        Returns:
        StorageResponse: A tuple indicating whether the operation was successful and a message describing the result.
        """

        response = requests.delete(
            f"{KeyValueStorage.__base_url}prefix/{prefix}")
        if response.status_code == 200:
            return StorageResponse(success=True, message="Key-value pairs deleted successfully")
        else:
            logging.error("error: " + response.text)
            return StorageResponse(success=False, message="error: " + response.text)


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
