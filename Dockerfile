FROM golang:1.25.3 AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o kingdomauth ./main

FROM debian:latest

WORKDIR /app

COPY --from=builder /app/kingdomauth .

CMD ["./kingdomauth"]