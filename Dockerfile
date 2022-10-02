# syntax=docker/dockerfile:1

## Build
FROM golang:1.19-bullseye AS build

WORKDIR /app
COPY . /app

RUN go mod tidy && go build -o /bthist cmd/main.go

## Deploy
FROM gcr.io/distroless/base-debian11:latest

WORKDIR /app

COPY --from=build /bthist /app/bthist

EXPOSE 9000

ENTRYPOINT ["./bthist"]