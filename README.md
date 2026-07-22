# URL Shortener

A full-stack URL shortener with **user-scoped link management**. Anyone can shorten a link publicly, but registered users get a personal dashboard to create, track, and own their links — backed by **PostgreSQL** as the source of truth and **Redis** as a caching + rate-limiting layer.

## 🚀 Features

- **Public URL Shortening** — anyone can shorten a link, no account needed
- **User Accounts** — open self-registration with JWT-based authentication and bcrypt password hashing
- **Per-User Link Ownership** — logged-in users create links tied to their account and see only their own links on the dashboard
- **Admin Dashboard** — create short URLs and view click-count analytics for links you own
- **Atomic Click Tracking** — SQL-level atomic increments (`click_count = click_count + 1`) eliminate race conditions under concurrent redirects
- **Cache-Aside Architecture** — Redis-first reads with automatic PostgreSQL fallback and cache backfill on a miss
- **Rate Limiting** — Redis-backed request throttling on registration and login endpoints
- **Dockerized** — backend, frontend, Postgres, and Redis all run via Docker Compose
- **Schema Migrations** — versioned SQL migrations via `golang-migrate`
- **Route-Aware Tab Titles & Favicons** — each page updates the browser tab title/icon dynamically

## 🏗️ Architecture

```mermaid
flowchart TD
    Client["React Frontend"]
    API["Go Backend (Gin)"]
    Redis[("Redis Cache")]
    Postgres[("PostgreSQL")]

    Client -->|"create / login / redirect"| API
    API -->|"check cache first"| Redis
    Redis -.->|"miss - fallback"| Postgres
    API -->|"writes always go here"| Postgres
    Postgres -.->|"backfill on miss"| Redis
    API -->|"302 redirect / JSON response"| Client
```

```
React Frontend  →  Go Backend (Gin)  →  Redis (cache + rate limiting)  →  PostgreSQL (source of truth)
```

**Write path (public):** short URL saved to Postgres (unowned, `user_id = NULL`), then cached in Redis.

**Write path (authenticated):** short URL saved to Postgres tagged with the creating user's ID, then cached in Redis.

**Read path (redirect):** Redis checked first. On a hit, redirect immediately and increment the click counter. On a miss, fall back to Postgres, backfill Redis, increment the counter, then redirect. Redis being empty or restarted never breaks a redirect — it only costs one extra Postgres lookup until the cache warms back up.

**Authorization:** JWTs carry the authenticated user's ID as a claim. The `/admin/urls` and `/admin/create-short-url` endpoints derive the caller's identity from the verified token — never from a client-supplied parameter — so one user can never see or create links attributed to another.

## 📋 Prerequisites

- Go 1.21 or higher
- Node.js 16 or higher (with npm)
- Docker Desktop (recommended — runs Postgres, Redis, and both services without installing anything natively)
- [`golang-migrate`](https://github.com/golang-migrate/migrate) CLI, for applying schema migrations

## 🏗️ Project Structure

```
url_shortener/
├── go_backend/
│   ├── main.go
│   ├── go.mod / go.sum
│   ├── Dockerfile
│   ├── .env.example
│   ├── auth/                     # JWT generation/validation, bcrypt hashing
│   │   └── auth.go
│   ├── middleware/
│   │   ├── auth_middleware.go     # Verifies JWT, injects user_id into context
│   │   └── rate_limit.go          # Redis-backed request throttling
│   ├── handler/
│   │   ├── handlers.go            # Public create + redirect + stats
│   │   ├── user_handlers.go       # Register, login, user-scoped create/list
│   │   └── errors.go              # Human-readable validation error formatting
│   ├── routes/
│   │   └── handle_urls.go
│   ├── shortener/
│   │   ├── shorturl_generator.go
│   │   └── shorturl_generator_test.go
│   ├── store/
│   │   ├── store_service.go        # Postgres (source of truth) + Redis (cache) core logic
│   │   ├── user_store.go            # User accounts + per-user URL queries
│   │   └── store_service_test.go
│   └── migrations/                  # Versioned SQL schema migrations
│       ├── 000001_create_urls_table.up.sql / .down.sql
│       ├── 000002_create_admins_table.up.sql / .down.sql
│       ├── 000003_add_admin_id_to_urls.up.sql / .down.sql
│       └── 000004_rename_admins_to_users.up.sql / .down.sql
│
├── react_frontend/
│   ├── src/
│   │   ├── pages/
│   │   │   ├── PublicShortener.tsx     # "/" — public shortening form
│   │   │   ├── Login.tsx                # "/admin/login"
│   │   │   ├── Register.tsx             # "/admin/register"
│   │   │   └── AdminDashboard.tsx       # "/admin" — protected, per-user
│   │   ├── components/
│   │   │   ├── Header.tsx / Footer.tsx
│   │   │   ├── ProtectedRoute.tsx       # Redirects logged-out users to /admin/login
│   │   │   └── PublicOnlyRoute.tsx      # Redirects logged-in users away from public pages
│   │   ├── api/client.ts                 # Axios instance, auto-attaches JWT
│   │   └── hooks/usePageMeta.ts           # Per-route tab title + favicon
│   ├── package.json
│   └── Dockerfile
│
└── docker-compose.yaml
```

## 🔧 Backend Setup (Go)

### Environment variables

```bash
cd go_backend
cp .env.example .env
```

`.env` (not committed to git):
```
DATABASE_URL=postgres://urlshortener:devpassword@localhost:5432/urlshortener?sslmode=disable
REDIS_ADDR=localhost:6379
BASE_URL=http://localhost:8000
JWT_SECRET=<a long random string>
```

> ⚠️ If your password contains `@` or other URL-reserved characters, URL-encode it in `DATABASE_URL` (`@` → `%40`) or avoid special characters in local dev passwords.

### Running Postgres and Redis locally (without Docker Compose)

```bash
docker run --name pg-urlshortener -e POSTGRES_USER=urlshortener -e POSTGRES_PASSWORD=devpassword -e POSTGRES_DB=urlshortener -p 5432:5432 -d postgres:16-alpine
docker run --name redis-urlshortener -p 6379:6379 -d redis:7-alpine
```

### Applying migrations

```bash
migrate -path migrations -database "postgres://urlshortener:devpassword@localhost:5432/urlshortener?sslmode=disable" up
```

### Installation & running

```bash
cd go_backend
go mod download
go mod tidy
go run main.go
```

The backend will be available at `http://localhost:8000`

### API Endpoints

**Public**
| Method | Path | Description |
|---|---|---|
| `GET` | `/` | Welcome message |
| `POST` | `/create-short-url` | Create an unowned short URL |
| `GET` | `/:shortUrl` | Redirect (cache-first, click count incremented) |
| `GET` | `/stats/:shortUrl` | Click count for a given short code |
| `POST` | `/admin/register` | Create a new user account (rate-limited: 5/hour/IP) |
| `POST` | `/admin/login` | Log in, returns a JWT (rate-limited: 10/15min/IP) |

**Authenticated** (`Authorization: Bearer <token>`)
| Method | Path | Description |
|---|---|---|
| `GET` | `/admin/urls` | List only the calling user's short URLs |
| `POST` | `/admin/create-short-url` | Create a short URL owned by the calling user |

### Testing

```bash
go test ./...
```
> The store tests are integration tests — they require a real, migrated Postgres and a running Redis reachable via `DATABASE_URL`/`REDIS_ADDR`.

## 🎨 Frontend Setup (React + TypeScript)

```bash
cd react_frontend
npm install
npm run dev
```
Available at `http://localhost:5173`. Routes:
- `/` — public shortener (redirects logged-in users to `/admin`)
- `/admin/login`, `/admin/register` — auth pages (same redirect behavior)
- `/admin` — protected dashboard (redirects logged-out users to `/admin/login`)

## 🐳 Running Everything with Docker Compose

```bash
docker compose up --build
```
- Frontend: `http://localhost:5173`
- Backend: `http://localhost:8000`
- Postgres: `localhost:5432`
- Redis: `localhost:6379`

Apply migrations against the Compose-managed Postgres:
```bash
migrate -path go_backend/migrations -database "postgres://urlshortener:devpassword@localhost:5432/urlshortener?sslmode=disable" up
```

Stop everything, keeping data: `docker compose down`
Stop and **wipe** the Postgres volume: `docker compose down -v`

> Inside Compose, services reach each other by service name (`postgres`, `redis`), not `localhost` — confirm your backend's `REDIS_ADDR` is set to `redis:6379` in the Compose environment, not inherited from a local `.env` meant for `go run`.

### Frontend dev workflow while iterating

The frontend's `Dockerfile` builds a static production bundle served by nginx — it has **no hot-reload**. While actively developing the UI, run Postgres/Redis/backend via Compose and the frontend locally instead:
```bash
docker compose up postgres redis backend
# in a second terminal:
cd react_frontend && npm run dev
```
Only rebuild the frontend container (`docker compose up --build frontend`) when testing the real production build.

## 🔍 Checking Data in Postgres

```bash
docker exec -it pg-urlshortener psql -U urlshortener -d urlshortener -c "SELECT * FROM urls;"
```

## 📦 Technology Stack

**Backend:** Gin, PostgreSQL (`pgx`/`pgxpool`), Redis (`go-redis`), `golang-migrate`, JWT (`golang-jwt`), bcrypt, `godotenv`

**Frontend:** React 19, TypeScript, Vite, React Router, Axios, ESLint

## 🚨 Error Handling

- `400` — invalid input or failed validation (human-readable messages, e.g. "password must be at least 8 characters")
- `401` — invalid credentials or missing/expired JWT
- `403` — attempting to act outside your own scope
- `404` — short URL not found in either Redis or Postgres
- `429` — rate limit exceeded on register/login
- `500` — Postgres/Redis connection or query failure

## 🐛 Troubleshooting

**Login "succeeds" but the dashboard never loads** — check `localStorage.getItem("admin_token")` in DevTools; a typo in the key name between where it's set (`Login.tsx`) and where it's read (`ProtectedRoute.tsx`, `Header.tsx`) causes a silent redirect loop with no console error.

**Frontend changes don't appear** — if running via `docker compose up --build`, remember the frontend container serves a static build with no hot-reload; rebuild after every change, or run the frontend locally via `npm run dev` instead.

**Postgres connection refused** — confirm the container is running (`docker ps`); restart with `docker start pg-urlshortener` or recreate it.

**Port conflicts** — remap the host side, e.g. `-p 5433:5432`, and update `DATABASE_URL` accordingly.

## 🔜 Roadmap

- [ ] CI/CD pipeline (GitHub Actions) — lint, test, build, push Docker images
- [ ] Kubernetes manifests
- [ ] `/healthz` endpoint for container/orchestrator probes
- [ ] Redis-backed JWT revocation (denylist) for "log out everywhere"

## 📝 License

This project is for educational purposes.

## 👨‍💻 Author

Created as a self-learning project in Go and React.

Happy URL Shortening! 🎉