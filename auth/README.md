# Auth Service

Микросервис аутентификации и авторизации для системы MediaLog. Отвечает за регистрацию пользователей, вход в систему, управление сессиями, выдачу и обновление токенов доступа, а также управление ролями пользователей.

## Содержание

- [Технологический стек](#технологический-стек)
- [Архитектура](#архитектура)
- [Быстрый старт](#быстрый-старт)
- [Конфигурация](#конфигурация)
- [API Endpoints](#api-endpoints)
- [Безопасность](#безопасность)

## Технологический стек

- **Язык**: Go 1.25+
- **Протоколы**: gRPC, HTTP/REST (через gRPC-Gateway)
- **База данных**: PostgreSQL 15+ (пользователи)
- **Кэш/Сессии**: Redis 7+ (сессии, rate limiting)
- **Аутентификация**: JWT (JSON Web Tokens)
- **Хеширование паролей**: bcrypt
- **Генерация кода**: Protocol Buffers (protoc)
- **Документация API**: Scalar API Reference UI

## Архитектура

### Слои сервиса

```
┌─────────────────────────────────────┐
│   HTTP/REST (gRPC-Gateway) :8081    │
│   gRPC :50052                        │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│      API Layer (auth_service_api)   │
│  - Валидация запросов               │
│  - Rate limiting                    │
│  - Обработка ошибок                 │
│  - Конвертация proto ↔ models       │
│  - Логирование (slog)              │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│   Service Layer (auth_service)      │
│  - Бизнес-логика                    │
│  - Валидация данных                 │
│  - Генерация токенов                │
│  - Управление сессиями              │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│      Storage Layer                  │
│  ├── AuthStorage (PostgreSQL)       │
│  └── SessionStorage (Redis)         │
└─────────────────────────────────────┘
```

### Порты и протоколы

- **gRPC**: `:50052` — для межсервисного взаимодействия
- **HTTP/REST**: `:8081` — для клиентских приложений (через gRPC-Gateway)
- **Scalar API Reference**: `http://localhost:8081/docs`

### Зависимости

- **PostgreSQL**: Хранение пользователей (таблица `auth_users`)
- **Redis**: Хранение сессий и rate limiting


### Запуск сервиса

```bash
make run
```

Сервис будет доступен на:
- **gRPC**: `localhost:50052`
- **HTTP/REST**: `http://localhost:8081`
- **API Documentation**: `http://localhost:8081/docs`

## Конфигурация

### Параметры конфигурации

| Параметр | Описание | По умолчанию |
|----------|----------|--------------|
| `database.host` | Хост PostgreSQL | `localhost` |
| `database.port` | Порт PostgreSQL | `5432` |
| `database.username` | Имя пользователя БД | `admin` |
| `database.password` | Пароль БД | `admin` |
| `database.name` | Имя базы данных | `medialog` |
| `database.ssl_mode` | Режим SSL | `disable` |
| `redis.host` | Хост Redis | `localhost` |
| `redis.port` | Порт Redis | `6379` |
| `redis.password` | Пароль Redis | `` |
| `redis.db` | Номер БД Redis | `0` |
| `auth.jwt_secret` | Секрет для подписи JWT | `CHANGE_ME_IN_PRODUCTION` |
| `auth.access_ttl_seconds` | TTL access token (сек) | `86400` (1 день) |
| `auth.refresh_ttl_seconds` | TTL refresh token (сек) | `604800` (7 дней) |
| `auth.rate_limit_login_per_minute` | Лимит запросов на логин | `2` |
| `auth.rate_limit_register_per_minute` | Лимит запросов на регистрацию | `2` |
| `auth.rate_limit_refresh_per_minute` | Лимит запросов на refresh | `2` |
| `server.grpc_addr` | Адрес gRPC сервера | `:50052` |
| `server.http_addr` | Адрес HTTP сервера | `:8081` |

## API Endpoints

### Регистрация пользователя

**Endpoint**: `POST /v1/auth/register`

**Описание**: Создаёт нового пользователя и возвращает пару access/refresh токенов.

**Request**:
```json
{
  "email": "user@example.com",
  "password": "Password123!"
}
```

**Response** (успех):
```json
{
  "user_id": 1,
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6..."
}
```

**Требования к паролю**:
- Минимум 8 символов
- Хотя бы одна заглавная буква
- Хотя бы одна строчная буква
- Хотя бы одна цифра

**Ошибки**:
- `400` - Неверный формат email или пароля
- `409` - Email уже зарегистрирован
- `429` - Превышен лимит запросов

---

### Вход в систему

**Endpoint**: `POST /v1/auth/login`

**Описание**: Проверяет учётные данные и возвращает пару access/refresh токенов.

**Request**:
```json
{
  "email": "user@example.com",
  "password": "Password123!"
}
```

**Response** (успех):
```json
{
  "user_id": 1,
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6..."
}
```

**Ошибки**:
- `400` - Неверный формат email или пароля
- `401` - Неверные учётные данные
- `429` - Превышен лимит запросов

---

### Обновление токенов

**Endpoint**: `POST /v1/auth/refresh`

**Описание**: Принимает refresh token и выдаёт новую пару access/refresh токенов (refresh токен ротируется).

**Request**:
```json
{
  "refresh_token": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6..."
}
```

**Response** (успех):
```json
{
  "user_id": 1,
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "x9y8z7w6v5u4t3s2r1q0p9o8n7m6l5k4..."
}
```

**Ошибки**:
- `400` - Refresh token не указан
- `401` - Неверный или истёкший refresh token
- `429` - Превышен лимит запросов

---

### Выход из сессии

**Endpoint**: `POST /v1/auth/logout`

**Описание**: Отзывает текущую сессию (удаляет refresh token из Redis).

**Headers**:
```
Authorization: Bearer <access_token>
```

**Request**:
```json
{
  "refresh_token": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6..."
}
```

**Response** (успех):
```json
{
  "success": true,
  "message": "Session revoked successfully."
}
```

**Ошибки**:
- `400` - Refresh token не указан
- `401` - Неверный refresh token или отсутствует access token

---

### Выход со всех устройств

**Endpoint**: `POST /v1/auth/logout-all`

**Описание**: Отзывает все сессии пользователя (удаляет все refresh tokens из Redis).

**Headers**:
```
Authorization: Bearer <access_token>
```

**Request**:
```json
{
  "refresh_token": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6..."
}
```

**Response** (успех):
```json
{
  "success": true,
  "message": "All sessions revoked successfully."
}
```

**Ошибки**:
- `400` - Refresh token не указан
- `401` - Неверный refresh token или отсутствует access token

---

### Изменение роли пользователя

**Endpoint**: `POST /v1/auth/users/{user_id}/role`

**Описание**: Изменяет роль пользователя. Доступно только администраторам.

**Headers**:
```
Authorization: Bearer <access_token>
```

**Request**:
```json
{
  "user_id": 2,
  "role": "admin"
}
```

**Response** (успех):
```json
{
  "success": true,
  "message": "User role updated successfully."
}
```

**Ошибки**:
- `400` - Неверная роль или попытка изменить свою роль
- `401` - Отсутствует или неверный access token
- `403` - Недостаточно прав (требуется роль admin)
- `404` - Пользователь не найден

---

#### Использование Scalar API Reference

1. Откройте `http://localhost:8081/docs`
2. Выберите нужный endpoint
3. Заполните параметры запроса
4. Нажмите "Execute"

## Безопасность

### Rate Limiting

Сервис защищён от брутфорс-атак через rate limiting на основе Redis:
- Лимит на логин: по умолчанию 2 запроса в минуту
- Лимит на регистрацию: по умолчанию 2 запроса в минуту
- Лимит на refresh: по умолчанию 2 запроса в минуту

Лимиты настраиваются в `config.yaml`.

### Хеширование паролей

Пароли хешируются с помощью bcrypt перед сохранением в базу данных.

### JWT токены

- **Access token**: JWT токен с коротким временем жизни (по умолчанию 1 день)
- **Refresh token**: Opaque токен с длинным временем жизни (по умолчанию 7 дней), хранится в Redis

### Управление сессиями

- Сессии хранятся в Redis с TTL равным времени жизни refresh token
- При logout сессии полностью удаляются из Redis
- При logout-all все сессии пользователя удаляются из Redis

## База данных

### Схема таблицы `auth_users`

```sql
CREATE TABLE IF NOT EXISTS auth_users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_auth_users_email ON auth_users(email);
```

Таблица создаётся автоматически при первом запуске сервиса.

## Redis

### Структура ключей

- Сессии: `session:<hex(refresh_hash)`
- Rate limiting: `ratelimit:<endpoint>:<ip>`
