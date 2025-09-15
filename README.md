# 📱 Subscription Service API

REST-сервис для агрегации данных об онлайн-подписках пользователей.

**Тестовое задание Junior Golang Developer — Effective Mobile**

## 🚀 Быстрый старт

```bash
# Клонирование репозитория
git clone https://github.com/BrikozO/testTaskEffectiveMobile.git
cd testTaskEffectiveMobile

# Создание .env файла (см. env.example)
cp env.example .env

# Запуск сервиса
docker compose up --build -d
```

**Swagger документация:** `http://localhost:8080/swagger/index.html`

## 📊 API Endpoints

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `POST` | `/api/v1/subscriptions` | Создать подписку |
| `GET` | `/api/v1/subscriptions/{user_id}` | Получить подписки пользователя |
| `GET` | `/api/v1/subscriptions/{user_id}/{subscription_id}` | Получить конкретную подписку |
| `PUT` | `/api/v1/subscriptions/{subscription_id}` | Обновить подписку |
| `DELETE` | `/api/v1/subscriptions/{subscription_id}` | Удалить подписку |
| `POST` | `/api/v1/calculate` | Рассчитать суммарную стоимость |

## 🔧 Структура данных

**Модель подписки:**
```json
{
  "service_name": "Netflix",
  "price": 999,
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "start_date": "01-2024",
  "end_date": "12-2024"
}
```

**Расчет стоимости (POST /api/v1/calculate):**
```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "service_name": "Netflix",
  "start_date": "01-2024", 
  "end_date": "12-2024"
}
```

> 💡 **Важно:** В ответах API дополнительно возвращается поле `id` записи из БД, необходимое для операций обновления и удаления конкретных подписок.

## 🛠 Технический стек

- **Go 1.25** — основной язык
- **PostgreSQL** — база данных
- **Swagger** — документация API
- **Docker Compose** — контейнеризация

## 📁 Архитектура проекта

```
├── main.go                    # Точка входа
├── handlers.go                # HTTP обработчики
├── routes.go                  # Маршрутизация
├── middlewares.go             # Middleware (логирование)
├── helpers.go                 # Вспомогательные функции
├── models/subscription.go     # Модели данных
├── dto/                       # Data Transfer Objects
├── postgres_db/               # Работа с БД
│   ├── connector.go           # Подключение к PostgreSQL
│   ├── repositories/          # Репозитории
│   └── migrations/v1/         # Миграции
├── docs/                      # Swagger документация
└── docker-compose.yml         # Конфигурация контейнеров
```

## 🔄 Особенности реализации

- **Автоматические миграции** при запуске приложения
- **UUID для пользователей** и автоинкремент ID для подписок
- **Формат дат MM-YYYY** (например, "01-2024")
- **Структурированное логирование** всех HTTP запросов
- **Graceful error handling** с соответствующими HTTP статусами

## 🎯 Примеры использования

**Создание подписки:**
```bash
curl -X POST http://localhost:8080/api/v1/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Spotify",
    "price": 299,
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "start_date": "03-2025"
  }'
```

**Расчет суммарной стоимости:**
```bash
curl -X POST http://localhost:8080/api/v1/calculate \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "start_date": "01-2025",
    "end_date": "12-2025"
  }'
```

## ⚙️ Конфигурация

Переменные окружения в `.env`:
```env
POSTGRES_HOST=postgres_db
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=your_password
POSTGRES_DB=subscriptions_db
```

## 👨‍💻 Автор

**Олег Якушев** — [GitHub](https://github.com/BrikozO) | [Email](mailto:oleg.yakushev.work@gmail.com)