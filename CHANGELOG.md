# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [2.0.0] - 2026-05-01

### Added
- `grove up` to start Docker Compose and then run `grove dev`.
- `grove start` to build and run the compiled binary.
- `grove setup` prompt for project name when omitted.
- Interactive observability selection in `grove setup` (Jaeger, Prometheus, Grafana, Loki, Promtail).
- Automatic `.env.example` updates for `OTEL_ENABLED` and `METRICS_ENABLED` during setup.
- Support for `infra/compose.yml` (and fallbacks) when configuring or running Compose.

### Changed
- Default `grove dev` watcher is now scoped to `cmd/` and `internal/` when present, with a larger exclude list for better performance.
- Default debounce increased to 100 ms to reduce rebuild thrash.
- Project structure documentation updated for the new `internal/app/` layout and `infra/compose.yml`.

### Fixed
- `grove up` now passes `.env` to Docker Compose with `--env-file` when present.
- Panic formatting in `grove dev` now detects stack traces embedded in structured JSON logs.
- Compose cleanup now removes unused volume definitions when observability services are disabled.

[2.0.0]: https://github.com/caiolandgraf/grove/releases/tag/v2.0.0
