FROM golang:1.23-alpine

# Устанавливаем зависимости
RUN apk add --no-cache git ca-certificates

# Рабочая директория
WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходники
COPY . .

# Собираем бинарник
RUN go build -o server ./cmd/server

# Запускаем
CMD ["./server"]