# user-api

A small API for managing users.

## Quick Start

Requirements: Docker (for Postgres), `migrate`, Go 1.24+.

Run:
```sh
docker compose up --build
```

Stop and remove data:
```sh
docker compose down -v
```

## Long Start

1) Start postgres in Docker and create the db:
```sh
go mod tidy
make postgres
make createdb
make migrateup
```
2) Run API server:
```sh
go run ./cmd/server
```
By default the config is taken from `config/app.env`. The API server listens on `0.0.0.0:8080`.

## API
- POST `/users` — create a user
- POST `/users/login` — login (response contains a `token`)
- GET `/users/:username` — user information (Authorization: `Bearer <token>`)
- POST `/users/:username` — update user (you can only update yourself) (Authorization: `Bearer <token>`)
