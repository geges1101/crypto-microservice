.PHONY: build run test clean docker-build docker-run

# Переменные
BINARY_NAME=crypto-microservice
DOCKER_IMAGE=crypto-microservice

# Сборка приложения
build:
	@echo "🔨 Сборка приложения..."
	go build -o $(BINARY_NAME) .

# Запуск приложения
run: build
	@echo "🚀 Запуск приложения..."
	./$(BINARY_NAME)

# Тестирование API
test:
	@echo "🧪 Запуск тестов API..."
	./test_api.sh

# Очистка
clean:
	@echo "🧹 Очистка..."
	rm -f $(BINARY_NAME)
	go clean

# Сборка Docker образа
docker-build:
	@echo "🐳 Сборка Docker образа..."
	docker build -t $(DOCKER_IMAGE) .

# Запуск с Docker Compose
docker-run:
	@echo "🐳 Запуск с Docker Compose..."
	docker-compose up -d

# Остановка Docker Compose
docker-stop:
	@echo "🛑 Остановка Docker Compose..."
	docker-compose down

# Просмотр логов
logs:
	@echo "📋 Просмотр логов..."
	docker-compose logs -f

# Установка зависимостей
deps:
	@echo "📦 Установка зависимостей..."
	go mod tidy

# Проверка кода
check:
	@echo "🔍 Проверка кода..."
	go vet ./...
	go fmt ./...

# Помощь
help:
	@echo "Доступные команды:"
	@echo "  build      - Сборка приложения"
	@echo "  run        - Запуск приложения"
	@echo "  test       - Тестирование API"
	@echo "  clean      - Очистка"
	@echo "  docker-build - Сборка Docker образа"
	@echo "  docker-run - Запуск с Docker Compose"
	@echo "  docker-stop - Остановка Docker Compose"
	@echo "  logs       - Просмотр логов"
	@echo "  deps       - Установка зависимостей"
	@echo "  check      - Проверка кода"
	@echo "  help       - Показать эту справку" 