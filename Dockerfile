# syntax=docker/dockerfile:1

## Build
FROM golang:1.19-bullseye AS build

WORKDIR /app
COPY . /app

RUN go mod tidy && go build -o /battery-history cmd/main.go

## Deploy
FROM gcr.io/distroless/base-debian11:latest

WORKDIR /

COPY --from=build /battery-history /battery-history

EXPOSE 9000

ENTRYPOINT ["/battery-history"]