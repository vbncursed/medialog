# MediaLog - Микросервисная система учёта контента

## Быстрый старт

### 1. Запуск зависимостей

```bash
docker-compose up -d postgres redis kafka
```

Это запустит:
- PostgreSQL на порту `5432`
- Redis на порту `6379`
- Kafka на порту `9092`
- Kafka UI на порту `8090`

### 2. Настройка конфигурации

**Важно:** 
- `jwt_secret` должен быть **одинаковым** во всех трёх сервисах (auth, library, metadata)
- Для работы поиска в metadata_service нужен API ключ TMDB (можно получить на https://www.themoviedb.org/settings/api)
- Если API ключа нет, поиск будет работать только по локальной БД

### 3. Запуск сервисов

Откройте три терминала и запустите каждый сервис:

#### Auth Service
```bash
cd auth_service
make run
```

#### Library Service
```bash
cd library_service
make run
```

#### Metadata Service
```bash
cd metadata_service
make run
```

### 4. Доступ к Swagger UI

После запуска всех сервисов доступны следующие Swagger UI:

- **Auth Service**: http://localhost:8081/docs
- **Library Service**: http://localhost:8082/docs
- **Metadata Service**: http://localhost:8083/docs
- **Kafka UI**: http://localhost:8090

## Порты сервисов

| Сервис | gRPC | HTTP | Swagger UI |
|--------|------|------|-----------|
| Auth Service | :50052 | :8081 | http://localhost:8081/docs |
| Library Service | :50053 | :8082 | http://localhost:8082/docs |
| Metadata Service | :50054 | :8083 | http://localhost:8083/docs |

## Порядок проверки через Swagger

### 1. Auth Service - Регистрация и вход

1. Откройте http://localhost:8081/docs
2. Найдите метод `POST /v1/auth/register` и нажмите "Try it out"
3. Введите данные:
   ```json
   {
     "email": "test@example.com",
     "password": "testpassword123"
   }
   ```
4. Нажмите "Execute"
5. Скопируйте `access_token` из поля в ответе

### 2. Metadata Service - Поиск контента

1. Откройте http://localhost:8083/docs
2. Нажмите кнопку **"Authorize"** (вверху справа)
3. В поле "Value" вставьте `access_token` из шага 1 (без префикса "Bearer")
4. Найдите метод `POST /v1/metadata/search`
5. Введите данные:
   ```json
   {
     "query": "Inception",
     "type": "MEDIA_TYPE_MOVIE",
     "page": 1,
     "page_size": 10
   }
   ```
6. Нажмите "Execute"
7. Скопируйте `mediaId` из первого результата в массиве `results` (например, `"1"`)

### 3. Library Service - Создание записи

1. Откройте http://localhost:8082/docs
2. Нажмите кнопку **"Authorize"** (вверху справа)
3. В поле "Value" вставьте тот же `access_token` из шага 1
4. Найдите метод `POST /v1/library/entries` и нажмите "Try it out"
5. Введите данные (замените `<media_id>` на значение из шага 2):
   ```json
   {
     "media_id": "1",
     "type": "MEDIA_TYPE_MOVIE",
     "status": "ENTRY_STATUS_PLANNED",
     "rating": 0
   }
   ```
6. Нажмите "Execute"
7. Проверьте, что запись создана успешно (статус 200, есть `entryId` в ответе)

**Примечание:** Если запись с таким `media_id` уже существует, вернётся ошибка `ENTRY_ALREADY_EXISTS`.

### 4. Проверка обогащения метаданных

1. Вернитесь в **Metadata Service Swagger** (http://localhost:8083/docs)
2. Убедитесь, что авторизация активна (токен из шага 1)
3. Найдите метод `GET /v1/metadata/media/{media_id}` и нажмите "Try it out"
4. В поле `media_id` введите значение из шага 2 (например, `1`)
5. Нажмите "Execute"
6. Проверьте ответ:
   - Должны быть заполнены поля: `title`, `genres` (не пустой массив), `posterUrl` (для фильмов)
   - Если метаданные обогатились через внешний API, должны быть `externalIds` с `source: "tmdb"`

**Примечание:** Обогащение метаданных происходит асинхронно через Kafka. Если метаданные ещё не обогатились, подождите несколько секунд и повторите запрос.

## Проверка Kafka событий

События можно проверить через Kafka UI:
- http://localhost:8090
- Топик: `library.entry.changed`

## Логи сервисов

Логи выводятся в консоль каждого сервиса. Для отладки проверяйте:
- Ошибки подключения к БД/Redis/Kafka
- Ошибки обработки запросов
- Логи обработки Kafka событий в metadata_service

