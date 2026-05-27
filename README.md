# 🏨 Hotel Booking — Backend

Два микросервиса на **Go 1.21**, чистая архитектура, PostgreSQL (без ORM), миграции, Swagger, zerolog, тесты с моками.

## Архитектура

```
hotel-backend/
├── common/             # Общий модуль: JWT, logger, response
├── auth-service/       # Порт 8081 — регистрация, вход, /me
├── booking-service/    # Порт 8082 — номера, брони, статусы
└── docker-compose.yml
```

### Слои (Clean Architecture)
```
Handler (delivery/http)  →  Usecase (бизнес-логика)  →  Repository (PostgreSQL)
                                                             ↑
                                                         Domain (сущности, ошибки)
```

## Быстрый старт

### Docker (рекомендуется)
```bash
make up
# auth-service:    http://localhost:8081/swagger/index.html
# booking-service: http://localhost:8082/swagger/index.html
```

### Локально
```bash
# 1. Запустить PostgreSQL
docker run -d --name authdb    -e POSTGRES_DB=authdb    -e POSTGRES_PASSWORD=postgres -p 5433:5432 postgres:16-alpine
docker run -d --name bookingdb -e POSTGRES_DB=bookingdb -e POSTGRES_PASSWORD=postgres -p 5434:5432 postgres:16-alpine

# 2. Генерация Swagger (нужен: go install github.com/swaggo/swag/cmd/swag@latest)
make swag-auth
make swag-booking

# 3. go mod tidy
make tidy

# 4. Запуск
cd auth-service    && go run ./cmd/server
cd booking-service && go run ./cmd/server
```

## Миграции

Применяются **автоматически** при старте сервиса (`golang-migrate`).

| Сервис          | # | Миграция                        |
|-----------------|---|---------------------------------|
| auth-service    | 1 | create_users_table              |
| auth-service    | 2 | create_refresh_tokens           |
| auth-service    | 3 | add_user_indexes                |
| booking-service | 1 | create_rooms                    |
| booking-service | 2 | create_bookings                 |
| booking-service | 3 | seed_rooms + booking_indexes    |

## API

### Auth Service `localhost:8081/api/v1`

| Метод | Путь              | Доступ | Описание                     |
|-------|-------------------|--------|------------------------------|
| POST  | /auth/register    | public | Регистрация → JWT            |
| POST  | /auth/login       | public | Вход → JWT                   |
| GET   | /auth/me          | 🔐     | Профиль текущего пользователя |

### Booking Service `localhost:8082/api/v1`

| Метод  | Путь                      | Доступ      | Описание                                  |
|--------|---------------------------|-------------|-------------------------------------------|
| GET    | /rooms                    | public      | Список номеров                            |
| GET    | /rooms/:id                | public      | Номер по ID                               |
| POST   | /bookings                 | 🔐 guest    | Создать бронь (status: pending)           |
| GET    | /bookings/my              | 🔐 guest    | Мои брони                                 |
| PUT    | /bookings/:id             | 🔐 guest    | Редактировать бронь (status: pending_edit) |
| DELETE | /bookings/:id             | 🔐 guest    | Удалить свою бронь                        |
| GET    | /bookings                 | 🔐 admin    | Все брони                                 |
| PATCH  | /bookings/:id/status      | 🔐 admin    | Изменить статус (confirmed/rejected/...)  |
| DELETE | /bookings/:id             | 🔐 admin    | Удалить любую бронь                       |

## Статусы брони

```
pending  →  confirmed  (admin подтверждает)
         →  rejected   (admin отклоняет)
confirmed → pending_edit (client редактирует)
pending_edit → confirmed (admin подтверждает изменение)
             → confirmed (admin отклоняет — возвращает confirmed без изменений)
```

## Тесты

```bash
make test-cover
```

Тесты покрывают:
- `usecase` — регистрация, вход, ошибки (auth); создание, редактирование, статусы, права (booking)
- `delivery/http` — все эндпоинты: 200/201/401/403/404/409

## Логирование (zerolog)

Каждый сервис логирует в stdout в формате JSON:
```json
{"level":"info","service":"booking-service","time":"2024-01-15T10:30:00Z","booking_id":"abc","msg":"booking created"}
```

## Переменные окружения

| Переменная       | Default                                      | Описание           |
|------------------|----------------------------------------------|--------------------|
| HTTP_PORT        | 8081 / 8082                                  | Порт сервиса       |
| DATABASE_URL     | postgres://postgres:postgres@localhost/...   | DSN PostgreSQL     |
| JWT_SECRET       | super-secret-key-change-in-production        | Секрет для JWT     |
| MIGRATIONS_PATH  | file://migrations                            | Путь к миграциям   |
