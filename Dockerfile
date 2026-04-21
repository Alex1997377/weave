# Используем официальный образ Go
FROM golang:1.23-alpine

# Устанавливаем необходимые инструменты
RUN apk add --no-cache git make bash

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum (для кэширования зависимостей)
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальной код (будет перемонтировано через volume)
COPY . .

# Команда по умолчанию – интерактивная оболочка
CMD ["/bin/sh"]