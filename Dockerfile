FROM golang:1.23-alpine AS builder

WORKDIR /app

# Копируем файлы с зависимостями
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./server/main.go

# Используем минимальный образ для запуска
FROM alpine:latest

WORKDIR /app

# Устанавливаем необходимые зависимости для работы приложения
RUN apk --no-cache add ca-certificates

# Копируем собранное приложение из предыдущего этапа
COPY --from=builder /app/main .

# Открываем порт (порт 8082 используется в main.go)
EXPOSE 8081

# Запускаем приложение
CMD ["./main"] 