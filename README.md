# Go-CRUD Generator

Генератор CRUD-проектов на Go, который автоматически создает структуру проекта, доменные модели, репозитории, usecase, контроллеры и миграции на основе JSON-конфигурации.

## Установка

1. Клонируй репозиторий:
   ```sh
   git clone https://github.com/KulikovAR/go-crud-generator.git
   cd go-crud-generator
   ```

2. Собери генератор:
   ```sh
   go build -o generator ./cmd/generator
   ```

## Использование

1. Создай JSON-конфигурацию (например, `config.json`):
   ```json
   {
     "name": "example-project",
     "module": "github.com/KulikovAR/example-project",
     "entities": [
       {
         "name": "User",
         "fields": [
           { "name": "ID", "type": "int64" },
           { "name": "Username", "type": "string" },
           { "name": "Email", "type": "string" },
           { "name": "CreatedAt", "type": "string" }
         ]
       }
     ],
     "repositories": ["postgres", "mongodb"],
     "features": {
       "rest": true,
       "grpc": false,
       "migrations": true
     }
   }
   ```

2. Запусти генератор:
   ```sh
   ./generator generate config.json
   ```

3. В результате появится папка `example-project` с готовой структурой:
   ```
   example-project/
     go.mod
     cmd/
     internal/
       domain/
       repository/
       usecase/
       controller/
     pkg/api/
     migrations/postgres/
     migrations/mongodb/
   ```

## Особенности

- Автоматическая генерация доменных моделей, репозиториев, usecase, контроллеров и миграций.
- Поддержка PostgreSQL и MongoDB.
- Гибкая конфигурация через JSON.

## Примеры

### Генерация проекта с одной сущностью

```json
{
  "name": "my-service",
  "module": "github.com/KulikovAR/my-service",
  "entities": [
    {
      "name": "Product",
      "fields": [
        { "name": "ID", "type": "int64" },
        { "name": "Name", "type": "string" },
        { "name": "Price", "type": "float64" }
      ]
    }
  ],
  "repositories": ["postgres"],
  "features": {
    "rest": true,
    "grpc": false,
    "migrations": true
  }
}
```

### Генерация проекта с несколькими сущностями

```json
{
  "name": "shop",
  "module": "github.com/KulikovAR/shop",
  "entities": [
    {
      "name": "Product",
      "fields": [
        { "name": "ID", "type": "int64" },
        { "name": "Name", "type": "string" },
        { "name": "Price", "type": "float64" }
      ]
    },
    {
      "name": "Order",
      "fields": [
        { "name": "ID", "type": "int64" },
        { "name": "ProductID", "type": "int64" },
        { "name": "Quantity", "type": "int" }
      ]
    }
  ],
  "repositories": ["postgres", "mongodb"],
  "features": {
    "rest": true,
    "grpc": true,
    "migrations": true
  }
}
```

## Лицензия

MIT 