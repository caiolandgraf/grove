# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/).
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

[2.1.0]: https://github.com/caiolandgraf/grove/releases/tag/v2.1.0
[2.0.0]: https://github.com/caiolandgraf/grove/releases/tag/v2.0.0
