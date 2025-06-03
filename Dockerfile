# Этап 1: Сборка приложения
FROM golang:latest AS builder

# Установка рабочей директории
WORKDIR /doc

# Копирование всех файлов проекта
COPY . .

# Загрузка зависимостей
RUN go mod tidy && go mod download

# Сборка приложения (основной файл находится в cmd/main.go)
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd

# Этап 2: Финальный образ
FROM alpine:latest

# Установка сертификатов
RUN apk --no-cache add ca-certificates

# Установка рабочей директории
WORKDIR /root/

# Копирование бинарника из первого этапа
COPY --from=builder /doc/app .

# Команда запуска
CMD ["./app"]