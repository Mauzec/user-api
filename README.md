# user-api

Небольшой api для управления пользователями (sqlc, pgxpool, Paseto токены).

## Быстрый старт

Требования: Docker (для Postgres), установленный `migrate` CLI, Go 1.24+.

1) Поднять Postgres в Docker и создать БД:
```sh
go mod tidy
make postgres
make createdb
make migrateup
```
2) Запустить API:
```sh
go run ./cmd/server
```
По умолчанию конфиг берётся из `config/app.env`. Порт: `0.0.0.0:8080`.

## Что, где?
- `db/` — логика Postgres: миграции, SQL-запросы (`db/query`), сгенерированный код sqlc (`db/sqlc`).
- `internal/api/` — логика API: хендлеры, роутинг, middleware, валидация.

## Запуск в Docker Compose

Всё будет в докере.
```sh
docker compose up --build
```

Остановка и очистка данных:
```sh
docker compose down -v
```

## HTTPS (опционально)
1) Сгенерировать самоподписанные сертификаты:
```sh
make gen-cert
```
2) В `config/app.env` расскоментировать TLS_CERT_FILE и TLS_KEY_FILE: