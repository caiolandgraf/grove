/**
 * Interactive project tree for the Architecture → Project Structure docs page.
 * Mirrors caiolandgraf/go-project-base (v2.0.0) layout and file order.
 */
export const projectStructureTree = [
  // ── cmd/ ──────────────────────────────────────────────────────────────
  {
    name: 'cmd/',
    isDir: true,
    level: 0,
    path: 'cmd/',
    desc: 'Application entry points. Each subdirectory builds a separate binary.',
    tip: 'Primary API: <code>cmd/api/</code>. Atlas schema loader: <code>cmd/atlas/</code>. Scalar UI: <code>cmd/scalar/</code>.'
  },
  {
    name: 'api/',
    isDir: true,
    level: 1,
    path: 'cmd/api/',
    desc: 'Main REST API entry point — boots config, connects infra, registers routes, starts Fuego.'
  },
  {
    name: 'main.go',
    isDir: false,
    level: 2,
    path: 'cmd/api/main.go',
    desc: 'Application entrypoint. Loads .env, initializes logger, OTel, metrics, DB, Redis, session, and starts the HTTP server.',
    code: `package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/caiolandgraf/go-project-base/cmd/scalar"
	"github.com/caiolandgraf/go-project-base/internal/app/config"
	"github.com/caiolandgraf/go-project-base/internal/routes"
	"github.com/go-fuego/fuego"
	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()
	config.InitLogger()

	ctx := context.Background()
	otelShutdown, _ := config.InitOtel(ctx)
	defer func() { _ = otelShutdown(ctx) }()

	metricsHandler, _ := config.InitMetrics()
	db, _ := config.InitDatabase()
	redisPool, _ := config.InitRedis()
	sessionManager := config.InitSessionManager(redisPool)

	s := fuego.NewServer(
		fuego.WithAddr("localhost:8080"),
		fuego.WithEngineOptions(
			fuego.WithOpenAPIConfig(fuego.OpenAPIConfig{
				UIHandler: scalar.NewUI,
			}),
		),
	)

	routes.SetupRoutes(s, db, redisPool, sessionManager, metricsHandler)
	_ = s.Run()
}`
  },
  {
    name: 'atlas/',
    isDir: true,
    level: 1,
    path: 'cmd/atlas/',
    desc: 'Atlas GORM schema loader invoked by atlas.hcl in program mode.',
    tip: 'Models self-register via <code>init()</code> in each module\'s <code>model.go</code>.'
  },
  {
    name: 'main.go',
    isDir: false,
    level: 2,
    path: 'cmd/atlas/main.go',
    desc: 'Loads all registered GORM models and outputs SQL schema for Atlas migrate diff.',
    code: `package main

import (
	"ariga.io/atlas-provider-gorm/gormschema"
	_ "github.com/caiolandgraf/go-project-base/internal/modules"
	"github.com/caiolandgraf/go-project-base/internal/app/database"
)

func main() {
	stmts, err := gormschema.New("postgres").Load(database.All()...)
	// writes schema to stdout for Atlas
	_ = stmts
}`
  },
  {
    name: 'scalar/',
    isDir: true,
    level: 1,
    path: 'cmd/scalar/',
    desc: 'Scalar OpenAPI documentation UI handler for Fuego.'
  },
  {
    name: 'scalar.go',
    isDir: false,
    level: 2,
    path: 'cmd/scalar/scalar.go',
    desc: 'Custom Scalar UI with themed CSS, wired into Fuego OpenAPI config.'
  },

  // ── internal/ ─────────────────────────────────────────────────────────
  {
    name: 'internal/',
    isDir: true,
    level: 0,
    path: 'internal/',
    desc: 'Private application code. Go prevents external packages from importing beyond this boundary.',
    tip: 'Shared infra in <code>internal/app/</code>. Domain logic in <code>internal/modules/</code>.'
  },

  // internal/app/
  {
    name: 'app/',
    isDir: true,
    level: 1,
    path: 'internal/app/',
    desc: 'Shared infrastructure: config, database, helpers, middleware, router, and types.'
  },
  {
    name: 'config/',
    isDir: true,
    level: 2,
    path: 'internal/app/config/',
    desc: 'Infrastructure initializers — PostgreSQL, Redis, OTel, slog, metrics, and SCS sessions.'
  },
  {
    name: 'database.go',
    isDir: false,
    level: 3,
    path: 'internal/app/config/database.go',
    desc: 'PostgreSQL + GORM setup with slog query logging and OTel tracing via otelgorm.'
  },
  {
    name: 'logger.go',
    isDir: false,
    level: 3,
    path: 'internal/app/config/logger.go',
    desc: 'Structured JSON slog handler writing to stdout and logs/app.log.'
  },
  {
    name: 'metrics.go',
    isDir: false,
    level: 3,
    path: 'internal/app/config/metrics.go',
    desc: 'Prometheus metrics via OTel Prometheus exporter.'
  },
  {
    name: 'otel.go',
    isDir: false,
    level: 3,
    path: 'internal/app/config/otel.go',
    desc: 'OpenTelemetry tracing — OTLP HTTP export to Jaeger.'
  },
  {
    name: 'redis.go',
    isDir: false,
    level: 3,
    path: 'internal/app/config/redis.go',
    desc: 'Redis connection pool via Redigo.'
  },
  {
    name: 'session.go',
    isDir: false,
    level: 3,
    path: 'internal/app/config/session.go',
    desc: 'SCS session manager backed by Redis store.'
  },
  {
    name: 'database/',
    isDir: true,
    level: 2,
    path: 'internal/app/database/',
    desc: 'Generic Repository[T] and Atlas model registry.'
  },
  {
    name: 'repository.go',
    isDir: false,
    level: 3,
    path: 'internal/app/database/repository.go',
    desc: 'Eloquent-like generic CRUD wrapper over GORM.',
    code: `type Repository[T any] struct {
	db *gorm.DB
}

func New[T any](db *gorm.DB) *Repository[T] {
	return &Repository[T]{db: db}
}`
  },
  {
    name: 'registry.go',
    isDir: false,
    level: 3,
    path: 'internal/app/database/registry.go',
    desc: 'Central registry where module models register via init() for Atlas.',
    code: `func Register(model schema.Tabler) {
	models = append(models, model)
}

func All() []schema.Tabler {
	return models
}`
  },
  {
    name: 'helpers/',
    isDir: true,
    level: 2,
    path: 'internal/app/helpers/',
    desc: 'Shared utility packages used across modules.'
  },
  {
    name: 'jsonutils/',
    isDir: true,
    level: 3,
    path: 'internal/app/helpers/jsonutils/',
    desc: 'JSON encoding/decoding utilities.'
  },
  {
    name: 'validator/',
    isDir: true,
    level: 3,
    path: 'internal/app/helpers/validator/',
    desc: 'Request validation helpers for DTO structs.'
  },
  {
    name: 'middleware/',
    isDir: true,
    level: 2,
    path: 'internal/app/middleware/',
    desc: 'HTTP middlewares applied globally in routes.SetupRoutes.'
  },
  {
    name: 'cors.go',
    isDir: false,
    level: 3,
    path: 'internal/app/middleware/cors.go',
    desc: 'CORS configuration middleware.'
  },
  {
    name: 'observability.go',
    isDir: false,
    level: 3,
    path: 'internal/app/middleware/observability.go',
    desc: 'Route tagging middleware for Prometheus metrics.'
  },
  {
    name: 'session.go',
    isDir: false,
    level: 3,
    path: 'internal/app/middleware/session.go',
    desc: 'SCS session load/save middleware.'
  },
  {
    name: 'types/',
    isDir: true,
    level: 2,
    path: 'internal/app/types/',
    desc: 'Shared HTTP response types (errors, messages, health checks).'
  },
  {
    name: 'common-dto.go',
    isDir: false,
    level: 3,
    path: 'internal/app/types/common-dto.go',
    desc: 'ErrorResponse, MessageResponse, and other shared DTOs.'
  },
  {
    name: 'health-dto.go',
    isDir: false,
    level: 3,
    path: 'internal/app/types/health-dto.go',
    desc: 'HealthCheckResponse and HealthCheckDetailedResponse types.'
  },
  {
    name: 'router/',
    isDir: true,
    level: 2,
    path: 'internal/app/router/',
    desc: 'Declarative OpenAPI documentation wrappers for Fuego routes.',
    tip: 'Each endpoint declares a <code>router.Doc</code> in the module\'s <code>docs.go</code>.'
  },
  {
    name: 'doc.go',
    isDir: false,
    level: 3,
    path: 'internal/app/router/doc.go',
    desc: 'Doc, Detail, Response, and QueryParam types for OpenAPI metadata.'
  },
  {
    name: 'options.go',
    isDir: false,
    level: 3,
    path: 'internal/app/router/options.go',
    desc: 'Converts router.Doc into Fuego route options.'
  },
  {
    name: 'register.go',
    isDir: false,
    level: 3,
    path: 'internal/app/router/register.go',
    desc: 'router.Get/Post/Put/Delete wrappers that attach OpenAPI docs to routes.'
  },
  {
    name: 'app.go',
    isDir: false,
    level: 2,
    path: 'internal/app/app.go',
    desc: 'Package documentation — groups shared infrastructure. Domain code lives in internal/modules/.',
    code: `// Package app groups shared infrastructure: config, database, helpers,
// middleware, router, and types. Domain code lives in internal/modules.
package app`
  },

  // internal/modules/
  {
    name: 'modules/',
    isDir: true,
    level: 1,
    path: 'internal/modules/',
    desc: 'Self-contained domain packages. Each module owns model, dto, service, controller, and docs.',
    tip: 'Run <code>grove make:resource Post</code> to scaffold a new module and register it automatically.'
  },
  {
    name: 'auth/',
    isDir: true,
    level: 2,
    path: 'internal/modules/auth/',
    desc: 'Authentication domain — login, logout, session, register (delegates user creation to users module).'
  },
  {
    name: 'controller.go',
    isDir: false,
    level: 3,
    path: 'internal/modules/auth/controller.go',
    desc: 'Auth HTTP handlers with Wire() and Mount() for /auth routes.'
  },
  {
    name: 'service.go',
    isDir: false,
    level: 3,
    path: 'internal/modules/auth/service.go',
    desc: 'Login validation and password checking logic.'
  },
  {
    name: 'dto.go',
    isDir: false,
    level: 3,
    path: 'internal/modules/auth/dto.go',
    desc: 'LoginRequest, LoginResponse, and auth-specific DTOs.'
  },
  {
    name: 'docs.go',
    isDir: false,
    level: 3,
    path: 'internal/modules/auth/docs.go',
    desc: 'OpenAPI router.Doc definitions for auth endpoints.'
  },
  {
    name: 'users/',
    isDir: true,
    level: 2,
    path: 'internal/modules/users/',
    desc: 'Users domain — model+repo, dto, service, controller, docs.'
  },
  {
    name: 'model.go',
    isDir: false,
    level: 3,
    path: 'internal/modules/users/model.go',
    desc: 'User GORM model with init() Atlas registration and typed Repo.',
    code: `func init() {
	database.Register(&User{})
}

type Repo struct {
	*database.Repository[User]
}

func Users(db *gorm.DB) *Repo {
	return &Repo{Repository: database.New[User](db)}
}`
  },
  {
    name: 'dto.go',
    isDir: false,
    level: 3,
    path: 'internal/modules/users/dto.go',
    desc: 'CreateUserRequest, UpdateUserRequest, UserResponse, UsersListResponse.'
  },
  {
    name: 'service.go',
    isDir: false,
    level: 3,
    path: 'internal/modules/users/service.go',
    desc: 'User business logic — CRUD, bcrypt hashing, pagination.'
  },
  {
    name: 'controller.go',
    isDir: false,
    level: 3,
    path: 'internal/modules/users/controller.go',
    desc: 'User HTTP handlers with Mount() registering /users routes via router.* helpers.'
  },
  {
    name: 'docs.go',
    isDir: false,
    level: 3,
    path: 'internal/modules/users/docs.go',
    desc: 'OpenAPI router.Doc definitions for user CRUD endpoints.'
  },
  {
    name: 'module.go',
    isDir: false,
    level: 2,
    path: 'internal/modules/module.go',
    desc: 'Module interface, Boot struct, and Factory type.',
    code: `type Module interface {
	Mount(api *fuego.Server, session *scs.SessionManager)
}

type Boot struct {
	DB      *gorm.DB
	Session *scs.SessionManager
}`
  },
  {
    name: 'register.go',
    isDir: false,
    level: 2,
    path: 'internal/modules/register.go',
    desc: 'Module registry — grove make:resource adds new Wire() factories here.',
    code: `var registry = []Factory{
	func(b Boot) Module { return users.Wire(b.DB) },
	func(b Boot) Module { return auth.Wire(b.DB, b.Session) },
}

func Mount(api *fuego.Server, boot Boot) {
	for _, factory := range registry {
		factory(boot).Mount(api, boot.Session)
	}
}`
  },

  // internal/routes/
  {
    name: 'routes/',
    isDir: true,
    level: 1,
    path: 'internal/routes/',
    desc: 'Global route setup — middleware, health checks, metrics, and module mounting.'
  },
  {
    name: 'health.go',
    isDir: false,
    level: 2,
    path: 'internal/routes/health.go',
    desc: 'GET / and GET /health handlers with DB and Redis status checks.'
  },
  {
    name: 'routes.go',
    isDir: false,
    level: 2,
    path: 'internal/routes/routes.go',
    desc: 'Applies global middleware and mounts all modules via modules.Mount().',
    code: `func SetupRoutes(
	s *fuego.Server,
	db *gorm.DB,
	redisPool *redis.Pool,
	session *scs.SessionManager,
	metricsHandler http.Handler,
) {
	fuego.Use(s, otelhttp.NewMiddleware("go-project-base"))
	fuego.Use(s, middleware.CORSMiddleware(middleware.DefaultCORSConfig()))
	fuego.Use(s, middleware.SessionMiddleware(session))

	fuego.Get(s, "/", healthCheck)
	fuego.Get(s, "/health", healthCheckDetailed(db, redisPool))
	fuego.GetStd(s, "/metrics", metricsHandler.ServeHTTP)

	api := fuego.Group(s, "/api/v1")
	modules.Mount(api, modules.Boot{DB: db, Session: session})
}`
  },

  // ── infra/ ────────────────────────────────────────────────────────────
  {
    name: 'infra/',
    isDir: true,
    level: 0,
    path: 'infra/',
    desc: 'Observability stack configuration — Prometheus, Grafana, Loki, Jaeger, Promtail.'
  },
  {
    name: 'grafana/',
    isDir: true,
    level: 1,
    path: 'infra/grafana/',
    desc: 'Pre-provisioned Grafana dashboards and datasource configs.'
  },
  {
    name: 'dashboards/',
    isDir: true,
    level: 2,
    path: 'infra/grafana/dashboards/',
    desc: 'Pre-built Go Project Base dashboard JSON.'
  },
  {
    name: 'provisioning/',
    isDir: true,
    level: 2,
    path: 'infra/grafana/provisioning/',
    desc: 'Auto-provisioned datasources (Prometheus, Loki, Jaeger) and dashboards on startup.'
  },
  {
    name: 'jeaguer-ui/',
    isDir: true,
    level: 1,
    path: 'infra/jeaguer-ui/',
    desc: 'Jaeger UI configuration (dark mode).'
  },
  {
    name: 'jaeger-ui.json',
    isDir: false,
    level: 2,
    path: 'infra/jeaguer-ui/jaeger-ui.json',
    desc: 'Jaeger UI theme and layout config mounted into the Jaeger container.'
  },
  {
    name: 'loki-config.yml',
    isDir: false,
    level: 1,
    path: 'infra/loki-config.yml',
    desc: 'Loki log aggregation storage configuration.'
  },
  {
    name: 'prometheus.yml',
    isDir: false,
    level: 1,
    path: 'infra/prometheus.yml',
    desc: 'Prometheus scrape config — scrapes GET /metrics from the Go app.'
  },
  {
    name: 'promtail-config.yml',
    isDir: false,
    level: 1,
    path: 'infra/promtail-config.yml',
    desc: 'Promtail config — tails logs/app.log and ships to Loki.'
  },

  // ── logs/ ─────────────────────────────────────────────────────────────
  {
    name: 'logs/',
    isDir: true,
    level: 0,
    path: 'logs/',
    desc: 'Application log files written by slog (tailed by Promtail → Loki → Grafana).'
  },
  {
    name: '.gitkeep',
    isDir: false,
    level: 1,
    path: 'logs/.gitkeep',
    desc: 'Keeps the logs directory in git. app.log is created at runtime.'
  },

  // ── migrations/ ───────────────────────────────────────────────────────
  {
    name: 'migrations/',
    isDir: true,
    level: 0,
    path: 'migrations/',
    desc: 'Atlas SQL migration files and integrity checksum.'
  },
  {
    name: '20260127143000_initial.sql',
    isDir: false,
    level: 1,
    path: 'migrations/20260127143000_initial.sql',
    desc: 'Initial migration creating baseline tables (users, etc.).'
  },
  {
    name: 'atlas.sum',
    isDir: false,
    level: 1,
    path: 'migrations/atlas.sum',
    desc: 'Atlas migration directory checksum — run grove migrate:hash after editing migrations.'
  },

  // ── doc/ ──────────────────────────────────────────────────────────────
  {
    name: 'doc/',
    isDir: true,
    level: 0,
    path: 'doc/',
    desc: 'Generated OpenAPI specification (produced by Fuego at runtime/build).'
  },
  {
    name: 'openapi.json',
    isDir: false,
    level: 1,
    path: 'doc/openapi.json',
    desc: 'Auto-generated OpenAPI 3 spec consumed by Scalar at /swagger.'
  },

  // ── Root files ────────────────────────────────────────────────────────
  {
    name: '.air.toml',
    isDir: false,
    level: 0,
    path: '.air.toml',
    desc: 'Air hot-reload configuration (used by make dev in go-project-base).'
  },
  {
    name: 'atlas.hcl',
    isDir: false,
    level: 0,
    path: 'atlas.hcl',
    desc: 'Atlas migration environments — program mode loads schema via go run ./cmd/atlas.',
    code: `data "external_schema" "gorm" {
  program = ["go", "run", "-mod=mod", "./cmd/atlas"]
}

env "local" {
  src = data.external_schema.gorm.url
  migration { dir = "file://migrations" }
}`
  },
  {
    name: 'docker-compose.yml',
    isDir: false,
    level: 0,
    path: 'docker-compose.yml',
    desc: 'Full infrastructure stack — PostgreSQL, Redis, Jaeger, Prometheus, Loki, Promtail, Grafana.'
  },
  {
    name: 'Makefile',
    isDir: false,
    level: 0,
    path: 'Makefile',
    desc: 'Dev commands: install, run, dev, test, migrate-create, migrate-up, db-reset.'
  },
  {
    name: 'go.mod',
    isDir: false,
    level: 0,
    path: 'go.mod',
    desc: 'Go module definition — Fuego, GORM, Atlas provider, OTel, SCS, Scalar.'
  },
  {
    name: 'm.go',
    isDir: false,
    level: 0,
    path: 'm.go',
    desc: 'Root guard — reminds developers to use make dev or Air instead of go run . at repo root.'
  },
  {
    name: '.env.example',
    isDir: false,
    level: 0,
    path: '.env.example',
    desc: 'Environment variable template — DB, Redis, OTel, logging, and app metadata.',
    code: `DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=mcs_dctfweb_sender
OTEL_SERVICE_NAME=go-project-base
LOG_LEVEL=info`
  },
  {
    name: 'grove.toml',
    isDir: false,
    level: 0,
    path: 'grove.toml',
    desc: 'Optional Grove CLI hot-reload config — added when using grove dev (not in go-project-base by default).',
    tip: 'Created by <code>grove setup</code> or manually. Configures <code>grove dev</code> watch paths and excludes.'
  },
  {
    name: 'README.md',
    isDir: false,
    level: 0,
    path: 'README.md',
    desc: 'Project documentation — tech stack, observability architecture, and getting started guide.'
  }
]

/** Plain-text tree matching go-project-base README for docs tables and copy blocks. */
export const projectStructureAscii = `my-api/
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
├── infra/
│   ├── grafana/                   # Dashboards + provisioning
│   ├── jeaguer-ui/                # Jaeger UI config
│   ├── loki-config.yml
│   ├── prometheus.yml
│   └── promtail-config.yml
├── logs/                          # App log files (Promtail → Loki)
├── migrations/                    # Atlas SQL migrations
├── doc/
│   └── openapi.json               # Generated OpenAPI spec
├── .air.toml                      # Air hot reload config
├── atlas.hcl                      # Atlas migration config
├── docker-compose.yml             # Full infrastructure stack
├── Makefile                       # Dev commands
└── go.mod                         # Go module definition`
