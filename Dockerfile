FROM golang:1.20.5-alpine3.18 AS builder
LABEL maintainer="Tom Keur <mail@tomkeur.net>"

ENV CGO_ENABLED=0
RUN apk add --no-cache git

WORKDIR /tmp/mysql-to-strict

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go test -v
RUN go build -ldflags="-s -w" -o ./mysql-to-strict .

FROM scratch

COPY --from=builder /tmp/mysql-to-strict/mysql-to-strict /usr/local/bin/mysql-to-strict
