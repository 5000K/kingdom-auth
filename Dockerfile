FROM golang:1.25.3 AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o kingdomauth ./main

FROM debian:latest

# ca-certificates required for openid
RUN apt-get update && apt-get install -y ca-certificates --no-install-recommends

WORKDIR /app

COPY --from=builder /app/kingdomauth .

CMD ["./kingdomauth"]