# Focus Auth

Локальный backend-сервис для личного кабинета в Focus Android: регистрация, логин, подтверждение email, JWT-авторизация и управление профилем.

## Что реализовано

- `POST /api/auth/register` - регистрация пользователя
- `POST /api/auth/login` - логин
- `POST /api/auth/verify-email` - подтверждение email кодом
- `POST /api/auth/resend-verification` - повторная отправка кода
- `GET /api/users/me` - получить профиль текущего пользователя
- `PATCH /api/users/me` - обновить профиль (имя, email, пароль, аватар)
- `GET /health` - health-check
- Swagger UI: `GET /swagger/index.html`
- OpenAPI JSON: `GET /openapi.json`

## Требования

- Go `1.22+`

## Быстрый старт

Из директории `account-cabinet`:

```bash
go mod download
go run ./cmd/server
```

По умолчанию сервис стартует на `http://127.0.0.1:8080`.

## Переменные окружения

| Переменная | По умолчанию | Описание |
|---|---|---|
| `ADDR` | `:8080` | Адрес/порт HTTP-сервера |
| `DATABASE_PATH` | `account-cabinet.db` | Путь к SQLite базе |
| `JWT_SECRET` | `dev-insecure-change-me` | Секрет подписи JWT (обязательно поменять в реальном окружении) |
| `JWT_EXPIRY_HOURS` | `168` | Срок жизни JWT в часах |
| `VERIFICATION_CODE_TTL_MINUTES` | `15` | Срок жизни кода подтверждения email в минутах |
| `SKIP_EMAIL_VERIFICATION` | `false` | Если `true`, email подтверждается автоматически (удобно для локальной разработки) |
| `GIN_MODE` | *(пусто)* | Если `release`, Gin работает в release-режиме |

Пример запуска для разработки без подтверждения email:

```bash
SKIP_EMAIL_VERIFICATION=true JWT_SECRET=local-dev-secret go run ./cmd/server
```

## Быстрая проверка API

### 1) Регистрация

```bash
curl -X POST http://127.0.0.1:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "userName": "dev-user",
    "email": "dev@example.com",
    "password": "strongpass123"
  }'
```

Если `SKIP_EMAIL_VERIFICATION=true`, сразу вернется `accessToken`.
Если `false`, сервис отправит код подтверждения в консоль.

### 2) Логин

```bash
curl -X POST http://127.0.0.1:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "dev@example.com",
    "password": "strongpass123"
  }'
```

Скопируйте `accessToken` из ответа.

### 3) Профиль текущего пользователя

```bash
curl http://127.0.0.1:8080/api/users/me \
  -H "Authorization: Bearer <accessToken>"
```

## Документация API

- Swagger UI: [http://127.0.0.1:8080/swagger/index.html](http://127.0.0.1:8080/swagger/index.html)
- OpenAPI JSON: [http://127.0.0.1:8080/openapi.json](http://127.0.0.1:8080/openapi.json)

## Примечания

- База SQLite создается автоматически при старте (по пути из `DATABASE_PATH`).
- Для любой среды кроме локальной разработки обязательно задайте свой `JWT_SECRET`.
