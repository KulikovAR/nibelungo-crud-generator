.PHONY: build run docker-build docker-run

# Переменные
BINARY_NAME=generator
DOCKER_IMAGE=go-crud-generator
DOCKER_TAG=latest

# Сборка локально
build:
	go build -o $(BINARY_NAME) ./cmd/generator

# Запуск локально
run: build
	./$(BINARY_NAME)

# Сборка Docker образа
docker-build:
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) -f build/Dockerfile.run .

# Запуск в Docker
docker-run: docker-build
	docker run -v $(PWD):/workspace $(DOCKER_IMAGE):$(DOCKER_TAG) generate config.json

# Очистка
clean:
	rm -f $(BINARY_NAME) 