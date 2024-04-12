# KeyValue Database API Documentation

## Introduction
This document outlines the RESTful API endpoints and functionality provided by the KeyValue Database API. The API allows users to store, retrieve, and delete key-value pairs.

## Base URL
The base URL for accessing the API is `http://localhost:5173`.

## Authentication
The API does not currently require authentication.

## Error Handling
Errors are returned in JSON format with appropriate HTTP status codes. The response body contains an `error` field describing the encountered issue.

## API Endpoints

### `GET /db/v1/:id`
- **Description**: Retrieve the value associated with the specified key.
- **Parameters**:
  - `id` (string, required): The unique identifier of the key.
- **Responses**:
  - `200 OK`: Successful retrieval. Returns JSON with the key, value, and existence status.
  - `400 Bad Request`: If the `id` parameter is missing.
  - `404 Not Found`: If the key does not exist in the database.

### `POST /db/v1/`
- **Description**: Store a new key-value pair in the database.
- **Request Body**:
  ```json
  {
    "key": "string",
    "val": "string"
  }
  ```
- **Responses**:
  - `200 OK`: Successful storage. Returns JSON with the stored key and value.
  - `400 Bad Request`: If the request body is missing or invalid.
  - `500 Internal Server Error`: If an error occurs while storing the key-value pair.

### `DELETE /db/v1/:id`
- **Description**: Delete the key-value pair associated with the specified key.
- **Parameters**:
  - `id` (string, required): The unique identifier of the key.
- **Responses**:
  - `200 OK`: Successful deletion.
  - `400 Bad Request`: If the `id` parameter is missing.
  - `500 Internal Server Error`: If an error occurs while deleting the key-value pair.

## Cross-Origin Resource Sharing (CORS)
The API supports Cross-Origin Resource Sharing (CORS) to allow requests from a specified origin. The following CORS settings are applied:
- **Allowed Headers**: Authorization, Origin, Content-Type, Accept, Content-Length, Accept-Language, Accept-Encoding, Connection, Access-Control-Allow-Origin.
- **Allowed Origins**: *
- **Allow Credentials**: true
- **Allowed Methods**: GET, POST, DELETE

## Usage
Below is an example of how to interact with the API endpoints using cURL:

1. **GET Request**:
   ```bash
   curl -X GET http://localhost:5173/db/v1/:id
   ```

2. **POST Request**:
   ```bash
   curl -X POST -H "Content-Type: application/json" -d '{"key": "example_key", "val": "example_value"}' http://localhost:5173/db/v1/
   ```

3. **DELETE Request**:
   ```bash
   curl -X DELETE http://localhost:5173/db/v1/:id
   ```

## Examples

Sure, here are examples of cURL commands with data for each endpoint:

1. **GET Request**:
   ```bash
   curl -X GET http://localhost:5173/db/v1/example_key
   ```

2. **POST Request**:
   ```bash
   curl -X POST -H "Content-Type: application/json" -d '{"key": "example_key", "val": "example_value"}' http://localhost:5173/db/v1/
   ```

3. **DELETE Request**:
   ```bash
   curl -X DELETE http://localhost:5173/db/v1/example_key
   ```

These examples demonstrate how to interact with the API endpoints using cURL commands. Adjust the `example_key` and `example_value` accordingly based on your actual data.