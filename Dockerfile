FROM golang:1.22-alpine as runner-builder
MAINTAINER isKONSTANTIN <me@knst.su>

ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /build/

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

RUN go build ./cmd/telegram-gpt


FROM golang:1.22-alpine as migrator-builder
MAINTAINER isKONSTANTIN <me@knst.su>

ENV GOOS=linux
ENV GOARCH=amd64

RUN apk update
RUN apk add git

WORKDIR /

RUN git clone -b v3.18.0 https://github.com/pressly/goose

WORKDIR /goose

RUN go mod download
RUN go build ./cmd/goose


FROM alpine:3.14 as runner
MAINTAINER isKONSTANTIN <me@knst.su>

WORKDIR /gpt-bot

COPY --from=runner-builder /build/telegram-gpt ./

ENV GOMAXPROCS=2

ENTRYPOINT ["./telegram-gpt"]


FROM alpine:3.14 as migrator
MAINTAINER isKONSTANTIN <me@knst.su>

WORKDIR /migrator

COPY build/migrate.sh ./
COPY migrations/* migrations/
COPY --from=migrator-builder /goose/goose ./

ENTRYPOINT ["./migrate.sh"]