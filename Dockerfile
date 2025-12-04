FROM golang:1.23-alpine

WORKDIR /app

# 1. Копируем go.mod/go.sum и тянем зависимости (кэшируется)
COPY go.mod go.sum ./
RUN go mod download

# 2. Копируем остальной код
COPY . .

# 3. Сборка бинаря из main.go в корне
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./main.go

# 4. netcat и wait-for-it
RUN apk add --no-cache netcat-openbsd bash

COPY wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

EXPOSE 3001

CMD ["/wait-for-it.sh", "deposits-db:5432", "--", "./server"]
