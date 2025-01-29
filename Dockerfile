# Используем официальный образ Golang с версией 1.22
FROM golang:1.22-alpine AS builder

# Устанавливаем зависимости для работы с git и PostgreSQL
RUN apk add --no-cache git build-base

# Устанавливаем утилиту goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы для кэширования зависимостей
COPY go.mod go.sum ./

# Устанавливаем зависимости Go
RUN go mod tidy

# Копируем исходный код проекта
COPY . .

# Сборка приложения
RUN go build -o /app/main ./cmd

# Используем легковесный образ для финального контейнера
FROM alpine:latest

# Устанавливаем необходимые зависимости
RUN apk add --no-cache libpq bash curl

# Копируем собранное приложение
COPY --from=builder /app /app

# Копируем миграции
COPY db/migrations /app/migrations

# Копируем утилиту goose
COPY --from=builder /go/bin/goose /usr/local/bin/goose

# Копируем скрипт wait-for-it
COPY wait-for-it.sh /usr/local/bin/wait-for-it

# Делаем wait-for-it исполняемым
RUN chmod +x /usr/local/bin/wait-for-it

# Устанавливаем рабочую директорию
WORKDIR /app

# Открываем порт
EXPOSE 8080

# Запуск миграций и приложения
CMD ["sh", "-c", "/usr/local/bin/wait-for-it db:5432 --strict --timeout=30 -- goose -dir /app/migrations postgres \"postgres://postgres:postgres@db:5432/testdb?sslmode=disable\" up && /app/main"]


