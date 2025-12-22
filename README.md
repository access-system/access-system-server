# Access System Server

Server for access control, written in Go.

## Status

DEV

## Features

- Layered Architecture
- Dockerized for easy deployment
- Database initialization scripts
- Unit and integration tests
- Mutual TLS (mTLS) authentication via Nginx

## Prerequisites

- Docker & Docker Compose
- Optional (local dev tooling):
  - Go 1.20+
  - Mockgen for generating test mocks: `go install github.com/golang/mock/mockgen@v1.6.0`

## Getting Started

1) Clone the repository

```
git clone https://github.com/access-system/access-system-server.git
cd access-system-server
```

2) Configure environment

- Copy `.env.example` to `.env` and adjust values if needed.
- Generate TLS certificates for Nginx and clients authentication. You can use the provided scripts to generate self-signed certs for development:

```
sudo chmod +x gen-test-nginx-certs.sh
./gen-test-nginx-certs.sh
```

```
sudo chmod +x gen-test-client-certs.sh
./gen-test-client-certs.sh
```

- Move the generated client certs to client project.

3) Run in Docker

```
go mod download
sudo chmod +x run.sh
./run.sh dev
```

- Reverse proxy (Nginx) exposes HTTPS on `https://localhost` (port 443)
- The upstream app listens on `:8081` inside the Docker network (not published directly)

4) Health check

- GET `https://localhost/health` → 200 OK

## Running Tests

To run all tests using Docker profile:

```
./run.sh test
```

Notes:
- The script runs `go generate ./...` on the host if `internal/mocks` is missing. Ensure `mockgen` is installed locally.
- Test Postgres is exposed on host port `5433`.

## API Reference

Vector size: All endpoints that receive a vector require exactly 512 `float32` values. Requests with other sizes return 400/500 depending on the layer; the service validates length.

### Main API
Base URL through Nginx:
- `https://localhost/api/v1`

#### POST /api/v1/embedding — Add embedding
Body:
- `name` (string, required)
- `vector` (array<float32>, required, length 512)

Responses:
- 201 Created
- 400 Bad Request (invalid body)
- 500 Internal Server Error

Example:
```
curl https://localhost/api/v1/embedding \
  --cert client_crt/client.crt --key client_crt/client.key -k \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","vector":[/* 512 floats */]}' \
  -i
```

#### POST /api/v1/embedding/validate — Validate embedding
Body:
- `vector` (array<float32>, required, length 512)

Responses:
- 200 OK with JSON body:
  - `id` (int64)
  - `name` (string)
  - `vector` (array<float32>)
  - `accuracy` (float32)
- 404 Not Found (no relevant match)
- 400 Bad Request (invalid body)
- 500 Internal Server Error

Example:
```
curl https://localhost/api/v1/embedding/validate \
  --cert client_crt/client.crt --key client_crt/client.key -k \
  -H "Content-Type: application/json" \
  -d '{"vector":[/* 512 floats */]}' \
  -i
```

#### DELETE /api/v1/embedding — Delete embedding
Body:
- `id` (int64, required)

Responses:
- 200 OK
- 400 Bad Request (invalid body)
- 500 Internal Server Error

Example:
```
curl https://localhost/api/v1/embedding \
  --cert client_crt/client.crt --key client_crt/client.key -k \
  -H "Content-Type: application/json" \
  -X DELETE \
  -d '{"id":1}' \
  -i
```

### Admin API
Base URL through Nginx:
- `https://localhost/api/admin` → proxies to `/api/v1/admin` upstream

Endpoints:
- POST `/embedding` — Add embedding
  - Body: `{ "name": string, "vector": float32[512] }`
  - 201, 400, 500
- GET `/embedding/:id` — Get embedding by ID
  - 200 with `{ id, name, vector }`, 400 (bad id), 500
- GET `/embeddings` — List all embeddings
  - 200 with `[{ id, name, vector }, ...]`, 500
- PUT `/embedding` — Update embedding
  - Body: `{ "id": int64, "name": string, "vector": float32[512] }`
  - 200, 400, 500
- DELETE `/embedding` — Delete embedding
  - Body: `{ "id": int64 }`
  - 200, 400, 500

Examples:
```
# List embeddings
curl https://localhost/api/admin/embeddings \
  --cert client_crt/client.crt --key client_crt/client.key -k -i

# Get embedding by ID
curl https://localhost/api/admin/embedding/1 \
  --cert client_crt/client.crt --key client_crt/client.key -k -i
```

## Environment

Environment variables (see `.env.example`):
- `POSTGRES_HOST` — Postgres hostname (container name in dev)
- `POSTGRES_PORT` — Postgres port
- `POSTGRES_DB` — Database name
- `POSTGRES_USER` — DB user
- `POSTGRES_PASSWORD` — DB password
- `POSTGRES_TEST_HOST`, `POSTGRES_TEST_PORT`, `POSTGRES_TEST_DB`, `POSTGRES_TEST_USER`, `POSTGRES_TEST_PASSWORD` — Test DB settings
- `PGADMIN_DEFAULT_EMAIL`, `PGADMIN_DEFAULT_PASSWORD` — PgAdmin (if enabled)

Database initialization runs from `docker/db/scripts/init.sql`.

## Project Structure

- `cmd/` — Entry point (main.go)
- `internal/` — Application logic
  - `cfg/` — Configuration
  - `client/` — External clients
  - `domain/` — Domain models
  - `handler/` — HTTP handlers
  - `mocks/` — Test mocks
  - `repository/` — Data access
  - `router/` — Routing
  - `service/` — Business logic
- `docker/` — Dockerfiles and DB scripts
- `scripts/` — Utility scripts

## Troubleshooting

- TLS/mTLS errors:
  - Ensure you pass `--cert client_crt/client.crt --key client_crt/client.key`.
  - For self-signed local certs, add `-k` or trust the CA.
- 404 on validate:
  - Means no relevant match found (similarity threshold is enforced at the DB layer: accuracy > 0.58).
- Vector length errors:
  - Vectors must be exactly 512 `float32` values.

## Copyright

@ 2025 NJSC "K.Zhubanov Aktobe regional university". All rights reserved.
Use, copying, modification, and distribution of this code are prohibited without the written permission of NJSC "K.Zhubanov Aktobe Regional University".