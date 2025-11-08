FROM golang:1.25.3 AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o kingdomauth ./main

FROM alpine:latest

COPY --from=builder /kingdomauth .

CMD ["/kingdomauth"]