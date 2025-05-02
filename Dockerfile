
# Этап 1: билд (builder)
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache ca-certificates git

WORKDIR /app

COPY go.mod go.sum ./

RUN update-ca-certificates
RUN go mod download

COPY . .
RUN go build -o proxy-server ./cmd/main.go


# Этап 2: финальный образ
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/proxy-server .
COPY config.yaml .

EXPOSE 8080

CMD ["./proxy-server"]

