# syntax=docker/dockerfile:1

## Build
FROM golang:1.19-bullseye AS build

# Copy files to WORKDIR
WORKDIR /app
COPY . /app

# Build BatteryHistory
RUN go mod tidy
RUN go build -o /battery-history cmd/main.go

## Deploy
FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /battery-history /battery-history

EXPOSE 9000

ENTRYPOINT ["/battery-history"]