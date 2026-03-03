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
[![Release](https://img.shields.io/badge/release-v1.3.0-c82838?style=flat-square)](https://github.com/caiolandgraf/grove/releases/tag/v1.3.0)
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
| [gest](https://github.com/caiolandgraf/gest) | Jest-inspired testing framework for Go |
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
| `grove make:model <Name> -m` | Scaffold model + migration (Atlas diff) |
| `grove make:model <Name> -c` | Scaffold model + controller |
| `grove make:model <Name> -d` | Scaffold model + DTO |
| `grove make:model <Name> -mc` | Scaffold model + migration + controller |
| `grove make:model <Name> -mcd` | Scaffold model + migration + controller + DTO |
| `grove make:model <Name> -r` | Full resource — shorthand for `-mcd` |
| `grove make:controller <Name>` | Scaffold a fuego controller in `internal/controllers/` |
| `grove make:dto <Name>` | Scaffold DTO request/response files in `internal/dto/` |
| `grove make:middleware <Name>` | Scaffold an HTTP middleware in `internal/middleware/` |
| `grove make:migration <name>` | Generate a SQL migration via Atlas diff |
| `grove make:resource <Name>` | Scaffold model + migration + controller + DTO in one shot |

> **Name singularization:** all generator commands accept plural or mixed-case names and convert them automatically. `Books`, `books`, and `Book` all produce the same `Book` model, `books` table, and `create_books_table` migration.

### Testing

| Command | Description |
|---|---|
| `grove make:test <Name>` | Scaffold a new [gest](https://github.com/caiolandgraf/gest) spec file in `internal/tests/` |
| `grove test` | Run all spec files via gest |
| `grove test -c` | Run specs and display a per-suite coverage report |
| `grove test -w` | Watch mode — re-run specs on every save (no external tools required) |
| `grove test -wc` | Watch mode + coverage report |

> `grove make:test` automatically creates `internal/tests/main.go` (the gest entrypoint) if it does not exist yet.

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
│   ├── app/                     # Shared singletons (DB, config)
│   ├── controllers/             # fuego route handlers
│   ├── dto/                     # Request and response types
│   ├── middleware/              # HTTP middlewares
│   ├── models/                  # GORM models
│   ├── routes/                  # Route registration
│   └── tests/                   # gest spec files
│       ├── main.go              # gest entrypoint (auto-created)
│       └── post_spec.go         # example spec
├── migrations/                  # Atlas SQL migrations
├── .env.example                 # Committed env template
├── atlas.hcl                    # Atlas configuration
├── grove.toml                   # Grove configuration (optional)
└── go.mod
```

The `internal/` boundary is intentional — it prevents external packages from importing your application internals, keeping the codebase clean as it grows.

---

## Typical Workflow

```bash
# 1. Scaffold a full resource (model + migration + controller + DTO)
grove make:resource Post

# 2. Add fields to the model and DTO, then apply the migration
grove migrate

# 3. Register routes in internal/routes/routes.go
#    fuego.Post(s, "/posts", controllers.CreatePost)

# 4. Write tests for your new resource
grove make:test Post

# 5. Run the test suite
grove test -c
```

---

## Hot Reload with `grove dev`

Grove ships a built-in hot reload watcher — no Air, no external tools required.

```bash
grove dev
```

On every `.go` save Grove recompiles and restarts your binary automatically. A debounce window collapses burst saves into a single rebuild, and newly created subdirectories are picked up at runtime without restarting the watcher.

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

> **Tip:** spec files (`*_spec.go`) and the `tests` directory are always excluded from the watcher so a test save never triggers an application rebuild.

---

## Testing with gest

Grove uses [gest](https://github.com/caiolandgraf/gest) — a Jest-inspired testing framework for Go with beautiful output and zero dependencies.

```bash
# Scaffold a spec file (creates internal/tests/main.go if needed)
grove make:test UserService

# Run all specs
grove test

# Run with per-suite coverage report
grove test -c

# Watch mode — re-run specs on every save (no Air or external tools)
grove test -w
```

Each spec file lives in `internal/tests/` and self-registers via `init()`:

```go
// internal/tests/user_service_spec.go
package main

import "github.com/caiolandgraf/gest/gest"

func init() {
    s := gest.Describe("UserService")

    s.It("should create a user", func(t *gest.T) {
        // ...
        t.Expect(user.ID).Not().ToBeNil()
    })

    gest.Register(s)
}
```

> **Note:** gest uses `_spec.go` instead of `_test.go` because the Go toolchain reserves `_test.go` for `go test`. gest runs via `go run`, so any other suffix works.

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

## Contributing

Contributions are welcome — bug fixes, new commands, documentation improvements and ideas alike.

1. Fork the repository
2. Make your change and build with `make grove-build`
3. Open a pull request with a clear description

See the full documentation at **[caiolandgraf.github.io/grove](https://caiolandgraf.github.io/grove/)**.

---

## License

MIT © [Caio Landgraf](https://github.com/caiolandgraf)