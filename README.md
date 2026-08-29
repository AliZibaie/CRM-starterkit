# CRM StarterKit

A production-ready **Modular Monolith** for customer relationship management, e-commerce, support, and internal operations. This starter kit includes a **Backend (Go)**, **Admin Panel (Next.js)**, **Storefront (Next.js)**, and **Docker/Nginx** deployment configuration.

The project is designed for businesses that sell **services**, **physical/digital products**, and **course licenses** (`course_license`). The catalog module supports `physical_product`, `service`, `course_license`, and `digital` types; orders, payments (Zarinpal/Stripe), CRM, CMS, and support are integrated on the same platform.

The architecture is **modular**: each capability (auth, catalog, crm, etc.) is an independent module with its own routes, migrations, permissions, and use cases. Modules are enabled or disabled via `backend/modules.config.go`; the core (`internal/core`) provides shared infrastructure (HTTP, DB, Cache, Storage, Event Bus, RBAC).

---

## Key Features (Active Modules)

The list below is extracted directly from `backend/modules.config.go`:

| Module | Description |
|--------|-------------|
| **rbac** | Roles, permissions, assign permission to role and role to user |
| **user** | User CRUD, role assignment, default admin user creation |
| **auth** | Login, refresh token, logout (JWT) |
| **cronbuilder** | Scheduled jobs via UI with condition/action engine on DB tables |
| **sms** | SMS sending (provider: Kavenegar) |
| **notification** | Multi-channel dispatch (in_app, email, sms, push) + event bus subscription |
| **announcement** | Internal announcements and broadcast |
| **support** | Departments, tickets, messages, WebSocket |
| **banner** | Promotional banners with scheduling and position |
| **catalog** | Categories, products, variants, attributes, inventory + public storefront API |
| **cart** | User shopping cart |
| **discount** | Discount codes (preview, CRUD) |
| **order** | Orders, checkout, status updates, refund |
| **payment** | Zarinpal and Stripe gateways, init/verify/callback/refund |
| **media** | Upload, thumbnail, file deletion (local/S3) |
| **cms** | Pages, blocks, menus, SEO, FAQ (public API + admin) |
| **review** | Product reviews (public + admin moderation) |
| **wishlist** | User wishlist |
| **crm** | Leads, contacts, deals, pipelines, tasks, segments, Customer 360 |
| **setting** | System key-value settings |
| **auditlog** | Change logging and viewing (event `audit.change`) |
| **requestlog** | HTTP request logs, stats, suspicious IPs |

---

## Tech Stack

| Layer | Technology | Version / Details |
|-------|------------|-------------------|
| Backend language | Go | `1.25.0` (`backend/go.mod`) |
| HTTP framework | Fiber | `v2.52.5` |
| ORM | GORM | `v1.25.12` |
| DB drivers | postgres, mysql, sqlserver, sqlite | GORM drivers + golang-migrate |
| Cache | Redis, Memcached | `go-redis/v9`, `gomemcache` |
| Queue / Worker | Asynq | `v0.24.1` (Redis) |
| Auth | JWT | `golang-jwt/jwt/v5` |
| Config | Viper | `v1.19.0` |
| Logging | Zap | `v1.27.0` |
| Validation | go-playground/validator | `v10.23.0` |
| Storage | Local filesystem, AWS S3 | aws-sdk-go-v2 |
| API docs | Swaggo | mounted at `/swagger/*` |
| Admin Panel | Next.js + React | `16.3.2` / `19.2.8` |
| Storefront | Next.js + React | `16.3.2` / `19.2.8` |
| UI (Admin) | Tailwind CSS 4, Radix UI, Recharts | — |
| State / Data (Front) | TanStack Query 5, Zustand | — |
| Reverse proxy | Nginx | `1.27-alpine` (Docker) |
| Container base | golang:1.25-alpine, node:22-alpine | — |

---

## Architecture

### Modular Monolith

- Each module implements the `module.Module` interface: `Init`, `RegisterRoutes`, `RegisterMigrations`, `RegisterPermissions`.
- `module.Loader` registers and initializes modules from `EnabledModules`.
- API routes are mounted under the **`/api/v1`** prefix (`server.APIGroup()`).
- Table schema: **GORM AutoMigrate** on each module startup + SQL migrations in `backend/migrations/` via `cmd/migrate`.

### Enabling / Disabling Modules

File `backend/modules.config.go`:

```go
var EnabledModules = []module.Module{
	rbac.New(),
	user.New(),
	auth.New(),
	// cronbuilder.New(),  // disable cron module
	catalog.New(),
	order.New(),
	// ...
}
```

A module removed from the slice no longer registers routes or permissions (note: existing DB tables remain).

### Folder Structure

```
CRM-starterkit/
├── backend/
│   ├── cmd/
│   │   ├── api/              # HTTP server
│   │   ├── worker/           # Asynq worker
│   │   └── migrate/          # golang-migrate CLI
│   ├── config/               # .env loading
│   ├── internal/
│   │   ├── core/
│   │   │   ├── cache/        # redis, memcached
│   │   │   ├── database/     # postgres, mysql, sqlserver, sqlite
│   │   │   ├── event/        # event bus
│   │   │   ├── http/         # server, middleware, response
│   │   │   ├── module/       # registry, loader, contract
│   │   │   ├── storage/      # local, s3
│   │   │   └── validator/
│   │   └── modules/          # 22 domain modules
│   │       ├── auth/
│   │       ├── catalog/
│   │       ├── crm/
│   │       └── ...
│   ├── migrations/           # SQL migrations (core)
│   ├── modules.config.go     # active module list
│   ├── Dockerfile
│   ├── Dockerfile.worker
│   └── .env.example
├── admin-panel/
│   ├── app/                  # Next.js App Router
│   ├── components/
│   ├── lib/
│   └── .env.example
├── storefront/
│   ├── app/
│   ├── lib/
│   └── .env.example
└── deploy/
    ├── docker-compose.yml
    ├── nginx/
    │   ├── default.conf      # unified reverse proxy
    │   ├── admin.conf
    │   └── storefront.conf
    └── .env.example
```

---

## Prerequisites

| Tool | Recommended Version | Purpose |
|------|---------------------|---------|
| Go | 1.25+ | Backend |
| Node.js | 22+ | Admin Panel and Storefront |
| npm | Bundled with Node | Frontend install and build |
| Docker + Docker Compose | v2+ | Recommended production-like setup |
| PostgreSQL / MySQL / SQL Server / SQLite | One of | Database (per `DB_DRIVER`) |
| Redis | 7+ | Cache, queue worker, rate limiting |
| Memcached | Optional | Cache alternative (`CACHE_DRIVER=memcached`) |

---

## Quick Start with Docker (Recommended)

```bash
git clone <repository-url>
cd CRM-starterkit/deploy
cp .env.example .env
# Edit JWT_SECRET and DB_PASSWORD in .env if needed
docker compose --profile postgres up -d --build
```

> **Note:** `deploy/.env.example` sets `COMPOSE_PROFILES=postgres,redis`. For MySQL: `COMPOSE_PROFILES=mysql,redis`. For SQL Server: `COMPOSE_PROFILES=sqlserver,redis`. For Memcached instead of Redis: set `CACHE_DRIVER=memcached` and use the `memcached` profile.

Migrations run automatically when the `api` container starts (`scripts/docker-entrypoint.sh` → `migrate -command up` → `api`).

### Ports and Paths (Docker + Nginx)

| Service | Internal Port | Host Access |
|---------|---------------|-------------|
| **nginx** | 80 | `http://localhost` |
| api | 8080 | `http://localhost/api/` |
| admin-panel | 3000 | `http://localhost/admin/` |
| storefront | 3000 | `http://localhost/` (root) |
| uploads (static) | — | `http://localhost/uploads/` |
| Swagger | — | `http://localhost/swagger/` |
| health / ready | — | `http://localhost/health`, `/ready` |
| postgres | 5432 | Docker network only (profile `postgres`) |
| mysql | 3306 | Docker network only (profile `mysql`) |
| redis | 6379 | Docker network only |

---

## Manual Setup (Development)

### Backend

```bash
cd backend
cp .env.example .env
# Configure DB, Redis, and JWT_SECRET for your local environment
go mod download
go run ./cmd/migrate -command up
go run ./cmd/api
```

The API is available at `http://localhost:8080` (or `APP_PORT`). Local uploads are served from `http://localhost:8080/uploads/`.

Worker (separate terminal):

```bash
cd backend
go run ./cmd/worker
```

Additional migrate commands:

```bash
go run ./cmd/migrate -command down
go run ./cmd/migrate -command version
go run ./cmd/migrate -command force -version 1
go run ./cmd/migrate -command up -steps 1
```

### Admin Panel

```bash
cd admin-panel
cp .env.example .env.local
npm install
npm run dev
```

Default: `http://localhost:3000`

**Default admin login** (created on user module startup):

- Email: `admin@example.com`
- Password: `ChangeMe123!`

### Storefront

```bash
cd storefront
cp .env.example .env.local
npm install
npm run dev
```

Default: `http://localhost:3001` (the `dev` script in `package.json` uses port 3001)

---

## Environment Variables

### Backend (`backend/.env.example`)

| Variable | Description | Default / Example |
|----------|-------------|-------------------|
| `APP_NAME` | Service name | `CRM-StarterKit` |
| `APP_ENV` | Environment (`development` / `production`) | `development` |
| `APP_DEBUG` | Debug mode | `true` |
| `APP_URL` | Backend base URL | `http://localhost:8080` |
| `APP_PORT` | HTTP port | `8080` |
| `DB_DRIVER` | DB driver: `postgres` \| `mysql` \| `sqlserver` \| `sqlite` | `postgres` |
| `DB_HOST` | DB host | `localhost` |
| `DB_PORT` | DB port | `5432` |
| `DB_NAME` | Database name | `starterkit` |
| `DB_USER` | DB user | `postgres` |
| `DB_PASSWORD` | DB password | `secret` |
| `DB_SSLMODE` | SSL (postgres) | `disable` |
| `DB_TIMEZONE` | Timezone | `UTC` |
| `DB_SQLITE_PATH` | SQLite file path | `./storage/database.sqlite` |
| `CACHE_DRIVER` | `redis` \| `memcached` | `redis` |
| `REDIS_HOST` | Redis host | `localhost` |
| `REDIS_PORT` | Redis port | `6379` |
| `REDIS_PASSWORD` | Redis password | *(empty)* |
| `REDIS_DB` | Redis DB number | `0` |
| `MEMCACHED_HOST` | Memcached address | `localhost:11211` |
| `CACHE_PREFIX` | Cache key prefix | `crm:` |
| `CACHE_DEFAULT_TTL` | Default TTL | `1h` |
| `JWT_SECRET` | JWT signing key | *(required)* |
| `JWT_ACCESS_TTL` | Access token lifetime | `15m` |
| `JWT_REFRESH_TTL` | Refresh token lifetime | `168h` |
| `JWT_ISSUER` | Issuer | `crm-starterkit` |
| `CORS_ALLOW_ORIGINS` | Allowed origins (CSV) | `*` |
| `CORS_ALLOW_METHODS` | Allowed methods | `GET,POST,PUT,PATCH,DELETE,OPTIONS` |
| `CORS_ALLOW_HEADERS` | Allowed headers | `Origin,Content-Type,Accept,Authorization` |
| `CORS_ALLOW_CREDENTIALS` | CORS credentials | `false` |
| `CORS_MAX_AGE` | Preflight max-age | `86400` |
| `RATE_LIMIT_ENABLED` | Enable rate limiting | `true` |
| `RATE_LIMIT_MAX` | Max requests per window | `100` |
| `RATE_LIMIT_WINDOW` | Rate limit window | `1m` |
| `STORAGE_DRIVER` | `local` \| `s3` | `local` |
| `STORAGE_LOCAL_PATH` | Local storage path | `./storage/uploads` |
| `STORAGE_PUBLIC_BASE_URL` | Public file URL | `http://localhost:8080/uploads` |
| `AWS_REGION` | S3 region | `us-east-1` |
| `AWS_ACCESS_KEY_ID` | AWS key | *(empty)* |
| `AWS_SECRET_ACCESS_KEY` | AWS secret | *(empty)* |
| `AWS_S3_BUCKET` | Bucket name | *(empty)* |
| `AWS_S3_ENDPOINT` | Custom S3 endpoint | *(empty)* |
| `AWS_S3_FORCE_PATH_STYLE` | S3 path-style | `false` |
| `AWS_S3_PUBLIC_BASE_URL` | Public S3 URL | *(empty)* |
| `LOG_LEVEL` | Log level | `info` |
| `LOG_ENCODING` | `json` \| `console` | `json` |
| `WORKER_CONCURRENCY` | Worker concurrency | `10` |
| `QUEUE_REDIS_ADDR` | Redis queue address | `localhost:6379` |
| `QUEUE_REDIS_PASSWORD` | Redis queue password | *(empty)* |
| `QUEUE_REDIS_DB` | Redis queue DB | `1` |
| `MIGRATIONS_PATH` | SQL migrations path | `./migrations` |
| `PAYMENT_CALLBACK_BASE_URL` | Payment callback base URL | `http://localhost:8080/api/v1/payments/callback` |
| `ZARINPAL_MERCHANT_ID` | Zarinpal merchant ID | *(empty)* |
| `ZARINPAL_SANDBOX` | Zarinpal sandbox mode | `true` |
| `STRIPE_SECRET_KEY` | Stripe secret key | *(empty)* |
| `STRIPE_WEBHOOK_SECRET` | Stripe webhook secret | *(empty)* |
| `STRIPE_CURRENCY` | Stripe currency | `usd` |
| `SMS_PROVIDER` | SMS provider | `kavenegar` |
| `KAVENEGAR_API_KEY` | Kavenegar API key | *(empty)* |
| `KAVENEGAR_SENDER` | Sender line | *(empty)* |
| `SMTP_HOST` | SMTP host | *(empty)* |
| `SMTP_PORT` | SMTP port | `587` |
| `SMTP_USER` | SMTP user | *(empty)* |
| `SMTP_PASSWORD` | SMTP password | *(empty)* |
| `SMTP_FROM` | From address | *(empty)* |
| `ENV_FILE` | Alternate env file path (optional) | — |

### Admin Panel (`admin-panel/.env.example`)

| Variable | Description | Default |
|----------|-------------|---------|
| `NEXT_PUBLIC_API_URL` | Backend API URL | `http://localhost:8080/api/v1` |
| `NEXT_PUBLIC_WS_URL` | WebSocket URL (support) | `ws://localhost:8080/api/v1` |

### Storefront (`storefront/.env.example`)

| Variable | Description | Default |
|----------|-------------|---------|
| `NEXT_PUBLIC_API_URL` | Backend API URL | `http://localhost:8080/api/v1` |
| `NEXT_PUBLIC_WS_URL` | WebSocket URL | `ws://localhost:8080/api/v1` |
| `NEXT_PUBLIC_SITE_URL` | Public storefront URL | `http://localhost:3001` |

### Deploy (`deploy/.env.example`)

| Variable | Description | Default |
|----------|-------------|---------|
| `COMPOSE_PROFILES` | Docker Compose profiles | `postgres,redis` |
| `JWT_SECRET` | Production JWT secret | *(change this)* |
| `DB_DRIVER` | DB driver | `postgres` |
| `DB_HOST` | DB host in Docker | `postgres` |
| `DB_PORT` | DB port | `5432` |
| `DB_NAME` | DB name | `starterkit` |
| `DB_USER` | DB user | `postgres` |
| `DB_PASSWORD` | DB password | `secret` |
| `DB_ROOT_PASSWORD` | MySQL root password | `rootsecret` |
| `CACHE_DRIVER` | Cache driver | `redis` |
| `REDIS_PASSWORD` | Redis password | *(empty)* |
| `ZARINPAL_MERCHANT_ID` | Merchant ID | *(empty)* |
| `STRIPE_SECRET_KEY` | Stripe key | *(empty)* |
| `STRIPE_WEBHOOK_SECRET` | Stripe webhook | *(empty)* |
| `KAVENEGAR_API_KEY` | SMS API key | *(empty)* |
| `KAVENEGAR_SENDER` | SMS sender | *(empty)* |
| `SMTP_*` | Email | *(empty)* |

---

## Database and Cache Selection

### Database (`DB_DRIVER`)

Allowed values in `backend/internal/core/database/driver.go` and `config.Validate()`:

| `DB_DRIVER` | Description |
|-------------|-------------|
| `postgres` | PostgreSQL (GORM + golang-migrate) |
| `mysql` | MySQL 8 |
| `sqlserver` | Microsoft SQL Server |
| `sqlite` | SQLite (file at `DB_SQLITE_PATH`) |

To switch: edit `.env` only; no code changes required. In Docker, use the appropriate profile (`postgres`, `mysql`, `sqlserver`).

### Cache (`CACHE_DRIVER`)

| `CACHE_DRIVER` | Description |
|----------------|-------------|
| `redis` | Default; host/port from `REDIS_*` |
| `memcached` | Uses `MEMCACHED_HOST` |

---

## API Documentation (Swagger)

Swagger UI is mounted via `github.com/gofiber/swagger`:

| Environment | URL |
|-------------|-----|
| Local (no Docker) | `http://localhost:8080/swagger/index.html` |
| Docker (via Nginx) | `http://localhost/swagger/index.html` |

API base path in annotations: `/`. Module endpoints are under `/api/v1/...`. Authentication: header `Authorization: Bearer <access_token>`.

---

## Main Endpoints

> Common prefix: **`/api/v1`**. **Public** endpoints require no JWT; others need a Bearer token and the appropriate permission.

### Core (outside `/api/v1`)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/` | Welcome message |
| GET | `/health` | Health check |
| GET | `/ready` | Readiness (DB ping) |
| GET | `/swagger/*` | Swagger UI |
| GET | `/uploads/*` | Static files (only `STORAGE_DRIVER=local`) |

### auth

| Method | Path | Auth |
|--------|------|------|
| POST | `/auth/login` | — |
| POST | `/auth/refresh` | — |
| POST | `/auth/logout` | JWT |

### rbac

| Method | Path |
|--------|------|
| POST | `/rbac/roles` |
| GET | `/rbac/roles` |
| GET | `/rbac/roles/:id` |
| POST | `/rbac/permissions` |
| GET | `/rbac/permissions` |
| PUT | `/rbac/permissions/:id` |
| POST | `/rbac/assign-permission` |
| POST | `/rbac/revoke-permission` |
| POST | `/rbac/assign-role` |
| POST | `/rbac/permissions/sync` |

### user

| Method | Path |
|--------|------|
| POST | `/users` |
| GET | `/users` |
| GET | `/users/:id` |
| PUT | `/users/:id` |
| DELETE | `/users/:id` |
| PUT | `/users/:id/roles` |

### cronbuilder

| Method | Path |
|--------|------|
| GET | `/cronbuilder/tables` |
| GET | `/cronbuilder/tables/:table/columns` |
| GET | `/cronbuilder/actions` |
| POST | `/cronbuilder/conditions/build` |
| POST | `/cronbuilder/jobs` |
| GET | `/cronbuilder/jobs` |
| GET | `/cronbuilder/jobs/:id` |
| PUT | `/cronbuilder/jobs/:id` |
| DELETE | `/cronbuilder/jobs/:id` |
| POST | `/cronbuilder/jobs/:id/run` |
| GET | `/cronbuilder/jobs/:id/logs` |

### sms

| Method | Path |
|--------|------|
| POST | `/sms/send` |

### notification

| Method | Path |
|--------|------|
| POST | `/notifications/dispatch` |
| GET | `/notifications/me` |
| POST | `/notifications/:id/read` |

### announcement

| Method | Path |
|--------|------|
| POST | `/announcements` |
| GET | `/announcements` |
| GET | `/announcements/:id` |
| PUT | `/announcements/:id` |
| DELETE | `/announcements/:id` |
| POST | `/announcements/:id/broadcast` |

### support

| Method | Path |
|--------|------|
| POST | `/support/departments` |
| GET | `/support/departments` |
| PUT | `/support/departments/:id` |
| DELETE | `/support/departments/:id` |
| POST | `/support/tickets` |
| GET | `/support/tickets` |
| GET | `/support/tickets/:id` |
| POST | `/support/tickets/:id/reply` |
| POST | `/support/tickets/:id/assign` |
| PUT | `/support/tickets/:id/status` |
| POST | `/support/tickets/:id/close` |
| GET | `/support/ws/:id` | WebSocket |

### banner

| Method | Path |
|--------|------|
| POST | `/banners` |
| GET | `/banners` |
| GET | `/banners/positions` |
| PUT | `/banners/:id` |
| DELETE | `/banners/:id` |
| PUT | `/banners/:id/schedule` |
| PUT | `/banners/:id/position` |

### catalog

| Method | Path | Auth |
|--------|------|------|
| GET | `/catalog/public/items` | Public |
| GET | `/catalog/public/items/:slug` | Public |
| GET | `/catalog/public/categories` | Public |
| GET | `/catalog/public/categories/:slug` | Public |
| POST | `/catalog/categories` | JWT |
| GET | `/catalog/categories` | JWT |
| PUT | `/catalog/categories/:id` | JWT |
| DELETE | `/catalog/categories/:id` | JWT |
| POST | `/catalog/items` | JWT |
| GET | `/catalog/items` | JWT |
| GET | `/catalog/items/:id` | JWT |
| PUT | `/catalog/items/:id` | JWT |
| DELETE | `/catalog/items/:id` | JWT |
| POST | `/catalog/items/:id/attributes` | JWT |
| POST | `/catalog/items/:id/variants` | JWT |
| GET | `/catalog/inventory/:variant_id` | JWT |
| PUT | `/catalog/inventory/:variant_id` | JWT |

### cart

| Method | Path |
|--------|------|
| GET | `/cart` |
| POST | `/cart/items` |
| DELETE | `/cart/items/:item_id` |
| DELETE | `/cart` |

### discount

| Method | Path |
|--------|------|
| POST | `/discounts` |
| GET | `/discounts` |
| POST | `/discounts/preview` |
| PUT | `/discounts/:id` |
| DELETE | `/discounts/:id` |

### order

| Method | Path |
|--------|------|
| POST | `/orders` |
| GET | `/orders` |
| GET | `/orders/:id` |
| PUT | `/orders/:id/status` |
| POST | `/orders/:id/checkout` |
| POST | `/orders/:id/discount` |
| POST | `/orders/:id/refund` |

### payment

| Method | Path | Auth |
|--------|------|------|
| GET | `/payments/gateways` | JWT |
| GET | `/payments/failed` | JWT |
| POST | `/payments/init` | JWT |
| POST | `/payments/verify` | JWT |
| GET | `/payments/callback/:gateway` | Public (callback) |
| POST | `/payments/refund/:order_id` | JWT |

Implemented gateways: **zarinpal** (with `ZARINPAL_MERCHANT_ID`) and **stripe** (with `STRIPE_SECRET_KEY`).

### media

| Method | Path |
|--------|------|
| POST | `/media/upload` |
| GET | `/media` |
| POST | `/media/:id/thumbnail` |
| DELETE | `/media/:id` |

### cms

| Method | Path | Auth |
|--------|------|------|
| GET | `/cms/pages/:slug` | Public |
| GET | `/cms/faqs` | Public |
| GET | `/cms/seo` | Public |
| GET | `/cms/menus` | Public |
| POST | `/cms/pages` | JWT |
| GET | `/cms/pages` | JWT |
| GET | `/cms/pages/id/:id` | JWT |
| PUT | `/cms/pages/:id` | JWT |
| PUT | `/cms/pages/:id/blocks` | JWT |
| POST | `/cms/menus` | JWT |
| GET | `/cms/menus/all` | JWT |
| POST | `/cms/seo` | JWT |
| POST | `/cms/faqs` | JWT |

### review

| Method | Path | Auth |
|--------|------|------|
| GET | `/reviews/items/:item_id` | Public |
| GET | `/reviews/admin` | JWT |
| POST | `/reviews` | JWT |
| PUT | `/reviews/:id/visibility` | JWT |
| DELETE | `/reviews/:id/admin` | JWT |
| DELETE | `/reviews/:id` | JWT |

### wishlist

| Method | Path |
|--------|------|
| GET | `/wishlist` |
| POST | `/wishlist` |
| DELETE | `/wishlist/:item_id` |

### crm

| Method | Path |
|--------|------|
| POST | `/crm/leads` |
| GET | `/crm/leads` |
| POST | `/crm/leads/:id/convert` |
| GET | `/crm/contacts` |
| POST | `/crm/contacts` |
| PUT | `/crm/contacts/:id` |
| DELETE | `/crm/contacts/:id` |
| GET | `/crm/deals` |
| POST | `/crm/deals` |
| POST | `/crm/deals/:id/stage` |
| GET | `/crm/activities` |
| POST | `/crm/assign-owner` |
| POST | `/crm/activities` |
| POST | `/crm/tasks/reminders` |
| GET | `/crm/customers/:contact_id/360` |
| POST | `/crm/pipelines` |
| GET | `/crm/pipelines` |
| POST | `/crm/tags` |
| POST | `/crm/tags/attach` |
| POST | `/crm/segments` |
| GET | `/crm/segments` |

### setting

| Method | Path |
|--------|------|
| GET | `/settings` |
| GET | `/settings/:key` |
| PUT | `/settings` |

### auditlog

| Method | Path |
|--------|------|
| GET | `/audit-logs` |

### requestlog

| Method | Path |
|--------|------|
| GET | `/request-logs` |
| GET | `/request-logs/stats` |
| GET | `/request-logs/suspicious/:ip` |

---

## RBAC (Roles and Permissions)

1. Each module defines permission keys in `RegisterPermissions()`.
2. On startup, the **rbac** module syncs permissions from the registry and creates the **admin** role with all permissions.
3. JWT access tokens include `roles` and `permissions`; the `RequireAnyPermission` middleware checks JWT first, then the DB (for stale tokens).
4. The default admin user has the admin role.

### Registering a New Permission (real example — media module)

```go
func (m *Module) RegisterPermissions() []module.Permission {
	return []module.Permission{
		{Key: "media.upload", Name: "Upload Media", Description: "Upload media files", Module: "media"},
		{Key: "media.read", Name: "Read Media", Description: "View media files", Module: "media"},
		{Key: "media.delete", Name: "Delete Media", Description: "Delete media files", Module: "media"},
	}
}
```

After adding permissions: restart the API (auto-sync) or call `POST /api/v1/rbac/permissions/sync`, then update the relevant role.

---

## Cron Builder (Non-Technical UI)

Admin panel path: **`/cron-jobs/builder`**

1. Log in to the admin panel (requires `cronbuilder.*` permissions).
2. Enter a **job name** and **cron expression** (e.g. `0 * * * *`).
3. Select a **table** from the tables list (API: `GET /cronbuilder/tables`).
4. Build conditions on columns with the **Condition Builder** (equals, greater than, etc.).
5. Click build condition → preview the condition JSON.
6. Select an **action type** (list: `GET /cronbuilder/actions`).
7. Click **Create Job** → job is saved in `GET /cronbuilder/jobs`.
8. Manual run: `POST /cronbuilder/jobs/:id/run` — logs at `/cron-jobs/logs`.

---

## Testing

There are currently no `*_test.go` files in the backend and no `*.test.ts(x)` files in the frontend. After adding tests:

```bash
cd backend && go test ./...
cd admin-panel && npm run lint
cd storefront && npm run lint
```

---

## Deployment / Production

### Sensitive Variables

- `JWT_SECRET`: long, unique random string
- `DB_PASSWORD`, `REDIS_PASSWORD`
- `ZARINPAL_MERCHANT_ID`, `STRIPE_SECRET_KEY`, `KAVENEGAR_API_KEY`
- `SMTP_*` for email notifications

### HTTPS and Nginx

`deploy/nginx/default.conf` currently only exposes **HTTP on port 80**. For production:

- Add TLS certificates (Let's Encrypt) and a `listen 443 ssl` block
- Update `APP_URL`, `STORAGE_PUBLIC_BASE_URL`, `PAYMENT_CALLBACK_BASE_URL`, and `NEXT_PUBLIC_*` to your real HTTPS domain
- `client_max_body_size 50m` is already set for media uploads

### CORS

Use `CORS_ALLOW_ORIGINS`; in production, specify admin and storefront domains instead of `*`. Set `CORS_ALLOW_CREDENTIALS` if you need cookies/credentials.

### Docker Production

```bash
cd deploy
docker compose --profile postgres up -d --build
```

The `backend_storage` volume persists uploaded files. Nginx mounts the same volume at `/var/www/uploads/`.

### Manual Image Builds

```bash
docker build -f backend/Dockerfile -t crm-api ../backend
docker build -f backend/Dockerfile.worker -t crm-worker ../backend
docker build -f admin-panel/Dockerfile -t crm-admin ../admin-panel
docker build -f storefront/Dockerfile -t crm-storefront ../storefront
```

---

## License

This project is a starter kit. Review the upstream repository license before commercial use.

## Contributing

1. Fork and create a feature branch
2. Verify changes with `go build ./...` and `npm run build` in admin/storefront
3. Open a Pull Request describing affected modules/endpoints

---

## UI Pages

### Admin Panel (sample routes)

`/dashboard` · `/users` · `/roles-permissions` · `/catalog/products` · `/orders` · `/crm/leads` · `/cms/pages` · `/media-library` · `/cron-jobs/builder` · `/audit-log` · `/request-monitor` · …

### Storefront

`/` · `/products/[slug]` · `/categories/[slug]` · `/cart` · `/checkout` · `/account/login` · `/account/orders` · `/account/tickets` · `/account/licenses` · `/[slug]` (CMS pages)
