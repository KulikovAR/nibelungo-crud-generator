# some-project

Автоматически сгенерированный CRUD сервис на Go.

## Возможности

- REST API с автоматической документацией
- gRPC API
- Поддержка PostgreSQL и MongoDB
- Автоматические тесты
- Docker окружение
- Миграции базы данных

## Быстрый старт

### С Docker

```bash
docker-compose up -d
```

Сервис будет доступен на порту 8080

### Без Docker

1. Установите зависимости:
```bash
go mod download
```

2. Настройте базу данных

3. Запустите сервис:
```bash
go run cmd/server/main.go
```

## API Endpoints

### user

- POST /api/v1/users - Создать user
- GET /api/v1/users/:id - Получить user по ID
- PUT /api/v1/users/:id - Обновить user
- DELETE /api/v1/users/:id - Удалить user
- GET /api/v1/users - Список всех user

### product

- POST /api/v1/products - Создать product
- GET /api/v1/products/:id - Получить product по ID
- PUT /api/v1/products/:id - Обновить product
- DELETE /api/v1/products/:id - Удалить product
- GET /api/v1/products - Список всех product

## Тестирование

```bash
go test ./tests/...
```

## Документация API

Swagger UI доступен по адресу: http://localhost:8080/swagger/index.html
