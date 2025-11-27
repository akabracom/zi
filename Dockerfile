FROM golang:1.25-alpine

WORKDIR /app

# Копируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект
COPY . .

# Показываем содержимое для отладки
RUN ls -la

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/main ./cmd/awesomeProject/main.go

# Проверяем, создан ли бинарник
RUN ls -la /app

# Делаем исполняемым
RUN chmod +x /app/main

EXPOSE 8080

CMD ["/app/main"]