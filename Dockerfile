# syntax=docker/dockerfile:1

FROM golang:1.24-alpine AS builder
WORKDIR /app

#odules caching
COPY go.mod go.sum ./
ENV GOTOOLCHAIN=auto
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM alpine:3.19
WORKDIR /app

#server binary
COPY --from=builder /app/server /usr/local/bin/server

# copy config for viper
COPY config/ /app/config/

RUN apk --no-cache add ca-certificates tzdata

EXPOSE 8080
CMD ["server"]
