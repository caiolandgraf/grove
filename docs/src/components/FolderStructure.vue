<template>
  <div class="folder-explorer" :class="{ 'is-fullscreen': isMaximized }">
    <!-- Overlay background when maximized -->
    <div v-if="isMaximized" class="fullscreen-overlay" @click="isMaximized = false" />

    <div class="explorer-grid" :class="{ 'maximized-grid': isMaximized }">
      <!-- Left side: The file tree -->
      <div class="tree-panel">
        <div class="panel-header">
          <div class="header-top">
            <span class="panel-title">Project Structure</span>
            <div class="panel-actions">
              <!-- Expand tree action -->
              <button 
                class="btn-action" 
                @click="toggleAllNodes" 
                :title="allExpanded ? 'Collapse All' : 'Expand All'"
              >
                <svg v-if="allExpanded" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="4" y1="12" x2="20" y2="12"/>
                </svg>
                <svg v-else width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/>
                </svg>
              </button>
              <!-- Maximize action -->
              <button 
                class="btn-action" 
                @click="isMaximized = !isMaximized" 
                :title="isMaximized ? 'Minimize' : 'Maximize to Fullscreen'"
              >
                <svg v-if="isMaximized" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M4 14h6v6M20 10h-6V4M14 10l7-7M10 14l-7 7"/>
                </svg>
                <svg v-else width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M8 3H5a2 2 0 0 0-2 2v3M21 8V5a2 2 0 0 0-2-2h-3M3 16v3a2 2 0 0 0 2 2h3M16 21h3a2 2 0 0 0 2-2v-3"/>
                </svg>
              </button>
            </div>
          </div>
          <span class="panel-subtitle">Click items to inspect their purpose</span>
        </div>
        <div class="tree-container">
          <div
            v-for="(node, index) in visibleNodes"
            :key="node.path"
            class="tree-node"
            :class="{
              'is-dir': node.isDir,
              'is-active': selectedNode.path === node.path,
              'is-collapsed': node.isDir && collapsedPaths.includes(node.path)
            }"
            :style="{ paddingLeft: `${node.level * 16 + 12}px` }"
            @click="selectNode(node)"
          >
            <!-- Indent connectors -->
            <div
              v-for="i in node.level"
              :key="i"
              class="tree-indent"
              :style="{ left: `${(i - 1) * 16 + 20}px` }"
            />

            <!-- Toggle arrow (only for folders) -->
            <span
              v-if="node.isDir"
              class="toggle-arrow"
              @click.stop="toggleFolder(node.path)"
            >
              <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                <polyline points="9 18 15 12 9 6"/>
              </svg>
            </span>
            <span v-else class="toggle-spacer" />

            <!-- Icon -->
            <span class="node-icon">
              <!-- Folder open icon -->
              <svg v-if="node.isDir && !collapsedPaths.includes(node.path)" width="16" height="16" viewBox="0 0 24 24" fill="currentColor" class="icon-folder-open">
                <path d="M20 6h-8l-2-2H4c-1.1 0-1.99.9-1.99 2L2 18c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V8c0-1.1-.9-2-2-2zm0 12H4V8h16v10z"/>
              </svg>
              <!-- Folder closed icon -->
              <svg v-else-if="node.isDir" width="16" height="16" viewBox="0 0 24 24" fill="currentColor" class="icon-folder">
                <path d="M10 4H4c-1.1 0-1.99.9-1.99 2L2 18c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V8c0-1.1-.9-2-2-2h-8l-2-2z"/>
              </svg>
              <!-- Go file icon -->
              <svg v-else-if="node.name.endsWith('.go')" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="icon-file-go">
                <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                <polyline points="14 2 14 8 20 8"/>
                <path d="M8 13h8M8 17h6"/>
              </svg>
              <!-- Toml/Config icon -->
              <svg v-else-if="node.name.endsWith('.toml') || node.name.endsWith('.hcl') || node.name.endsWith('.yml')" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="icon-file-config">
                <path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/>
                <polyline points="14 2 14 8 20 8"/>
                <rect x="8" y="12" width="8" height="6" rx="1"/>
              </svg>
              <!-- General file icon -->
              <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="icon-file">
                <path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z"/>
                <polyline points="13 2 13 9 20 9"/>
              </svg>
            </span>

            <!-- Label -->
            <span class="node-label">{{ node.name }}</span>
          </div>
        </div>
      </div>

      <!-- Right side: Explanation box -->
      <div class="explanation-panel">
        <transition name="fade-slide" mode="out-in">
          <div :key="selectedNode.path" class="explanation-content card">
            <div class="explanation-header">
              <span class="path-badge">{{ selectedNode.path }}</span>
              <h3 class="node-title">{{ selectedNode.name.replace('/', '') }}</h3>
            </div>
            <div class="explanation-body">
              <p class="node-desc">{{ selectedNode.desc }}</p>

              <!-- Conditionally show structural helper info or tips -->
              <div v-if="selectedNode.tip" class="tip-box">
                <span class="tip-icon">💡</span>
                <p class="tip-text" v-html="selectedNode.tip"></p>
              </div>

              <!-- Extra code examples or stubs if clicking files -->
              <div v-if="selectedNode.code" class="code-preview">
                <span class="code-preview-title">Quick Glance:</span>
                <pre class="mono-code"><code v-html="highlightCode(selectedNode.code, selectedNode.name)" /></pre>
              </div>
            </div>
          </div>
        </transition>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import hljs from 'highlight.js/lib/core'
import go from 'highlight.js/lib/languages/go'
import yaml from 'highlight.js/lib/languages/yaml'
import ini from 'highlight.js/lib/languages/ini'
import markdown from 'highlight.js/lib/languages/markdown'

hljs.registerLanguage('go', go)
hljs.registerLanguage('yaml', yaml)
hljs.registerLanguage('toml', ini)
hljs.registerLanguage('ini', ini)
hljs.registerLanguage('markdown', markdown)

const treeNodes = [
  {
    name: 'cmd/',
    isDir: true,
    level: 0,
    path: 'cmd/',
    desc: 'Application entry points. Each subdirectory under cmd/ builds a separate executable binary. The primary service is cmd/api/; cmd/atlas/ loads GORM models for Atlas migrations.',
    tip: 'Add commands or other sub-applications (like CLI tasks, background workers) by creating folders inside <code>cmd/</code>.'
  },
  {
    name: 'api/',
    isDir: true,
    level: 1,
    path: 'cmd/api/',
    desc: 'The directory enclosing the main entry point for the REST API service.',
    tip: 'This holds your primary API build target. You usually run this target using <code>grove dev</code>.'
  },
  {
    name: 'main.go',
    isDir: false,
    level: 2,
    path: 'cmd/api/main.go',
    desc: 'The starting point of the REST API application. It boots config loaders, connects to the database, initiates session states, registers middlewares, attaches routes, and fires up Fuego.',
    code: `package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"your-app/cmd/scalar"
	"your-app/internal/app"
	"your-app/internal/app/config"
	"your-app/internal/routes"
	"github.com/go-fuego/fuego"
)

func main() {
	config.Load()
	config.InitLogger()

	ctx := context.Background()
	otelShutdown, _ := config.InitOtel(ctx)
	defer func() { _ = otelShutdown(ctx) }()

	_ = app.Boot()
	defer app.Shutdown()

	s := fuego.NewServer(
		fuego.WithAddr("localhost:8080"),
		fuego.WithEngineOptions(
			fuego.WithOpenAPIConfig(fuego.OpenAPIConfig{
				JSONFilePath: "./internal/app/docs/openapi.json",
				SwaggerURL:   "/api-doc",
				UIHandler:    scalar.NewUI,
			}),
		),
	)

	routes.SetupRoutes(s)
	_ = s.Run()
}`
  },
  {
    name: 'atlas/',
    isDir: true,
    level: 1,
    path: 'cmd/atlas/',
    desc: 'Atlas GORM schema loader invoked by atlas.hcl. Exports all registered models for migration diffing.',
    tip: 'Models self-register via <code>init()</code> in each module\'s <code>model.go</code>. This binary loads them all for Atlas.'
  },
  {
    name: 'main.go',
    isDir: false,
    level: 2,
    path: 'cmd/atlas/main.go',
    desc: 'Loads all GORM models registered in internal/app/database and outputs SQL schema for Atlas.',
    code: `package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	_ "github.com/caiolandgraf/go-project-base/internal/modules"
	"github.com/caiolandgraf/go-project-base/internal/app/database"
)

func main() {
	stmts, err := gormschema.New("postgres").Load(database.All()...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\\n", err)
		os.Exit(1)
	}
	_, _ = io.WriteString(os.Stdout, stmts)
}`
  },
  {
    name: 'scalar/',
    isDir: true,
    level: 1,
    path: 'cmd/scalar/',
    desc: 'Adapter package wrapping the Scalar UI layout to serve OpenAPI docs.',
  },
  {
    name: 'scalar.go',
    isDir: false,
    level: 2,
    path: 'cmd/scalar/scalar.go',
    desc: 'Registers UI Handlers for the modern Scalar OpenAPI documentation interface.',
    code: `package scalar

import (
	"net/http"
	"github.com/bdpiprava/scalar-go"
)

func NewUI(htmlRoute string, specRoute string) http.Handler {
	handler, _ := scalar.ApiReferenceHandler(&scalar.Options{
		SpecURL: specRoute,
		Theme:   scalar.ThemePurple,
	})
	return handler
}`
  },
  {
    name: 'internal/',
    isDir: true,
    level: 0,
    path: 'internal/',
    desc: 'Holds all private project source code. Go prevents external projects from importing files within the internal/ boundary.',
    tip: 'Shared infrastructure lives in <code>internal/app/</code>. Domain logic is organized as self-contained modules under <code>internal/modules/</code>.'
  },
  {
    name: 'app/',
    isDir: true,
    level: 1,
    path: 'internal/app/',
    desc: 'Application setup, globals, bootstrap wiring, config initializers, middlewares, and repository structures.',
  },
  {
    name: 'app.go',
    isDir: false,
    level: 2,
    path: 'internal/app/app.go',
    desc: 'Declares and exposes application global singletons like DB (GORM) and Session (SCS manager). Handles clean Boot and Shutdown routines.',
    code: `package app

import (
	"log/slog"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/caiolandgraf/go-project-base/internal/app/config"
	"github.com/gomodule/redigo/redis"
	"gorm.io/gorm"
)

var (
	DB      *gorm.DB
	Redis   *redis.Pool
	Session *scs.SessionManager
	Metrics http.Handler
)

func Boot() error {
	slog.Info("Booting application...")
	config.InitLogger()

	db, err := config.InitDatabase()
	if err != nil { return err }
	DB = db

	redisPool, err := config.InitRedis()
	if err != nil { return err }
	Redis = redisPool

	Session = config.InitSessionManager(redisPool)
	return nil
}

func Shutdown() {
	slog.Info("Shutting down application...")
	if DB != nil {
		if sqlDB, err := DB.DB(); err == nil {
			_ = sqlDB.Close()
		}
	}
	if Redis != nil { _ = Redis.Close() }
}`
  },
  {
    name: 'config/',
    isDir: true,
    level: 2,
    path: 'internal/app/config/',
    desc: 'Bootstrap initializers for environmental configurations and integrations.',
  },
  {
    name: 'database.go',
    isDir: false,
    level: 3,
    path: 'internal/app/config/database.go',
    desc: 'Initializes and connects to the PostgreSQL database using GORM.',
    code: `package config

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabase() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		Env.DBHost, Env.DBPort, Env.DBUser, Env.DBPassword, Env.DBName, Env.DBSSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return db, nil
}`
  },
  {
    name: 'env.go',
    isDir: false,
    level: 3,
    path: 'internal/app/config/env.go',
    desc: 'Loads and parses environment configurations from the .env file.',
    code: `package config

import (
	"log/slog"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	AppPort      int    \`env:"APP_PORT" envDefault:"8080"\`
	DBHost       string \`env:"DB_HOST" envDefault:"localhost"\`
	DBPort       int    \`env:"DB_PORT" envDefault:"5432"\`
	DBUser       string \`env:"DB_USER"\`
	DBPassword   string \`env:"DB_PASSWORD"\`
	DBName       string \`env:"DB_NAME"\`
	DBSSLMode    string \`env:"DB_SSLMODE" envDefault:"disable"\`
}

var Env Config

func Load() {
	_ = godotenv.Load()
	if err := env.Parse(&Env); err != nil {
		slog.Error("Failed to parse environment variables", "err", err)
	}
}`
  },
  {
    name: 'logger.go',
    isDir: false,
    level: 3,
    path: 'internal/app/config/logger.go',
    desc: 'Configures slog to log debug-friendly values locally or JSON outputs in production.',
    code: `package config

import (
	"log/slog"
	"os"
)

func InitLogger() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	slog.SetDefault(slog.New(handler))
}`
  },
  {
    name: 'metrics.go',
    isDir: false,
    level: 3,
    path: 'internal/app/config/metrics.go',
    desc: 'Initializes the Prometheus collector tracker.',
    code: `package config

import (
	"net/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func InitMetrics() (http.Handler, error) {
	return promhttp.Handler(), nil
}`
  },
  {
    name: 'otel.go',
    isDir: false,
    level: 3,
    path: 'internal/app/config/otel.go',
    desc: 'Configures OpenTelemetry tracer exporters.',
    code: `package config

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

func InitOtel(ctx context.Context) (func(context.Context) error, error) {
	tp := trace.NewTracerProvider()
	otel.SetTracerProvider(tp)
	return tp.Shutdown, nil
}`
  },
  {
    name: 'redis.go',
    isDir: false,
    level: 3,
    path: 'internal/app/config/redis.go',
    desc: 'Initializes connection client for Redis cache.',
    code: `package config

import (
	"github.com/gomodule/redigo/redis"
)

func InitRedis() (*redis.Pool, error) {
	return &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}, nil
}`
  },
  {
    name: 'session.go',
    isDir: false,
    level: 3,
    path: 'internal/app/config/session.go',
    desc: 'Configures alexedwards SCS session manager settings.',
    code: `package config

import (
	"time"
	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
)

func InitSessionManager(pool *redis.Pool) *scs.SessionManager {
	sm := scs.New()
	sm.Store = redisstore.New(pool)
	sm.Lifetime = 24 * time.Hour
	return sm
}`
  },
  {
    name: 'database/',
    isDir: true,
    level: 2,
    path: 'internal/app/database/',
    desc: 'Generic repository layer and Atlas model registry used by all domain modules.',
  },
  {
    name: 'registry.go',
    isDir: false,
    level: 3,
    path: 'internal/app/database/registry.go',
    desc: 'Central registry where module models register themselves via init() for Atlas schema loading.',
    code: `package database

import "gorm.io/schema"

var models []schema.Tabler

func Register(model schema.Tabler) {
	models = append(models, model)
}

func All() []schema.Tabler {
	return models
}`
  },
  {
    name: 'repository.go',
    isDir: false,
    level: 3,
    path: 'internal/app/database/repository.go',
    desc: 'A generic database CRUD wrapper over GORM Repository[T].',
    code: `package database

import (
	"errors"
	"gorm.io/gorm"
)

type Repository[T any] struct {
	db *gorm.DB
}

func New[T any](db *gorm.DB) *Repository[T] {
	return &Repository[T]{db: db}
}

func (r *Repository[T]) Find(id any) (*T, error) {
	var entity T
	if err := r.db.First(&entity, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *Repository[T]) Create(entity *T) error {
	return r.db.Create(entity).Error
}

func (r *Repository[T]) Update(entity *T) error {
	return r.db.Save(entity).Error
}

func (r *Repository[T]) Delete(id any) error {
	return r.db.Delete(new(T), "id = ?", id).Error
}`
  },
  {
    name: 'helpers/',
    isDir: true,
    level: 2,
    path: 'internal/app/helpers/',
    desc: 'Contains helper packages and utilities.',
  },
  {
    name: 'response-json.go',
    isDir: false,
    level: 3,
    path: 'internal/app/helpers/response-json.go',
    desc: 'Wraps output JSON payloads consistently.',
    code: `package helpers

type StandardResponse struct {
	Message string \`json:"message"\`
	Success bool   \`json:"success"\`
}`
  },
  {
    name: 'middleware/',
    isDir: true,
    level: 2,
    path: 'internal/app/middleware/',
    desc: 'Global route middlewares.',
  },
  {
    name: 'cors.go',
    isDir: false,
    level: 3,
    path: 'internal/app/middleware/cors.go',
    desc: 'Global CORS configuration.',
    code: `package middleware

import "net/http"

type CORSConfig struct {
	AllowedOrigins []string
}

func DefaultCORSConfig() CORSConfig {
	return CORSConfig{AllowedOrigins: []string{"*"}}
}

func CORSMiddleware(cfg CORSConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			next.ServeHTTP(w, r)
		})
	}
}`
  },
  {
    name: 'observability.go',
    isDir: false,
    level: 3,
    path: 'internal/app/middleware/observability.go',
    desc: 'Traces incoming requests.',
    code: `package middleware

import "net/http"

func RouteTagMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Adds custom metrics route tagging logic here
		next.ServeHTTP(w, r)
	})
}`
  },
  {
    name: 'session.go',
    isDir: false,
    level: 3,
    path: 'internal/app/middleware/session.go',
    desc: 'Extracts session tokens.',
    code: `package middleware

import (
	"net/http"
	"github.com/alexedwards/scs/v2"
)

func SessionMiddleware(sm *scs.SessionManager) func(http.Handler) http.Handler {
	return sm.LoadAndSave
}

func AuthRequired(sm *scs.SessionManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !sm.Exists(r.Context(), "user_id") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}`
  },
  {
    name: 'modules/',
    isDir: true,
    level: 1,
    path: 'internal/modules/',
    desc: 'Self-contained domain packages. Each module owns its model, DTO, service, controller, and OpenAPI docs.',
    tip: 'Run <code>grove make:resource Post</code> to scaffold a new module and auto-register it in <code>register.go</code>.'
  },
  {
    name: 'module.go',
    isDir: false,
    level: 2,
    path: 'internal/modules/module.go',
    desc: 'Defines the Module interface, Boot dependencies, and Factory type used by the registry.',
    code: `package modules

import (
	"github.com/alexedwards/scs/v2"
	"github.com/go-fuego/fuego"
	"gorm.io/gorm"
)

type Boot struct {
	DB      *gorm.DB
	Session *scs.SessionManager
}

type Module interface {
	Mount(api *fuego.Server, session *scs.SessionManager)
}

type Factory func(Boot) Module`
  },
  {
    name: 'register.go',
    isDir: false,
    level: 2,
    path: 'internal/modules/register.go',
    desc: 'Central module registry. New domains are wired here by grove make:resource.',
    code: `package modules

import (
	"github.com/caiolandgraf/go-project-base/internal/modules/auth"
	"github.com/caiolandgraf/go-project-base/internal/modules/users"
	"github.com/go-fuego/fuego"
)

var registry = []Factory{
	func(b Boot) Module { return users.Wire(b.DB) },
	func(b Boot) Module { return auth.Wire(b.DB, b.Session) },
}

func Mount(api *fuego.Server, boot Boot) {
	for _, factory := range registry {
		factory(boot).Mount(api, boot.Session)
	}
}`
  },
  {
    name: 'users/',
    isDir: true,
    level: 2,
    path: 'internal/modules/users/',
    desc: 'Users domain module — model, DTO, service, controller, and docs in one package.',
  },
  {
    name: 'model.go',
    isDir: false,
    level: 3,
    path: 'internal/modules/users/model.go',
    desc: 'User GORM model with init() registration for Atlas and typed repository.',
    code: `package users

import (
	"time"
	"github.com/caiolandgraf/go-project-base/internal/app/database"
	"gorm.io/gorm"
)

func init() {
	database.Register(&User{})
}

type User struct {
	ID        string         \`gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"\`
	Name      string         \`gorm:"type:varchar(255);not null"                     json:"name"\`
	Email     string         \`gorm:"type:varchar(255);uniqueIndex;not null"         json:"email"\`
	Password  string         \`gorm:"type:varchar(255);not null"                     json:"-"\`
	CreatedAt time.Time      \`gorm:"autoCreateTime"                                 json:"created_at"\`
	UpdatedAt time.Time      \`gorm:"autoUpdateTime"                                 json:"updated_at"\`
	DeletedAt gorm.DeletedAt \`gorm:"index"                                          json:"-"\`
}

type Repo struct {
	*database.Repository[User]
}

func Users(db *gorm.DB) *Repo {
	return &Repo{Repository: database.New[User](db)}
}`
  },
  {
    name: 'controller.go',
    isDir: false,
    level: 3,
    path: 'internal/modules/users/controller.go',
    desc: 'HTTP handlers and route mounting for the users domain.',
    code: `package users

import (
	"github.com/alexedwards/scs/v2"
	"github.com/caiolandgraf/go-project-base/internal/app/router"
	"github.com/go-fuego/fuego"
	"gorm.io/gorm"
)

type Controller struct {
	service Service
}

func Wire(db *gorm.DB) *Controller {
	return NewController(WireService(db))
}

func (ctrl *Controller) Mount(api *fuego.Server, session *scs.SessionManager) {
	group := fuego.Group(api, "/users")
	router.Get(group, "/", ctrl.ListUsers, ListUsersDoc, session)
}`
  },
  {
    name: 'auth/',
    isDir: true,
    level: 2,
    path: 'internal/modules/auth/',
    desc: 'Authentication module — login, session management, and auth routes.',
  },
  {
    name: 'routes/',
    isDir: true,
    level: 1,
    path: 'internal/routes/',
    desc: 'Global router configuration — middleware, health checks, and module mounting.',
  },
  {
    name: 'routes.go',
    isDir: false,
    level: 2,
    path: 'internal/routes/routes.go',
    desc: 'Applies global middleware and mounts all domain modules via modules.Mount().',
    code: `package routes

import (
	"net/http"
	"github.com/alexedwards/scs/v2"
	"github.com/caiolandgraf/go-project-base/internal/modules"
	"github.com/caiolandgraf/go-project-base/internal/app/middleware"
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
	fuego.Use(s, otelhttp.NewMiddleware("go-project-base"))
	fuego.Use(s, middleware.CORSMiddleware(middleware.DefaultCORSConfig()))
	fuego.Use(s, middleware.SessionMiddleware(session))

	api := fuego.Group(s, "/api/v1")
	modules.Mount(api, modules.Boot{DB: db, Session: session})
}`
  },
  {
    name: 'tests/',
    isDir: true,
    level: 1,
    path: 'internal/tests/',
    desc: 'Contains test suites executed by gest.',
  },
  {
    name: 'service_test.go',
    isDir: false,
    level: 2,
    path: 'internal/tests/service_test.go',
    desc: 'Service assertions using mocked repositories.',
    code: `package tests

import (
	"testing"
	"github.com/caiolandgraf/gest/v2/gest"
	"github.com/caiolandgraf/go-project-base/internal/modules/users"
)

func TestUserModel(t *testing.T) {
	s := gest.Describe("User model")
	s.It("should have valid fields", func(t *gest.T) {
		user := users.User{
			Name:  "John Doe",
			Email: "john@example.com",
		}
		t.Expect(user.Name).ToBe("John Doe")
		t.Expect(user.Email).ToBe("john@example.com")
	})
	s.Run(t)
}`
  },
  {
    name: 'migrations/',
    isDir: true,
    level: 0,
    path: 'migrations/',
    desc: 'Atlas-managed SQL migration files and integrity checksum.',
  },
  {
    name: '20260127143000_initial.sql',
    isDir: false,
    level: 1,
    path: 'migrations/20260127143000_initial.sql',
    desc: 'Default initial migration creating baseline tables (like users).',
    code: `-- create "users" table
CREATE TABLE "public"."users" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(255) NOT NULL,
  "email" character varying(255) NOT NULL,
  "password" character varying(255) NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
CREATE UNIQUE INDEX "idx_users_email" ON "public"."users" ("email");`
  },
  {
    name: 'atlas.sum',
    isDir: false,
    level: 1,
    path: 'migrations/atlas.sum',
    desc: 'Integrity checksum file maintaining migration history sequence correctness.',
    code: `h1:n9H/0c00bC8eD6eC9aC7fB3eA8eA9eA8eB5eC2eD1eD=
20260127143000_initial.sql h1:B8c7aF6eC5eD9aC3fB8eA2eD9eA8eC6eB5eA3eA1eD=`
  },
  {
    name: 'infra/',
    isDir: true,
    level: 0,
    path: 'infra/',
    desc: 'Observability stack configuration — Prometheus, Grafana, Loki, Jaeger, and Promtail.',
  },
  {
    name: 'prometheus.yml',
    isDir: false,
    level: 1,
    path: 'infra/prometheus.yml',
    desc: 'Prometheus scraper configuration.',
    code: `global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'grove-app'
    static_configs:
      - targets: ['host.docker.internal:8080']`
  },
  {
    name: 'loki-config.yml',
    isDir: false,
    level: 1,
    path: 'infra/loki-config.yml',
    desc: 'Grafana Loki configuration.',
    code: `auth_enabled: false
server:
  http_listen_port: 3100
common:
  path_prefix: /tmp/loki
  storage:
    filesystem:
      chunks_directory: /tmp/loki/chunks
      rules_directory: /tmp/loki/rules`
  },
  {
    name: 'promtail-config.yml',
    isDir: false,
    level: 1,
    path: 'infra/promtail-config.yml',
    desc: 'Promtail collector configuration.',
    code: `server:
  http_listen_port: 9080
clients:
  - url: http://loki:3100/loki/api/v1/push
scrape_configs:
  - job_name: local
    static_configs:
      - targets: [localhost]
        labels:
          job: varlogs
          __path__: /var/log/*log`
  },
  {
    name: 'grafana/',
    isDir: true,
    level: 1,
    path: 'infra/grafana/',
    desc: 'Grafana provisionings and telemetry dashboard panels.'
  },
  {
    name: 'docker-compose.yml',
    isDir: false,
    level: 0,
    path: 'docker-compose.yml',
    desc: 'Full infrastructure stack — PostgreSQL, Redis, Jaeger, Prometheus, Loki, and Grafana.',
    code: `services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_USER: \${DB_USER}
      POSTGRES_PASSWORD: \${DB_PASSWORD}
      POSTGRES_DB: \${DB_NAME}
    ports:
      - "\${DB_PORT}:5432"

  redis:
    image: redis:7
    ports:
      - "6379:6379"

  jaeger:
    image: jaegertracing/all-in-one:latest
    ports:
      - "16686:16686"`
  },
  {
    name: '.env.example',
    isDir: false,
    level: 0,
    path: '.env.example',
    desc: 'Committed templates for local variables.',
    code: `APP_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=grove_db
DB_SSLMODE=disable`
  },
  {
    name: 'atlas.hcl',
    isDir: false,
    level: 0,
    path: 'atlas.hcl',
    desc: 'Atlas GORM migration environment setups.',
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
  url = "postgres://\${getenv("DB_USER")}:\${getenv("DB_PASSWORD")}@\${getenv("DB_HOST")}:\${getenv("DB_PORT")}/\${getenv("DB_NAME")}?sslmode=disable"
  dev = "docker://postgres/15/dev"
  migration {
    dir = "file://migrations"
  }
}`
  },
  {
    name: 'grove.toml',
    isDir: false,
    level: 0,
    path: 'grove.toml',
    desc: 'CLI hot-reload settings configuration.',
    code: `[dev]
root        = "."
bin         = ".grove/tmp/app"
build_cmd   = "go build -o .grove/tmp/app ./cmd/api/"
watch_dirs  = ["cmd", "internal"]
exclude     = [".grove", "vendor", "node_modules", ".git", "infra", "migrations", "logs"]`
  },
  {
    name: 'go.mod',
    isDir: false,
    level: 0,
    path: 'go.mod',
    desc: 'Standard Go modules file.',
    code: `module github.com/caiolandgraf/go-project-base

go 1.25.7

require (
	github.com/alexedwards/scs/redisstore v0.0.0
	github.com/alexedwards/scs/v2 v2.9.0
	github.com/bdpiprava/scalar-go v0.13.0
	github.com/caiolandgraf/gest/v2 v2.0.3
	github.com/go-fuego/fuego v0.19.0
	github.com/gomodule/redigo v1.9.3
	gorm.io/driver/postgres v1.6.0
	gorm.io/gorm v1.31.1
)`
  },
  {
    name: 'README.md',
    isDir: false,
    level: 0,
    path: 'README.md',
    desc: 'General project readme.',
    code: `# Go Project Base
Modular Grove template with self-contained domain modules under internal/modules/.
Pre-configured with PostgreSQL, Redis, sessions, OpenTelemetry, and Atlas migrations.`
  }
]

const collapsedPaths = ref([])
const selectedNode = ref(treeNodes[2]) // default to main.go selected
const isMaximized = ref(false)

const visibleNodes = computed(() => {
  return treeNodes.filter(node => {
    // Check if any ancestor path is collapsed
    for (const collapsedPath of collapsedPaths.value) {
      if (node.path.startsWith(collapsedPath) && node.path !== collapsedPath) {
        return false
      }
    }
    return true
  })
})

function selectNode(node) {
  selectedNode.value = node
}

function toggleFolder(path) {
  const idx = collapsedPaths.value.indexOf(path)
  if (idx === -1) {
    collapsedPaths.value.push(path)
  } else {
    collapsedPaths.value.splice(idx, 1)
  }
}

const allExpanded = computed(() => collapsedPaths.value.length === 0)

function toggleAllNodes() {
  if (allExpanded.value) {
    // Collapse all directories in the tree
    const allDirs = treeNodes.filter(n => n.isDir).map(n => n.path)
    collapsedPaths.value = allDirs
  } else {
    // Expand all directories
    collapsedPaths.value = []
  }
}

function highlightCode(code, name) {
  let lang = 'go'
  if (name.endsWith('.go')) lang = 'go'
  else if (name.endsWith('.toml')) lang = 'toml'
  else if (name.endsWith('.yaml') || name.endsWith('.yml')) lang = 'yaml'
  else if (name.endsWith('.hcl')) lang = 'ini'
  else if (name.endsWith('.md')) lang = 'markdown'
  
  try {
    return hljs.highlight(code ?? '', { language: lang }).value
  } catch {
    return code ?? ''
  }
}
</script>

<style scoped>
.folder-explorer {
  margin: 1.5rem 0;
  width: 100%;
}

.explorer-grid {
  display: grid;
  grid-template-columns: 280px 1fr;
  gap: 1.25rem;
  align-items: start;
}

@media (max-width: 768px) {
  .explorer-grid {
    grid-template-columns: 1fr;
  }
}

/* Tree Panel */
.tree-panel {
  background: var(--bg-elevated);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  overflow: hidden;
}

.panel-header {
  padding: 0.85rem 1rem;
  background: rgba(255, 255, 255, 0.02);
  border-bottom: 1px solid var(--border);
}

.panel-title {
  display: block;
  font-size: 0.82rem;
  font-weight: 600;
  color: var(--text);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.panel-subtitle {
  font-size: 0.72rem;
  color: var(--text-muted);
}

.tree-container {
  padding: 0.75rem 0;
  max-height: 480px;
  overflow-y: auto;
}

/* Tree Nodes */
.tree-node {
  display: flex;
  align-items: center;
  position: relative;
  padding: 0.35rem 0.5rem;
  cursor: pointer;
  user-select: none;
  font-family: var(--font-mono);
  font-size: 0.8rem;
  color: rgba(226, 228, 239, 0.82);
  transition: background 0.15s, color 0.15s;
}

.tree-node:hover {
  background: var(--bg-hover);
  color: var(--text);
}

.tree-node.is-active {
  background: var(--red-dim);
  color: var(--red-hover);
  font-weight: 500;
}

.tree-indent {
  position: absolute;
  top: 0;
  bottom: 0;
  width: 1px;
  background: rgba(255, 255, 255, 0.03);
}

.toggle-arrow {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  margin-right: 4px;
  color: var(--text-muted);
  transition: transform 0.2s var(--ease);
  cursor: pointer;
  z-index: 2;
}

.toggle-arrow:hover {
  color: var(--text);
}

.tree-node.is-collapsed .toggle-arrow {
  transform: rotate(0deg);
}

.tree-node:not(.is-collapsed) .toggle-arrow {
  transform: rotate(90deg);
}

.toggle-spacer {
  width: 16px;
  margin-right: 4px;
}

.node-icon {
  display: inline-flex;
  align-items: center;
  margin-right: 6px;
  flex-shrink: 0;
}

.icon-folder { color: #f2a93b; }
.icon-folder-open { color: #e59728; }
.icon-file-go { color: #00add8; }
.icon-file-config { color: #8f92b2; }
.icon-file { color: var(--text-muted); }

.node-label {
  white-space: nowrap;
}

/* Explanation Panel */
.explanation-panel {
  min-height: 250px;
  min-width: 0; /* Prevents CSS grid column expansion from wide child code blocks */
}

.card {
  background: var(--bg-card);
  border: 1px solid var(--border-md);
  border-radius: var(--radius);
  padding: 1.25rem;
  box-shadow: 0 4px 30px rgba(0, 0, 0, 0.2);
  backdrop-filter: blur(5px);
  min-width: 0; /* Keeps card content confined to the grid cell width */
}

.explanation-header {
  border-bottom: 1px solid var(--border);
  padding-bottom: 0.75rem;
  margin-bottom: 0.85rem;
}

.path-badge {
  font-family: var(--font-mono);
  font-size: 0.7rem;
  background: rgba(255, 255, 255, 0.04);
  color: var(--text-muted);
  padding: 0.15rem 0.45rem;
  border-radius: 4px;
  border: 1px solid var(--border);
}

.node-title {
  font-size: 1.15rem;
  font-weight: 600;
  margin-top: 0.45rem;
  color: var(--text);
}

.node-desc {
  font-size: 0.88rem;
  color: rgba(226, 228, 239, 0.85);
  line-height: 1.6;
  margin-bottom: 1rem;
}

.tip-box {
  display: flex;
  gap: 0.65rem;
  background: rgba(200, 40, 56, 0.04);
  border-left: 2px solid var(--red);
  padding: 0.75rem;
  border-radius: 0 var(--radius-sm) var(--radius-sm) 0;
  margin-bottom: 1rem;
}

.tip-icon {
  font-size: 0.85rem;
  flex-shrink: 0;
}

.tip-text {
  font-size: 0.78rem;
  color: rgba(226, 228, 239, 0.72);
  line-height: 1.5;
  margin: 0;
}

.tip-text :deep(code) {
  font-family: var(--font-mono);
  background: rgba(255, 255, 255, 0.06);
  padding: 0.05rem 0.25rem;
  border-radius: 3px;
  color: var(--text);
}

.code-preview {
  margin-top: 1rem;
}

.code-preview-title {
  display: block;
  font-size: 0.7rem;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-bottom: 0.4rem;
}

.mono-code {
  background: #06060a;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  padding: 0.75rem;
  margin: 0;
  overflow-x: auto;
  max-width: 100%; /* Prevents code blocks from breaking flex/grid limits */
}

.mono-code code {
  font-family: var(--font-mono);
  font-size: 0.76rem;
  color: #a9b2d3;
  line-height: 1.5;
  display: block;
}

/* Animations */
.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: opacity 0.15s ease, transform 0.15s ease;
}

.fade-slide-enter-from {
  opacity: 0;
  transform: translateY(6px);
}

.fade-slide-leave-to {
  opacity: 0;
  transform: translateY(-6px);
}

/* Maximized & Fullscreen Styles */
.folder-explorer.is-fullscreen .explorer-grid.maximized-grid {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: min(1200px, 92vw);
  height: min(720px, 86vh);
  background: #18181c;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  box-shadow: 0 50px 100px rgba(0, 0, 0, 0.9);
  z-index: 9999;
  display: grid;
  grid-template-columns: 320px 1fr;
  grid-template-rows: 1fr;
  overflow: hidden;
  gap: 0;
}

.folder-explorer.is-fullscreen .tree-panel {
  border: none;
  border-right: 1px solid rgba(255, 255, 255, 0.08);
  height: 100%;
  display: flex;
  flex-direction: column;
  border-radius: 12px 0 0 12px;
}

.folder-explorer.is-fullscreen .tree-container {
  flex: 1;
  max-height: none;
  overflow-y: auto;
}

.folder-explorer.is-fullscreen .explanation-panel {
  height: 100%;
  overflow-y: auto;
  border-radius: 0 12px 12px 0;
  box-sizing: border-box;
}

.folder-explorer.is-fullscreen .card {
  height: 100%;
  border: none;
  border-radius: 0;
  box-shadow: none;
  overflow-y: auto;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
}

.folder-explorer.is-fullscreen .explanation-body {
  flex: 1;
  overflow-y: auto;
}

.fullscreen-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background: rgba(0, 0, 0, 0.7);
  backdrop-filter: blur(8px);
  z-index: 9998;
}

.header-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2px;
}

.panel-actions {
  display: flex;
  gap: 6px;
}

.btn-action {
  background: none;
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: var(--text-muted);
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s;
}

.btn-action:hover {
  color: var(--text);
  background: rgba(255, 255, 255, 0.05);
  border-color: rgba(255, 255, 255, 0.2);
}

/* Custom syntax highlighting colors for code in FolderStructure */
.mono-code code :deep(.hljs-keyword)   { color: #569cd6; }
.mono-code code :deep(.hljs-built_in)  { color: #4ec9b0; }
.mono-code code :deep(.hljs-type)      { color: #4ec9b0; }
.mono-code code :deep(.hljs-string)    { color: #ce9178; }
.mono-code code :deep(.hljs-number)    { color: #b5cea8; }
.mono-code code :deep(.hljs-comment)   { color: #6a9955; font-style: italic; }
.mono-code code :deep(.hljs-title)     { color: #dcdcaa; }
.mono-code code :deep(.hljs-params)    { color: #9cdcfe; }
.mono-code code :deep(.hljs-attr)      { color: #9cdcfe; }
.mono-code code :deep(.hljs-literal)   { color: #569cd6; }
</style>
