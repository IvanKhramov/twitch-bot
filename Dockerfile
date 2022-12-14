# syntax=docker/dockerfile:1

FROM golang:1.19-alpine

WORKDIR /app

COPY . .
RUN go mod download


RUN go build -o twitch-bot


CMD [ "/app/twitch-bot", "start" ]