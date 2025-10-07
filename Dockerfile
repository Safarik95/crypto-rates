# Используем версию Go 1.24
FROM golang:1.24-alpine AS builder

# Устанавливаем зависимости для сборки
RUN apk add --no-cache git

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем ВСЕ файлы проекта
COPY . .

# Скачиваем зависимости
RUN go mod download

# Собираем приложение (правильный путь)
RUN go build -o server ./cmd/service

# Финальный образ
FROM alpine:latest

# Устанавливаем зависимости для runtime
RUN apk add --no-cache ca-certificates tzdata

# Создаем пользователя для безопасности
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем бинарник из стадии builder
COPY --from=builder /app/server .

# Копируем миграции
COPY --from=builder /app/migrations ./migrations

# Меняем владельца файлов
RUN chown -R appuser:appgroup /app

# Переключаемся на непривилегированного пользователя
USER appuser

# Экспортируем порт
EXPOSE 8080

# Запускаем приложение
CMD ["./server"]