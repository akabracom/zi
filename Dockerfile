FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# Настройка Go proxy для обхода проблем с TLS
# Используем direct режим для обхода проблем с proxy
ENV GOPROXY=direct
ENV GOSUMDB=off
ENV CGO_ENABLED=0

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs

EXPOSE 3001

CMD ["./main"]
