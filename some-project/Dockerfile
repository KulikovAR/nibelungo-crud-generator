FROM golang:1.24-alpine AS builder

WORKDIR /app

# Установка зависимостей
RUN apk add --no-cache git

# Копирование go mod файлов
COPY go.mod go.sum ./
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копирование бинарного файла
COPY --from=builder /app/main .

# Копирование миграций
COPY --from=builder /app/migrations ./migrations

# Открытие порта
EXPOSE 8080

# Запуск приложения
CMD ["./main"]
