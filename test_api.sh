#!/bin/bash

# Тестовый скрипт для проверки API микросервиса

BASE_URL="http://localhost:8080"

echo "🧪 Тестирование Crypto Microservice API"
echo "========================================"

# Проверяем, что сервер запущен
echo "1. Проверка доступности сервера..."
if curl -s "$BASE_URL/swagger/index.html" > /dev/null; then
    echo "✅ Сервер доступен"
else
    echo "❌ Сервер недоступен. Убедитесь, что он запущен на порту 8080"
    exit 1
fi

echo ""
echo "2. Добавление Bitcoin в список наблюдения..."
curl -X POST "$BASE_URL/currency/add" \
  -H "Content-Type: application/json" \
  -d '{"symbol": "bitcoin"}' \
  -w "\nHTTP Status: %{http_code}\n"

echo ""
echo "3. Добавление Ethereum в список наблюдения..."
curl -X POST "$BASE_URL/currency/add" \
  -H "Content-Type: application/json" \
  -d '{"symbol": "ethereum"}' \
  -w "\nHTTP Status: %{http_code}\n"

echo ""
echo "4. Получение цены Bitcoin (текущее время)..."
TIMESTAMP=$(date +%s)
curl "$BASE_URL/currency/price?coin=bitcoin&timestamp=$TIMESTAMP" \
  -w "\nHTTP Status: %{http_code}\n"

echo ""
echo "5. Получение цены Ethereum (текущее время)..."
curl "$BASE_URL/currency/price?coin=ethereum&timestamp=$TIMESTAMP" \
  -w "\nHTTP Status: %{http_code}\n"

echo ""
echo "6. Удаление Bitcoin из списка наблюдения..."
curl -X DELETE "$BASE_URL/currency/remove" \
  -H "Content-Type: application/json" \
  -d '{"symbol": "bitcoin"}' \
  -w "\nHTTP Status: %{http_code}\n"

echo ""
echo "✅ Тестирование завершено!"
echo ""
echo "📖 Документация API доступна по адресу: $BASE_URL/swagger/index.html" 