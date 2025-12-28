# Auth Service

Микросервис авторизации для MediaLog. Предоставляет регистрацию, вход, управление сессиями и токенами.

## Описание

Auth Service реализует:
- Регистрацию пользователей
- Аутентификацию (логин)
- Обновление токенов (refresh)
- Выход из сессии (logout)
- Выход со всех устройств (logout-all)
- Rate limiting для защиты от брутфорса
- Хранение сессий в Redis
- JWT access tokens и opaque refresh tokens

## Архитектура

```
Frontend/BFF
    ↓ HTTP/REST (gRPC-Gateway)
Auth Service API Layer
    ↓
Auth Service (Business Logic)
    ↓
Storage Layer
    ├── PostgreSQL (users)
    └── Redis (sessions, rate limiting)
```

### Протоколы

- **gRPC**: `:50052` - для межсервисного взаимодействия
- **HTTP/REST**: `:8081` - через gRPC-Gateway для frontend/BFF

## Требования

- Go 1.25+
- PostgreSQL 15+
- Redis 7+
- protoc (для генерации кода)

## Быстрый старт

### 1. Запуск зависимостей

```bash
# Из корня проекта
docker-compose up -d postgres redis
```

### 2. Конфигурация

Скопируйте и отредактируйте `config.yaml`:

```yaml
database:
  host: "localhost"
  port: 5432
  username: "admin"
  password: "admin"
  name: "medialog"
  ssl_mode: "disable"

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0

auth:
  jwt_secret: "CHANGE_ME_IN_PRODUCTION"  # В проде через env/secret manager
  access_ttl_seconds: 3600               # 1 час
  refresh_ttl_seconds: 2592000           # 30 дней
  rate_limit_login_per_minute: 5
  rate_limit_register_per_minute: 5
  rate_limit_refresh_per_minute: 5

server:
  grpc_addr: ":50052"
  http_addr: ":8081"
```

### 3. Генерация кода

```bash
make generate-api
```

Генерирует:
- gRPC код (`internal/pb/auth_api/`)
- gRPC-Gateway (`internal/pb/auth_api/auth.pb.gw.go`)
- Swagger/OpenAPI (`internal/pb/swagger/auth_api/auth.swagger.json`)

### 4. Запуск сервиса

```bash
make run
```

Или вручную:

```bash
configPath=./config.yaml \
swaggerPath=./internal/pb/swagger/auth_api/auth.swagger.json \
go run ./cmd/app
```

## API Документация

После запуска сервиса доступна Swagger UI:

- **Swagger UI**: http://localhost:8081/docs/
- **Swagger JSON**: http://localhost:8081/swagger.json

### Endpoints

- `POST /v1/auth/register` - Регистрация
- `POST /v1/auth/login` - Вход
- `POST /v1/auth/refresh` - Обновление токенов
- `POST /v1/auth/logout` - Выход
- `POST /v1/auth/logout-all` - Выход со всех устройств

## Переменные окружения

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `configPath` | Путь к файлу конфигурации | `./config.yaml` |
| `swaggerPath` | Путь к swagger.json | `./internal/pb/swagger/auth_api/auth.swagger.json` |

## Конфигурация

### Database (PostgreSQL)

- **users** таблица создается автоматически при первом запуске
- Схема: `id`, `email` (UNIQUE), `password_hash`, `created_at`

### Redis

#### Ключи сессий

Формат: `session:{hex(refresh_hash)}`

- **TTL**: Устанавливается при создании сессии, равен `refresh_ttl_seconds`
- **Структура**: JSON с полями:
  ```json
  {
    "UserID": 1,
    "RefreshHash": "...",
    "ExpiresAt": "2026-01-27T...",
    "RevokedAt": null,
    "CreatedAt": "2025-12-28T...",
    "UserAgent": "...",
    "IP": "127.0.0.1"
  }
  ```

#### Ключи rate limiting

Формат: `rl:{kind}:{ip}`

Где `kind`:
- `login` - ограничение попыток входа
- `register` - ограничение регистраций
- `refresh` - ограничение обновлений токенов

- **TTL**: 60 секунд (скользящее окно)
- **Алгоритм**: Sliding window через Lua script

## Структура проекта

```
auth_service/
├── api/                    # Proto файлы
│   ├── auth_api/
│   │   └── auth.proto      # gRPC сервис
│   └── models/
│       └── auth_model.proto
├── cmd/app/                # Точка входа
│   └── main.go
├── config/                 # Конфигурация
│   └── config.go
├── internal/
│   ├── api/                # API слой (gRPC handlers)
│   │   └── auth_service_api/
│   ├── bootstrap/          # Инициализация зависимостей
│   ├── models/             # Доменные модели
│   ├── pb/                 # Сгенерированный код
│   │   ├── auth_api/       # gRPC + Gateway
│   │   └── swagger/        # OpenAPI документация
│   ├── services/           # Бизнес-логика
│   │   └── auth_service/
│   └── storage/            # Хранилища
│       ├── auth_storage/   # PostgreSQL
│       └── session_storage/ # Redis
├── scripts/                # Скрипты генерации
└── config.yaml             # Конфигурационный файл
```

## Тестирование

### Unit тесты

```bash
make cov
```

Запускает тесты с покрытием для сервисного слоя.

### Генерация моков

```bash
make mock
```

Генерирует моки для интерфейсов (использует `mockery`).

## Безопасность

### Пароли

- Хранятся в виде bcrypt хэшей (cost 10)
- Валидация: минимум 8 символов, обязательны: заглавные, строчные, цифры

### Токены

- **Access Token**: JWT с коротким TTL (по умолчанию 1 час)
- **Refresh Token**: Opaque токен, хранится в Redis с TTL (по умолчанию 30 дней)
- При refresh старый токен отзывается, выдается новая пара (ротация)

### Rate Limiting

- **Login**: 10 попыток в минуту (настраивается)
- **Register**: 10 попыток в минуту (настраивается)
- **Refresh**: 30 попыток в минуту (настраивается)
- Ограничение по IP адресу
- Sliding window алгоритм

### Валидация

- Email: RFC-совместимая проверка
- Пароль: сложность (8+ символов, верхний/нижний регистр, цифры)
- Структурированные ошибки для frontend

## Обработка ошибок

Все ошибки возвращаются в структурированном JSON формате:

```json
{
  "code": "ERROR_CODE",
  "message": "Human-readable message",
  "field": "field_name"  // опционально, для валидационных ошибок
}
```

Коды ошибок:
- `INVALID_EMAIL` - Неверный формат email
- `INVALID_PASSWORD` - Неверный формат пароля
- `INVALID_CREDENTIALS` - Неверный email или пароль
- `EMAIL_ALREADY_EXISTS` - Email уже зарегистрирован
- `INVALID_TOKEN` - Неверный токен
- `SESSION_EXPIRED` - Сессия истекла
- `SESSION_REVOKED` - Сессия отозвана
- `RATE_LIMIT_EXCEEDED` - Превышен лимит запросов
- `INTERNAL_ERROR` - Внутренняя ошибка сервера

## Разработка

### Генерация кода из proto

```bash
make generate-api
```

Требует установленные:
- `protoc`
- `protoc-gen-go`
- `protoc-gen-go-grpc`
- `protoc-gen-grpc-gateway`
- `protoc-gen-openapiv2`

### Зависимости

```bash
go mod download
go mod tidy
```

## Production

### Рекомендации

1. **JWT Secret**: Использовать сильный секретный ключ через переменные окружения или secret manager
2. **HTTPS**: Настроить TLS на уровне reverse proxy (nginx/traefik)
3. **Миграции**: Использовать систему миграций (golang-migrate) вместо `initTables()`
4. **Мониторинг**: Добавить метрики (Prometheus) и трейсинг (OpenTelemetry)
5. **Логирование**: Настроить централизованное логирование (ELK, Loki)

### Переменные окружения для production

```bash
export CONFIG_PATH=/etc/auth-service/config.yaml
export JWT_SECRET=$(cat /run/secrets/jwt_secret)
```

## Лицензия

См. LICENSE в корне проекта.
