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
            text: 'Grove is distributed as a single binary via <code>go install</code>. Make sure you have Go 1.22 or newer and that <code>$GOPATH/bin</code> (or <code>$HOME/go/bin</code>) is in your <code>$PATH</code>.'
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

# 2. Enter the project directory
cd my-api

# 3. Configure your environment
cp .env.example .env
# edit .env with your database credentials

# 4. Start the development server with built-in hot reload
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
                text: 'Edit <code>internal/models/post.go</code> to add your fields (e.g. <code>Title</code>, <code>Body</code>). Also fill in <code>internal/dto/post-dto.go</code> with your request/response fields.'
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
                text: "Add your routes to <code>internal/routes/</code> and you're done."
              }
            ]
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
                '1.22+',
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
            text: 'Downloads and scaffolds a complete Grove project from the official template repository on GitHub.'
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'terminal',
            code: `grove setup <project-name> [--module <go-module-path>]`
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
grove setup my-api --module github.com/acme/my-api`
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
watch_dirs  = ["."]
exclude     = [".grove", "vendor", "node_modules", "tests"]
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
                '<code>["."]</code>',
                'Directories to watch for file changes (recursive)'
              ],
              [
                '<code>exclude</code>',
                '<code>[".grove", "vendor", "node_modules", "tests"]</code>',
                'Directory names to ignore (the <code>internal/tests/</code> directory is always excluded so test saves never trigger a rebuild)'
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
        id: 'cmd-make-model',
        title: 'grove make:model',
        blocks: [
          {
            type: 'paragraph',
            text: 'Scaffolds a new GORM model in <code>internal/models/</code>. The name is automatically converted to PascalCase and the file to snake_case. Combine flags to scaffold additional layers in the same step — or use <code>-r</code> as a shorthand for the full resource.'
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
            text: 'The generated model includes UUID primary key, <code>CreatedAt</code>, <code>UpdatedAt</code>, soft-delete (<code>DeletedAt</code>) and a typed repository accessor:'
          },
          {
            type: 'code',
            lang: 'go',
            label: 'internal/models/post.go',
            code: `package models

import (
	"time"

	"your/module/internal/app"
	"your/module/internal/database"
	"gorm.io/gorm"
)

type Post struct {
	ID        string         \`gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"\`
	CreatedAt time.Time      \`gorm:"autoCreateTime"  json:"created_at"\`
	UpdatedAt time.Time      \`gorm:"autoUpdateTime"  json:"updated_at"\`
	DeletedAt gorm.DeletedAt \`gorm:"index"           json:"-"\`
}

func (Post) TableName() string { return "posts" }

// Posts returns a repository scoped to the Post model.
func Posts() *database.Repository[Post] {
	return database.New[Post](app.DB)
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
            text: 'Scaffolds a new fuego controller with full CRUD handlers in <code>internal/controllers/</code>. If the file already exists the command prints <strong>SKIPPED</strong> and exits cleanly.'
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'terminal',
            code: `grove make:controller <Name>`
          },
          {
            type: 'code',
            lang: 'go',
            label: 'internal/controllers/post-controller.go (generated)',
            code: `package controllers

import (
	"net/http"

	"your/module/internal/dto"
	"your/module/internal/models"
	"github.com/go-fuego/fuego"
)

func GetPost(c fuego.ContextNoBody) (*dto.PostResponse, error) {
	id := c.PathParam("post_id")

	item, err := models.Posts().Find(id)
	if err != nil {
		return nil, fuego.HTTPError{Status: http.StatusNotFound, Err: err}
	}
	return toPostDTO(item), nil
}

func ListPosts(c fuego.ContextNoBody) (*dto.PostsListResponse, error) {
	items, err := models.Posts().All()
	if err != nil {
		return nil, fuego.HTTPError{Status: http.StatusInternalServerError, Err: err}
	}
	result := make([]dto.PostResponse, len(items))
	for i, item := range items {
		result[i] = *toPostDTO(&item)
	}
	return &dto.PostsListResponse{Items: result, Total: len(result)}, nil
}

func CreatePost(c fuego.ContextWithBody[dto.CreatePostRequest]) (*dto.PostResponse, error) {
	body, err := c.Body()
	if err != nil {
		return nil, fuego.HTTPError{Status: http.StatusBadRequest, Err: err}
	}
	item := &models.Post{
		// TODO: map fields from body
	}
	_ = body
	if err := models.Posts().Create(item); err != nil {
		return nil, fuego.HTTPError{Status: http.StatusBadRequest, Err: err}
	}
	return toPostDTO(item), nil
}

func toPostDTO(m *models.Post) *dto.PostResponse {
	return &dto.PostResponse{ID: m.ID}
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
            text: 'Scaffolds DTO request/response structs in <code>internal/dto/</code>. If the file already exists the command prints <strong>SKIPPED</strong> and exits cleanly.'
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
            label: 'internal/dto/post-dto.go (generated)',
            code: `package dto

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
            text: 'Scaffolds a new HTTP middleware in <code>internal/middleware/</code>. The generated file follows the standard <code>func(http.Handler) http.Handler</code> signature compatible with fuego and the standard library. If the file already exists the command prints <strong>SKIPPED</strong> and exits cleanly.'
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
            label: 'internal/middleware/auth-middleware.go (generated)',
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
            label: 'cmd/api/main.go',
            code: `s := fuego.NewServer(fuego.WithAddr(":8080"))

// Apply globally
s.Use(middleware.Auth)

routes.Register(s)
s.Run()`
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
        id: 'cmd-make-resource',
        title: 'grove make:resource',
        blocks: [
          {
            type: 'paragraph',
            text: 'Scaffolds a model, controller and DTO in one shot. Equivalent to running <code>grove make:model &lt;Name&gt; -r</code>. Every file respects the <strong>SKIPPED</strong> rule — existing files are never overwritten.'
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
            head: ['Input', 'Resolved name', 'File', 'Table'],
            rows: [
              [
                '<code>Post</code>',
                '<code>Post</code>',
                '<code>post.go</code>',
                '<code>posts</code>'
              ],
              [
                '<code>Posts</code>',
                '<code>Post</code>',
                '<code>post.go</code>',
                '<code>posts</code>'
              ],
              [
                '<code>BlogPost</code>',
                '<code>BlogPost</code>',
                '<code>blog_post.go</code>',
                '<code>blog_posts</code>'
              ],
              [
                '<code>order_items</code>',
                '<code>OrderItem</code>',
                '<code>order_item.go</code>',
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
                'Add fields to the model in <code>internal/models/&lt;name&gt;.go</code>'
              ],
              [
                '2',
                'Add request/response fields to the DTO in <code>internal/dto/&lt;name&gt;-dto.go</code>'
              ],
              [
                '3',
                'Run <code>grove make:migration create_&lt;table&gt;_table</code> to generate the SQL diff'
              ],
              ['4', 'Run <code>grove migrate</code> to apply it'],
              ['5', 'Register routes in <code>internal/routes/</code>']
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
	"your/module/internal/models"
)

func TestPost(t *testing.T) {
	s := gest.Describe("Post")

	s.It("TableName should return 'posts'", func(t *gest.T) {
		t.Expect(models.Post{}.TableName()).ToBe("posts")
	})

	s.It("slug should contain the title", func(t *gest.T) {
		post := models.Post{Title: "Hello World"}
		t.Expect(post.Slug()).ToContain("hello")
	})

	s.It("draft post should not be published", func(t *gest.T) {
		post := models.Post{Published: false}
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
            text: 'Every Grove project follows a clean, layered directory layout. The separation of concerns is intentional — each layer has one job.'
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'project tree',
            code: `my-api/
├── cmd/
│   └── api/
│       └── main.go              # Entry point — wires everything together
├── internal/
│   ├── app/                     # Shared singletons (DB, Redis, Session, Metrics)
│   ├── config/                  # Infrastructure initializers (DB, Redis, OTel, etc.)
│   ├── controllers/             # fuego route handlers
│   ├── database/                # Generic GORM repository
│   ├── dto/                     # Request and response types
│   ├── middleware/              # HTTP middlewares (CORS, session, observability)
│   ├── models/                  # GORM models
│   ├── routes/                  # Route registration
│   └── tests/                   # gest v2 test files
│       └── user_test.go         # Example test (generated by grove make:test)
├── migrations/                  # Atlas SQL migrations
├── infra/                       # Observability stack config (Prometheus, Grafana, Loki, Jaeger)
├── .env.example                 # Committed environment template
├── atlas.hcl                    # Atlas configuration
├── docker-compose.yml           # Full observability stack
└── grove.toml                   # Grove dev server configuration`
          },
          {
            type: 'paragraph',
            text: 'The <code>internal/</code> package boundary is intentional — it prevents external packages from importing your application internals, keeping the codebase clean as it grows.'
          },
          {
            type: 'table',
            head: ['Directory', 'Purpose'],
            rows: [
              [
                '<code>cmd/api/</code>',
                'Application entry point — wires singletons, routes and starts the server'
              ],
              [
                '<code>internal/app/</code>',
                'Shared singletons: DB, Redis, session store, metrics — initialized once at startup'
              ],
              [
                '<code>internal/config/</code>',
                'Infrastructure initializers for DB, Redis, OpenTelemetry and other external services'
              ],
              [
                '<code>internal/controllers/</code>',
                'fuego route handlers — one file per resource, OpenAPI inferred automatically'
              ],
              [
                '<code>internal/database/</code>',
                'Generic GORM repository (<code>Repository[T]</code>) used by all models'
              ],
              [
                '<code>internal/dto/</code>',
                'Request and response structs — decoupled from GORM models'
              ],
              [
                '<code>internal/middleware/</code>',
                'HTTP middlewares: CORS, session, observability, auth, etc.'
              ],
              [
                '<code>internal/models/</code>',
                'GORM models with typed repository accessors'
              ],
              [
                '<code>internal/routes/</code>',
                'Route registration — fuego typed routes wired to controllers'
              ],
              [
                '<code>internal/tests/</code>',
                'gest spec files — <code>main.go</code> is auto-created by <code>grove make:test</code>'
              ],
              [
                '<code>migrations/</code>',
                'Versioned Atlas SQL migration files + <code>atlas.sum</code>'
              ],
              [
                '<code>infra/</code>',
                'Observability stack configuration: Prometheus, Grafana, Loki, Jaeger'
              ],
              [
                '<code>docker-compose.yml</code>',
                'Spins up the full observability stack locally with a single command'
              ],
              [
                '<code>grove.toml</code>',
                'Optional Grove configuration — <code>[dev]</code> section for <code>grove dev</code>'
              ]
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
            text: 'Models are plain Go structs with GORM tags. Every generated model ships with a typed <strong>Repository</strong> accessor that gives you a clean, generic CRUD interface without writing boilerplate.'
          },
          {
            type: 'code',
            lang: 'go',
            label: 'internal/models/post.go',
            code: `package models

import (
	"time"

	"your/module/internal/app"
	"your/module/internal/database"
	"gorm.io/gorm"
)

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

func Posts() *database.Repository[Post] {
	return database.New[Post](app.DB)
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
            code: `// Find by primary key
post, err := models.Posts().Find(id)

// List all
posts, err := models.Posts().All()

// Create
err := models.Posts().Create(&post)

// Update
err := models.Posts().Update(&post)

// Soft-delete
err := models.Posts().Delete(id)

// Custom query — drop down to GORM when needed
var result []models.Post
err := models.Posts().DB().Where("published = ?", true).Find(&result).Error`
          }
        ]
      },
      {
        id: 'arch-controllers',
        title: 'Controllers',
        blocks: [
          {
            type: 'paragraph',
            text: 'Controllers are plain functions — no struct receivers, no dependency injection frameworks. fuego automatically generates OpenAPI 3.1 documentation from your handler signatures.'
          },
          {
            type: 'code',
            lang: 'go',
            label: 'internal/controllers/post-controller.go',
            code: `package controllers

import (
	"net/http"

	"your/module/internal/dto"
	"your/module/internal/models"
	"github.com/go-fuego/fuego"
)

// GetPost handles GET /posts/:post_id
// fuego infers the response type for OpenAPI automatically.
func GetPost(c fuego.ContextNoBody) (*dto.PostResponse, error) {
	id := c.PathParam("post_id")

	post, err := models.Posts().Find(id)
	if err != nil {
		return nil, fuego.HTTPError{
			Status: http.StatusNotFound,
			Err:    err,
		}
	}
	return toPostDTO(post), nil
}

// CreatePost handles POST /posts
func CreatePost(c fuego.ContextWithBody[dto.CreatePostRequest]) (*dto.PostResponse, error) {
	body, err := c.Body()
	if err != nil {
		return nil, fuego.HTTPError{Status: http.StatusBadRequest, Err: err}
	}

	post := &models.Post{
		Title:    body.Title,
		Content:  body.Content,
		AuthorID: body.AuthorID,
	}

	if err := models.Posts().Create(post); err != nil {
		return nil, fuego.HTTPError{Status: http.StatusUnprocessableEntity, Err: err}
	}
	return toPostDTO(post), nil
}

// toPostDTO maps a model to its response DTO.
func toPostDTO(m *models.Post) *dto.PostResponse {
	return &dto.PostResponse{
		ID:        m.ID,
		Title:     m.Title,
		Content:   m.Content,
		Published: m.Published,
		CreatedAt: m.CreatedAt.Format("2006-01-02T15:04:05Z"),
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
            text: 'DTOs (Data Transfer Objects) live in <code>internal/dto/</code> and define the exact shape of your API requests and responses. Keeping them separate from models means your API contract is independent of your database schema.'
          },
          {
            type: 'code',
            lang: 'go',
            label: 'internal/dto/post-dto.go',
            code: `package dto

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
            text: "All routes are registered in <code>internal/routes/routes.go</code>. fuego's route functions are fully typed — the compiler catches mismatched handler signatures before runtime."
          },
          {
            type: 'code',
            lang: 'go',
            label: 'internal/routes/routes.go',
            code: `package routes

import (
	"your/module/internal/controllers"
	"github.com/go-fuego/fuego"
)

func Register(s *fuego.Server) {
	// Posts
	fuego.Get(s,    "/posts",           controllers.ListPosts)
	fuego.Post(s,   "/posts",           controllers.CreatePost)
	fuego.Get(s,    "/posts/{post_id}", controllers.GetPost)
	fuego.Put(s,    "/posts/{post_id}", controllers.UpdatePost)
	fuego.Delete(s, "/posts/{post_id}", controllers.DeletePost)

	// Auth (example grouping)
	authGroup := fuego.Group(s, "/auth")
	fuego.Post(authGroup, "/login",    controllers.Login)
	fuego.Post(authGroup, "/register", controllers.Register)
	fuego.Post(authGroup, "/logout",   controllers.Logout)
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
#    e.g. add "PublishedAt *time.Time" to models/post.go

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
            text: '<code>atlas.hcl</code> sits at the project root and tells Atlas how to connect to your database and where to find your GORM models for diffing.'
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'atlas.hcl',
            code: `data "composite_schema" "app" {
  # Load your GORM models into Atlas via the Go loader
  schema "public" {
    url = "ent://internal/models"
  }
}

env "local" {
  src = data.composite_schema.app.url
  dev = "docker://postgres/15/dev?search_path=public"
  url = getenv("DATABASE_URL")
  migration {
    dir = "file://migrations"
  }
}

env "dev" {
  src = data.composite_schema.app.url
  url = getenv("DATABASE_URL")
  migration {
    dir = "file://migrations"
  }
}`
          },
          {
            type: 'note',
            kind: 'tip',
            text: 'The <code>local</code> environment uses a Docker-based shadow database for diffing. Make sure Docker is running when executing <code>grove make:migration</code> locally.'
          }
        ]
      },
      {
        id: 'app-singleton',
        title: 'App Singleton',
        blocks: [
          {
            type: 'paragraph',
            text: 'The <code>internal/app/</code> package holds shared singletons — the database connection, config, and logger. Everything is initialized once at startup and reused across the application.'
          },
          {
            type: 'code',
            lang: 'go',
            label: 'internal/app/app.go',
            code: `package app

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the shared database connection.
var DB *gorm.DB

// Init connects to the database and stores the handle in DB.
// Call this once from main() before starting the server.
func Init() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	DB = db
}`
          },
          {
            type: 'code',
            lang: 'go',
            label: 'cmd/api/main.go',
            code: `package main

import (
	"your/module/internal/app"
	"your/module/internal/routes"
	"github.com/go-fuego/fuego"
)

func main() {
	// Initialize shared singletons
	app.Init()

	// Create fuego server
	s := fuego.NewServer(
		fuego.WithAddr(":8080"),
	)

	// Register all routes
	routes.Register(s)

	// Start — fuego also serves /swagger automatically
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
