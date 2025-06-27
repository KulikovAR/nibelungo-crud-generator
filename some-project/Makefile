# Makefile для some-project

.PHONY: build run test clean docker-build docker-run

# Сборка приложения
build:
	go build -o bin/server cmd/server/main.go

# Запуск приложения
run:
	go run cmd/server/main.go

# Запуск тестов
test:
	go test -v ./tests/...

# Очистка
clean:
	rm -rf bin/
	go clean

# Сборка Docker образа
docker-build:
	docker build -t some-project .

# Запуск в Docker
docker-run:
	docker-compose up -d

# Остановка Docker
docker-stop:
	docker-compose down

# Просмотр логов
logs:
	docker-compose logs -f

# Миграции
migrate-up:
	migrate -path migrations/postgres -database "postgres://postgres:password@localhost:5432/some-project?sslmode=disable" up

migrate-down:
	migrate -path migrations/postgres -database "postgres://postgres:password@localhost:5432/some-project?sslmode=disable" down
