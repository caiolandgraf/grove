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
            code: `grove --help`
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

# 4. Start the development server
grove serve`
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
                title: 'Generate a migration',
                text: 'Run <code>grove make:migration create_posts_table</code> to diff your models against the current schema.'
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
                '<a href="https://github.com/caiolandgraf/gest" target="_blank">gest</a>',
                'latest',
                'Jest-inspired testing framework used by <code>grove test</code>'
              ],
              [
                '<a href="https://github.com/air-verse/air" target="_blank">air</a> (optional)',
                'latest',
                'Hot-reload via <code>grove serve</code>'
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
            text: 'Install <code>air</code> with <code>go install github.com/air-verse/air@latest</code> to enable hot-reload. Grove will detect it automatically.'
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
        id: 'cmd-serve',
        title: 'grove serve',
        blocks: [
          {
            type: 'paragraph',
            text: 'Starts the development HTTP server. If <code>air</code> is installed it will be used for hot-reload; otherwise falls back to <code>go run ./cmd/api/main.go</code>.'
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'terminal',
            code: `grove serve`
          },
          {
            type: 'note',
            kind: 'tip',
            text: 'SIGINT (<kbd>Ctrl+C</kbd>) is forwarded to the child process so it can shut down gracefully.'
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
            code: `grove make:model <Name> [-m] [-c] [-d] [-r]`
          },
          {
            type: 'table',
            head: ['Flag', 'Description'],
            rows: [
              [
                '<code>-m</code>, <code>--migration</code>',
                'Also run <code>atlas migrate diff</code> to generate a migration'
              ],
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
                'Full resource — shorthand for <code>-m -c -d</code> combined'
              ]
            ]
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'examples',
            code: `grove make:model Post            # model only
grove make:model Post -m         # model + migration
grove make:model Post -c         # model + controller
grove make:model Post -d         # model + DTO
grove make:model Post -mc        # model + migration + controller
grove make:model Post -mcd       # model + migration + controller + DTO
grove make:model Post -r         # full resource (same as -mcd)
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
            text: 'Generates a new SQL migration file by diffing your GORM models against the current database schema using Atlas. Make sure your models are up-to-date before running this.'
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
            kind: 'info',
            text: 'When you run <code>grove make:model Post -m</code> the migration is generated automatically with the name <code>create_posts_table</code> — you do not need to run this command separately.'
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
            text: 'Scaffolds a model, controller and DTO in one shot. Equivalent to running <code>make:model</code>, <code>make:controller</code> and <code>make:dto</code> together. Every file respects the <strong>SKIPPED</strong> rule — existing files are never overwritten.'
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
            code: `grove make:resource Post
grove make:resource BlogPost
grove make:resource order_item`
          },
          {
            type: 'note',
            kind: 'tip',
            text: 'This is the fastest way to bootstrap a new feature. After running it, add your fields to the model and DTO, run <code>grove make:migration</code>, then register your routes.'
          }
        ]
      },
      {
        id: 'cmd-migrate',
        title: 'grove migrate',
        blocks: [
          {
            type: 'paragraph',
            text: 'Applies all pending migrations to the database using <code>atlas migrate apply</code>.'
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'terminal',
            code: `grove migrate [--env <atlas-env>]`
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
            text: 'Grove uses <a href="https://github.com/caiolandgraf/gest" target="_blank">gest</a> as its testing framework — a Jest-inspired library for Go with beautiful colored output, descriptive failure messages and a fluent assertion API. Zero external dependencies, zero config files.'
          },
          {
            type: 'paragraph',
            text: 'All spec files live in <code>internal/tests/</code>. Each file self-registers its suite via <code>init()</code> and the single <code>main.go</code> entrypoint calls <code>gest.RunRegistered()</code>. Grove manages both files for you.'
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'project layout',
            code: `internal/
└── tests/
    ├── main.go            # gest entrypoint — auto-created by make:test
    ├── post_spec.go       # grove make:test Post
    └── user_service_spec.go`
          },
          {
            type: 'note',
            kind: 'info',
            text: 'gest uses <code>_spec.go</code> instead of <code>_test.go</code> because the Go toolchain reserves <code>_test.go</code> for <code>go test</code>. gest runs via <code>go run</code>, so any other suffix works fine.'
          }
        ]
      },
      {
        id: 'cmd-make-test',
        title: 'grove make:test',
        blocks: [
          {
            type: 'paragraph',
            text: 'Scaffolds a new gest spec file in <code>internal/tests/</code>. If <code>internal/tests/main.go</code> does not exist yet it is created automatically as the gest entrypoint — you never have to write it by hand. If the spec file already exists the command prints <strong>SKIPPED</strong> and exits cleanly.'
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
            text: 'The first time you run <code>make:test</code> grove creates the entrypoint and the spec:'
          },
          {
            type: 'code',
            lang: 'go',
            label: 'internal/tests/main.go (auto-created)',
            code: `package main

import (
	"os"

	"github.com/caiolandgraf/gest/gest"
)

func main() {
	if !gest.RunRegistered() {
		os.Exit(1)
	}
}`
          },
          {
            type: 'code',
            lang: 'go',
            label: 'internal/tests/user_spec.go (generated)',
            code: `package main

import "github.com/caiolandgraf/gest/gest"

func init() {
	s := gest.Describe("User")

	s.It("should work", func(t *gest.T) {
		// TODO: write your test here
		t.Expect(true).ToBeTrue()
	})

	gest.Register(s)
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
            text: 'Compiles and runs every <code>*_spec.go</code> file in <code>internal/tests/</code>. Pass <code>-c</code> to display a per-suite coverage report after the run.'
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'terminal',
            code: `grove test [-c]`
          },
          {
            type: 'table',
            head: ['Flag', 'Description'],
            rows: [
              [
                '<code>-c</code>, <code>--coverage</code>',
                'Display a per-suite pass-rate bar after the run'
              ]
            ]
          },
          {
            type: 'code',
            lang: 'bash',
            label: 'examples',
            code: `grove test         # run all specs
grove test -c      # run all specs + coverage report`
          },
          {
            type: 'note',
            kind: 'tip',
            text: 'Coverage bar colours: <strong>green</strong> ≥ 80% · <strong>yellow</strong> ≥ 50% · <strong>red</strong> &lt; 50%.'
          }
        ]
      },
      {
        id: 'testing-writing-specs',
        title: 'Writing Specs',
        blocks: [
          {
            type: 'paragraph',
            text: 'Each spec file calls <code>gest.Describe()</code> to create a suite, adds cases with <code>s.It()</code>, then registers with <code>gest.Register()</code>. The <code>init()</code> function ensures self-registration at startup.'
          },
          {
            type: 'code',
            lang: 'go',
            label: 'internal/tests/post_spec.go',
            code: `package main

import (
	"testing/quick"

	"github.com/caiolandgraf/gest/gest"
	"your/module/internal/models"
)

func init() {
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

	gest.Register(s)
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
│   ├── app/                     # Shared singletons (DB, config, logger)
│   │   └── app.go
│   ├── controllers/             # fuego route handlers
│   │   └── post-controller.go
│   ├── database/                # Generic repository + DB init
│   │   ├── database.go
│   │   └── repository.go
│   ├── dto/                     # Request and response types
│   │   └── post-dto.go
│   ├── middleware/              # HTTP middlewares
│   │   └── auth-middleware.go
│   ├── models/                  # GORM models
│   │   └── post.go
│   ├── routes/                  # Route registration
│   │   └── routes.go
│   └── tests/                   # gest spec files
│       ├── main.go              # gest entrypoint (auto-created by make:test)
│       └── post_spec.go
├── migrations/                  # Atlas SQL migrations
│   ├── 20240101120000_create_posts_table.sql
│   └── atlas.sum
├── .env                         # Local secrets (git-ignored)
├── .env.example                 # Committed env template
├── atlas.hcl                    # Atlas configuration
├── go.mod
└── go.sum`
          },
          {
            type: 'paragraph',
            text: 'The <code>internal/</code> package boundary is intentional — it prevents external packages from importing your application internals, keeping the codebase clean as it grows.'
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
