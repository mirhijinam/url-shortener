FROM golang:1.20.4-alpine as builder 

WORKDIR /go/src/url-shortener

RUN apk update && \
    apk upgrade && \
    apk add \
        --no-cache \
            git \
            gcc \
            libc-dev

COPY ./go.* ./
RUN go mod download

COPY .env .env
COPY . .
RUN go build -o .bin/main ./cmd/url-shortener


FROM alpine

WORKDIR /go/src/url-shortener

RUN adduser -D appuser
USER appuser

COPY --from=builder /go/src/url-shortener/.bin/main ./.bin/main
COPY --from=builder /go/src/url-shortener/config ./config
COPY --from=builder /go/src/url-shortener/.env .env


CMD ["./.bin/main"]
