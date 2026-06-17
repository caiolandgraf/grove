# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/).
## [2.3.0] - 2026-06-16

### Added
- Quality commands replace Makefile targets in Grove projects:
  - `grove prepare` — format, lint, test, and build before committing
  - `grove check` — lint, test, and build without modifying files
  - `grove fmt`, `grove lint`, `grove lint:fix`
  - `grove test:unit`, `grove test:all` — `go test ./...` with or without `-short`
  - `grove build:binaries` — compile api and atlas to `.grove/bin/`

### Changed
- `grove setup` now scaffolds from the latest [`grove-base`](https://github.com/caiolandgraf/grove-base) template (modular MSC architecture with rate limiting).
- Generators align with grove-base: controllers wire `ratelimiter.Settings`, services use a `Store` interface, and seeders live under `internal/app/database/seeders/`.
- `make:resource` registers new modules with `Wire(b.DB, b.RateLimit)` in `internal/modules/register.go`.
- Seed runner scaffold uses `config.Load()` and `config.InitDatabase()` instead of legacy app globals.
- Documentation updated for rate limiting, module registration, and the seed workflow.

## [2.2.0] - 2026-06-10

### Changed
- `grove setup` now scaffolds from [`grove-base`](https://github.com/caiolandgraf/grove-base) v2.0.0 (modular MSC architecture).
- All generators (`make:model`, `make:service`, `make:controller`, `make:dto`, `make:resource`) now scaffold domain modules under `internal/modules/<domain>/`.
- `make:resource` registers new modules in `internal/modules/register.go` instead of wiring routes manually.
- `make:middleware` scaffolds into `internal/app/middleware/`.
- `make:relations` defaults to scanning `internal/modules/` (`model.go` per domain).
- Documentation updated for the modular project layout, Atlas program mode, and observability stack.

## [2.1.0] - 2026-06-05

### Added
- `grove make:service` command to scaffold Go services.
- Automatic routing wiring for `grove make:resource` (and `grove make:model -r`) inside `internal/routes/routes.go`.

### Changed
- `grove make:controller` to use service dependency injection under the MCS (Model-Controller-Service) architecture.
- `grove make:model` repository struct to satisfy the repository interfaces.
- `grove make:resource` to also scaffold the corresponding service layer.

## [2.0.0] - 2026-05-01

### Added
- `grove up` to start Docker Compose and then run `grove dev`.
- `grove down` to stop Docker Compose services.
- `grove start` to build and run the compiled binary.
- `grove setup` prompt for project name when omitted.
- Interactive observability selection in `grove setup` (Jaeger, Prometheus, Grafana, Loki, Promtail).
- Automatic `.env.example` updates for `OTEL_ENABLED` and `METRICS_ENABLED` during setup.
- Support for `infra/compose.yml` (and fallbacks) when configuring or running Compose.
- Auto-build and caching of the `atlas-provider-gorm` binary for faster `grove make:migration`.
- Port guard support in `grove dev` (auto-detected from `.env`, configurable via `port_guard`).

### Changed
- Default `grove dev` watcher now aligns with the grove.toml template (`cmd/`, `internal/`, and path-based excludes like `internal/tests`).
- `grove migrate` and `grove migrate:status` now print a completion badge on success.
- Project structure documentation updated for the new `internal/app/` layout and `infra/compose.yml`.

### Fixed
- `grove up` now passes `.env` to Docker Compose with `--env-file` when present.
- Panic formatting in `grove dev` now detects stack traces embedded in structured JSON logs.
- Compose cleanup now removes unused volume definitions when observability services are disabled.

[2.3.0]: https://github.com/caiolandgraf/grove/releases/tag/v2.3.0
[2.2.0]: https://github.com/caiolandgraf/grove/releases/tag/v2.2.0
[2.1.0]: https://github.com/caiolandgraf/grove/releases/tag/v2.1.0
[2.0.0]: https://github.com/caiolandgraf/grove/releases/tag/v2.0.0
