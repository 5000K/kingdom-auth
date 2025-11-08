FROM golang:1.25.3 AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -o kingdomauth ./main

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/kingdomauth .

CMD ["./kingdomauth"]