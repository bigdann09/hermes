# Hermes — Notification Service

Hermes is a Go-based notification platform that accepts notification requests over HTTP, queues them through Kafka, and dispatches them across multiple delivery channels (database, email, push). It uses PostgreSQL for persistence, Redis for caching, and Mailhog for local email testing.

## Architecture

```
Client
  │
  ▼
┌─────────────┐     publish      ┌───────┐     consume     ┌──────────────────┐
│  API (HTTP) │ ───────────────► │ Kafka │ ──────────────► │ Notification     │
│  hermes     │                  └───────┘                 │ Worker           │
└──────┬──────┘                                            └────────┬─────────┘
       │                                                             │
       │ read/write                                                  │ dispatch
       ▼                                                             ▼
┌─────────────┐                                            ┌──────────────────┐
│ PostgreSQL  │ ◄────────────────────────────────────────│ Channels:        │
└─────────────┘                                            │ database, email, │
       ▲                                                   │ push             │
       │ cache                                              └──────────────────┘
┌─────────────┐                                                    │
│ Redis       │                                                    ▼
└─────────────┘                                              Mailhog (SMTP)
```

**Request flow**

1. Client sends `POST /api/v1/notifications/` with a user ID, type, title, and message.
2. API validates the request, loads the user and their notification preferences, and publishes a message to the `notifications` Kafka topic.
3. The notification worker consumes the message and dispatches it to enabled channels:
   - **database** — persists the notification in PostgreSQL
   - **email** — sends an HTML email via SMTP (Mailhog in development)
   - **push** — logs the push delivery (stub)
   - **websocket** — queued when enabled in preferences (not yet implemented)

## Tech Stack

| Layer        | Technology                          |
| ------------ | ----------------------------------- |
| Language     | Go 1.25+                            |
| HTTP         | Gin                                 |
| Database     | PostgreSQL 15, GORM                 |
| Cache        | Redis                               |
| Messaging    | Apache Kafka (KRaft mode)           |
| Email        | gomail + Mailhog (local dev)        |
| Logging      | Zap                                 |
| Config       | Viper                               |
| Migrations   | [golang-migrate](https://github.com/golang-migrate/migrate) |
| Containers   | Docker Compose                      |

## Project Structure

```
.
├── cmd/
│   ├── api/                    # HTTP API entrypoint
│   ├── cli/                    # CLI (seed command)
│   └── workers/notifications/  # Kafka consumer / dispatcher worker
├── internal/
│   ├── config/                 # Environment configuration
│   ├── controllers/            # HTTP handlers
│   ├── dtos/                   # Request/response DTOs
│   ├── infrastructure/
│   │   ├── cache/              # Redis client
│   │   ├── database/           # DB connection, migrations, seeders
│   │   └── kafka/              # Kafka producer & consumer
│   ├── models/                 # GORM models
│   ├── repositories/           # Data access layer
│   ├── routes/                 # Route registration
│   ├── server/                 # HTTP server lifecycle
│   └── services/
│       └── notification/
│           ├── channels/       # Delivery channel implementations
│           ├── dispatcher.service.go
│           └── notification.service.go
├── pkgs/                       # Shared packages (logger, pagination, validators)
├── docker-compose.yml
├── Dockerfile
├── Makefile
└── .env.example
```

## Prerequisites

- **Go** 1.25 or later
- **Docker** and **Docker Compose**
- **golang-migrate** CLI (optional, for local migrations outside Docker)
  ```bash
  go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
  ```

## Environment Variables

Copy the example file and fill in the values:

```bash
cp .env.example .env
```

| Variable          | Description                              | Example (local Docker)        |
| ----------------- | ---------------------------------------- | ----------------------------- |
| `PORT`            | API port                                 | `3002`                        |
| `ENVIRONMENT`     | Runtime environment                      | `development`                 |
| `NAME`            | Application name                         | `hermes`                      |
| `DOMAIN`          | Application domain                       | `localhost`                   |
| `URL`             | Base application URL                     | `http://localhost:3002`       |
| `ALLOWED_ORIGINS` | CORS allowed origins                     | `*`                           |
| `LOG_LEVEL`       | Log level (`debug`, `info`, `warn`)      | `debug`                       |
| `DB_HOST`         | PostgreSQL host                          | `localhost` (host) / `postgres` (container) |
| `DB_PORT`         | PostgreSQL port                          | `5432`                        |
| `DB_USER`         | PostgreSQL user                          | `postgres`                    |
| `DB_PASS`         | PostgreSQL password                      | `postgres`                    |
| `DB_NAME`         | PostgreSQL database name                 | `notifications`               |
| `DB_SSL`          | SSL mode                                 | `disable`                     |
| `REDIS_HOST`      | Redis host                               | `localhost` / `redis`         |
| `REDIS_PORT`      | Redis port                               | `6379`                        |
| `REDIS_USER`      | Redis username (optional)                |                               |
| `REDIS_PASS`      | Redis password (optional)                |                               |
| `KAFKA_BROKERS`   | Kafka broker addresses (comma-separated) | `localhost:9093` / `kafka:9092` |
| `KAFKA_RETRIES`   | Kafka producer retry count               | `3`                           |
| `MAIL_HOST`       | SMTP host                                | `localhost` / `mailhog`       |
| `MAIL_PORT`       | SMTP port                                | `1025`                        |
| `MAIL_USER`       | SMTP username (optional)                 |                               |
| `MAIL_PASS`       | SMTP password (optional)                 |                               |
| `MAIL_FROM`       | Sender email address                     | `no-reply@hermes.com`       |

> **Note:** Inside Docker Compose, service hostnames (`postgres`, `redis`, `kafka`, `mailhog`) are used automatically. When running the API or CLI on your host machine against Docker services, use the published host ports (check with `docker compose ps`).

## Docker Setup

### 1. Configure environment

```bash
cp .env.example .env
# Edit .env with your values (see table above)
```

### 2. Start infrastructure and services

```bash
docker compose up --build -d
```

This starts:

| Service               | Container name          | Purpose                        |
| --------------------- | ----------------------- | ------------------------------ |
| `backend`             | `hermes`                | HTTP API                       |
| `notifications`       | `notification-worker`   | Kafka consumer / dispatcher    |
| `postgres`            | `notification-postgres` | Database                       |
| `redis`               | `redis`                 | Cache                          |
| `kafka`               | `kafka`                 | Message broker                 |
| `mailhog`             | `mailhog`               | Local SMTP + web UI            |
| `migrate`             | `migrate`               | One-shot migration runner      |

Verify everything is running:

```bash
docker compose ps
curl http://localhost:3002/api/v1/health
# {"message":"OK"}
```

### 3. Run database migrations (Docker)

Migrations live in `internal/infrastructure/database/migrations/`:

| Version | File                                      | Creates                          |
| ------- | ----------------------------------------- | -------------------------------- |
| 000001  | `000001_users.up.sql`                     | `users` table                    |
| 000002  | `000002_notifications.up.sql`             | `notifications` table + enum     |
| 000003  | `000003_notification_preferences.up.sql`  | `notification_preferences` table |

**Option A — Makefile (recommended)**

```bash
make docker_migrate
```

This runs the dedicated `migrate` service defined in `docker-compose.yml`, which connects to the `postgres` container on the internal Docker network.

**Option B — Docker Compose directly**

```bash
docker compose run --rm migrate
```

**Option C — Manual one-off container**

```bash
docker run --rm \
  -v "$(pwd)/internal/infrastructure/database/migrations:/migrations" \
  --network notification_notification \
  migrate/migrate \
  -path=/migrations \
  -database "postgresql://postgres:postgres@notification-postgres:5432/notifications?sslmode=disable" \
  up
```

> Replace credentials and network name if your `.env` values differ. Find the network name with `docker network ls | grep notification`.

**Rollback one migration**

```bash
docker compose run --rm migrate \
  -database "postgres://${DB_USER}:${DB_PASS}@postgres:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL}" \
  -path ./migrations down 1
```

**Check migration status**

```bash
docker exec notification-postgres psql -U postgres -d notifications -c "\dt"
```

Expected tables: `users`, `notifications`, `notification_preferences`, `schema_migrations`.

### 4. Seed the database (Docker)

The seed command creates a demo user and matching notification preferences.

**Option A — Makefile (recommended)**

```bash
make docker_seed
```

**Option B — Docker Compose directly**

```bash
docker compose run --rm --no-deps backend ./cli seed
```

The `--no-deps` flag avoids restarting dependent services; the backend container already has `DB_HOST=postgres` configured.

**What gets seeded**

| Data                     | Details                                              |
| ------------------------ | ---------------------------------------------------- |
| User                     | John Doe — `johndoe@gmail.com`                       |
| Notification preferences | Email, push, and websocket enabled for the seeded user |

After seeding, retrieve the user ID for API testing:

```bash
docker exec notification-postgres psql -U postgres -d notifications \
  -c "SELECT id, email FROM users;"
```

### 5. Full Docker bootstrap (first-time setup)

Run these commands in order:

```bash
cp .env.example .env          # configure values
docker compose up --build -d  # start all services
make docker_migrate           # apply migrations
make docker_seed              # seed demo data
```

## Local Development (without Docker for the app)

You can run the API and worker on your host while using Docker only for infrastructure:

```bash
# Start dependencies only
docker compose up -d postgres redis kafka mailhog

# Point .env at host-mapped ports (check with docker compose ps)
# DB_HOST=localhost  DB_PORT=<mapped postgres port>
# REDIS_HOST=localhost  KAFKA_BROKERS=localhost:9093

make migrate_up   # requires golang-migrate installed locally
make seed         # builds and runs ./cli seed
make run          # starts the API
```

In a second terminal, run the worker:

```bash
go run cmd/workers/notifications/main.go
```

## Makefile Commands

| Command            | Description                                      |
| ------------------ | ------------------------------------------------ |
| `make run`         | Run the API locally                              |
| `make seed`        | Build CLI and seed database (local)              |
| `make migrate_up`  | Apply all pending migrations (local)             |
| `make migrate_down`| Roll back one migration (local)                  |
| `make migrate_create name=<name>` | Create a new migration file       |
| `make migrate_force version=<n>`  | Force migration version           |
| `make docker_run`  | Start stack with rebuild + file watch            |
| `make docker_down` | Stop and remove containers                       |
| `make docker_migrate` | Run migrations inside Docker                |
| `make docker_seed` | Seed database inside Docker                      |

## API Reference

Base URL: `http://localhost:3002/api/v1`

### Health Check

```
GET /health
```

**Response `200`**
```json
{ "message": "OK" }
```

---

### Create Notification

Queues a notification for asynchronous delivery.

```
POST /notifications/
Content-Type: application/json
```

**Request body**

| Field      | Type   | Required | Description                              |
| ---------- | ------ | -------- | ---------------------------------------- |
| `user_id`  | UUID   | yes      | Target user ID                           |
| `type`     | string | yes      | `system` or `social`                     |
| `title`    | string | yes      | Notification title                       |
| `message`  | string | yes      | Notification body                        |
| `metadata` | object | no       | Arbitrary JSON metadata                  |

**Example**

```bash
curl -X POST http://localhost:3002/api/v1/notifications/ \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "<USER_ID>",
    "type": "system",
    "title": "Welcome",
    "message": "Hello from Hermes"
  }'
```

**Response `202`**
```json
{
  "code": 202,
  "status": "Accepted",
  "message": "notification queued successfully",
  "data": {
    "user_id": "...",
    "email": "johndoe@gmail.com",
    "title": "Welcome",
    "message": "Hello from Hermes",
    "type": "system",
    "channels": ["database", "email", "push", "websocket"]
  }
}
```

---

### List Notifications

```
GET /notifications/?user_id=<UUID>&page=1&limit=10
```

**Query parameters**

| Parameter  | Required | Default      | Description                    |
| ---------- | -------- | ------------ | ------------------------------ |
| `user_id`  | yes      | —            | Filter by user UUID            |
| `page`     | no       | `1`          | Page number                    |
| `limit`    | no       | `10`         | Items per page                 |
| `type`     | no       | —            | Filter by `system` or `social` |
| `is_read`  | no       | `false`      | Filter by read status          |
| `order_by` | no       | `created_at` | Sort column                    |
| `order`    | no       | `desc`       | Sort direction                 |

**Example**

```bash
curl "http://localhost:3002/api/v1/notifications/?user_id=<USER_ID>&type=system"
```

**Response `200`**
```json
{
  "code": 200,
  "status": "OK",
  "message": "notifications retrieved successfully",
  "data": {
    "data": [
      {
        "id": "...",
        "user_id": "...",
        "type": "system",
        "title": "Welcome",
        "metadata": "{}",
        "read_at": "0001-01-01T00:00:00Z",
        "created_at": "2026-05-31T20:01:31.035085Z"
      }
    ],
    "metadata": {
      "total": 1,
      "has_next_page": false,
      "has_previous_page": false,
      "page": 1,
      "limit": 10
    }
  }
}
```

## End-to-End Test

After completing the Docker bootstrap steps:

```bash
# 1. Get seeded user ID
USER_ID=$(docker exec notification-postgres psql -U postgres -d notifications -t -c "SELECT id FROM users LIMIT 1;" | tr -d ' \n')

# 2. Send a notification
curl -X POST http://localhost:3002/api/v1/notifications/ \
  -H "Content-Type: application/json" \
  -d "{\"user_id\":\"$USER_ID\",\"type\":\"system\",\"title\":\"Test\",\"message\":\"Hello\"}"

# 3. Confirm it was persisted
curl "http://localhost:3002/api/v1/notifications/?user_id=$USER_ID"

# 4. Check email in Mailhog UI
open http://localhost:8025

# 5. Check worker logs
docker logs notification-worker --tail 20
```

## Useful Docker Commands

```bash
# View logs
docker compose logs -f backend
docker compose logs -f notifications

# Restart a single service
docker compose restart backend notifications

# Stop everything
docker compose down

# Stop and remove volumes (fresh database)
docker compose down -v
```

## Troubleshooting

**Migrations fail with "connection refused"**
- Ensure Postgres is healthy: `docker compose ps postgres`
- Run migrations after Postgres is up: `make docker_migrate`

**Seed fails with "relation does not exist"**
- Migrations have not been applied yet. Run `make docker_migrate` first.

**Notification returns `404 user not found`**
- Run `make docker_seed` or verify the user exists:
  ```bash
  docker exec notification-postgres psql -U postgres -d notifications -c "SELECT * FROM users;"
  ```

**Email channel errors in worker logs**
- Confirm Mailhog is running: `docker compose ps mailhog`
- Confirm worker has mail env vars: `docker inspect notification-worker --format '{{range .Config.Env}}{{println .}}{{end}}' | grep MAIL`
- Recreate the worker if env vars are missing: `docker compose up -d --force-recreate notifications`

**Kafka consumer not processing messages**
- Check Kafka health: `docker compose ps kafka`
- Inspect worker logs: `docker logs notification-worker`
- Verify the topic exists: `docker exec kafka kafka-topics --bootstrap-server kafka:9092 --list`

**API cannot connect to Redis/Kafka/Postgres when run locally**
- Use the host-mapped ports from `docker compose ps`, not the internal container ports.

## License

Private project — all rights reserved.
