<div align="center">

<br />

<pre><code style="color:#c82838">
█▀▀ █▀█ █▀█ █░█ █▀▀
█▄█ █▀▄ █▄█ ▀▄▀ ██▄
</code></pre>

<br />

**Grove is an opinionated Go foundation for building structured, observable, production-ready applications.**

<br />

[![Go Version](https://img.shields.io/badge/go-1.22+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev)
[![Release](https://img.shields.io/badge/release-v2.2.0-c82838?style=flat-square)](https://github.com/caiolandgraf/grove/releases/tag/v2.2.0)
[![License](https://img.shields.io/badge/license-MIT-c82838?style=flat-square)](LICENSE)
[![Docs](https://img.shields.io/badge/docs-caiolandgraf.github.io%2Fgrove-c82838?style=flat-square)](https://caiolandgraf.github.io/grove/)

<br />

[**Documentation**](https://caiolandgraf.github.io/grove/) · [**Quick Start**](#quick-start) · [**Commands**](#commands) · [**Contributing**](#contributing)

<br />

</div>

---

## Overview

Grove is a CLI that scaffolds and manages Go applications following a clean, layered project layout. It wires together [GORM](https://gorm.io), [fuego](https://github.com/go-fuego/fuego) and [Atlas](https://atlasgo.io) so you can generate models, controllers, DTOs, middlewares and migrations in seconds — and focus entirely on your business logic.

| Tool | Role |
|---|---|
| [GORM](https://gorm.io) | ORM & typed repository layer |
| [fuego](https://github.com/go-fuego/fuego) | HTTP router + automatic OpenAPI 3.1 |
| [Atlas](https://atlasgo.io) | Schema migration engine |
| [gest](https://github.com/caiolandgraf/gest) | Jest-inspired testing framework for Go (v2) |
| [air](https://github.com/air-verse/air) _(optional)_ | Hot-reload via `grove dev:air` (not needed for `grove dev`) |

---

## Installation

```bash
go install github.com/caiolandgraf/grove@latest
```

Verify:

```bash
grove -v        # print version
grove --help    # full command reference
```

> **Requirements:** Go 1.22+, [Atlas CLI](https://atlasgo.io/docs) for migration commands.

---

## Quick Start

```bash
# 1. Scaffold a new project from the official template
grove setup my-api
# or: grove setup (prompts for project name + observability)

# 2. Enter the project and configure your environment
cd my-api && cp .env.example .env

# 3. Start infra + dev server (docker compose + hot reload)
grove up

# or start only the dev server if infra is already running
grove dev
```

Your API is running at `http://localhost:8080`.  
The OpenAPI / Swagger UI is available at `http://localhost:8080/swagger` automatically.

---

## Commands

### Generators

| Command | Description |
|---|---|
| `grove make:model <Name>` | Scaffold a GORM model in `internal/modules/<domain>/model.go` |
| `grove make:model <Name> -c` | Scaffold model + controller + OpenAPI docs |
| `grove make:model <Name> -d` | Scaffold model + DTO |
| `grove make:model <Name> -cd` | Scaffold model + controller + DTO + docs |
| `grove make:model <Name> -r` | Full resource — service, controller, DTO, docs + module registration |
| `grove make:controller <Name>` | Scaffold a module controller in `internal/modules/<domain>/` |
| `grove make:service <Name>` | Scaffold a service in `internal/modules/<domain>/service.go` |
| `grove make:seeder <Name>` | Scaffold a database seeder in `internal/database/seeders/` |
| `grove make:seed` | Scaffold a seed runner entrypoint in `cmd/seed/main.go` (used by `grove db:seed`) |
| `grove make:dto <Name>` | Scaffold DTO types in `internal/modules/<domain>/dto.go` |
| `grove make:middleware <Name>` | Scaffold an HTTP middleware in `internal/app/middleware/` |
| `grove make:migration <name>` | Generate a SQL migration via Atlas diff (after editing your model) |
| `grove db:seed` | Run database seeders via `go run ./cmd/seed` |
| `grove make:resource <Name>` | Scaffold a full domain module and register it in `internal/modules/register.go` |
| `grove make:relations` | Infer and add GORM relationships from foreign keys |

> **Name singularization:** all generator commands accept plural or mixed-case names and convert them automatically. `Books`, `books`, and `Book` all produce the same `Book` model and `books` table.

> **Migration workflow:** migrations are **not** generated automatically when scaffolding a model or resource. Add your fields to the model first, then run `grove make:migration <name>` to let Atlas diff your schema and generate the correct SQL. This ensures the migration reflects the fields you actually defined.

#### `grove make:relations`

Infers model relationships from foreign key fields (for example, `UserID`) and adds GORM relation fields automatically.

#### Seeders (`make:seeder`, `make:seed`, `db:seed`)

Grove supports a simple, explicit seeding workflow compatible with `go-project-base`.

- `grove make:seeder <Name>` scaffolds a new seeder file in `internal/database/seeders/` following the interface:

  - `Name() string`
  - `Seed(db *gorm.DB) error`

- `grove make:seed` scaffolds the seed runner entrypoint at `cmd/seed/main.go`. This runner should initialize your app (DB, session, etc.) and call:

  - `internal/database/seeders.Run(app.DB)`

- `grove db:seed` loads `.env` (if present) and runs seeders by executing:

  - `go run ./cmd/seed`

Examples:

```bash
# 1) Create a new seeder
grove make:seeder Users

# 2) (Once per project) create the seed runner entrypoint
grove make:seed

# 3) Run all seeders
grove db:seed

# Optional: customize runner package or env file
grove db:seed --package ./cmd/seed
grove db:seed --env-file .env
```

> Keep seeders idempotent: running `grove db:seed` multiple times should be safe.

By default, it generates only the **has-many** side on the target model (for example, `PaymentMethods []PaymentMethod` on `User`).  
Use `--with-belongs-to` to also generate the **belongs-to** side on the source model (for example, `User *User` on `PaymentMethod`).

| Flag | Description |
|---|---|
| `--path` | Modules directory (default: `internal/modules`) |
| `--dry-run` | Preview inferred relations without writing files |
| `--verbose` | Print detailed relation inference logs |
| `--with-belongs-to` | Also generate belongs-to fields on source models |
| `--model` | Limit processing to specific source model(s). Can be repeated or comma-separated |

Examples:

```bash
# infer relations for all models (has-many only)
grove make:relations

# preview changes with detailed logs
grove make:relations --dry-run --verbose

# generate both sides (has-many + belongs-to)
grove make:relations --with-belongs-to

# process only specific source models
grove make:relations --model PaymentMethod
grove make:relations --model PaymentMethod,Order
grove make:relations --model PaymentMethod --model Order
```

### Testing

| Command | Description |
|---|---|
| `grove make:test <Name>` | Scaffold a new [gest](https://github.com/caiolandgraf/gest) v2 test file in `internal/tests/` |
| `grove test` | Run all tests via the gest CLI (falls back to `go test -v` if gest is not installed) |
| `grove test -c` | Run tests and display a per-suite coverage report |
| `grove test -w` | Watch mode — re-run tests on every save |
| `grove test -wc` | Watch mode + coverage report |

> `grove make:test` generates standard `*_test.go` files with a `func Test<Name>(t *testing.T)` entry point. You can also run `go test ./internal/tests/...` directly at any time.

### Server & Build

| Command | Description |
|---|---|
| `grove up` | Start docker compose and run the dev server |
| `grove down` | Stop docker compose services |
| `grove dev` | Hot reload — watch, build & restart on every save (no external tools required) |
| `grove dev:air` | Start the development server using Air for hot-reload |
| `grove build` | Compile the application binary to `./bin/app` |
| `grove start` | Build and start the application binary |
| `grove setup [project-name]` | Scaffold a new project from the official template (prompts if omitted) |

### Database

| Command | Description |
|---|---|
| `grove migrate` | Apply all pending migrations |
| `grove migrate:rollback` | Rollback the last applied migration |
| `grove migrate:status` | Show migration status |
| `grove migrate:fresh` | Drop all tables and re-apply every migration ⚠️ |
| `grove migrate:hash` | Rehash the `atlas.sum` file |

`grove migrate` formats the Atlas output with Grove's colour palette — each migration version gets a `MIGRATE` badge, SQL statements are syntax-highlighted with the keyword in cyan, and a final summary line shows total time, migrations and statements applied:

```
  Running migrations (atlas migrate apply --env local)

   MIGRATE   20260127143000

    CREATE TABLE  users ( … )
    CREATE INDEX  idx_users_deleted_at ON users(deleted_at)
   OK   18.4ms

   MIGRATE   20260304122639

    ALTER TABLE  "public"."users" DROP CONSTRAINT …
    CREATE TABLE  "public"."books" ( … )
    CREATE INDEX  "idx_books_deleted_at" ON "public"."books" …
   OK   5.2ms

  ────────────────────────────────────────
  81.473ms
  2 migrations
  9 sql statements
```

If all migrations are already applied, Grove prints an `UP TO DATE` badge instead.

### Maintenance

| Command | Description |
|---|---|
| `grove update` | Update Grove project dependencies (gest) to their latest versions |

---

### Shell Completion

```bash
grove completion [bash|zsh|fish|powershell]
```

```bash
# Zsh — persist
echo 'source <(grove completion zsh)' >> ~/.zshrc

# Fish — persist
grove completion fish > ~/.config/fish/completions/grove.fish
```

---

## Project Structure

```
my-api/
├── cmd/
│   ├── api/
│   │   └── main.go                # Application entrypoint
│   ├── atlas/
│   │   └── main.go                # Atlas GORM schema loader
│   └── scalar/
│       └── scalar.go              # Scalar API docs UI handler
├── internal/
│   ├── app/                       # Shared infrastructure
│   │   ├── config/                # PostgreSQL, Redis, OTel, slog, metrics, session
│   │   ├── database/              # Repository[T] + Atlas model registry
│   │   ├── helpers/               # jsonutils, validator
│   │   ├── middleware/            # CORS, session, observability
│   │   ├── types/                 # Shared HTTP response types
│   │   ├── router/                # Declarative OpenAPI route docs
│   │   └── app.go                 # Package documentation
│   ├── modules/
│   │   ├── auth/                  # Auth domain (dto, service, controller, docs)
│   │   ├── users/                 # Users domain (model, dto, service, controller, docs)
│   │   ├── module.go              # Module interface + Boot
│   │   └── register.go            # Module registry
│   └── routes/
│       ├── health.go              # Health check routes
│       └── routes.go              # Global routes + module mounting
├── infra/                         # Prometheus, Grafana, Loki, Jaeger configs
├── logs/                          # App log files (Promtail → Loki)
├── migrations/                    # Atlas SQL migrations
├── doc/
│   └── openapi.json               # Generated OpenAPI spec
├── .air.toml                      # Air hot reload config
├── atlas.hcl                      # Atlas migration config
├── docker-compose.yml             # Full infrastructure stack
├── Makefile                       # Dev commands
└── go.mod                         # Go module definition
```

Each domain under `internal/modules/` is self-contained: **model + repo**, **dto**, **service**, **controller**, and **docs**. Modules wire themselves via `Wire` and register HTTP routes via `Mount`. Add new domains in `internal/modules/register.go`. Models auto-register for Atlas via `init()` in each module's `model.go`.

| Directory | Purpose |
|---|---|
| `cmd/api/` | Application entrypoint — boots config, connects infra, starts the server |
| `cmd/atlas/` | Loads all registered GORM models for Atlas migrations |
| `cmd/scalar/` | Scalar OpenAPI documentation UI handler |
| `internal/app/config/` | DB, Redis, OTel, slog, metrics, and session initializers |
| `internal/app/database/` | Generic `Repository[T]` and Atlas model registry |
| `internal/app/helpers/` | JSON utilities and request validation |
| `internal/app/middleware/` | CORS, session, and observability middlewares |
| `internal/app/types/` | Shared HTTP response types (errors, health, messages) |
| `internal/app/router/` | Declarative OpenAPI docs per endpoint (`router.Doc`) |
| `internal/modules/` | Self-contained domain packages (auth, users, your resources) |
| `internal/modules/register.go` | Module registry — `grove make:resource` adds new domains here |
| `internal/routes/` | Global middleware, health checks, mounts all modules |
| `infra/` | Observability stack configuration |
| `logs/` | Structured JSON log files tailed by Promtail |
| `migrations/` | Versioned Atlas SQL migration files + `atlas.sum` |
| `docker-compose.yml` | PostgreSQL, Redis, Jaeger, Prometheus, Loki, Grafana |
| `atlas.hcl` | Atlas configuration — program mode via `cmd/atlas` |
| `Makefile` | Dev commands: install, run, dev, migrate-up |
| `grove.toml` | Optional `[dev]` section for `grove dev` |

---

## Typical Workflow

```bash
# 1. Scaffold a full domain module
grove make:resource Post

# 2. Add your fields to the model
#    edit internal/modules/posts/model.go → add Title, Body, etc.

# 3. Add request/response fields to the DTO
#    edit internal/modules/posts/dto.go

# 4. Generate the migration — Atlas diffs your model against the DB
grove make:migration create_posts_table

# 5. Apply the migration
grove migrate

# 6. Module is registered automatically in internal/modules/register.go
#    Routes are mounted via posts.Wire(db).Mount(api, session)

# 7. Write tests for your new resource
grove make:test Post

# 8. Run the test suite
grove test -c
```

> **Updating a model later?** Add the new fields to your struct, then run `grove make:migration add_<field>_to_posts` — Atlas will generate an `ALTER TABLE` migration with exactly the diff between the current DB schema and your updated model.

---

## Hot Reload with `grove dev`

Grove ships a built-in hot reload watcher — no Air, no external tools required.

```bash
grove dev
```

On every `.go` save Grove recompiles and restarts your binary automatically. A debounce window collapses burst saves into a single rebuild, and newly created subdirectories are picked up at runtime without restarting the watcher.

> **Tip:** the `internal/tests/` directory is always excluded from the dev watcher so a test save never triggers an application rebuild.

### Output formatting

`grove dev` processes your application's stdout/stderr and formats it intelligently:

**Structured JSON logs** (slog, zap, zerolog) are parsed and rendered as human-readable coloured lines:

```
  08:38:28  INF  Booting application...
  08:38:28  INF  OpenTelemetry initialized  service=grove-app  endpoint=localhost:4318
  08:38:28  ERR  Failed to boot application  error=failed to connect to database: ...
```

**Panics** are captured and rendered as a styled block with the stack trace clearly formatted instead of raw text.

### Startup hints

Grove detects common startup errors and prints an actionable `HINT` immediately below the error:

| Error detected | Hint shown |
|---|---|
| `.env not found` | `cp .env.example .env` |
| `connection refused` / `dial error` / `failed to connect` | `docker compose up -d` |

Each hint is shown once per rebuild — if the error persists after the next file save, the hint appears again.

### Configuration

Configure behaviour via the optional `[dev]` section in `grove.toml` at the project root:

```toml
[dev]
root        = "."
bin         = ".grove/tmp/app"
build_cmd   = "go build -o .grove/tmp/app ./cmd/api/"
watch_dirs  = ["cmd", "internal"]
exclude     = [".grove", "vendor", "node_modules", ".git", "infra", "migrations", "bin", "internal/tests"]
extensions  = [".go"]
debounce_ms = 50
```

All fields are optional. When `grove.toml` is absent or the `[dev]` section is omitted, sensible defaults are applied and `grove dev` works out of the box.

---

## Testing with gest

Grove uses [gest v2](https://github.com/caiolandgraf/gest) — a Jest-inspired testing framework for Go that runs on top of the native `go test` engine for full IDE support, caching and coverage.

```bash
# Scaffold a test file
grove make:test UserService

# Run all tests (beautiful gest CLI output)
grove test

# Run with per-suite coverage report
grove test -c

# Watch mode — re-run tests on every save
grove test -w
```

Each test file lives in `internal/tests/` and follows the standard Go test convention:

```go
// internal/tests/user_service_test.go
package myapp

import (
    "testing"

    "github.com/caiolandgraf/gest/v2/gest"
)

func TestUserService(t *testing.T) {
    s := gest.Describe("UserService")

    s.It("should create a user", func(t *gest.T) {
        // ...
        t.Expect(user.ID).Not().ToBeNil()
    })

    s.Run(t)
}
```

You can also run the tests directly with the standard Go toolchain at any time:

```bash
go test ./internal/tests/...
```

Install the gest CLI globally for the full Jest-style output:

```bash
go install github.com/caiolandgraf/gest/v2/cmd/gest@latest
```

> **Note:** gest v2 uses standard `*_test.go` files and integrates with `go test` — no separate `main.go` entrypoint needed. `grove test` falls back to `go test -v` automatically if the gest CLI is not installed.

---

## Generator Name Singularization

All generator commands automatically singularize the entity name before generating files. This means you can type the name in any form and Grove will always produce consistent output:

| Input | Resolved name | File | Table |
|---|---|---|---|
| `Book` | `Book` | `book.go` | `books` |
| `Books` | `Book` | `book.go` | `books` |
| `books` | `Book` | `book.go` | `books` |
| `BlogPost` | `BlogPost` | `blog_post.go` | `blog_posts` |
| `order_items` | `OrderItem` | `order_item.go` | `order_items` |

---

## Updating dependencies

Run `grove update` inside your project to update Grove-managed dependencies to their latest versions and tidy the module graph:

```bash
grove update
```

This updates the [gest](https://github.com/caiolandgraf/gest) library in your `go.mod`, installs the latest `gest` CLI binary globally, and runs `go mod tidy` automatically. Use `grove update` whenever you want to pull in a newer version.

---

## Contributing

Contributions are welcome — bug fixes, new commands, documentation improvements and ideas alike.

1. Fork the repository
2. Make your change and build with `make grove-build`
3. Open a pull request with a clear description

See the full documentation at **[caiolandgraf.github.io/grove](https://caiolandgraf.github.io/grove/)**.

---

## License

MIT © [Caio Landgraf](https://github.com/caiolandgraf)