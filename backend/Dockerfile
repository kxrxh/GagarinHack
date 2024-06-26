FROM golang:latest AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Copy app.json file from the root directory
COPY app.json ./
# COPY database/sql/schema.sql ./

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o main .

FROM alpine:latest
# RUN apk --no-cache add ca-certificates

COPY --from=builder ["/build/main", "/"]
COPY --from=builder ["/build/app.json", "/"]
# COPY --from=builder ["/build/database/sql/schema.sql", "/"]
CMD ["openssl req -config example-com.conf -new -x509 -sha256 -newkey rsa:2048 -nodes \
    -keyout example-com.key.pem -days 365 -out example-com.cert.pem"]

ENTRYPOINT ["./main"]