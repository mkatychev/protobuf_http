FROM golang:1.14

ENV GOFLAGS=-mod=vendor
ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /src/notebook/
COPY . .

ENTRYPOINT go run .
