# Tojtech API

A REST API backend for an e-commerce store specializing in IT hardware and computer products. Written in Go with a Clean Architecture layout, PostgreSQL storage, JWT authentication with rotating refresh tokens, and role-based access control.

## Features

- Categories and products CRUD, with UUID public identifiers
- Product images (URL or base64 data URL), with multiple images per product
- Registration, login, token refresh and logout
- Admin-only routes protected by JWT middleware (`RequireAuth` + `RequireRole`)
- SQL migrations via `golang-migrate`
- Structured, colored logs in local mode and JSON logs in production

## Tech stack

- Go 1.26.5 (see `go.mod`)
- PostgreSQL 17
- [chi](https://github.com/go-chi/chi) router, [cors](https://github.com/go-chi/cors)
- [pgx v5](https://github.com/jackc/pgx) with `pgxpool`
- [golang-jwt](https://github.com/golang-jwt/jwt) (HS256) and `bcrypt` for passwords
- [cleanenv](https://github.com/ilyakaznacheev/cleanenv) for configuration
- Docker and Docker Compose (optional)

## Project structure

```
cmd/api/                  entry point: wiring of config, DB, usecases, handlers
config/                   local.yaml (local development), docker.yaml (container)
internal/
  config/                 configuration loading (YAML + environment variables)
  domain/                 entities, validation, sentinel errors
  usecase/                business logic, depends on interfaces it declares
  repository/postgres/    pgx implementations of the repository interfaces
  delivery/http_server/   routes, middleware, HTTP handlers
  domain/dto/             request and response DTOs
pkg/                      shared helpers: logger, HTTP responses, request logging
migrations/               SQL migrations (up and down)
tests/api/                .http files for manual API testing
```

Dependencies point inward: handlers call usecases, usecases call repositories through interfaces, and the domain knows nothing about HTTP or SQL.

## Prerequisites

- Go 1.26.5 or newer
- PostgreSQL 17 (or a local install of any recent version), **or** Docker with Docker Compose
- `golang-migrate` CLI, if you run migrations locally (see below)

## Quick start with Docker (backend only)

The compose file starts PostgreSQL, applies the migrations, and runs the API. The frontend is not containerized and runs separately with `npm run dev`.

```bash
git clone https://github.com/shodruzhoshimzoda/Tojtech.git
cd Tojtech

cp .env.example .env     # then edit DB_PASSWORD and JWT_SECRET
docker compose up -d --build
docker compose logs -f api
```

The API is available at `http://localhost:8080/api/v1`. PostgreSQL is exposed on `localhost:5433` so you can connect with `psql` or a GUI client.

Stop the stack with `docker compose down`. To also delete the database volume, use `docker compose down -v`.

## Local setup (without Docker)

### 1. Create the database

```sql
CREATE DATABASE myshop;
```

### 2. Download dependencies

```bash
go mod tidy
```

### 3. Install the migration tool

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### 4. Configure the database connection

Set `user` and `database` in `config/local.yaml` to match your PostgreSQL setup:

```yaml
db-conn:
  host: "localhost"
  port: 5432
  user: your_username
  database: your_database
  sslmode: disable
```

The password is never stored in the config file. Provide it at runtime (see below).

### 5. Apply migrations

```bash
DB_PASSWORD='your_postgres_password' make migrate-up
```

To roll back the last migration: `DB_PASSWORD=... make migrate-down`.

### 6. Run the server

Either pass the password as a flag:

```bash
go run ./cmd/api/ -db-password=your_postgres_password
```

or as an environment variable:

```bash
DB_PASSWORD=your_postgres_password make server-run
```

The server listens on `localhost:8080` by default.

## Configuration

Settings are read from `config/local.yaml` (or `config/docker.yaml` in the container). Secrets and environment-specific values come from environment variables or flags.

| Variable / flag | Required | Default | Description |
|---|---|---|---|
| `DB_PASSWORD` / `-db-password` | yes | none | PostgreSQL password |
| `JWT_SECRET` | yes | from `local.yaml` | Secret used to sign access tokens |
| `DB_USER` | no | `shodruz` (Docker) | Database user (Docker only) |
| `DB_NAME` | no | `myshop` (Docker only) | Database name (Docker only) |
| `CORS_ORIGIN` | no | `http://localhost:5173` | Origin allowed to call the API from a browser |
| `CONFIG_PATH` | no | `./config/local.yaml` | Path to the config file, if not passed with `-config-file` |

Token lifetimes are set in the config file: `jwt.token-ttl` (access token, 24h) and `jwt.refresh-token-ttl` (refresh token, 720h).

> **Note:** generate a strong secret with `openssl rand -hex 32` and never commit a real `JWT_SECRET` or database password.

## API overview

All routes are prefixed with `/api/v1`. Protected routes require an `Authorization: Bearer <access_token>` header and an admin role.

### Authentication

| Method | Path | Auth | Description |
|---|---|---|---|
| POST | `/auth/register` | public | Create an account (role `customer`) |
| POST | `/auth/login` | public | Returns access and refresh tokens |
| POST | `/auth/refresh` | public | Exchanges a refresh token for a new pair (the old one is revoked) |
| POST | `/auth/logout` | public | Revokes a refresh token |

### Categories

| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/categories` | public | List categories |
| GET | `/categories/{uuid}` | public | Get one category |
| POST | `/categories` | admin | Create a category |
| PATCH | `/categories/{uuid}` | admin | Update a category |
| DELETE | `/categories/{uuid}` | admin | Delete a category |

### Products

| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/products` | public | List products with category and images |
| GET | `/products/{uuid}` | public | Get one product |
| POST | `/products` | admin | Create a product (with optional images) |
| PATCH | `/products/{uuid}` | admin | Update a product |
| DELETE | `/products/{uuid}` | admin | Delete a product |
| POST | `/products/{uuid}/images` | admin | Add an image |
| DELETE | `/products/{uuid}/images/{image_uuid}` | admin | Delete an image |

Products reference their category by UUID:

```json
{
  "name": "Mechanical keyboard",
  "slug": "mechanical-keyboard",
  "description": "75% layout, hot-swappable switches",
  "price": 89.99,
  "stock": 12,
  "category": { "uuid": "2f1c8e4a-0c1e-4f7e-9a8b-3d2a1b0c9e11" },
  "images": [
    { "image_url": "https://example.com/keyboard.jpg", "is_main": true }
  ]
}
```

### Errors

Errors use a consistent shape. `error` is safe to show to clients, and `error_description` (when present) gives validation details:

```json
{ "error": "category is invalid", "error_description": "name is required" }
```

## Creating the first admin

New accounts are always created as `customer`. To promote an account to admin, run this once against the database:

```bash
# Docker
docker compose exec db psql -U shodruz -d myshop \
  -c "UPDATE users SET role='admin' WHERE email='you@example.com';"

# Local PostgreSQL
psql -U your_username -d your_database \
  -c "UPDATE users SET role='admin' WHERE email='you@example.com';"
```

Log in again afterwards so the new role is included in your token.

## Testing

Manual API tests live in `tests/api/` as `.http` files (for example, the REST Client extension for VS Code or the HTTP client in JetBrains IDEs). Run the requests in order: create a category, then a product, then the protected routes with an admin token.

Automated unit and integration tests are not written yet.

## Common problems

- **`DB Password is not set`**: pass `-db-password` or set `DB_PASSWORD`.
- **`JWT token is not set`**: the `jwt.secret` value is empty. Set `JWT_SECRET` (Docker) or `jwt.secret` in `local.yaml`.
- **`403 Forbidden` on admin routes after a role change**: log in again, because the role is stored in the access token.
- **`value too long for type character varying`** when uploading images: migration `000005` must be applied (`make migrate-up`).
- **Port 8080 already in use**: stop the local `go run` process, or set `API_PORT` when using Docker.
- **Changes not visible after editing code**: restart the server. A stale binary is the most common cause of confusing results.