FROM golang:1.18.0-alpine3.14

WORKDIR /workspaces

RUN apk add --no-cache \
    pkgconfig \
    gcc \
    libc-dev \
    git