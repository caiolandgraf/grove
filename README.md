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
[![Release](https://img.shields.io/badge/release-v1.5.0-c82838?style=flat-square)](https://github.com/caiolandgraf/grove/releases/tag/v1.5.0)
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

# 2. Enter the project and configure your environment
cd my-api && cp .env.example .env

# 3. Start the development server with built-in hot reload
grove dev
```

Your API is running at `http://localhost:8080`.  
The OpenAPI / Swagger UI is available at `http://localhost:8080/swagger` automatically.

---

## Commands

### Generators

| Command | Description |
|---|---|
| `grove make:model <Name>` | Scaffold a GORM model in `internal/models/` |
| `grove make:model <Name> -c` | Scaffold model + controller |
| `grove make:model <Name> -d` | Scaffold model + DTO |
| `grove make:model <Name> -cd` | Scaffold model + controller + DTO |
| `grove make:model <Name> -r` | Full resource — shorthand for `-cd` |
| `grove make:controller <Name>` | Scaffold a fuego controller in `internal/controllers/` |
| `grove make:dto <Name>` | Scaffold DTO request/response files in `internal/dto/` |
| `grove make:middleware <Name>` | Scaffold an HTTP middleware in `internal/middleware/` |
| `grove make:migration <name>` | Generate a SQL migration via Atlas diff (after editing your model) |
| `grove make:resource <Name>` | Scaffold model + controller + DTO in one shot |

> **Name singularization:** all generator commands accept plural or mixed-case names and convert them automatically. `Books`, `books`, and `Book` all produce the same `Book` model and `books` table.

> **Migration workflow:** migrations are **not** generated automatically when scaffolding a model or resource. Add your fields to the model first, then run `grove make:migration <name>` to let Atlas diff your schema and generate the correct SQL. This ensures the migration reflects the fields you actually defined.

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
| `grove dev` | Hot reload — watch, build & restart on every save (no external tools required) |
| `grove dev:air` | Start the development server using Air for hot-reload |
| `grove build` | Compile the application binary to `./bin/app` |
| `grove setup <project-name>` | Scaffold a new project from the official template |

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
│   └── api/
│       └── main.go              # Entry point
├── internal/
│   ├── app/                     # Shared singletons (DB, Redis, Session, Metrics)
│   ├── config/                  # Infrastructure initializers (DB, Redis, OTel, etc.)
│   ├── controllers/             # fuego route handlers
│   ├── database/                # Generic GORM repository
│   ├── dto/                     # Request and response types
│   ├── middleware/              # HTTP middlewares (CORS, session, observability)
│   ├── models/                  # GORM models
│   ├── routes/                  # Route registration
│   └── tests/                   # gest test files
│       └── user_test.go         # Example test (generated by grove make:test)
├── migrations/                  # Atlas SQL migrations
├── infra/                       # Observability stack config (Prometheus, Grafana, Loki, Jaeger)
├── .env.example                 # Committed environment template
├── atlas.hcl                    # Atlas configuration
├── docker-compose.yml           # Full observability stack
└── grove.toml                   # Grove dev server configuration
```

The `internal/` boundary is intentional — it prevents external packages from importing your application internals, keeping the codebase clean as it grows.

| Directory | Purpose |
|---|---|
| `cmd/api/` | Application entry point — wires singletons, routes and starts the server |
| `internal/app/` | Shared singletons: DB, Redis, session store, metrics — initialized once at startup |
| `internal/config/` | Infrastructure initializers for DB, Redis, OpenTelemetry and other external services |
| `internal/controllers/` | fuego route handlers — one file per resource, OpenAPI inferred automatically |
| `internal/database/` | Generic GORM repository (`Repository[T]`) used by all models |
| `internal/dto/` | Request and response structs — decoupled from GORM models |
| `internal/middleware/` | HTTP middlewares: CORS, session, observability, auth, etc. |
| `internal/models/` | GORM models with typed repository accessors |
| `internal/routes/` | Route registration — fuego typed routes wired to controllers |
| `internal/tests/` | gest v2 test files — standard `*_test.go` files, auto-created by `grove make:test` |
| `migrations/` | Versioned Atlas SQL migration files + `atlas.sum` |
| `infra/` | Observability stack configuration: Prometheus, Grafana, Loki, Jaeger |
| `docker-compose.yml` | Spins up the full observability stack locally with a single command |
| `grove.toml` | Optional Grove configuration — `[dev]` section for `grove dev` |

---

## Typical Workflow

```bash
# 1. Scaffold a full resource (model + controller + DTO)
grove make:resource Post

# 2. Add your fields to the model
#    edit internal/models/post.go → add Title, Body, etc.

# 3. Add request/response fields to the DTO
#    edit internal/dto/post-dto.go

# 4. Generate the migration — Atlas diffs your model against the DB
grove make:migration create_posts_table

# 5. Apply the migration
grove migrate

# 6. Register routes in internal/routes/routes.go
#    fuego.Post(s, "/posts", controllers.CreatePost)

# 7. Write tests for your new resource
grove make:test Post

# 8. Run the test suite
grove test -c

# or with plain go test
go test ./internal/tests/...
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
watch_dirs  = ["."]
exclude     = [".grove", "vendor", "node_modules", "tests"]
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