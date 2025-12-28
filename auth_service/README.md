## auth_service

Микросервис авторизации для MediaLog.

### Команды

- Генерация protobuf/gateway/swagger:

```bash
make generate-api
```

- Запуск:

```bash
make run
```

- Тесты:

```bash
make cov
```

### Порты (по умолчанию)

- **gRPC**: `:50052`
- **HTTP (grpc-gateway)**: `:8080`


