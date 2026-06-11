import { projectStructureAscii } from './projectStructure.js'

export const sections = [
  // ─────────────────────────────────────────────
  // GETTING STARTED
  // ─────────────────────────────────────────────
  {
    id: 'getting-started',
    title: 'Getting Started',
    items: [
      {
        id: 'installation',
        title: 'Installation',
        blocks: [
          {
            type: 'paragraph',
            text: 'Grove is distributed as a single binary via <code>go install</code>. Make sure you have Go 1.25 or newer and that <code>$GOPATH/bin</code> (or <code>$HOME/go/bin</code>) is in your <code>$PATH</code>.'
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'terminal',
            code: `go install github.com/caiolandgraf/grove@latest`
          },
          {
            type: 'paragraph',
            text: 'Verify the installation by running:'
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'terminal',
            code: `grove --version   # print version
grove --help      # full command reference`
          },
          {
            type: 'note',
            kind: 'info',
            text: 'Documentation version: <strong>v2.2.0</strong>.'
          },
          {
            type: 'note',
            kind: 'info',
            text: 'Grove requires the <strong>Atlas CLI</strong> for all migration-related commands. Install it from <a href="https://atlasgo.io/docs" target="_blank">atlasgo.io</a>.'
          }
        ]
      },
      {
        id: 'quick-start',
        title: 'Quick Start',
        blocks: [
          {
            type: 'paragraph',
            text: 'Scaffold a new project with <code>grove setup</code>, then configure your environment and start the server.'
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'terminal',
            code: `# 1. Scaffold a new project
grove setup my-api
# or: grove setup (prompts for name + observability)

# 2. Enter the project directory
cd my-api

# 3. Configure your environment
cp .env.example .env
# edit .env with your database credentials

# 4. Start infra + dev server (docker compose + hot reload)
grove up

# — or start only the dev server —
grove dev

# — or, if you prefer Air for hot-reload —
grove dev:air`
          },
          {
            type: 'paragraph',
            text: 'Your API is now running at <code>http://localhost:8080</code>. The OpenAPI docs are available at <code>http://localhost:8080/swagger</code> automatically via fuego.'
          },
          {
            type: 'steps',
            items: [
              {
                title: 'Scaffold a resource',
                text: 'Run <code>grove make:resource Post</code> to generate a model, controller and DTO in one shot.'
              },
              {
                title: 'Add your fields',
                text: 'Edit <code>internal/modules/posts/model.go</code> to add your fields (e.g. <code>Title</code>, <code>Body</code>). Also fill in <code>internal/modules/posts/dto.go</code> with your request/response fields.'
              },
              {
                title: 'Generate the migration',
                text: 'Run <code>grove make:migration create_posts_table</code> — Atlas diffs your updated model against the current DB schema and generates the exact SQL.'
              },
              {
                title: 'Apply it',
                text: 'Run <code>grove migrate</code> to apply pending migrations to your database.'
              },
              {
                title: 'Register routes',
                text: "Verify your module is registered in <code>internal/modules/register.go</code> — <code>grove make:resource</code> does this automatically."
              }
            ]
          }
        ]
      },
      {
        id: 'first-project-guide',
        title: 'First Project Guide',
        blocks: [
          {
            type: 'paragraph',
            text: 'Ready to build your first Go API? Follow this step-by-step interactive playground guide to scaffold, develop, and launch a complete micro-blogging API using Grove.'
          },
          {
            type: 'first-project-guide'
          }
        ]
      },
      {
        id: 'requirements',
        title: 'Requirements',
        blocks: [
          {
            type: 'table',
            head: ['Tool', 'Version', 'Purpose'],
            rows: [
              [
                '<a href="https://go.dev/dl" target="_blank">Go</a>',
                '1.25+',
                'Build and install grove and your application'
              ],
              [
                '<a href="https://atlasgo.io/docs" target="_blank">Atlas CLI</a>',
                'latest',
                'Migration generation and application'
              ],
              [
                '<a href="https://github.com/caiolandgraf/gest" target="_blank">gest library</a>',
                'v2+',
                "Jest-inspired testing library — added to your project's <code>go.mod</code> automatically by <code>grove make:test</code>"
              ],
              [
                '<a href="https://github.com/caiolandgraf/gest" target="_blank">gest CLI</a> (optional)',
                'v2+',
                'Renders beautiful Jest-style output. Install with <code>go install github.com/caiolandgraf/gest/v2/cmd/gest@latest</code>. <code>grove test</code> falls back to <code>go test -v</code> when absent.'
              ],
              [
                '<a href="https://github.com/air-verse/air" target="_blank">air</a> (optional)',
                'latest',
                'Hot-reload via <code>grove dev:air</code> — not needed for <code>grove dev</code>'
              ],
              [
                'PostgreSQL (or any GORM driver)',
                '14+',
                'Default database — swap driver freely'
              ]
            ]
          },
          {
            type: 'note',
            kind: 'tip',
            text: '<code>grove dev</code> provides built-in hot reload without any external tools. <code>air</code> is only needed if you use <code>grove dev:air</code>.'
          }
        ]
      }
    ]
  },

  // ─────────────────────────────────────────────
  // COMMANDS
  // ─────────────────────────────────────────────
  {
    id: 'commands',
    title: 'Commands',
    items: [
      {
        id: 'cmd-setup',
        title: 'grove setup',
        blocks: [
          {
            type: 'paragraph',
            text: 'Downloads and scaffolds a complete Grove project from the official <a href="https://github.com/caiolandgraf/go-project-base" target="_blank">go-project-base</a> template repository on GitHub (v2.0.0 modular MSC architecture).'
          },
          {
            type: 'paragraph',
            text: 'If you omit the project name, Grove will prompt for it and ask which observability services to enable (Jaeger, Prometheus, Grafana, Loki, Promtail). The selections update <code>.env.example</code> (<code>OTEL_ENABLED</code> / <code>METRICS_ENABLED</code>) and <code>docker-compose.yml</code> at the project root (or <code>infra/compose.yml</code> when present).'
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'terminal',
            code: `grove setup [project-name] [--module <go-module-path>]`
          },
          {
            type: 'table',
            head: ['Flag', 'Default', 'Description'],
            rows: [
              [
                '<code>--module</code>',
                'project name',
                'Go module path (e.g. <code>github.com/acme/my-api</code>)'
              ]
            ]
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'examples',
            code: `grove setup my-api
grove setup my-api --module github.com/acme/my-api
grove setup  # prompt for project name + observability`
          }
        ]
      },
      {
        id: 'cmd-up',
        title: 'grove up',
        blocks: [
          {
            type: 'paragraph',
            text: 'Starts infrastructure via <code>docker compose</code> and then launches <code>grove dev</code>. It detects <code>docker-compose.yml</code> at the project root (or <code>infra/compose.yml</code> as fallback) and passes <code>.env</code> to Compose when present.'
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'terminal',
            code: `grove up`
          },
          {
            type: 'note',
            kind: 'tip',
            text: 'If you only want the dev server (no infra), use <code>grove dev</code>.'
          }
        ]
      },
      {
        id: 'cmd-down',
        title: 'grove down',
        blocks: [
          {
            type: 'paragraph',
            text: 'Stops the Docker Compose stack started by <code>grove up</code>. It uses the same compose file detection (<code>docker-compose.yml</code> at root, or <code>infra/compose.yml</code> as fallback) and passes <code>.env</code> when present.'
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'terminal',
            code: `grove down`
          }
        ]
      },
      {
        id: 'cmd-dev',
        title: 'grove dev',
        blocks: [
          {
            type: 'paragraph',
            text: 'Compiles and runs your application, then watches for file changes and automatically recompiles and restarts the binary on every save. No external tools required — hot reload is built directly into Grove.'
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'terminal',
            code: `grove dev`
          },
          {
            type: 'paragraph',
            text: "<code>grove dev</code> also processes your application's stdout/stderr and formats it intelligently — structured JSON logs are rendered as human-readable coloured lines, and panics are captured into a styled block."
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'JSON log output (slog / zap / zerolog)',
            code: `  08:38:28  INF  Booting application...
  08:38:28  INF  OpenTelemetry initialized  service=grove-app  endpoint=localhost:4318
  08:38:28  ERR  Failed to boot application  error=failed to connect to database: ...`
          },
          {
            type: 'table',
            head: ['Level', 'Badge', 'Colour'],
            rows: [
              ['<code>DEBUG</code>', '<code>DBG</code>', 'Gray'],
              ['<code>INFO</code>', '<code>INF</code>', 'Green'],
              ['<code>WARN</code>', '<code>WRN</code>', 'Yellow'],
              ['<code>ERROR</code>', '<code>ERR</code>', 'Red']
            ]
          },
          {
            type: 'note',
            kind: 'info',
            text: "Compatible with any structured logger that emits <code>level</code>, <code>msg</code> and <code>time</code> fields — including Go's standard <code>slog</code>, <code>zap</code> and <code>zerolog</code>. The timestamp is trimmed to <code>HH:MM:SS</code> and extra fields are shown inline, dimmed."
          },
          {
            type: 'paragraph',
            text: '<strong>Startup hints</strong> — Grove detects common startup errors and prints an actionable <code>HINT</code> block immediately below the error line:'
          },
          {
            type: 'table',
            head: ['Error detected', 'Hint shown'],
            rows: [
              [
                '<code>.env not found</code>',
                '<code>cp .env.example .env</code>'
              ],
              [
                '<code>connection refused</code> · <code>dial error</code> · <code>failed to connect</code>',
                '<code>docker compose up -d</code> and check <code>.env</code>'
              ]
            ]
          },
          {
            type: 'note',
            kind: 'tip',
            text: 'Each hint is shown once per rebuild. If the same error persists after the next file save, the hint appears again — no spam within the same run.'
          },
          {
            type: 'paragraph',
            text: 'Behaviour is fully configurable via the optional <code>[dev]</code> section in <code>grove.toml</code>. All fields are optional — when the file is absent or the section is omitted, sensible defaults are used and <code>grove dev</code> works out of the box.'
          },
          {
            type: 'code',
            lang: 'toml',
            label: 'grove.toml',
            code: `[dev]
root        = "."
bin         = ".grove/tmp/app"
build_cmd   = "go build -o .grove/tmp/app ./cmd/api/"
watch_dirs  = ["cmd", "internal"]
exclude     = [".grove", "vendor", "node_modules", ".git", "infra", "migrations", "bin", "internal/tests"]
extensions  = [".go"]
debounce_ms = 50`
          },
          {
            type: 'table',
            head: ['Field', 'Default', 'Description'],
            rows: [
              [
                '<code>root</code>',
                '<code>.</code>',
                'Working directory for build commands'
              ],
              [
                '<code>bin</code>',
                '<code>.grove/tmp/app</code>',
                'Path to the compiled binary'
              ],
              [
                '<code>build_cmd</code>',
                '<code>go build -o .grove/tmp/app ./cmd/api/</code>',
                'Command used to compile the project'
              ],
              [
                '<code>watch_dirs</code>',
                '<code>["cmd", "internal"]</code>',
                'Directories to watch for file changes (recursive)'
              ],
              [
                '<code>exclude</code>',
                '<code>[".grove", "vendor", "node_modules", ".git", "infra", "migrations", "bin", "internal/tests"]</code>',
                'Directories or path fragments to ignore (the <code>internal/tests/</code> directory is always excluded so test saves never trigger a rebuild)'
              ],
              [
                '<code>extensions</code>',
                '<code>[".go"]</code>',
                'File extensions that trigger a rebuild'
              ],
              [
                '<code>debounce_ms</code>',
                '<code>50</code>',
                'Milliseconds to wait after the last change before rebuilding'
              ],
              [
                '<code>port_guard</code>',
                '<code>auto</code>',
                'If set, Grove will attempt to free the port before starting the app. When omitted, it tries to infer the port from <code>.env</code> (<code>APP_PORT</code>, <code>PORT</code>, <code>BASE_URL</code>, etc.). Set to <code>0</code> to disable.'
              ]
            ]
          },
          {
            type: 'note',
            kind: 'tip',
            text: 'Newly created subdirectories are picked up automatically at runtime — no restart of <code>grove dev</code> required. The <code>internal/tests/</code> directory is always excluded so a test save never triggers an application rebuild.'
          }
        ]
      },
      {
        id: 'cmd-dev-air',
        title: 'grove dev:air',
        blocks: [
          {
            type: 'paragraph',
            text: 'Starts the development HTTP server using <a href="https://github.com/air-verse/air" target="_blank">Air</a> for hot-reload. If <code>air</code> is not installed it falls back to <code>go run ./cmd/api/main.go</code>.'
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'terminal',
            code: `grove dev:air`
          },
          {
            type: 'note',
            kind: 'info',
            text: 'For a zero-dependency hot-reload experience use <code>grove dev</code> instead — it has a built-in watcher that requires no external tools whatsoever.'
          },
          {
            type: 'note',
            kind: 'tip',
            text: 'Install Air with <code>go install github.com/air-verse/air@latest</code>. SIGINT (<kbd>Ctrl+C</kbd>) is forwarded to the child process so it can shut down gracefully.'
          }
        ]
      },
      {
        id: 'cmd-build',
        title: 'grove build',
        blocks: [
          {
            type: 'paragraph',
            text: 'Compiles the application and writes the binary to <code>./bin/app</code> by default. The <code>bin/</code> directory is created automatically if it does not exist.'
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'terminal',
            code: `grove build [--output <path>]`
          },
          {
            type: 'table',
            head: ['Flag', 'Default', 'Description'],
            rows: [
              [
                '<code>-o</code>, <code>--output</code>',
                '<code>./bin/app</code>',
                'Output path for the compiled binary'
              ]
            ]
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'examples',
            code: `grove build
grove build -o ./bin/my-api`
          }
        ]
      },
      {
        id: 'cmd-start',
        title: 'grove start',
        blocks: [
          {
            type: 'paragraph',
            text: 'Builds the application and runs the resulting binary. Useful for a quick local smoke test without hot reload.'
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'terminal',
            code: `grove start [--output <path>]`
          },
          {
            type: 'table',
            head: ['Flag', 'Default', 'Description'],
            rows: [
              [
                '<code>-o</code>, <code>--output</code>',
                '<code>./bin/app</code>',
                'Output path for the compiled binary'
              ]
            ]
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'examples',
            code: `grove start
grove start -o ./bin/my-api`
          }
        ]
      },
      {
        id: 'cmd-make-model',
        title: 'grove make:model',
        blocks: [
          {
            type: 'paragraph',
            text: 'Scaffolds a new GORM model in <code>internal/modules/&lt;domain&gt;/model.go</code> (e.g. <code>internal/modules/posts/model.go</code>). The entity name is automatically singularized and the domain package pluralized. Combine flags to scaffold additional layers in the same step — or use <code>-r</code> as a shorthand for the full resource.'
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'terminal',
            code: `grove make:model <Name> [-c] [-d] [-r]`
          },
          {
            type: 'table',
            head: ['Flag', 'Description'],
            rows: [
              [
                '<code>-c</code>, <code>--controller</code>',
                'Also scaffold a fuego controller'
              ],
              [
                '<code>-d</code>, <code>--dto</code>',
                'Also scaffold a DTO request/response file'
              ],
              [
                '<code>-r</code>, <code>--resource</code>',
                'Full resource — shorthand for <code>-c -d</code> combined'
              ]
            ]
          },
          {
            type: 'note',
            kind: 'warning',
            text: 'Migrations are <strong>not</strong> generated automatically when scaffolding a model. Add your fields to the model first, then run <code>grove make:migration &lt;name&gt;</code> to let Atlas diff your schema and generate the correct SQL. This ensures the migration reflects the fields you actually defined — not an empty struct.'
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'examples',
            code: `grove make:model Post            # model only
grove make:model Post -c         # model + controller
grove make:model Post -d         # model + DTO
grove make:model Post -cd        # model + controller + DTO
grove make:model Post -r         # full resource (same as -cd)
grove make:model order_item      # → OrderItem`
          },
          {
            type: 'paragraph',
            text: 'The generated model includes UUID primary key, <code>CreatedAt</code>, <code>UpdatedAt</code>, soft-delete (<code>DeletedAt</code>), Atlas auto-registration via <code>database.Register</code>, and a typed <code>Repo</code> accessor:'
          },
          {
            type: 'code',
            lang: 'go',
            label: 'internal/modules/posts/model.go',
            code: `package posts

import (
	"time"

	"your/module/internal/app/database"
	"gorm.io/gorm"
)

func init() {
	database.Register(&Post{})
}

type Post struct {
	ID        string         \`gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"\`
	CreatedAt time.Time      \`gorm:"autoCreateTime" json:"created_at"\`
	UpdatedAt time.Time      \`gorm:"autoUpdateTime" json:"updated_at"\`
	DeletedAt gorm.DeletedAt \`gorm:"index"          json:"-"\`
	// TODO: add fields
}

func (Post) TableName() string { return "posts" }

type Repo struct {
	*database.Repository[Post]
}

func Posts(db *gorm.DB) *Repo {
	return &Repo{
		Repository: database.New[Post](db),
	}
}`
          }
        ]
      },
      {
        id: 'cmd-make-controller',
        title: 'grove make:controller',
        blocks: [
          {
            type: 'paragraph',
            text: 'Scaffolds a new fuego controller with <code>Wire</code> and <code>Mount</code> helpers in <code>internal/modules/&lt;domain&gt;/controller.go</code>. By default, the generated controller uses a struct injected with the corresponding service layer. Pass <code>--no-auth</code> to generate the legacy function-based stub instead. If the file already exists the command prints <strong>SKIPPED</strong> and exits cleanly.'
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'terminal',
            code: `grove make:controller <Name> [--no-auth]`
          },
          {
            type: 'code',
            lang: 'go',
            label: 'internal/modules/posts/controller.go (generated)',
            code: `package posts

import (
	"github.com/alexedwards/scs/v2"
	"your/module/internal/app/router"
	"github.com/go-fuego/fuego"
	"gorm.io/gorm"
)

type Controller struct {
	service Service
}

func NewController(service Service) *Controller {
	return &Controller{service: service}
}

func Wire(db *gorm.DB) *Controller {
	return NewController(WireService(db))
}

func (ctrl *Controller) Mount(api *fuego.Server, session *scs.SessionManager) {
	group := fuego.Group(api, "/posts")

	router.Get(group, "/", ctrl.ListPosts, ListPostsDoc, session)
	router.Post(group, "/", ctrl.CreatePost, CreatePostDoc, session)
	router.Get(group, "/{post_id}", ctrl.GetPost, GetPostDoc, session)
	router.Put(group, "/{post_id}", ctrl.UpdatePost, UpdatePostDoc, session)
	router.Delete(group, "/{post_id}", ctrl.DeletePost, DeletePostDoc, session)
}

// TODO: implement HTTP handlers delegating to ctrl.service`
          }
        ]
      },
      {
        id: 'cmd-make-service',
        title: 'grove make:service',
        blocks: [
          {
            type: 'paragraph',
            text: 'Scaffolds a new Go service struct in <code>internal/modules/&lt;domain&gt;/service.go</code>. The service receives the module repository directly and exposes a <code>WireService</code> helper for dependency injection. If the file already exists the command prints <strong>SKIPPED</strong> and exits.'
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'terminal',
            code: `grove make:service <Name>`
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'examples',
            code: `grove make:service Post
grove make:service Posts        # same as Post (singularized)
grove make:service BlogPost
grove make:service user_profile`
          },
          {
            type: 'code',
            lang: 'go',
            label: 'internal/modules/posts/service.go (generated)',
            code: `package posts

import "gorm.io/gorm"

type Service interface {
	// TODO: add business methods (CreatePost, GetPostByID, etc.)
}

type service struct {
	repo *Repo
}

func NewService(repo *Repo) Service {
	return &service{repo: repo}
}

func WireService(db *gorm.DB) Service {
	return NewService(Posts(db))
}`
          }
        ]
      },
      {
        id: 'cmd-make-dto',
        title: 'grove make:dto',
        blocks: [
          {
            type: 'paragraph',
            text: 'Scaffolds DTO request/response structs in <code>internal/modules/&lt;domain&gt;/dto.go</code>. If the file already exists the command prints <strong>SKIPPED</strong> and exits cleanly.'
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'terminal',
            code: `grove make:dto <Name>`
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'examples',
            code: `grove make:dto Post
grove make:dto BlogPost
grove make:dto order_item   # → OrderItem`
          },
          {
            type: 'code',
            lang: 'go',
            label: 'internal/modules/posts/dto.go (generated)',
            code: `package posts

type CreatePostRequest struct {
	// TODO: add fields
}

type UpdatePostRequest struct {
	// TODO: add fields
}

type PostResponse struct {
	ID string \`json:"id"\`
}

type PostsListResponse struct {
	Items []PostResponse \`json:"items"\`
	Total int            \`json:"total"\`
}`
          }
        ]
      },
      {
        id: 'cmd-make-middleware',
        title: 'grove make:middleware',
        blocks: [
          {
            type: 'paragraph',
            text: 'Scaffolds a new HTTP middleware in <code>internal/app/middleware/</code>. The generated file follows the standard <code>func(http.Handler) http.Handler</code> signature compatible with fuego and the standard library. If the file already exists the command prints <strong>SKIPPED</strong> and exits cleanly.'
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'terminal',
            code: `grove make:middleware <Name>`
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'examples',
            code: `grove make:middleware Auth
grove make:middleware RateLimit
grove make:middleware cors_headers   # → CorsHeaders`
          },
          {
            type: 'code',
            lang: 'go',
            label: 'internal/app/middleware/auth-middleware.go (generated)',
            code: `package middleware

import "net/http"

// Auth is an HTTP middleware that handles Auth logic.
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: implement Auth middleware logic

		next.ServeHTTP(w, r)
	})
}`
          },
          {
            type: 'paragraph',
            text: 'Register the middleware in your server setup by passing it to <code>s.Use()</code>:'
          },
          {
            type: 'code',
            lang: 'go',
            label: 'internal/routes/routes.go',
            code: `func SetupRoutes(s *fuego.Server, session *scs.SessionManager) {
	// Apply globally
	fuego.Use(s, middleware.Auth)

	// ... health checks, metrics ...

	api := fuego.Group(s, "/api/v1")
	modules.Mount(api, modules.Boot{DB: db, Session: session})
}`
          }
        ]
      },
      {
        id: 'cmd-make-migration',
        title: 'grove make:migration',
        blocks: [
          {
            type: 'paragraph',
            text: 'Generates a new SQL migration file by diffing your GORM models against the current database schema using Atlas. <strong>Always edit your model first</strong>, then run this command — Atlas will produce the exact SQL diff between your updated struct and the current DB schema.'
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'terminal',
            code: `grove make:migration <name> [--env <atlas-env>]`
          },
          {
            type: 'table',
            head: ['Flag', 'Default', 'Description'],
            rows: [
              [
                '<code>--env</code>',
                '<code>local</code>',
                'Atlas environment defined in <code>atlas.hcl</code>'
              ]
            ]
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'examples',
            code: `grove make:migration create_posts_table
grove make:migration add_title_to_posts
grove make:migration create_orders_table --env dev`
          },
          {
            type: 'note',
            kind: 'tip',
            text: 'Need to update an existing table? Add the new fields to your struct, then run <code>grove make:migration add_&lt;field&gt;_to_&lt;table&gt;</code> — Atlas will generate an <code>ALTER TABLE</code> migration with exactly the diff between the current DB schema and your updated model.'
          },
          {
            type: 'note',
            kind: 'warning',
            text: 'Always review the generated SQL in <code>migrations/</code> before applying. Atlas is thorough but your business logic may require manual adjustments.'
          }
        ]
      },
      {
        id: 'cmd-make-relations',
        title: 'grove make:relations',
        blocks: [
          {
            type: 'paragraph',
            text: 'Infers model relationships from foreign key fields (for example, <code>UserID</code>) and adds GORM relation fields automatically.'
          },
          {
            type: 'note',
            kind: 'tip',
            text: 'Foreign key suffixes <code>ID</code>, <code>Id</code>, and <code>id</code> are supported for inference. The generated GORM tag always references the original struct field name (for example, <code>foreignKey:BillCategoryId</code>).'
          },
          {
            type: 'paragraph',
            text: 'By default, it generates only the <strong>has-many</strong> side on the target model (for example, <code>PaymentMethods []PaymentMethod</code> on <code>User</code>).'
          },
          {
            type: 'note',
            kind: 'info',
            text: 'Use <code>--with-belongs-to</code> to also generate the <strong>belongs-to</strong> side on the source model (for example, <code>User *User</code> on <code>PaymentMethod</code>).'
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'terminal',
            code: `grove make:relations [--path <models-dir>] [--dry-run] [--verbose] [--with-belongs-to] [--model <Model>]`
          },
          {
            type: 'table',
            head: ['Flag', 'Default', 'Description'],
            rows: [
              [
                '<code>--path</code>',
                '<code>internal/modules</code>',
                'Modules directory to scan (<code>model.go</code> per domain)'
              ],
              [
                '<code>--dry-run</code>',
                '<code>false</code>',
                'Preview inferred relations without writing files'
              ],
              [
                '<code>--verbose</code>',
                '<code>false</code>',
                'Print detailed relation inference logs'
              ],
              [
                '<code>--with-belongs-to</code>',
                '<code>false</code>',
                'Also generate belongs-to fields on source models'
              ],
              [
                '<code>--model</code>',
                'all models',
                'Limit processing to specific source model(s). Can be repeated or comma-separated.'
              ]
            ]
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'examples',
            code: `# infer relations for all models (has-many only)
grove make:relations

# preview inferred changes with detailed logs
grove make:relations --dry-run --verbose

# generate both sides (has-many + belongs-to)
grove make:relations --with-belongs-to

# process only specific source models
grove make:relations --model PaymentMethod
grove make:relations --model PaymentMethod,Order
grove make:relations --model PaymentMethod --model Order`
          },
          {
            type: 'note',
            kind: 'tip',
            text: 'The command avoids duplicate relation fields and only infers relationships from supported foreign key types.'
          }
        ]
      },
      {
        id: 'cmd-make-seeder',
        title: 'grove make:seeder',
        blocks: [
          {
            type: 'paragraph',
            text: 'Scaffolds a database seeder in <code>internal/database/seeders/</code>. Seeders implement a simple interface with <code>Name() string</code> and <code>Seed(db *gorm.DB) error</code>.'
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'terminal',
            code: `grove make:seeder <Name> [--path <dir>] [--register]`
          },
          {
            type: 'table',
            head: ['Flag', 'Default', 'Description'],
            rows: [
              [
                '<code>--path</code>',
                '<code>internal/database/seeders</code>',
                'Target seeders directory'
              ],
              [
                '<code>--register</code>',
                '<code>false</code>',
                'Register the new seeder in <code>internal/database/seeders/seeder.go</code> (adds it to the seeders list) if the file exists'
              ]
            ]
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'examples',
            code: `grove make:seeder Users
grove make:seeder BillCategories
grove make:seeder order_items`
          },
          {
            type: 'note',
            kind: 'tip',
            text: 'Keep seeders idempotent — running them multiple times should be safe.'
          }
        ]
      },
      {
        id: 'cmd-make-seed-runner',
        title: 'grove make:seed',
        blocks: [
          {
            type: 'paragraph',
            text: 'Scaffolds a dedicated seed runner entrypoint at <code>cmd/seed/main.go</code>. This runner initializes app globals and calls <code>internal/database/seeders.Run(app.DB)</code>.'
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'terminal',
            code: `grove make:seed [--path <dir>]`
          },
          {
            type: 'table',
            head: ['Flag', 'Default', 'Description'],
            rows: [
              [
                '<code>--path</code>',
                '<code>cmd/seed</code>',
                'Target directory (creates <code>main.go</code> inside it)'
              ]
            ]
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'examples',
            code: `grove make:seed
grove make:seed --path ./cmd/seed`
          }
        ]
      },
      {
        id: 'cmd-db-seed',
        title: 'grove db:seed',
        blocks: [
          {
            type: 'paragraph',
            text: 'Runs your project seeders by executing the dedicated seed runner via <code>go run</code>. By default it runs <code>go run ./cmd/seed</code>.'
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'terminal',
            code: `grove db:seed [--package <path>] [--env-file <file>]`
          },
          {
            type: 'table',
            head: ['Flag', 'Default', 'Description'],
            rows: [
              [
                '<code>--package</code>',
                '<code>./cmd/seed</code>',
                'Go package/path to execute (seed runner)'
              ],
              [
                '<code>--env-file</code>',
                '<code>.env</code>',
                'Env file to load before running seeders (only if present; does not override existing env vars)'
              ]
            ]
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'examples',
            code: `grove db:seed
grove db:seed --package ./cmd/seed
grove db:seed --env-file .env`
          },
          {
            type: 'note',
            kind: 'info',
            text: 'If you have not scaffolded the seed runner yet, run <code>grove make:seed</code> first.'
          }
        ]
      },
      {
        id: 'cmd-make-resource',
        title: 'grove make:resource',
        blocks: [
          {
            type: 'paragraph',
            text: 'Scaffolds a full domain module (model, service, controller, DTO, and OpenAPI docs) in <code>internal/modules/&lt;domain&gt;/</code> and registers it in <code>internal/modules/register.go</code>. Equivalent to running <code>grove make:model &lt;Name&gt; -r</code>. Every file respects the <strong>SKIPPED</strong> rule — existing files are never overwritten.'
          },
          {
            type: 'paragraph',
            text: 'The entity name is <strong>automatically singularized</strong> before generating files, so you can pass the name in any form and Grove will always produce consistent output.'
          },
          {
            type: 'note',
            kind: 'warning',
            text: '<strong>Migrations are not generated automatically.</strong> After scaffolding, add your fields to the model, then run <code>grove make:migration create_&lt;table&gt;_table</code> to generate the SQL. This ensures the migration reflects the actual fields you defined.'
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'terminal',
            code: `grove make:resource <Name>`
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'examples',
            code: `grove make:resource Post        # → Post model, posts table
grove make:resource Posts       # → Post model, posts table (singularized)
grove make:resource BlogPost    # → BlogPost model, blog_posts table
grove make:resource order_items # → OrderItem model, order_items table`
          },
          {
            type: 'table',
            head: ['Input', 'Resolved name', 'Module path', 'Table'],
            rows: [
              [
                '<code>Post</code>',
                '<code>Post</code>',
                '<code>internal/modules/posts/</code>',
                '<code>posts</code>'
              ],
              [
                '<code>Posts</code>',
                '<code>Post</code>',
                '<code>internal/modules/posts/</code>',
                '<code>posts</code>'
              ],
              [
                '<code>BlogPost</code>',
                '<code>BlogPost</code>',
                '<code>internal/modules/blog_posts/</code>',
                '<code>blog_posts</code>'
              ],
              [
                '<code>order_items</code>',
                '<code>OrderItem</code>',
                '<code>internal/modules/order_items/</code>',
                '<code>order_items</code>'
              ]
            ]
          },
          {
            type: 'paragraph',
            text: 'After running <code>make:resource</code>, follow this workflow:'
          },
          {
            type: 'table',
            head: ['Step', 'Action'],
            rows: [
              [
                '1',
                'Add fields to the model in <code>internal/modules/&lt;domain&gt;/model.go</code>'
              ],
              [
                '2',
                'Add request/response fields to the DTO in <code>internal/modules/&lt;domain&gt;/dto.go</code>'
              ],
              [
                '3',
                'Run <code>grove make:migration create_&lt;table&gt;_table</code> to generate the SQL diff'
              ],
              ['4', 'Run <code>grove migrate</code> to apply it'],
              [
                '5',
                'Verify module registered in <code>internal/modules/register.go</code>'
              ]
            ]
          },
          {
            type: 'note',
            kind: 'tip',
            text: 'This is the fastest way to bootstrap a new feature — scaffolding takes seconds, and the explicit migration step ensures the SQL always matches the fields you actually defined.'
          }
        ]
      },
      {
        id: 'cmd-migrate',
        title: 'grove migrate',
        blocks: [
          {
            type: 'paragraph',
            text: 'Applies all pending migrations to the database using <code>atlas migrate apply</code>. Grove parses the Atlas output and renders it with badges, syntax-highlighted SQL keywords and a formatted summary.'
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'terminal',
            code: `grove migrate [--env <atlas-env>]`
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'output example',
            code: `  Running migrations (atlas migrate apply --env local)

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
  9 sql statements`
          },
          {
            type: 'paragraph',
            text: 'Each migration version gets a <code>MIGRATE</code> badge. SQL keywords (<code>CREATE TABLE</code>, <code>ALTER TABLE</code>, <code>DROP INDEX</code>, etc.) are highlighted in cyan. Continuation lines of multi-line statements are dimmed. If all migrations are already applied, Grove prints an <code>UP TO DATE</code> badge instead.'
          },
          {
            type: 'table',
            head: ['Subcommand', 'Description'],
            rows: [
              ['<code>grove migrate</code>', 'Apply all pending migrations'],
              [
                '<code>grove migrate:rollback</code>',
                'Roll back the last applied migration'
              ],
              [
                '<code>grove migrate:status</code>',
                'Show which migrations are applied / pending'
              ],
              [
                '<code>grove migrate:fresh</code>',
                'Drop all tables and re-apply every migration'
              ],
              [
                '<code>grove migrate:hash</code>',
                'Rehash the <code>atlas.sum</code> file'
              ]
            ]
          },
          {
            type: 'note',
            kind: 'warning',
            text: '<code>migrate:fresh</code> is a <strong>destructive operation</strong>. It will drop all tables. Only use it on development databases.'
          }
        ]
      },
      {
        id: 'cmd-completion',
        title: 'grove completion',
        blocks: [
          {
            type: 'paragraph',
            text: 'Generates a shell completion script so you get tab-completion for all grove commands and flags.'
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'terminal',
            code: `grove completion [bash|zsh|fish|powershell]`
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'zsh — persist',
            code: `echo 'source <(grove completion zsh)' >> ~/.zshrc`
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'fish — persist',
            code: `grove completion fish > ~/.config/fish/completions/grove.fish`
          }
        ]
      },
      {
        id: 'cmd-update',
        title: 'grove update',
        blocks: [
          {
            type: 'paragraph',
            text: 'Updates Grove-managed project dependencies to their latest versions and runs <code>go mod tidy</code> to clean up the module graph.'
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'terminal',
            code: `grove update`
          },
          {
            type: 'table',
            head: ['Dependency', 'Action'],
            rows: [
              [
                '<a href="https://github.com/caiolandgraf/gest" target="_blank">gest library</a>',
                'Updated to <code>@latest</code> in <code>go.mod</code> via <code>go get</code>'
              ],
              [
                '<a href="https://github.com/caiolandgraf/gest" target="_blank">gest CLI binary</a>',
                'Installed globally to <code>$GOPATH/bin</code> via <code>go install</code>'
              ],
              [
                'Module graph',
                '<code>go mod tidy</code> is run automatically after updates'
              ]
            ]
          },
          {
            type: 'note',
            kind: 'info',
            text: '<code>grove test</code> does not update gest automatically on every run. Use <code>grove update</code> whenever you want to pull in a newer version of both the library and the CLI.'
          }
        ]
      }
    ]
  },

  // ─────────────────────────────────────────────
  // TESTING
  // ─────────────────────────────────────────────
  {
    id: 'testing',
    title: 'Testing',
    items: [
      {
        id: 'testing-overview',
        title: 'Overview',
        blocks: [
          {
            type: 'paragraph',
            text: 'Grove uses <a href="https://github.com/caiolandgraf/gest" target="_blank">gest v2</a> as its testing framework — a Jest-inspired library for Go that runs on top of the native <code>go test</code> engine. You get beautiful colored output, descriptive failure messages and a fluent assertion API, while keeping full IDE support, caching, <code>-race</code> detection and real coverage for free.'
          },
          {
            type: 'paragraph',
            text: 'All test files live in <code>internal/tests/</code> as standard <code>*_test.go</code> files. Each file has a <code>func Test&lt;Name&gt;(t *testing.T)</code> entry point that calls <code>s.Run(t)</code> — no separate <code>main.go</code>, no <code>init()</code> registration.'
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'project layout',
            code: `internal/
└── tests/
    ├── post_test.go            # grove make:test Post
    └── user_service_test.go    # grove make:test UserService`
          },
          {
            type: 'note',
            kind: 'info',
            text: 'gest v2 uses standard <code>*_test.go</code> files — the same convention as <code>go test</code>. You can run <code>go test ./internal/tests/...</code> at any time without the gest CLI.'
          }
        ]
      },
      {
        id: 'testing-installation',
        title: 'Installing the gest CLI',
        blocks: [
          {
            type: 'paragraph',
            text: 'The gest CLI renders beautiful Jest-style output by wrapping <code>go test -v -json</code>. It is optional — <code>grove test</code> falls back to plain <code>go test -v</code> automatically when it is not installed.'
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'terminal',
            code: `go install github.com/caiolandgraf/gest/v2/cmd/gest@latest`
          },
          {
            type: 'note',
            kind: 'tip',
            text: 'Run <code>grove update</code> at any time to update both the gest library in your <code>go.mod</code> and the global gest CLI binary to their latest versions.'
          }
        ]
      },
      {
        id: 'cmd-make-test',
        title: 'grove make:test',
        blocks: [
          {
            type: 'paragraph',
            text: "Scaffolds a new gest v2 test file in <code>internal/tests/</code>. The generated file is a standard <code>*_test.go</code> file with a <code>func Test&lt;Name&gt;(t *testing.T)</code> entry point. If the file already exists the command prints <strong>SKIPPED</strong> and exits cleanly. On the first call, gest is added to the project's <code>go.mod</code> automatically via <code>go get</code>."
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'terminal',
            code: `grove make:test <Name>`
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'examples',
            code: `grove make:test User
grove make:test AuthService
grove make:test order_calculations   # → OrderCalculations`
          },
          {
            type: 'paragraph',
            text: 'The generated file follows the gest v2 convention:'
          },
          {
            type: 'code',
            lang: 'go',
            label: 'internal/tests/user_test.go (generated)',
            code: `package myapp

import (
	"testing"

	"github.com/caiolandgraf/gest/v2/gest"
)

func TestUser(t *testing.T) {
	s := gest.Describe("User")

	s.It("should work", func(t *gest.T) {
		// TODO: write your test here
		t.Expect(true).ToBeTrue()
	})

	s.Run(t)
}`
          }
        ]
      },
      {
        id: 'cmd-test',
        title: 'grove test',
        blocks: [
          {
            type: 'paragraph',
            text: 'Runs every <code>*_test.go</code> file in <code>internal/tests/</code> using the gest CLI for beautiful Jest-style output. If the gest CLI is not installed, grove falls back to <code>go test -v</code> automatically. Pass <code>-c</code> for a per-suite coverage report. Pass <code>-w</code> to enter watch mode. Combine both as <code>-wc</code>.'
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'terminal',
            code: `grove test [-c] [-w]`
          },
          {
            type: 'table',
            head: ['Flag', 'Description'],
            rows: [
              [
                '<code>-c</code>, <code>--coverage</code>',
                'Display a per-suite pass-rate bar after the run'
              ],
              [
                '<code>-w</code>, <code>--watch</code>',
                'Watch mode: re-run tests on file changes. Delegates to <code>gest --watch</code> when the CLI is installed; falls back to a polling loop otherwise.'
              ],
              [
                '<code>-wc</code>',
                'Watch mode with coverage report (shorthand for <code>-w -c</code>)'
              ]
            ]
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'examples',
            code: `grove test         # run all tests
grove test -c      # run all tests + coverage report
grove test -w      # watch mode — re-run tests on every save
grove test -wc     # watch mode + coverage report

# you can also use go test directly at any time:
go test ./internal/tests/...`
          },
          {
            type: 'note',
            kind: 'tip',
            text: 'Coverage bar colours: <strong>green</strong> ≥ 80% · <strong>yellow</strong> ≥ 50% · <strong>red</strong> &lt; 50%. Pressing <kbd>Ctrl+C</kbd> in watch mode stops cleanly without printing a failure message.'
          },
          {
            type: 'note',
            kind: 'info',
            text: '<code>grove test</code> does not update gest automatically. Run <code>grove update</code> to pull the latest version of both the gest library and the gest CLI.'
          }
        ]
      },
      {
        id: 'testing-writing-tests',
        title: 'Writing Tests',
        blocks: [
          {
            type: 'paragraph',
            text: 'Each test file creates a suite with <code>gest.Describe()</code>, adds cases with <code>s.It()</code>, then hands off to <code>go test</code> via <code>s.Run(t)</code>. The wrapping <code>func Test&lt;Name&gt;(t *testing.T)</code> is a standard Go test function — IDEs, <code>go test</code> and the gest CLI all discover it automatically.'
          },
          {
            type: 'code',
            lang: 'go',
            label: 'internal/tests/post_test.go',
            code: `package myapp

import (
	"testing"

	"github.com/caiolandgraf/gest/v2/gest"
	"your/module/internal/modules/posts"
)

func TestPost(t *testing.T) {
	s := gest.Describe("Post")

	s.It("TableName should return 'posts'", func(t *gest.T) {
		t.Expect(posts.Post{}.TableName()).ToBe("posts")
	})

	s.It("slug should contain the title", func(t *gest.T) {
		post := posts.Post{Title: "Hello World"}
		t.Expect(post.Slug()).ToContain("hello")
	})

	s.It("draft post should not be published", func(t *gest.T) {
		post := posts.Post{Published: false}
		t.Expect(post.Published).ToBeFalse()
	})

	s.Run(t)
}`
          },
          {
            type: 'paragraph',
            text: 'Available matchers:'
          },
          {
            type: 'table',
            head: ['Matcher', 'Description'],
            rows: [
              ['<code>.ToBe(v)</code>', 'Strict equality (<code>==</code>)'],
              [
                '<code>.ToEqual(v)</code>',
                'Deep equality (<code>reflect.DeepEqual</code>)'
              ],
              ['<code>.ToBeNil()</code>', 'Value is <code>nil</code>'],
              ['<code>.ToBeTrue()</code>', 'Value is <code>true</code>'],
              ['<code>.ToBeFalse()</code>', 'Value is <code>false</code>'],
              ['<code>.ToContain(s)</code>', 'String contains substring'],
              [
                '<code>.ToHaveLength(n)</code>',
                'Length of string, slice or map'
              ],
              [
                '<code>.ToBeGreaterThan(n)</code>',
                'Number greater than <code>n</code>'
              ],
              [
                '<code>.ToBeLessThan(n)</code>',
                'Number less than <code>n</code>'
              ],
              [
                '<code>.ToBeCloseTo(n, delta?)</code>',
                'Float approximately equal (default ±0.001)'
              ],
              [
                '<code>.Not()</code>',
                'Negates any matcher — e.g. <code>.Not().ToBeNil()</code>'
              ]
            ]
          }
        ]
      }
    ]
  },

  // ─────────────────────────────────────────────
  // ARCHITECTURE
  // ─────────────────────────────────────────────
  {
    id: 'architecture',
    title: 'Architecture',
    items: [
      {
        id: 'project-structure',
        title: 'Project Structure',
        blocks: [
          {
            type: 'paragraph',
            text: 'Every Grove project follows the <strong>go-project-base</strong> modular MSC layout. Each domain is a self-contained module under <code>internal/modules/</code>; shared infrastructure lives under <code>internal/app/</code>.'
          },
          {
            type: 'code',
            lang: 'text',
            label: 'project layout',
            code: projectStructureAscii
          },
          {
            type: 'folder-structure'
          },
          {
            type: 'paragraph',
            text: 'Request flow: <code>Routes → Middlewares → Module (Controller → Service → Repository) → Database</code>. The <code>internal/</code> package boundary prevents external packages from importing your application internals.'
          },
          {
            type: 'table',
            head: ['Directory', 'Purpose'],
            rows: [
              ['<code>cmd/api/</code>', 'Application entrypoint — boots config, connects infra, starts Fuego'],
              ['<code>cmd/atlas/</code>', 'Atlas GORM schema loader (program mode via <code>atlas.hcl</code>)'],
              ['<code>cmd/scalar/</code>', 'Scalar OpenAPI documentation UI handler'],
              ['<code>internal/app/config/</code>', 'DB, Redis, OTel, slog, metrics, and session initializers'],
              ['<code>internal/app/database/</code>', 'Generic <code>Repository[T]</code> and Atlas model registry'],
              ['<code>internal/app/helpers/</code>', 'JSON utilities (<code>jsonutils/</code>) and request validation (<code>validator/</code>)'],
              ['<code>internal/app/middleware/</code>', 'CORS, session, and observability middlewares'],
              ['<code>internal/app/types/</code>', 'Shared HTTP types — errors, messages, health responses'],
              ['<code>internal/app/router/</code>', 'Declarative OpenAPI docs per endpoint (<code>router.Doc</code>)'],
              ['<code>internal/modules/</code>', 'Self-contained domains — model, dto, service, controller, docs'],
              ['<code>internal/modules/module.go</code>', 'Module interface, Boot struct, and Factory type'],
              ['<code>internal/modules/register.go</code>', 'Module registry — <code>grove make:resource</code> adds new domains here'],
              ['<code>internal/routes/</code>', 'Global middleware, health checks (<code>health.go</code>), module mounting'],
              ['<code>infra/</code>', 'Prometheus, Grafana, Loki, Jaeger, and Promtail configuration'],
              ['<code>logs/</code>', 'Structured JSON log files tailed by Promtail → Loki'],
              ['<code>migrations/</code>', 'Versioned Atlas SQL files + <code>atlas.sum</code>'],
              ['<code>doc/openapi.json</code>', 'Generated OpenAPI 3 spec served at <code>/swagger</code>'],
              ['<code>docker-compose.yml</code>', 'PostgreSQL, Redis, Jaeger, Prometheus, Loki, Grafana stack'],
              ['<code>atlas.hcl</code>', 'Atlas environments — loads schema via <code>go run ./cmd/atlas</code>'],
              ['<code>Makefile</code>', 'Dev commands: install, run, dev, migrate-up, db-reset'],
              ['<code>grove.toml</code>', 'Optional — Grove CLI hot-reload config for <code>grove dev</code>']
            ]
          }
        ]
      },
      {
        id: 'arch-models',
        title: 'Models',
        blocks: [
          {
            type: 'paragraph',
            text: 'Models live in each domain module (<code>internal/modules/&lt;domain&gt;/model.go</code>) as plain Go structs with GORM tags. Every model registers itself for Atlas via <code>database.Register</code> in <code>init()</code> and ships with a typed <strong>Repo</strong> accessor.'
          },
          {
            type: 'code',
            lang: 'go',
            label: 'internal/modules/posts/model.go',
            code: `package posts

import (
	"time"

	"your/module/internal/app/database"
	"gorm.io/gorm"
)

func init() {
	database.Register(&Post{})
}

type Post struct {
	ID        string         \`gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"\`
	Title     string         \`gorm:"not null"      json:"title"\`
	Content   string         \`gorm:"type:text"     json:"content"\`
	Published bool           \`gorm:"default:false" json:"published"\`
	AuthorID  string         \`gorm:"not null"      json:"author_id"\`
	CreatedAt time.Time      \`gorm:"autoCreateTime" json:"created_at"\`
	UpdatedAt time.Time      \`gorm:"autoUpdateTime" json:"updated_at"\`
	DeletedAt gorm.DeletedAt \`gorm:"index"          json:"-"\`
}

func (Post) TableName() string { return "posts" }

type Repo struct {
	*database.Repository[Post]
}

func Posts(db *gorm.DB) *Repo {
	return &Repo{Repository: database.New[Post](db)}
}`
          },
          {
            type: 'paragraph',
            text: 'The repository exposes typed methods so you never write raw GORM queries for standard operations:'
          },
          {
            type: 'code',
            lang: 'go',
            label: 'usage',
            code: `repo := posts.Posts(db)

// Find by primary key
post, err := repo.Find(id)

// List all
all, err := repo.All()

// Create
err := repo.Create(&post)

// Update
err := repo.Update(&post)

// Soft-delete
err := repo.Delete(id)

// Custom query — drop down to GORM when needed
var result []posts.Post
err := repo.DB().Where("published = ?", true).Find(&result).Error`
          }
        ]
      },
      {
        id: 'arch-services',
        title: 'Services',
        blocks: [
          {
            type: 'paragraph',
            text: 'Services live in <code>internal/modules/&lt;domain&gt;/service.go</code> and house your core business logic, orchestrate database transactions, and handle operations independent of HTTP requests. Each module exposes <code>WireService</code> for dependency injection from the controller.'
          },
          {
            type: 'code',
            lang: 'go',
            label: 'internal/modules/posts/service.go',
            code: `package posts

import (
	"context"
	"gorm.io/gorm"
)

type Service interface {
	CreatePost(ctx context.Context, item *Post) error
	GetPostByID(ctx context.Context, id string) (*Post, error)
}

type service struct {
	repo *Repo
}

func NewService(repo *Repo) Service {
	return &service{repo: repo}
}

func WireService(db *gorm.DB) Service {
	return NewService(Posts(db))
}

func (s *service) CreatePost(ctx context.Context, item *Post) error {
	// Implement validations, triggers, or custom logic here before database persistence
	return s.repo.Create(item)
}

func (s *service) GetPostByID(ctx context.Context, id string) (*Post, error) {
	return s.repo.Find(id)
}`
          }
        ]
      },
      {
        id: 'arch-controllers',
        title: 'Controllers',
        blocks: [
          {
            type: 'paragraph',
            text: 'Controllers live in <code>internal/modules/&lt;domain&gt;/controller.go</code> and bind HTTP endpoints to service layers under the modular MCS architecture. In v2.2.0, each module exposes <code>Wire</code> and <code>Mount</code> helpers — the controller wires its own dependencies and registers routes on the API group.'
          },
          {
            type: 'code',
            lang: 'go',
            label: 'internal/modules/posts/controller.go',
            code: `package posts

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"your/module/internal/app/router"
	"github.com/go-fuego/fuego"
	"gorm.io/gorm"
)

type Controller struct {
	service Service
}

func NewController(service Service) *Controller {
	return &Controller{service: service}
}

func Wire(db *gorm.DB) *Controller {
	return NewController(WireService(db))
}

func (ctrl *Controller) Mount(api *fuego.Server, session *scs.SessionManager) {
	group := fuego.Group(api, "/posts")

	router.Get(group, "/{post_id}", ctrl.GetPost, GetPostDoc, session)
	router.Post(group, "/", ctrl.CreatePost, CreatePostDoc, session)
}

// GetPost handles GET /posts/:post_id
func (ctrl *Controller) GetPost(c fuego.ContextNoBody) (*PostResponse, error) {
	id := c.PathParam("post_id")

	post, err := ctrl.service.GetPostByID(c.Context(), id)
	if err != nil {
		return nil, fuego.HTTPError{
			Status: http.StatusNotFound,
			Err:    err,
		}
	}
	return toPostDTO(post), nil
}

// CreatePost handles POST /posts
func (ctrl *Controller) CreatePost(c fuego.ContextWithBody[CreatePostRequest]) (*PostResponse, error) {
	body, err := c.Body()
	if err != nil {
		return nil, fuego.HTTPError{Status: http.StatusBadRequest, Err: err}
	}

	post := &Post{
		Title:   body.Title,
		Content: body.Content,
	}

	if err := ctrl.service.CreatePost(c.Context(), post); err != nil {
		return nil, fuego.HTTPError{Status: http.StatusUnprocessableEntity, Err: err}
	}
	return toPostDTO(post), nil
}

func toPostDTO(m *Post) *PostResponse {
	return &PostResponse{
		ID:        m.ID,
		Title:     m.Title,
		Content:   m.Content,
		Published: m.Published,
	}
}`
          }
        ]
      },
      {
        id: 'arch-dto',
        title: 'DTOs',
        blocks: [
          {
            type: 'paragraph',
            text: 'DTOs (Data Transfer Objects) live in <code>internal/modules/&lt;domain&gt;/dto.go</code> and define the exact shape of your API requests and responses. Keeping them in the same module as the model means your API contract is colocated with the domain it serves.'
          },
          {
            type: 'code',
            lang: 'go',
            label: 'internal/modules/posts/dto.go',
            code: `package posts

type CreatePostRequest struct {
	Title    string \`json:"title"     validate:"required,min=1,max=255"\`
	Content  string \`json:"content"   validate:"required"\`
	AuthorID string \`json:"author_id" validate:"required,uuid"\`
}

type UpdatePostRequest struct {
	Title   *string \`json:"title"   validate:"omitempty,min=1,max=255"\`
	Content *string \`json:"content" validate:"omitempty"\`
}

type PostResponse struct {
	ID        string \`json:"id"\`
	Title     string \`json:"title"\`
	Content   string \`json:"content"\`
	Published bool   \`json:"published"\`
	CreatedAt string \`json:"created_at"\`
}

type PostsListResponse struct {
	Items []PostResponse \`json:"items"\`
	Total int            \`json:"total"\`
}`
          },
          {
            type: 'note',
            kind: 'tip',
            text: 'Use pointer fields (<code>*string</code>) in update requests so you can distinguish between "field not sent" and "field sent as empty string". Grove\'s generated stubs leave this to you.'
          }
        ]
      },
      {
        id: 'arch-routes',
        title: 'Routes',
        blocks: [
          {
            type: 'paragraph',
            text: 'Global routes and middleware are configured in <code>internal/routes/routes.go</code>. Domain modules register themselves via <code>modules.Mount</code>, which iterates the registry in <code>internal/modules/register.go</code>. Each module\'s <code>Mount</code> method wires its own routes — no manual per-resource registration needed.'
          },
          {
            type: 'code',
            lang: 'go',
            label: 'internal/routes/routes.go',
            code: `package routes

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"your/module/internal/app/middleware"
	"your/module/internal/modules"
	"github.com/go-fuego/fuego"
	"github.com/gomodule/redigo/redis"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"gorm.io/gorm"
)

func SetupRoutes(
	s *fuego.Server,
	db *gorm.DB,
	redisPool *redis.Pool,
	session *scs.SessionManager,
	metricsHandler http.Handler,
) {
	fuego.Use(s, otelhttp.NewMiddleware("my-api"))
	fuego.Use(s, middleware.RouteTagMiddleware)
	fuego.Use(s, middleware.CORSMiddleware(middleware.DefaultCORSConfig()))
	fuego.Use(s, middleware.SessionMiddleware(session))

	fuego.Get(s, "/", healthCheck)
	fuego.Get(s, "/health", healthCheckDetailed(db, redisPool))
	fuego.GetStd(s, "/metrics", metricsHandler.ServeHTTP)

	api := fuego.Group(s, "/api/v1")
	modules.Mount(api, modules.Boot{DB: db, Session: session})
}`
          },
          {
            type: 'code',
            lang: 'go',
            label: 'internal/modules/register.go',
            code: `package modules

import (
	"your/module/internal/modules/auth"
	"your/module/internal/modules/posts"
	"github.com/go-fuego/fuego"
)

var registry = []Factory{
	func(b Boot) Module { return posts.Wire(b.DB) },
	func(b Boot) Module { return auth.Wire(b.DB, b.Session) },
}

func Mount(api *fuego.Server, boot Boot) {
	for _, factory := range registry {
		factory(boot).Mount(api, boot.Session)
	}
}`
          },
          {
            type: 'paragraph',
            text: 'fuego automatically generates and serves an OpenAPI 3.1 spec from your route registrations. Visit <code>/swagger</code> in development to explore your API interactively.'
          }
        ]
      },
      {
        id: 'arch-migrations',
        title: 'Migrations',
        blocks: [
          {
            type: 'paragraph',
            text: 'Grove uses <strong>Atlas</strong> for schema migrations. Atlas diffs your GORM models against the live database schema and produces precise, versioned SQL files — no manual SQL writing required.'
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'typical workflow',
            code: `# 1. Add a new field to your model
#    e.g. add "PublishedAt *time.Time" to internal/modules/posts/model.go

# 2. Generate the migration
grove make:migration add_published_at_to_posts

# 3. Review the generated SQL
cat migrations/20240801120000_add_published_at_to_posts.sql

# 4. Apply it
grove migrate`
          },
          {
            type: 'code',
            lang: 'sql',
            label: 'migrations/20240801120000_add_published_at_to_posts.sql',
            code: `-- Add column
ALTER TABLE "posts" ADD COLUMN "published_at" timestamptz NULL;

-- Create index for efficient filtering by published_at
CREATE INDEX "idx_posts_published_at" ON "posts" ("published_at");`
          },
          {
            type: 'paragraph',
            text: "<code>grove migrate</code> parses the Atlas output and renders it with Grove's colour palette — each migration version gets a <code>MIGRATE</code> badge, SQL keywords are highlighted in cyan, and a summary line shows total time, migrations and statements applied:"
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'grove migrate output',
            code: `  Running migrations (atlas migrate apply --env local)

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
  9 sql statements`
          },
          {
            type: 'note',
            kind: 'tip',
            text: 'If all migrations are already applied, Grove prints an <code>UP TO DATE</code> badge instead of a migration list.'
          },
          {
            type: 'note',
            kind: 'info',
            text: 'If you see a checksum error after editing a migration file manually, run <code>grove migrate:hash</code> to re-hash the <code>atlas.sum</code> file.'
          }
        ]
      }
    ]
  },

  // ─────────────────────────────────────────────
  // CONFIGURATION
  // ─────────────────────────────────────────────
  {
    id: 'configuration',
    title: 'Configuration',
    items: [
      {
        id: 'env-vars',
        title: 'Environment Variables',
        blocks: [
          {
            type: 'paragraph',
            text: 'Grove projects use a <code>.env</code> file for local configuration. The <code>.env.example</code> file is committed to the repository as a template.'
          },
          {
            type: 'code',
            lang: 'bash',
            label: '.env.example',
            code: `# Server
PORT=8080
ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=my_api_dev
DB_USER=postgres
DB_PASSWORD=postgres
DB_SSL_MODE=disable

# Computed — used by Atlas
DATABASE_URL=postgresql://\${DB_USER}:\${DB_PASSWORD}@\${DB_HOST}:\${DB_PORT}/\${DB_NAME}?sslmode=\${DB_SSL_MODE}`
          },
          {
            type: 'note',
            kind: 'info',
            text: 'The <code>.env</code> file is git-ignored by default. Never commit real credentials. Use secrets managers (e.g. AWS Secrets Manager, Doppler) in production.'
          }
        ]
      },
      {
        id: 'atlas-config',
        title: 'Atlas Configuration',
        blocks: [
          {
            type: 'paragraph',
            text: '<code>atlas.hcl</code> sits at the project root and tells Atlas how to connect to your database. GORM models are loaded via <strong>program mode</strong> — Atlas runs <code>go run ./cmd/atlas</code>, which imports every module and reads models registered via <code>database.Register</code> in each <code>model.go</code>.'
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'atlas.hcl',
            code: `data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "./cmd/atlas",
  ]
}

env "local" {
  src = data.external_schema.gorm.url
  url = getenv("DATABASE_URL")
  dev = "docker://postgres/15/dev"
  migration {
    dir = "file://migrations"
  }
}

env "dev" {
  src = data.external_schema.gorm.url
  url = getenv("DATABASE_URL")
  dev = "docker://postgres/15/dev"
  migration {
    dir = "file://migrations"
  }
}`
          },
          {
            type: 'note',
            kind: 'tip',
            text: 'The <code>local</code> environment uses a Docker-based shadow database for diffing. Make sure Docker is running when executing <code>grove make:migration</code> locally. Grove may cache the Atlas GORM provider binary under <code>.grove/bin/atlas-gorm</code> for faster subsequent runs.'
          }
        ]
      },
      {
        id: 'app-singleton',
        title: 'App Singleton',
        blocks: [
          {
            type: 'paragraph',
            text: 'The <code>internal/app/</code> package holds shared infrastructure — database, Redis, session, metrics, and config. Everything is initialized once at startup in <code>cmd/api/main.go</code> and passed to modules via <code>modules.Boot</code>.'
          },
          {
            type: 'code',
            lang: 'go',
            label: 'cmd/api/main.go (simplified)',
            code: `package main

import (
	"your/module/internal/app/config"
	"your/module/internal/routes"
	"github.com/go-fuego/fuego"
)

func main() {
	db := config.InitDB()
	redisPool := config.InitRedis()
	session := config.InitSession(redisPool)
	metrics := config.InitMetrics()

	s := fuego.NewServer(fuego.WithAddr(":8080"))

	routes.SetupRoutes(s, db, redisPool, session, metrics)

	s.Run()
}`
          }
        ]
      }
    ]
  }
]

// ─────────────────────────────────────────────
// Flat search index built from sections above
// ─────────────────────────────────────────────
export const searchIndex = sections.flatMap(section =>
  section.items.map(item => ({
    id: item.id,
    title: item.title,
    section: section.title,
    sectionId: section.id,
    text: item.blocks
      .filter(b => b.type === 'paragraph' || b.type === 'note')
      .map(b => b.text || b.text)
      .join(' ')
      .replace(/<[^>]+>/g, ''),
    url: `/docs#${item.id}`
  }))
)
