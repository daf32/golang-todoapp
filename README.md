# Golang Todo App

A Go REST API for a todo application (**users**, **tasks**, **statistics**) with a minimal web UI on `/`.

- **API base**: `/api/v1`
- **Swagger UI**: `/swagger/` (spec at `/swagger/doc.json`)
- **Default HTTP addr**: `:5050`

## Tech stack

- **Language**: Go (recommended `1.22+`; Docker image builds with Go `1.25.6`)
- **HTTP**: standard library `net/http` (no external framework)
- **Database**: PostgreSQL
- **DB driver**: `jackc/pgx/v5`
- **Config**: `kelseyhightower/envconfig` (+ `TIME_ZONE`)
- **Validation**: `go-playground/validator/v10`
- **Logging**: `go.uber.org/zap` (writes to files under `LOGGER_FOLDER`)
- **Migrations**: `migrate/migrate` (via Docker)
- **API docs**: Swagger (`swaggo/swag`, generated to `docs/`)
- **Runtime/deploy**: Docker / docker compose

## Architecture

The project follows **Clean Architecture** and keeps dependencies pointing inward.

Each feature (`users`, `tasks`, `statistics`, `web`) is split into 3 layers:

```
Transport (HTTP handler)  →  Service (business logic)  →  Repository (data access)
```

Shared, cross-cutting primitives live under `internal/core/` (HTTP server/middleware, logger, config, postgres pool adapter, domain helpers).

Dependency wiring is done in `cmd/todoapp/main.go`:

```
Repository → Service → HTTP Handler → Routes()
```

## Project structure (high level)

```
.
├── cmd/todoapp/                 # entrypoint + Dockerfile
├── internal/
│   ├── core/                    # shared infrastructure & primitives
│   └── features/                # users / tasks / statistics / web
├── migrations/                  # SQL migrations (golang-migrate)
├── public/                      # static files (served by web feature)
├── docs/                        # generated Swagger docs (checked in)
├── docker-compose.yaml          # local infra + app container
└── Makefile                     # canonical task runner (loads .env)
```

## Local run

### Prerequisites

- Docker + docker compose
- Go (if running via `go run`)
- `make`

### Steps

```bash
# 1) Create local env file
cp .env.example .env

# 2) Fill in required values
$EDITOR .env

# 3) Start Postgres
make env-up

# 4) Run migrations
make migrate-up

# 5) Expose Postgres on 127.0.0.1:5432 (the DB container itself does not publish ports)
make env-port-forward

# 6) Run the app locally
make todoapp-run
```

After start:

- **Main page**: `http://127.0.0.1:5050/`
- **Swagger UI**: `http://127.0.0.1:5050/swagger/`
- **API**: `http://127.0.0.1:5050/api/v1/`

## Run via Docker (compose)

This runs the app as a container and exposes it on port `5050`.

```bash
make env-up
make migrate-up
make todoapp-deploy
```

To stop it:

```bash
make todoapp-undeploy
```

## Deploy via GitHub Actions

The repository includes `.github/workflows/deploy.yml`. It deploys on pushes to `main` or on manual `workflow_dispatch`.

### Server prerequisites

The target host must have:

- Docker with `docker compose`
- `rsync`
- `base64`
- SSH access for the user from `DEPLOY_USER`
- write access to `DEPLOY_PATH`

The workflow copies the repository to the server with `rsync`, writes `.env` from a GitHub secret, starts Postgres, waits for readiness, runs migrations, and rebuilds `todoapp`.

### Required GitHub secrets

Set these in the repository or in the `production` environment:

| Secret | Description | Example |
|---|---|---|
| `DEPLOY_HOST` | Public server IP or domain | `203.0.113.10` |
| `DEPLOY_PORT` | SSH port | `22` |
| `DEPLOY_USER` | SSH user on the server | `root` |
| `DEPLOY_PATH` | Absolute deploy directory on the server | `/root/apps/golang-todoapp` |
| `DEPLOY_SSH_KEY` | Private SSH key content used by Actions | `-----BEGIN OPENSSH PRIVATE KEY-----...` |
| `DEPLOY_ENV_B64` | Base64-encoded `.env` contents | `SFRUUF9BRERSPTo1MDUw...` |

Generate `DEPLOY_ENV_B64` from your production `.env`:

```bash
# macOS
base64 < .env | tr -d '\n'

# Linux
base64 -w 0 .env
```

## Common Makefile commands

- **`make env-up`**: start Postgres container
- **`make env-port-forward`**: forward DB to `127.0.0.1:5432` via a socat sidecar
- **`make env-port-close`**: stop the port forwarder
- **`make env-cleanup`**: wipe `out/pgdata` (destructive)
- **`make migrate-up`** / **`make migrate-down`**: apply / rollback migrations (via container)
- **`make migrate-create seq=<name>`**: create a new migration file pair
- **`make todoapp-run`**: run locally with `go run` (forces `POSTGRES_HOST=localhost` and sets `LOGGER_FOLDER=$PROJECT_ROOT/out/logs`)
- **`make todoapp-deploy`** / **`make todoapp-undeploy`**: build & run / stop app via docker compose
- **`make swagger-gen`**: regenerate Swagger docs to `docs/`

## Environment variables

The Makefile includes `.env` (`include .env` + `export`), so the app picks up variables from it.

### Required in `.env` for local development

| Variable | Description | Example |
|---|---|---|
| `HTTP_ADDR` | HTTP server address | `:5050` |
| `HTTP_ALLOWED_ORIGINS` | CORS allowed origins (comma-separated) | `http://127.0.0.1:5050,null` |
| `POSTGRES_USER` | DB user | `todoapp` |
| `POSTGRES_PASSWORD` | DB password | `todoapp` |
| `POSTGRES_DB` | DB name | `todoapp` |
| `POSTGRES_TIMEOUT` | DB op timeout | `10s` |

### Notes

- `POSTGRES_HOST` is set by:
  - `make todoapp-run` to `localhost` (expects port-forwarder to be up)
  - docker compose to `todoapp-postgres`
- `POSTGRES_PORT` is optional (defaults to `5432`)
- `LOGGER_FOLDER` is required by the logger config and is set by:
  - `make todoapp-run` to `$PROJECT_ROOT/out/logs`
  - docker compose to `/app/out/logs`
- `TIME_ZONE` defaults to `UTC` if not set.

## API

### Users (`/api/v1/users`)

- **POST** `/users` — create user
- **GET** `/users` — list users (pagination)
- **GET** `/users/{id}` — get user by id
- **PATCH** `/users/{id}` — patch user
- **DELETE** `/users/{id}` — delete user

### Tasks (`/api/v1/tasks`)

- **POST** `/tasks` — create task
- **GET** `/tasks` — list tasks (pagination + optional `user_id` filter)
- **GET** `/tasks/{id}` — get task by id
- **PATCH** `/tasks/{id}` — patch task
- **DELETE** `/tasks/{id}` — delete task

### Statistics (`/api/v1/statistics`)

- **GET** `/statistics` — task statistics (filters: `user_id`, `from`, `to` in `YYYY-MM-DD`)

Full interactive documentation is available at **Swagger UI**: `/swagger/`.
