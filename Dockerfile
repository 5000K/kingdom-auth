FROM golang:1.25.1 AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o kingdom-auth ./main

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/kingdom-auth .

ENV CONFIG_PATH=/app/config.yml

CMD ["./kingdom-auth"]