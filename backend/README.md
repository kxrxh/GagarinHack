# README for API Version 1 Documentation

## Overview
This README details the configuration, deployment, and usage of the API Version 1 developed in Go using the Fiber framework. The documentation is crafted to guide through setup, endpoint functionalities, and operational details.

## Prerequisites
- Docker installed on your system.
- Go environment set up for local development if necessary.

## Build Instructions
The application is containerized using Docker for consistent deployment environments. Below are the steps to build and run the container:

### Docker Build
1. **Setting Up the Build Environment**
   - Start with a base image from the official Golang repository:
     ```dockerfile
     FROM golang:latest AS builder
     ```

2. **Preparing the Build Directory**
   - Set the working directory and copy necessary Go module files:
     ```dockerfile
     WORKDIR /build
     COPY go.mod go.sum ./
     RUN go mod download
     ```

3. **Copying Application Files**
   - Copy all source files and the application configuration:
     ```dockerfile
     COPY . .
     COPY app.json ./
     ```

4. **Building the Application**
   - Environment setup and compilation:
     ```dockerfile
     ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
     RUN go build -ldflags="-s -w" -o main .
     ```

5. **Preparing the Runtime Environment**
   - Using a lightweight Alpine Linux image:
     ```dockerfile
     FROM alpine:latest
     ```

6. **Copying the Built Executable and Configurations**
   - Transfer the compiled binary and config files:
     ```dockerfile
     COPY --from=builder ["/build/main", "/"]
     COPY --from=builder ["/build/app.json", "/"]
     ```

7. **Configuring SSL/TLS**
   - Optional: Generate SSL certificates if needed for secure connections:
     ```dockerfile
     CMD ["openssl req -config example-com.conf -new -x509 -sha256 -newkey rsa:2048 -nodes \
          -keyout example-com.key.pem -days 365 -out example-com.cert.pem"]
     ```

8. **Setting the Entry Point**
   - Specify the application's entry point:
     ```dockerfile
     ENTRYPOINT ["./main"]
     ```

### Running the Container
Execute the following command to build and run the Docker container:
```bash
docker build -t myapi:v1 .
docker run -p 3033:3033 myapi:v1
```

## API Endpoints
Documentation for the endpoints provided in Version 1 of the API, including routes for user authentication, data completions, and external API integrations.

### Route Initialization
```go
func SetupRoutesV1(v1 *fiber.Router)
```
Sets up all the API endpoints under the `/v1` base path, organizing the services offered by this version.

### Completion Requests
Various completion endpoints are grouped under `/v1/completion`. These include:
- `POST /yandex`
- `POST /yandex/questions`
- `POST /yandex/questions/biography`
- `POST /yandex/epitaph`
- `POST /yandex/biography`
- `POST /yandex/biography/short`
- `POST /gigachat`
- `POST /gigachat/questions`
- `POST /gigachat/questions/biography`
- `POST /gigachat/epitaph`
- `POST /gigachat/biography`
- `POST /gigachat/biography/short`

### Completion Services
Endpoints under `/v1/completion` handle various data generation requests, categorized by type such as `yandex` and `gigachat`.

### Token Services
External token retrieval is handled by:
```go
POST /v1/external/get-access-token
```
Used for authenticating with external APIs by obtaining necessary access tokens.

## Configuration Details
Configuration files like `app.json` should be set according to your deployment environment specifics. Parameters include port settings, API keys, and external URLs.

## Error Handling
Proper error responses are configured to provide clear and actionable information, aiding in troubleshooting and integration efforts.

## Logging
The API uses structured logging to facilitate easy monitoring and debugging. Configuration toggles between development and production log settings as specified in `app.json`.

## Running the Application
To start the server, execute the following command:
```bash
go run main.go
```
