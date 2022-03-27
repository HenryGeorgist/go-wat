FROM osgeo/gdal:alpine-normal-3.2.1 as build

COPY --from=golang:1.18-alpine3.14 /usr/local/go/ /usr/local/go/

RUN apk add --no-cache \
    pkgconfig \
    gcc \
    libc-dev \
    git

ENV GOROOT=/usr/local/go
ENV GOPATH=/go
ENV GO111MODULE="on"
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

WORKDIR /workspaces