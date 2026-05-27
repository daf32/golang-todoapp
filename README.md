# Golang Todo App

A Go REST API for a todo application (**auth**, **users**, **tasks**, **statistics**) with a React (Vite) frontend served on `/`.

- **API base**: `/api/v1`
- **Swagger UI**: `/swagger/` (spec at `/swagger/doc.json`)
- **Default HTTP addr**: `:5050`

## Tech stack

- **Language**: Go (recommended `1.22+`; Docker image builds with Go `1.25.6`)
- **Frontend**: React 19 + Vite 8 (built during Docker image build; output served from `frontend/dist`)
- **HTTP**: standard library `net/http` (no external framework)
- **Database**: PostgreSQL
- **DB driver**: `jackc/pgx/v5`
- **Auth**: JWT access tokens (`golang-jwt/jwt/v5`) + opaque refresh tokens; bcrypt password hashing
- **OAuth / OIDC**: `golang.org/x/oauth2` + `coreos/go-oidc/v3` (Google; provider registry ready for Apple etc.)
- **Email**: `net/smtp` for sending confirmation links
- **Rate limiting**: `golang.org/x/time/rate` (in-memory token bucket per IP)
- **Config**: `kelseyhightower/envconfig` (+ `TIME_ZONE`)
- **Validation**: `go-playground/validator/v10`
- **Logging**: `go.uber.org/zap` (writes to files under `LOGGER_FOLDER`)
- **Migrations**: `migrate/migrate` (via Docker)
- **Mocks**: `vektra/mockery` (testify template), generated via `make generate`
- **API docs**: Swagger (`swaggo/swag`, generated to `docs/`)
- **Runtime/deploy**: Docker / docker compose

## Architecture

The project follows **Clean Architecture** and keeps dependencies pointing inward.

Each feature (`auth`, `users`, `tasks`, `statistics`, `web`) is split into 3 layers:

```
Transport (HTTP handler)  →  Service (business logic)  →  Repository (data access)
```

Shared, cross-cutting primitives live under `internal/core/` (HTTP server/middleware, logger, config, postgres pool adapter, domain helpers, OAuth providers, mailer, rate-limit middleware, cookie manager).

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
│   └── features/                # auth / users / tasks / statistics / web
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

The repository includes `.github/workflows/deploy.yml`. It is triggered **manually** via the GitHub Actions UI — go to **Actions → Deploy → Run workflow**. There is no automatic deploy on push.

The `verify` job runs Go tests, builds the Go binary, and builds the React frontend (`npm ci && npm run build`) before any deployment proceeds. The `deploy` job then rsyncs the source (excluding `node_modules`) to the server and rebuilds the Docker image — the Dockerfile handles the full frontend build on the server side via a Node multi-stage build.

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
- **`make generate`**: regenerate mockery mocks (run after changing any repo interface)
- **`make admin-promote email=user@example.com`**: promote a user to the `admin` role

## Environment variables

The Makefile includes `.env` (`include .env` + `export`), so the app picks up variables from it.

### Required in `.env` for local development

#### HTTP & Postgres

| Variable | Description | Example |
|---|---|---|
| `HTTP_ADDR` | HTTP server address | `:5050` |
| `HTTP_ALLOWED_ORIGINS` | CORS allowed origins (comma-separated) | `http://127.0.0.1:5050,null` |
| `POSTGRES_USER` | DB user | `todoapp` |
| `POSTGRES_PASSWORD` | DB password | `todoapp` |
| `POSTGRES_DB` | DB name | `todoapp` |
| `POSTGRES_TIMEOUT` | DB op timeout | `10s` |
| `APP_BASE_URL` | Public base URL of the app (used in confirmation links and OAuth redirect URIs) | `http://localhost:5050` |
| `COOKIE_SECURE` | Set `Secure` flag on cookies. Use `false` for local HTTP dev, `true` in HTTPS production | `false` |

#### Auth (JWT + email confirmation)

| Variable | Description | Example |
|---|---|---|
| `AUTH_JWT_SECRET` | Signing key for access tokens | `change-me-32+-chars` |
| `AUTH_ACCESS_TOKEN_EXPIRY` | Access token TTL | `15m` |
| `AUTH_REFRESH_TOKEN_EXPIRY` | Refresh token TTL | `168h` |
| `AUTH_EMAIL_CONFIRMATION_TOKEN_EXPIRY` | Email confirmation token TTL | `24h` |

#### SMTP (email confirmation delivery)

| Variable | Description | Example |
|---|---|---|
| `SMTP_HOST` | SMTP server host | `sandbox.smtp.mailtrap.io` |
| `SMTP_PORT` | SMTP server port (defaults to `587`) | `587` |
| `SMTP_USERNAME` | SMTP auth user | |
| `SMTP_PASSWORD` | SMTP auth password | |
| `SMTP_FROM` | `From:` address | `no-reply@example.com` |

#### Google OAuth

| Variable | Description |
|---|---|
| `GOOGLE_OAUTH_CLIENT_ID` | OAuth 2.0 Client ID from Google Cloud Console |
| `GOOGLE_OAUTH_CLIENT_SECRET` | OAuth 2.0 Client Secret |

The redirect URI to register in Google Cloud Console is `{APP_BASE_URL}/api/v1/auth/oauth/google/callback`.

#### Unverified user cleanup scheduler

| Variable | Description | Default |
|---|---|---|
| `USERS_UNVERIFIED_CLEANUP_INTERVAL` | How often the background loop runs | `24h` |
| `USERS_UNVERIFIED_CLEANUP_MIN_AGE` | Grace period before an unverified user is deleted | `168h` (7 days) |

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

All routes under `/api/v1/` except those listed in **Auth** require a valid `Authorization: Bearer <access_token>` header.

### Auth (`/api/v1/auth`)

| Route | Auth | Description |
|---|---|---|
| **POST** `/auth/register` | public | Create a new user; sends a confirmation email |
| **POST** `/auth/login` | public | Email + password → access + refresh tokens (rejects unverified emails with 403) |
| **POST** `/auth/refresh` | public | Exchange a refresh token for a new access token |
| **POST** `/auth/logout` | bearer | Revoke the supplied refresh token |
| **POST** `/auth/logout-all` | bearer | Revoke **all** refresh tokens for the user |
| **GET**  `/auth/confirm-email?token=...` | public | Confirm email; redirects to `/email-confirmed?status=...` |
| **POST** `/auth/resend-confirmation` | public | Resend a confirmation email (always returns 204, no enumeration) |
| **GET**  `/auth/oauth/{provider}` | public | Start an OAuth flow (must be a top-level browser navigation, not XHR) |
| **GET**  `/auth/oauth/{provider}/callback` | public | OAuth provider's callback — issues tokens |

Sensitive auth routes (`register`, `login`, `refresh`, `resend-confirmation`) are rate-limited per-IP via in-memory token buckets.

### Users (`/api/v1/users`)

- **GET** `/users` — list users (pagination + optional `?email_verified=true|false` filter; admin)
- **GET** `/users/{id}` — get user by id
- **PATCH** `/users/{id}` — patch user
- **POST** `/users/{id}/change-password` — change password (revokes all of the user's refresh tokens on success)
- **DELETE** `/users/{id}` — delete user

User registration is done via **POST** `/auth/register`, not `/users`.

### Tasks (`/api/v1/tasks`)

- **POST** `/tasks` — create task
- **GET** `/tasks` — list tasks (pagination + optional `user_id` filter)
- **GET** `/tasks/{id}` — get task by id
- **PATCH** `/tasks/{id}` — patch task
- **DELETE** `/tasks/{id}` — delete task

### Statistics (`/api/v1/statistics`)

- **GET** `/statistics` — task statistics (filters: `user_id`, `from`, `to` in `YYYY-MM-DD`)

Full interactive documentation is available at **Swagger UI**: `/swagger/`.

## Background jobs

- **Unverified-user cleanup loop** — runs in-process, started in `main.go`. On each tick it deletes users with `email_verified = false` whose `created_at` is older than `USERS_UNVERIFIED_CLEANUP_MIN_AGE`. Tied to the root signal context, so shuts down cleanly on `SIGINT`/`SIGTERM`.
