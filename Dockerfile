
FROM golang:1.24 AS builder

WORKDIR /app

# Копируем модули и скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -o todo-app main.go

# Финальный образ
FROM debian:bullseye-slim

WORKDIR /app

# Устанавливаем CA-сертификаты для работы с HTTPS
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Копируем бинарник из builder
COPY --from=builder /app/todo-app .
# Копируем папку с фронтендом
COPY --from=builder /app/web ./web

# Порт приложения
EXPOSE 7540

# Запускаем
CMD ["./todo-app"]