# Todo Application

Приложение для управления задачами с поддержкой повторений.

## Технологии

- Go 1.24
- PostgreSQL
- Docker / Docker Compose

## Запуск

### Локально с PostgreSQL:

```bash
# Запустить PostgreSQL
docker-compose up -d postgres

# Запустить приложение
go run main.go