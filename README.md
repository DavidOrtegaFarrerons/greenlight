# Greenlight

A JSON API for managing a movie database. Built in Go while working through [*Let's Go Further*](https://lets-go-further.alexedwards.net/) by Alex Edwards.

The project follows the book closely (Greenlight is the example application used throughout) but I treated it as a proper engineering exercise: real Docker setup and enough attention to the surrounding tooling that I'd feel comfortable extending it.

---

## What it does

The API lets you create, read, update and delete movie records. Each movie has a title, release year, runtime, and a list of genres. Beyond the basic CRUD, it supports full-text search, filtering by genre, and sorting and pagination on the listing endpoint — the kind of thing you'd actually need in a real product.

Users can register, activate their account via email, log in to get a bearer token, and use that token to make authenticated requests. Permissions are scoped per-user: a user can have read-only or read-write access to movies.

---

## Features

**Movies**
- `GET /v1/movies` — list with filtering (genre, title search), sorting (any field, asc/desc), and cursor-based pagination
- `POST /v1/movies` — create a movie (requires write permission)
- `GET /v1/movies/:id` — fetch a single movie
- `PATCH /v1/movies/:id` — partial update (requires write permission)
- `DELETE /v1/movies/:id` — delete (requires write permission)

**Users & auth**
- `POST /v1/users` — register; triggers an account activation email in the background
- `PUT /v1/users/activated` — activate account via emailed token
- `POST /v1/tokens/authentication` — exchange credentials for a stateful bearer token
- `PUT /v1/users/password` — reset password via emailed token

**System**
- `GET /v1/healthcheck` — returns version, environment, and status
- `GET /debug/vars` — Go's `expvar` metrics (goroutines, DB pool stats, request counts)

---

## Under the hood

- **PostgreSQL** with `lib/pq`, raw SQL, and `golang-migrate` for schema migrations
- **Token-based auth** — tokens are hashed (SHA-256) before storage, never kept in plaintext
- **Permission model** — per-user read/write scopes checked via middleware
- **Background jobs** — email sending runs in goroutines with a `sync.WaitGroup` so the server doesn't shut down mid-send
- **Rate limiting** — per-IP token bucket via `golang.org/x/time/rate`, configurable at startup
- **Graceful shutdown** — listens for `SIGINT`/`SIGTERM` and drains in-flight requests before exiting
- **CORS** — configurable trusted origins with preflight support
- **Structured JSON logging** — all logs go to stdout as JSON; no log files
- **Input validation** — custom `Validator` type accumulates field errors and returns them as a structured JSON response
- **Centralised error handling** — consistent JSON error envelopes across all failure modes

---

## Project layout

```
cmd/api/            → main package: server setup, routing, handlers, middleware
internal/           → data models, mailer, validator, DB helpers
migrations/         → SQL migration files (golang-migrate)
Makefile            → dev shortcuts (run, migrate, audit, build)
Dockerfile          → production image
Dockerfile.debug    → image with Delve attached for remote debugging
docker-compose.yaml → production compose
```

---

## Running locally

You need Docker.

```bash
cp .env.example .env
# fill in DB credentials and SMTP settings

docker compose up -d # starts postgres + api
```

The API will be available at `http://localhost:4000`.

---