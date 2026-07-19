# URL Shortener

A full-stack URL shortener application built with Go (backend) and React (frontend). Converts long URLs into short, shareable links, backed by **PostgreSQL** as the source of truth with **Redis** as a caching layer for fast redirects.

## 🚀 Features

- **URL Shortening** — convert long URLs into compact, unique short codes
- **URL Redirect** — access original URLs through short codes, served from cache when possible
- **RESTful API** — simple backend API for URL operations
- **Modern UI** — React + TypeScript frontend
- **CORS Enabled** — cross-origin requests supported for local development
- **PostgreSQL Storage** — durable, queryable persistence for URL mappings
- **Redis Caching** — cache-aside layer in front of Postgres for low-latency redirects
- **Dockerized** — backend, frontend, Postgres, and Redis all run via Docker Compose
- **Schema Migrations** — versioned SQL migrations via `golang-migrate`

## 🏗️ Architecture

```
React Frontend  →  Go Backend (Gin)  →  Redis (cache)  →  PostgreSQL (source of truth)
```

Write path: a new short URL is saved to **Postgres first**, then written to **Redis** as a cache entry.

Read path (redirect): the backend checks **Redis first**. On a cache hit, it redirects immediately. On a miss, it falls back to **Postgres**, backfills Redis with the result, and then redirects. This means Redis being empty or restarted never causes a broken redirect — it just costs one extra Postgres lookup until the cache warms back up.

## 📋 Prerequisites

- Go 1.21 or higher
- Node.js 16 or higher (with npm)
- Docker Desktop (recommended — runs Postgres, Redis, and both services without installing them natively)
- [`golang-migrate`](https://github.com/golang-migrate/migrate) CLI, for applying schema migrations

## 🏗️ Project Structure

```
url_shortener/
├── go_backend/                # Go backend application
│   ├── main.go                 # Entry point — loads .env, initializes store, starts server
│   ├── go.mod / go.sum
│   ├── Dockerfile               # Multi-stage build for the backend image
│   ├── .env.example              # Template for local environment variables
│   ├── handler/                 # Request handlers
│   │   └── handlers.go
│   ├── routes/                  # API route definitions
│   │   └── handle_urls.go
│   ├── shortener/                # URL shortening logic (SHA-256 + Base58)
│   │   ├── shorturl_generator.go
│   │   └── shorturl_generator_test.go
│   ├── store/                    # Storage layer — Postgres + Redis
│   │   ├── store_service.go        # Postgres (source of truth) + Redis (cache) logic
│   │   └── store_service_test.go
│   └── migrations/                # Versioned SQL schema migrations
│       ├── 000001_create_urls_table.up.sql
│       └── 000001_create_urls_table.down.sql
│
├── react_frontend/              # React TypeScript frontend
│   ├── src/
│   ├── public/
│   ├── package.json
│   ├── Dockerfile                 # Multi-stage build, served via nginx
│   ├── vite.config.ts
│   └── tsconfig.json
│
└── docker-compose.yml            # Runs postgres, redis, backend, frontend together
```

## 🔧 Backend Setup (Go)

### Environment variables

Copy the example file and fill in your local values:
```bash
cd go_backend
cp .env.example .env
```

`.env` (not committed to git):
```
DATABASE_URL=postgres://urlshortener:devpassword@localhost:5432/urlshortener?sslmode=disable
```

### Running Postgres and Redis locally (without Docker Compose)

If you want to run the Go server directly with `go run` instead of via Compose, start Postgres and Redis as standalone containers first:
```bash
docker run --name pg-urlshortener -e POSTGRES_USER=urlshortener -e POSTGRES_PASSWORD=devpassword -e POSTGRES_DB=urlshortener -p 5432:5432 -d postgres:16-alpine
docker run --name redis-urlshortener -p 6379:6379 -d redis:7-alpine
```

### Applying migrations

```bash
migrate -path migrations -database "postgres://urlshortener:devpassword@localhost:5432/urlshortener?sslmode=disable" up
```

### Installation

```bash
cd go_backend
go mod download
go mod tidy
```

### Running the Server

```bash
go run main.go
```

The backend will be available at `http://localhost:8000`

### API Endpoints

**1. Home Endpoint**
- Method: `GET`
- URL: `/`
- Description: Returns a welcome message

**2. Create Short URL**
- Method: `POST`
- URL: `/create-short-url`
- Request Body:
  ```json
  { "long_url": "https://example.com/very/long/url" }
  ```
- Response:
  ```json
  {
    "message": "short url created successfully",
    "short_url": "http://localhost:8000/abc123"
  }
  ```

**3. Redirect to Original URL**
- Method: `GET`
- URL: `/:shortUrl`
- Description: Checks Redis, falls back to Postgres on a cache miss, then redirects
- Response: `302` redirect to the original URL, or `404` if the short code doesn't exist in either store

### Testing

```bash
go test ./...
```

> ⚠️ The store tests are integration tests — they require a real, migrated Postgres and a running Redis, reachable via `DATABASE_URL`. Make sure both containers are up and the migration has been applied before running tests, or the test binary's `init()` will panic on connection failure.

Run with coverage:
```bash
go test -cover ./...
```

## 🎨 Frontend Setup (React + TypeScript)

```bash
cd react_frontend
npm install
npm run dev
```

The frontend will be available at `http://localhost:5173`

Build for production:
```bash
npm run build
```

Lint:
```bash
npm run lint
```

## 🐳 Running Everything with Docker Compose

The simplest way to run the full stack — Postgres, Redis, backend, and frontend — together:

```bash
docker compose up --build
```

- Frontend: `http://localhost:5173`
- Backend: `http://localhost:8000`
- Postgres: `localhost:5432` (still reachable from your host for `psql`/migrations)
- Redis: `localhost:6379`

Apply migrations against the Compose-managed Postgres (same command as local dev, since the port is still exposed to the host):
```bash
migrate -path go_backend/migrations -database "postgres://urlshortener:devpassword@localhost:5432/urlshortener?sslmode=disable" up
```

Stop everything, keeping data:
```bash
docker compose down
```

Stop everything and **wipe** the Postgres data volume:
```bash
docker compose down -v
```

> Inside the Compose network, services reach each other by service name, not `localhost` — the backend's `DATABASE_URL` in `docker-compose.yml` points at host `postgres`, not `localhost`, for exactly this reason.

## 🔍 Checking Data in Postgres

```bash
docker exec -it pg-urlshortener psql -U urlshortener -d urlshortener -c "SELECT * FROM urls;"
```

Or drop into an interactive shell:
```bash
docker exec -it pg-urlshortener psql -U urlshortener -d urlshortener
```
```sql
\dt                     -- list tables
\d urls                 -- describe the urls table
SELECT * FROM urls ORDER BY created_at DESC LIMIT 10;
```

## 📦 Technology Stack

**Backend**
- Gin — web framework for Go
- PostgreSQL (via `pgx`/`pgxpool`) — durable storage, source of truth
- Redis (`go-redis`) — cache-aside layer for fast redirects
- `golang-migrate` — versioned SQL schema migrations
- Base58 — URL-safe encoding for short codes
- `godotenv` — loads local `.env` for development

**Frontend**
- React 19
- TypeScript
- Vite
- Axios
- ESLint

## 🔄 Workflow

1. User enters a long URL in the React frontend
2. Frontend sends a `POST` request to `/create-short-url`
3. Backend generates a unique short code (SHA-256 hash + Base58 encoding)
4. Backend writes the mapping to **Postgres** first, then caches it in **Redis**
5. User receives the short URL and can share it
6. On visiting the short URL, the backend checks **Redis** first; on a miss it reads **Postgres**, backfills Redis, then redirects

## 🚨 Error Handling

- `400 Bad Request` — invalid input (missing or malformed JSON)
- `404 Not Found` — short URL does not exist in either Redis or Postgres
- `500 Internal Server Error` — Postgres or Redis connection/query failure

## 🐛 Troubleshooting

**Postgres connection refused**
Confirm the container is actually running:
```bash
docker ps
```
If it's not listed, start it (`docker start pg-urlshortener`) or recreate it with the `docker run` command above.

**Password contains special characters (e.g. `@`)**
URL-encode it in the connection string — `@` becomes `%40` — or avoid special characters in local dev passwords entirely.

**Port already in use**
- Postgres: map to a different host port, e.g. `-p 5433:5432`, and update `DATABASE_URL` accordingly
- Backend: change the port in `go_backend/main.go` (default: `8000`)
- Frontend: Vite will prompt to use a different port if `5173` is in use

**CORS Issues**
CORS is already enabled in the backend for local development. If issues persist, verify the frontend's URL is allowed in the backend's CORS configuration.

## 🔜 Roadmap

- [ ] CI/CD pipeline (GitHub Actions) — lint, test, build, and push Docker images
- [ ] Kubernetes manifests — Deployments, Services, ConfigMaps/Secrets, Ingress
- [ ] `/healthz` endpoint for container/orchestrator liveness and readiness probes

## 📝 License

This project is for educational purposes.

## 👨‍💻 Author

Created as a self-learning project in Go and React.

Happy URL Shortening! 🎉