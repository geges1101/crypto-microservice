# Crypto Microservice

Микросервис для сбора, хранения и отображения стоимости криптовалют.

## Описание

Этот микросервис предоставляет API для работы с криптовалютами. Основная идея - создать систему, которая автоматически отслеживает цены выбранных криптовалют и позволяет получать их исторические данные.

### Основные возможности:
- Добавление криптовалют в список наблюдения
- Удаление криптовалют из списка наблюдения  
- Получение цен криптовалют в определенный момент времени
- Автоматический сбор цен каждые N секунд через CoinGecko API
- Поиск ближайшей цены, если точное время не найдено

Микросервис использует PostgreSQL для хранения данных и Gin для HTTP API. Все цены сохраняются с временными метками для последующего анализа.

## Технологии

- **Go 1.24** - основной язык разработки
- **Gin** - веб-фреймворк
- **GORM** - ORM для работы с базой данных
- **PostgreSQL** - база данных
- **Docker Compose** - для развертывания
- **Swagger** - документация API

## API Endpoints

### 1. Добавить криптовалюту в список наблюдения
```
POST /currency/add
Content-Type: application/json

{
  "symbol": "bitcoin"
}
```

### 2. Удалить криптовалюту из списка наблюдения
```
DELETE /currency/remove
Content-Type: application/json

{
  "symbol": "bitcoin"
}
```

### 3. Получить цену криптовалюты
```
GET /currency/price?coin=bitcoin&timestamp=1736500490
```

## Запуск

### С помощью Docker Compose (рекомендуется)

1. Клонируйте репозиторий:
```bash
git clone <repository-url>
cd crypto-microservice
```

2. Запустите микросервис:
```bash
docker-compose up -d
```

3. Проверьте, что сервис запущен:
```bash
curl http://localhost:8080/swagger/index.html
```

### Локальный запуск

1. Установите зависимости:
```bash
go mod tidy
```

2. Запустите PostgreSQL (или используйте Docker):
```bash
docker run -d --name postgres \
  -e POSTGRES_DB=crypto_db \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=password \
  -p 5432:5432 \
  postgres:15
```

3. Установите переменные окружения:
```bash
export DATABASE_URL="postgres://postgres:password@localhost:5432/crypto_db?sslmode=disable"
export PRICE_UPDATE_INTERVAL=30
```

4. Запустите приложение:
```bash
go run main.go
```

## Переменные окружения

- `DATABASE_URL` - строка подключения к PostgreSQL (по умолчанию: postgres://postgres:password@localhost:5432/crypto_db?sslmode=disable)
- `PRICE_UPDATE_INTERVAL` - интервал обновления цен в секундах (по умолчанию: 30)
- `PORT` - порт для запуска сервера (по умолчанию: 8080)

## Документация API

После запуска сервиса документация Swagger доступна по адресу:
```
http://localhost:8080/swagger/index.html
```

## Примеры использования

### Добавление Bitcoin в список наблюдения:
```bash
curl -X POST http://localhost:8080/currency/add \
  -H "Content-Type: application/json" \
  -d '{"symbol": "bitcoin"}'
```

### Получение цены Bitcoin:
```bash
curl "http://localhost:8080/currency/price?coin=bitcoin&timestamp=1736500490"
```

### Удаление Bitcoin из списка наблюдения:
```bash
curl -X DELETE http://localhost:8080/currency/remove \
  -H "Content-Type: application/json" \
  -d '{"symbol": "bitcoin"}'
```

## Структура проекта

```
crypto-microservice/
├── main.go                 # Точка входа приложения
├── go.mod                  # Зависимости Go
├── docker-compose.yml      # Docker Compose конфигурация
├── Dockerfile             # Docker образ
├── README.md              # Документация
└── internal/
    ├── config/            # Конфигурация
    ├── database/          # Модели и инициализация БД
    ├── handlers/          # HTTP обработчики
    └── services/          # Бизнес-логика
```

## Особенности реализации

1. **Автоматический сбор цен**: Микросервис автоматически собирает цены активных криптовалют каждые N секунд через CoinGecko API
2. **Поиск ближайшей цены**: Если цена в точное время не найдена, возвращается цена в ближайший момент времени (используется SQL ORDER BY с ABS)
3. **Обработка ошибок**: Подробное логирование ошибок API, валидация входных данных, timeout для HTTP запросов
4. **Swagger документация**: Полная документация API доступна через Swagger UI
5. **Docker Compose**: Готовое решение для развертывания с PostgreSQL
6. **Валидация данных**: Проверка символов валют, обработка дубликатов, защита от некорректных данных

### Решенные проблемы:
- **Проблема**: CoinGecko API может быть недоступен → **Решение**: Добавлен timeout и обработка ошибок
- **Проблема**: Дублирование валют при повторном добавлении → **Решение**: Проверка существования и реактивация
- **Проблема**: Удаление валют с связанными ценами → **Решение**: Каскадное удаление через транзакции
- **Проблема**: Неточные временные метки → **Решение**: Поиск ближайшей цены с индексацией

## Поддерживаемые криптовалюты

Микросервис использует CoinGecko API для получения цен. Поддерживаются все криптовалюты, доступные в CoinGecko API.

Примеры популярных символов:
- `bitcoin` - Bitcoin
- `ethereum` - Ethereum
- `binancecoin` - Binance Coin
- `cardano` - Cardano
- `solana` - Solana 