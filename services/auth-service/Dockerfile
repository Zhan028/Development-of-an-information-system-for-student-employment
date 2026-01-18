# Этап сборки
FROM golang:1.23-alpine AS builder

# Установка необходимых пакетов для сборки
RUN apk add --no-cache git ca-certificates tzdata

# Установка рабочей директории
WORKDIR /app

# Копирование файлов зависимостей
COPY go.mod go.sum ./

# Загрузка зависимостей
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /auth-service ./cmd/main.go

# Этап запуска
FROM alpine:3.19

# Установка сертификатов CA и временных зон
RUN apk --no-cache add ca-certificates tzdata

# Создание непривилегированного пользователя
RUN adduser -D -g '' appuser

# Установка рабочей директории
WORKDIR /app

# Копирование бинарного файла из этапа сборки
COPY --from=builder /auth-service .

# Смена владельца файлов
RUN chown -R appuser:appuser /app

# Переключение на непривилегированного пользователя
USER appuser

# Порт приложения
EXPOSE 8081

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8081/health || exit 1

# Точка входа
ENTRYPOINT ["./auth-service"]
