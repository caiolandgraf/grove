<template>
  <div class="tutorial-layout">
    
    <!-- ─── Left Pane: Tutorial Guide ──────────── -->
    <aside class="tutorial-guide-pane">
      <div class="guide-header">
        <router-link to="/docs" class="back-link">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
            <line x1="19" y1="12" x2="5" y2="12"/><polyline points="12 19 5 12 12 5"/>
          </svg>
          Back to Docs
        </router-link>
        <span class="guide-brand">Grove CLI</span>
      </div>

      <div class="guide-scroll-content">
        <!-- Progress section -->
        <div class="guide-progress-section">
          <div class="progress-info">
            <span class="progress-title">First Project Guide</span>
            <span class="progress-percent">{{ currentStep + 1 }} / {{ steps.length }}</span>
          </div>
          <div class="progress-bar-wrap">
            <div class="progress-bar-fill" :style="{ width: ((currentStep + 1) / steps.length * 100) + '%' }" />
          </div>
        </div>

        <!-- Steps Pipeline -->
        <div class="steps-pipeline">
          <div
            v-for="(step, i) in steps"
            :key="i"
            class="pipeline-node"
            :class="{ 
              'is-active': currentStep === i, 
              'is-done': currentStep > i, 
              'is-locked': i > currentStep && !stepActionsCompleted[i - 1] 
            }"
            @click="goToStep(i)"
          >
            <div class="pipeline-circle">
              <svg v-if="currentStep > i" width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                <polyline points="20 6 9 17 4 12"/>
              </svg>
              <span v-else>{{ i + 1 }}</span>
            </div>
            <span class="pipeline-label">{{ step.shortTitle }}</span>
          </div>
        </div>

        <!-- Step Detail Card -->
        <div class="step-card">
          <div class="step-card-badge">Step {{ currentStep + 1 }} of {{ steps.length }}</div>
          <h2 class="step-card-title">{{ steps[currentStep].title }}</h2>
          <p class="step-card-desc" v-html="steps[currentStep].description" />

          <!-- Action execution triggers -->
          <div class="step-card-actions">
            <button
              v-if="steps[currentStep].actionText && !stepActionsCompleted[currentStep]"
              class="btn-trigger-action"
              @click="performStepAction"
              :disabled="typingCommand"
            >
              <svg width="12" height="12" viewBox="0 0 24 24" fill="currentColor">
                <polygon points="5 3 19 12 5 21 5 3"/>
              </svg>
              {{ steps[currentStep].actionText }}
            </button>

            <div v-if="stepActionsCompleted[currentStep]" class="action-success-badge">
              <div class="success-icon-wrap">
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                  <polyline points="20 6 9 17 4 12"/>
                </svg>
              </div>
              <div class="success-text-wrap">
                <strong>Completed!</strong> {{ steps[currentStep].successMessage }}
              </div>
            </div>
          </div>

          <!-- HTTP Sandbox (Step 6 only, after grove up is run) -->
          <div v-if="currentStep === 5 && stepActionsCompleted[5]" class="http-sandbox">
            <h3 class="sandbox-title">API Endpoint Tester</h3>
            <p class="sandbox-subtitle">Your API is running in the sandbox. Send mock HTTP requests below.</p>

            <div class="http-controls">
              <select v-model="httpMethod" class="http-method-select">
                <option value="GET">GET</option>
                <option value="POST">POST</option>
              </select>
              <input type="text" readonly value="/api/v1/posts" class="http-url-input" />
              <button class="btn-send-http" @click="sendMockHttpRequest" :disabled="sendingHttp">
                {{ sendingHttp ? 'Sending...' : 'Send' }}
              </button>
            </div>

            <!-- POST body input -->
            <div v-if="httpMethod === 'POST'" class="http-post-body">
              <label class="post-body-label">JSON Request Body:</label>
              <textarea v-model="httpRequestBody" class="post-body-textarea" rows="4"></textarea>
            </div>

            <!-- HTTP Response output -->
            <div v-if="httpResponse" class="http-response-panel">
              <div class="panel-tab">RESPONSE STATUS: {{ httpResponse.status }}</div>
              <pre class="response-body"><code v-html="highlightCode(httpResponse.body, 'json')" /></pre>
            </div>
          </div>

        </div>
      </div>

      <!-- Navigation footer -->
      <footer class="guide-nav-footer">
        <button class="btn-step-prev" :disabled="currentStep === 0" @click="prevStep">
          ← Back
        </button>
        <button
          class="btn-step-next"
          :disabled="currentStep === steps.length - 1 || !stepActionsCompleted[currentStep]"
          @click="nextStep"
        >
          Next Step →
        </button>
      </footer>
    </aside>

    <!-- ─── Right Pane: VS Code Window ─────────── -->
    <main class="vscode-pane">
      <div class="vscode-window">
        
        <!-- Menu/Title Bar -->
        <div class="vscode-menubar">
          <div class="menubar-traffic">
            <router-link to="/docs" class="traffic-dot dot-red" title="Exit to Docs" />
            <span class="traffic-dot dot-yellow" />
            <span class="traffic-dot dot-green" />
          </div>
          <div class="menubar-title">
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="3" y="3" width="18" height="18" rx="2"/><path d="M9 9l3 3-3 3M13 15h3"/>
            </svg>
            blog-api — VS Code
          </div>
          <router-link to="/docs" class="menubar-close" title="Exit to Docs">
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </router-link>
        </div>

        <div class="vscode-body">
          <!-- Activity Bar -->
          <div class="vscode-activity-bar">
            <button
              class="activity-icon"
              :class="{ active: sidebarView === 'explorer' }"
              @click="sidebarView = sidebarView === 'explorer' ? null : 'explorer'"
              title="Explorer"
            >
              <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
                <rect x="3" y="3" width="18" height="18" rx="2"/><line x1="9" y1="3" x2="9" y2="21"/>
              </svg>
            </button>
          </div>

          <!-- Primary Sidebar (Explorer) -->
          <div v-if="sidebarView === 'explorer'" class="vscode-sidebar">
            <div class="sidebar-header">EXPLORER: BLOG-API</div>
            <div class="file-tree">
              <div
                v-for="item in explorerTree"
                :key="item.path"
                class="tree-item"
                :class="{ 
                  folder: item.isFolder, 
                  'is-highlighted': item.highlight, 
                  'is-active': activeEditorTab === item.tabId 
                }"
                :style="{ paddingLeft: (item.depth * 12 + 10) + 'px' }"
                @click="item.tabId && !item.isFolder && (activeEditorTab = item.tabId)"
              >
                <span class="tree-icon">
                  <!-- Folder icon -->
                  <svg v-if="item.isFolder" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/>
                  </svg>
                  <!-- File Go icon -->
                  <svg v-else-if="item.lang === 'go'" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="#82aaff" stroke-width="2">
                    <circle cx="12" cy="12" r="9"/><path d="M9 12h6M12 9v6"/>
                  </svg>
                  <!-- File SQL icon -->
                  <svg v-else-if="item.lang === 'sql'" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="#f78c6c" stroke-width="2">
                    <ellipse cx="12" cy="5" rx="9" ry="3"/><path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5"/>
                  </svg>
                  <!-- General File icon -->
                  <svg v-else width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="#89ddff" stroke-width="2">
                    <path d="M14 2H6a2 2 0 0 0-2 2v16h16V8z"/>
                  </svg>
                </span>
                <span class="tree-name" :class="item.lang">{{ item.name }}</span>
                <span v-if="item.highlight" class="tree-badge">new</span>
              </div>
            </div>
          </div>

          <!-- Editor / Content Workspace -->
          <div class="editor-area">
            <!-- Tabs Row -->
            <div class="vscode-tabs">
              <div
                v-for="tab in openTabs"
                :key="tab.id"
                class="editor-tab"
                :class="{ 'is-active': activeEditorTab === tab.id }"
                @click="activeEditorTab = tab.id"
              >
                <span class="tab-icon" :class="tab.lang">
                  <svg v-if="tab.lang === 'go'" width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="9"/><path d="M8 12h8M12 8v8"/></svg>
                  <svg v-else-if="tab.lang === 'sql'" width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><ellipse cx="12" cy="5" rx="9" ry="3"/><path d="M21 12c0 1.66-4 3-9 3s-9-1.34-9-3"/><path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5"/></svg>
                  <svg v-else width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="4 17 10 11 4 5"/><line x1="12" y1="19" x2="20" y2="19"/></svg>
                </span>
                {{ tab.label }}
              </div>
            </div>

            <!-- Code Editor Pane -->
            <div v-if="activeEditorTab !== 'terminal'" class="editor-pane">
              <div class="editor-breadcrumb">
                <span class="bc-seg">blog-api</span>
                <span class="bc-sep">›</span>
                <span class="bc-seg">{{ currentTabMeta.path }}</span>
                <span class="bc-sep">›</span>
                <span class="bc-seg bc-file">{{ currentTabMeta.label }}</span>
              </div>

              <div class="editor-scroll">
                <div class="editor-gutter">
                  <span v-for="n in codeLineCount" :key="n" class="line-num">{{ n }}</span>
                </div>
                <pre class="editor-code-body"><code v-html="highlightCode(currentFileCode, currentTabMeta.lang)" /></pre>
              </div>
            </div>

            <!-- Terminal Pane -->
            <div v-else class="terminal-pane" ref="terminalRef">
              <div class="terminal-topbar">
                <span class="tbar-label">
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                    <polyline points="4 17 10 11 4 5"/><line x1="12" y1="19" x2="20" y2="19"/>
                  </svg>
                  zsh (terminal)
                </span>
                <button class="tbar-close" @click="activeEditorTab = lastEditorTab">✕</button>
              </div>
              <div class="terminal-body">
                <div v-for="(log, i) in terminalLogs" :key="i" class="t-line" :class="log.type" v-html="log.text" />
                <div v-if="typingCommand" class="t-line t-cmd">
                  <span class="t-prompt">❯</span> <span class="t-typing">{{ typedText }}</span><span class="t-cursor">▋</span>
                </div>
                <div v-else-if="stepActionsCompleted[currentStep]" class="t-line t-cmd">
                  <span class="t-prompt">❯</span>
                </div>
              </div>
            </div>

            <!-- Bottom tab selector bar -->
            <div class="editor-bottom-bar">
              <div class="bottom-tabs">
                <button
                  class="bottom-tab"
                  :class="{ active: activeEditorTab === 'terminal' }"
                  @click="toggleTerminal"
                >
                  <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                    <polyline points="4 17 10 11 4 5"/><line x1="12" y1="19" x2="20" y2="19"/>
                  </svg>
                  TERMINAL
                </button>
              </div>
            </div>

          </div>
        </div>

        <!-- Status Bar -->
        <div class="vscode-statusbar">
          <div class="status-left">
            <span class="status-item status-branch">
              <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                <line x1="6" y1="3" x2="6" y2="15"/><circle cx="18" cy="6" r="3"/><circle cx="6" cy="18" r="3"/><path d="M18 9a9 9 0 0 1-9 9"/>
              </svg>
              main
            </span>
            <span class="status-item">
              <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                <circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/>
              </svg>
              {{ activeEditorTab === 'terminal' ? 'bash' : (currentTabMeta?.lang ?? 'go') }}
            </span>
          </div>
          <div class="status-right">
            <span class="status-item status-step">
              Step {{ currentStep + 1 }}/{{ steps.length }} — {{ steps[currentStep].shortTitle }}
            </span>
            <span class="status-item status-grove">grove v2.1.0</span>
          </div>
        </div>

      </div>
    </main>

  </div>
</template>

<script setup>
import { ref, computed, watch, nextTick } from 'vue'
import hljs from 'highlight.js/lib/core'
import go from 'highlight.js/lib/languages/go'
import bash from 'highlight.js/lib/languages/bash'
import json from 'highlight.js/lib/languages/json'
import sql from 'highlight.js/lib/languages/sql'

hljs.registerLanguage('go', go)
hljs.registerLanguage('bash', bash)
hljs.registerLanguage('json', json)
hljs.registerLanguage('sql', sql)

// ── State ──────────────────────────────────────
const currentStep         = ref(0)
const sidebarView         = ref('explorer')
const activeEditorTab     = ref('main.go')
const lastEditorTab       = ref('main.go')
const typingCommand       = ref(false)
const typedText           = ref('')
const terminalLogs        = ref([])
const terminalRef         = ref(null)
const stepActionsCompleted = ref([false, false, false, false, false, false])

// HTTP Tester Sandbox State
const httpMethod          = ref('GET')
const httpRequestBody     = ref('{\n  "title": "Grove framework is awesome",\n  "content": "Building APIs in Go is now simpler than ever with Models, Controllers and Services.",\n  "published": true\n}')
const sendingHttp         = ref(false)
const httpResponse        = ref(null)

const mockPostsList       = ref([
  {
    id: "777c07b6-1721-4d92-8096-7d1a1b415a77",
    title: "Introducing Grove CLI",
    content: "Grove scaffolds modular REST APIs inside the Go ecology.",
    published: true,
    created_at: "2026-06-10T16:00:00Z",
    updated_at: "2026-06-10T16:00:00Z"
  }
])

// ── Steps definitions ──────────────────────────
const steps = [
  {
    shortTitle: 'Scaffold',
    title: '1. Create a New Project',
    description: `Run <code>grove setup</code> to scaffold a clean Grove codebase. It downloads the template from <code>caiolandgraf/grove-base</code>, configures environment vars, and installs Go modules.`,
    actionText: 'grove setup blog-api',
    successMessage: 'Project scaffolded successfully!'
  },
  {
    shortTitle: 'Resource',
    title: '2. Generate a Resource Layer',
    description: `Grove follows a strict **Model-Controller-Service (MCS)** pattern. Running <code>grove make:resource Post</code> generates the database model, DTO, service, controller, and routes wireup in one sweep.`,
    actionText: 'grove make:resource Post',
    successMessage: 'Model, Service, Controller & routes created!'
  },
  {
    shortTitle: 'Fields',
    title: '3. Add Fields to Model & DTO',
    description: `Open <code>post.go</code> and <code>post-dto.go</code> in the editor. We will add <code>Title</code>, <code>Content</code>, and <code>Published</code> fields along GORM metadata and validation tags.`,
    actionText: 'Apply custom fields to Post structs',
    successMessage: 'Fields & validation tags applied.'
  },
  {
    shortTitle: 'Migration',
    title: '4. Generate SQL Migration',
    description: `Grove integrates **Atlas** for schema migrations. Run <code>grove make:migration create_posts_table</code> to compare the new model struct with the local PostgreSQL database and generate an SQL diff.`,
    actionText: 'grove make:migration create_posts_table',
    successMessage: 'Migration file generated!'
  },
  {
    shortTitle: 'Migrate',
    title: '5. Apply Schema to Database',
    description: `Run <code>grove migrate</code> to push the generated Atlas SQL files into the database. Grove handles the migration safely inside a transaction schema.`,
    actionText: 'grove migrate',
    successMessage: 'Schema applied to PostgreSQL!'
  },
  {
    shortTitle: 'Run API',
    title: '6. Start Server & Test API',
    description: `Run <code>grove up</code> to launch dependencies (Postgres, Jaeger, Redis) via Docker and spin up Fuego's hot-reload server. The API and Swagger doc are served at <code>localhost:8080</code>.`,
    actionText: 'grove up',
    successMessage: 'API running at http://localhost:8080!'
  }
]

// ── Simulated file contents (computed based on steps) ──
const fileContents = computed(() => {
  const customized = currentStep.value >= 2
  return {
    'main.go': {
      lang: 'go',
      path: 'cmd/api',
      label: 'main.go',
      code: `package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"blog-api/cmd/scalar"
	"blog-api/internal/app"
	"blog-api/internal/app/config"
	"blog-api/internal/routes"
	"github.com/go-fuego/fuego"
)

func main() {
	// Load .env and parse all environment variables into config.Env
	config.Load()

	// Initialize structured logger (uses config.Env.LogLevel)
	config.InitLogger()

	// Initialize OpenTelemetry
	ctx := context.Background()
	otelShutdown, err := config.InitOtel(ctx)
	if err != nil {
		slog.Error("Failed to initialize OpenTelemetry", "error", err)
		os.Exit(1)
	}
	defer func() { _ = otelShutdown(ctx) }()

	// Boot application (DB, Redis, Session, Metrics)
	if err := app.Boot(); err != nil {
		slog.Error("Failed to boot application", "error", err)
		os.Exit(1)
	}
	defer app.Shutdown()

	// Initialize server
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

	// Configure routes
	routes.SetupRoutes(s)

	slog.Info("Server starting", "addr", ":8080")

	// Graceful shutdown
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint
		slog.Info("Shutting down server...")
		app.Shutdown()
		os.Exit(0)
	}()

	// Start server
	if err := s.Run(); err != nil {
		slog.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}`
    },
    'app.go': {
      lang: 'go',
      path: 'internal/app',
      label: 'app.go',
      code: `package app

import (
	"log/slog"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"blog-api/internal/app/config"
	"github.com/gomodule/redigo/redis"
	"gorm.io/gorm"
)

// Global application dependencies.
// Accessible from anywhere via app.DB, app.Redis, app.Session.
var (
	DB      *gorm.DB
	Redis   *redis.Pool
	Session *scs.SessionManager
	Metrics http.Handler
)

// Boot initializes all infrastructure dependencies.
// Call this once at application startup.
func Boot() error {
	slog.Info("Booting application...")

	// Logger
	config.InitLogger()

	// Database
	db, err := config.InitDatabase()
	if err != nil {
		return err
	}
	DB = db

	// Redis
	redisPool, err := config.InitRedis()
	if err != nil {
		return err
	}
	Redis = redisPool

	// Session Manager
	Session = config.InitSessionManager(redisPool)

	// Metrics
	if config.Env.MetricsEnabled {
		metricsHandler, err := config.InitMetrics()
		if err != nil {
			return err
		}
		Metrics = metricsHandler
	} else {
		slog.Info("Prometheus metrics disabled")
	}

	slog.Info("Application booted successfully")
	return nil
}

// Shutdown gracefully closes all connections.
func Shutdown() {
	slog.Info("Shutting down application...")

	if DB != nil {
		if sqlDB, err := DB.DB(); err == nil {
			_ = sqlDB.Close()
		}
	}

	if Redis != nil {
		_ = Redis.Close()
	}

	slog.Info("All connections closed")
}

// HealthCheck returns the status of all infrastructure dependencies.
func HealthCheck() map[string]string {
	status := make(map[string]string)

	// Database
	if DB != nil {
		if sqlDB, err := DB.DB(); err == nil {
			if err := sqlDB.Ping(); err == nil {
				status["database"] = "healthy"
			} else {
				status["database"] = "unhealthy"
			}
		}
	} else {
		status["database"] = "not initialized"
	}

	// Redis
	if Redis != nil {
		conn := Redis.Get()
		defer func() {
			_ = conn.Close()
		}()
		if _, err := conn.Do("PING"); err == nil {
			status["redis"] = "healthy"
		} else {
			status["redis"] = "unhealthy"
		}
	} else {
		status["redis"] = "not initialized"
	}

	return status
}`
    },
    'routes.go': {
      lang: 'go',
      path: 'internal/routes',
      label: 'routes.go',
      code: currentStep.value >= 1
        ? `package routes

import (
	"blog-api/internal/app"
	"blog-api/internal/app/config"
	"blog-api/internal/app/middleware"
	"blog-api/internal/controllers"
	"blog-api/internal/models"
	"blog-api/internal/services"
	"github.com/go-fuego/fuego"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type HealthCheckResponse struct {
	Status string \`json:"status"\`
}

type HealthCheckDetailedResponse struct {
	Status   string            \`json:"status"\`
	Services map[string]string \`json:"services"\`
}

// SetupRoutes configures all routes using app globals
func SetupRoutes(s *fuego.Server) {
	// OpenTelemetry / Metrics Middleware
	if config.Env.OtelEnabled || config.Env.MetricsEnabled {
		fuego.Use(s, otelhttp.NewMiddleware(config.Env.OtelServiceName))
		fuego.Use(s, middleware.RouteTagMiddleware)
	}

	// CORS Middleware global
	fuego.Use(s, middleware.CORSMiddleware(middleware.DefaultCORSConfig()))

	// Session Middleware global
	fuego.Use(s, middleware.SessionMiddleware(app.Session))

	// Health check
	fuego.Get(s, "/", healthCheck)
	fuego.Get(s, "/health", healthCheckDetailed)

	// API v1
	api := fuego.Group(s, "/api/v1")

	// Services instantiation
	userRepo := models.Users()
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo)

	// ========== AUTH ROUTES ==========
	auth := fuego.Group(api, "/auth")
	authController := controllers.NewAuthController(app.Session, authService, userService)
	fuego.Post(auth, "/login", authController.Login)
	fuego.Get(auth, "/me", authController.Me)

	// ========== USERS ROUTES ==========
	users := fuego.Group(api, "/users")
	userController := controllers.NewUserController(app.Session, userService)
	fuego.Use(users, middleware.AuthRequired(app.Session))
	fuego.Get(users, "/", userController.List)

	// ── Post routes (generated by grove make:resource Post) ─
	postRepo := models.Posts()
	postService := services.NewPostService(postRepo)
	postController := controllers.NewPostController(app.Session, postService)

	postsGroup := fuego.Group(s, "/api/v1/posts")
	fuego.Get(postsGroup, "/", postController.List)
	fuego.Post(postsGroup, "/", postController.Create)
	fuego.Get(postsGroup, "/{post_id}", postController.Get)
	fuego.Put(postsGroup, "/{post_id}", postController.Update)
	fuego.Delete(postsGroup, "/{post_id}", postController.Delete)
}

func healthCheck(c fuego.ContextNoBody) (*HealthCheckResponse, error) {
	return &HealthCheckResponse{Status: "Ok"}, nil
}

func healthCheckDetailed(c fuego.ContextNoBody) (*HealthCheckDetailedResponse, error) {
	return &HealthCheckDetailedResponse{
		Status:   "OK",
		Services: app.HealthCheck(),
	}, nil
}`
        : `package routes

import (
	"blog-api/internal/app"
	"blog-api/internal/app/config"
	"blog-api/internal/app/middleware"
	"blog-api/internal/controllers"
	"blog-api/internal/models"
	"blog-api/internal/services"
	"github.com/go-fuego/fuego"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type HealthCheckResponse struct {
	Status string \`json:"status"\`
}

type HealthCheckDetailedResponse struct {
	Status   string            \`json:"status"\`
	Services map[string]string \`json:"services"\`
}

// SetupRoutes configures all routes using app globals
func SetupRoutes(s *fuego.Server) {
	// OpenTelemetry / Metrics Middleware
	if config.Env.OtelEnabled || config.Env.MetricsEnabled {
		fuego.Use(s, otelhttp.NewMiddleware(config.Env.OtelServiceName))
		fuego.Use(s, middleware.RouteTagMiddleware)
	}

	// CORS Middleware global
	fuego.Use(s, middleware.CORSMiddleware(middleware.DefaultCORSConfig()))

	// Session Middleware global
	fuego.Use(s, middleware.SessionMiddleware(app.Session))

	// Health check
	fuego.Get(s, "/", healthCheck)
	fuego.Get(s, "/health", healthCheckDetailed)

	// API v1
	api := fuego.Group(s, "/api/v1")

	// Services instantiation
	userRepo := models.Users()
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo)

	// ========== AUTH ROUTES ==========
	auth := fuego.Group(api, "/auth")
	authController := controllers.NewAuthController(app.Session, authService, userService)
	fuego.Post(auth, "/login", authController.Login)
	fuego.Get(auth, "/me", authController.Me)

	// ========== USERS ROUTES ==========
	users := fuego.Group(api, "/users")
	userController := controllers.NewUserController(app.Session, userService)
	fuego.Use(users, middleware.AuthRequired(app.Session))
	fuego.Get(users, "/", userController.List)
}

func healthCheck(c fuego.ContextNoBody) (*HealthCheckResponse, error) {
	return &HealthCheckResponse{Status: "Ok"}, nil
}

func healthCheckDetailed(c fuego.ContextNoBody) (*HealthCheckDetailedResponse, error) {
	return &HealthCheckDetailedResponse{
		Status:   "OK",
		Services: app.HealthCheck(),
	}, nil
}`
    },
    'user.go': {
      lang: 'go',
      path: 'internal/models',
      label: 'user.go',
      code: `package models

import (
	"time"

	"blog-api/internal/app"
	"blog-api/internal/app/database"
	"gorm.io/gorm"
)

type User struct {
	ID        string         \`gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"\`
	Name      string         \`gorm:"type:varchar(255);not null"                     json:"name"\`
	Email     string         \`gorm:"type:varchar(255);uniqueIndex;not null"         json:"email"\`
	Password  string         \`gorm:"type:varchar(255);not null"                     json:"-"\`
	CreatedAt time.Time      \`gorm:"autoCreateTime"                                 json:"created_at"\`
	UpdatedAt time.Time      \`gorm:"autoUpdateTime"                                 json:"updated_at"\`
	DeletedAt gorm.DeletedAt \`gorm:"index"                                          json:"-"\`
}

func (User) TableName() string { return "users" }

type usersRepository struct {
	*database.Repository[User]
}

func Users() *usersRepository {
	return &usersRepository{
		Repository: database.New[User](app.DB),
	}
}

func (r *usersRepository) FindUserByEmail(email string) (*User, error) {
	return Users().FirstBy("email", email)
}

func (r *usersRepository) Exists(email string) (bool, error) {
	return r.Where("email", email).Exists()
}`
    },
    'post.go': {
      lang: 'go',
      path: 'internal/models',
      label: 'post.go',
      code: customized
        ? `package models

import (
	"time"

	"blog-api/internal/app"
	"blog-api/internal/app/database"
	"gorm.io/gorm"
)

type Post struct {
	ID        string         \`gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"\`
	Title     string         \`gorm:"type:varchar(255);not null"                     json:"title"\`
	Content   string         \`gorm:"type:text"                                      json:"content"\`
	Published bool           \`gorm:"default:false"                                  json:"published"\`
	CreatedAt time.Time      \`gorm:"autoCreateTime"                                 json:"created_at"\`
	UpdatedAt time.Time      \`gorm:"autoUpdateTime"                                 json:"updated_at"\`
	DeletedAt gorm.DeletedAt \`gorm:"index"                                          json:"-"\`
}

func (Post) TableName() string { return "posts" }

type postsRepository struct {
	*database.Repository[Post]
}

func Posts() *postsRepository {
	return &postsRepository{
		Repository: database.New[Post](app.DB),
	}
}`
        : `package models

import (
	"time"

	"blog-api/internal/app"
	"blog-api/internal/app/database"
	"gorm.io/gorm"
)

// Post — generated by grove make:resource Post.
// Add your domain fields below.
type Post struct {
	ID        string         \`gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"\`
	CreatedAt time.Time      \`gorm:"autoCreateTime"                                 json:"created_at"\`
	UpdatedAt time.Time      \`gorm:"autoUpdateTime"                                 json:"updated_at"\`
	DeletedAt gorm.DeletedAt \`gorm:"index"                                          json:"-"\`
	// TODO: add your domain fields here
}

func (Post) TableName() string { return "posts" }

type postsRepository struct {
	*database.Repository[Post]
}

func Posts() *postsRepository {
	return &postsRepository{
		Repository: database.New[Post](app.DB),
	}
}`
    },
    'post-dto.go': {
      lang: 'go',
      path: 'internal/dto',
      label: 'post-dto.go',
      code: customized
        ? `package dto

type CreatePostRequest struct {
	Title     string \`json:"title"     validate:"required,min=3,max=255"\`
	Content   string \`json:"content"   validate:"required"\`
	Published bool   \`json:"published"\`
}

type UpdatePostRequest struct {
	Title     *string \`json:"title"     validate:"omitempty,min=3"\`
	Content   *string \`json:"content"   validate:"omitempty"\`
	Published *bool   \`json:"published" validate:"omitempty"\`
}

type PostResponse struct {
	ID        string \`json:"id"\`
	Title     string \`json:"title"\`
	Content   string \`json:"content"\`
	Published bool   \`json:"published"\`
}

type PostsListResponse struct {
	Items []PostResponse \`json:"items"\`
	Total int            \`json:"total"\`
}`
        : `package dto

// CreatePostRequest — generated stub.
// Add validated fields matching your model.
type CreatePostRequest struct {
	// TODO: add fields with validate tags
	// Example: Title string \`json:"title" validate:"required,min=3"\`
}

type UpdatePostRequest struct {
	// TODO: add optional pointer fields
}

type PostResponse struct {
	ID string \`json:"id"\`
}

type PostsListResponse struct {
	Items []PostResponse \`json:"items"\`
	Total int            \`json:"total"\`
}`
    },
    'post_service.go': {
      lang: 'go',
      path: 'internal/services',
      label: 'post_service.go',
      code: `package services

import (
	"context"
	"blog-api/internal/models"
)

// PostRepository defines the data access contract for Post.
// Implemented by postsRepository in internal/models/post.go.
type PostRepository interface {
	Create(entity *models.Post) error
	Update(entity *models.Post) error
	Delete(id any) error
	Find(id any) (*models.Post, error)
	Paginate(page, perPage int) ([]models.Post, int64, error)
}

type PostService struct {
	repo PostRepository
}

func NewPostService(repo PostRepository) *PostService {
	return &PostService{repo: repo}
}

// TODO: Adicionar métodos de negócio aqui (ex: CreatePost, GetPost, etc.)`
    },
    'post-controller.go': {
      lang: 'go',
      path: 'internal/controllers',
      label: 'post-controller.go',
      code: `package controllers

import (
	"github.com/alexedwards/scs/v2"
	"blog-api/internal/services"
	"github.com/go-fuego/fuego"
)

type PostController struct {
	session     *scs.SessionManager
	postService *services.PostService
}

func NewPostController(
	session *scs.SessionManager,
	postService *services.PostService,
) *PostController {
	return &PostController{
		session:     session,
		postService: postService,
	}
}

// TODO: Adicionar os handlers HTTP delegando a lógica para ctl.postService`
    },
    'migration.sql': {
      lang: 'sql',
      path: 'internal/app/database/migrations',
      label: '20260610_create_posts_table.sql',
      code: `-- Atlas generated migration
-- grove make:migration create_posts_table
-- Diff: models.Post → current DB schema

-- create "posts" table
CREATE TABLE "public"."posts" (
  "id"         uuid        NOT NULL DEFAULT gen_random_uuid(),
  "title"      character varying(255) NOT NULL,
  "content"    text,
  "published"  boolean     NOT NULL DEFAULT false,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);

-- Index for GORM soft-deletes
CREATE INDEX "idx_posts_deleted_at"
  ON "public"."posts" ("deleted_at");`
    },
    'terminal': {
      lang: 'bash',
      path: '',
      label: 'terminal'
    }
  }
})

// ── Explorer tree (changes per step) ──────────
const explorerTree = computed(() => {
  const s = currentStep.value
  const items = [
    { name: 'blog-api/', isFolder: true, depth: 0, path: '' },
    { name: 'cmd/', isFolder: true, depth: 1, path: 'cmd' },
    { name: 'api/', isFolder: true, depth: 2, path: 'cmd/api' },
    { name: 'main.go', depth: 3, lang: 'go', tabId: 'main.go', path: 'cmd/api/main.go', highlight: s === 0 },
    { name: 'scalar/', isFolder: true, depth: 2, path: 'cmd/scalar' },
    { name: 'seed/', isFolder: true, depth: 2, path: 'cmd/seed' },
    { name: 'internal/', isFolder: true, depth: 1, path: 'internal' },
    { name: 'app/', isFolder: true, depth: 2, path: 'internal/app' },
    { name: 'app.go', depth: 3, lang: 'go', tabId: 'app.go', path: 'internal/app/app.go', highlight: s === 0 },
    { name: 'config/', isFolder: true, depth: 3, path: 'internal/app/config' },
    { name: 'database/', isFolder: true, depth: 3, path: 'internal/app/database' },
    { name: 'migrations/', isFolder: true, depth: 4, path: 'internal/app/database/migrations' },
    ...(s >= 3 ? [{ name: '20260610_create_posts_table.sql', depth: 5, lang: 'sql', tabId: 'migration.sql', path: 'internal/app/database/migrations/20260610_create_posts_table.sql', highlight: s === 3 }] : []),
    { name: 'middleware/', isFolder: true, depth: 3, path: 'internal/app/middleware' },
    { name: 'controllers/', isFolder: true, depth: 2, path: 'internal/controllers' },
    { name: 'auth-controller.go', depth: 3, lang: 'go', path: 'internal/controllers/auth-controller.go' },
    { name: 'user-controller.go', depth: 3, lang: 'go', path: 'internal/controllers/user-controller.go' },
    ...(s >= 1 ? [{ name: 'post-controller.go', depth: 3, lang: 'go', tabId: 'post-controller.go', path: 'internal/controllers/post-controller.go', highlight: s === 1 }] : []),
    { name: 'dto/', isFolder: true, depth: 2, path: 'internal/dto' },
    { name: 'auth-dto.go', depth: 3, lang: 'go', path: 'internal/dto/auth-dto.go' },
    { name: 'user-dto.go', depth: 3, lang: 'go', path: 'internal/dto/user-dto.go' },
    ...(s >= 1 ? [{ name: 'post-dto.go', depth: 3, lang: 'go', tabId: 'post-dto.go', path: 'internal/dto/post-dto.go', highlight: s === 1 || s === 2 }] : []),
    { name: 'models/', isFolder: true, depth: 2, path: 'internal/models' },
    { name: 'user.go', depth: 3, lang: 'go', tabId: 'user.go', path: 'internal/models/user.go' },
    ...(s >= 1 ? [{ name: 'post.go', depth: 3, lang: 'go', tabId: 'post.go', path: 'internal/models/post.go', highlight: s === 1 || s === 2 }] : []),
    { name: 'routes/', isFolder: true, depth: 2, path: 'internal/routes' },
    { name: 'routes.go', depth: 3, lang: 'go', tabId: 'routes.go', path: 'internal/routes/routes.go', highlight: s === 1 },
    { name: 'services/', isFolder: true, depth: 2, path: 'internal/services' },
    { name: 'auth_service.go', depth: 3, lang: 'go', path: 'internal/services/auth_service.go' },
    { name: 'user_service.go', depth: 3, lang: 'go', path: 'internal/services/user_service.go' },
    ...(s >= 1 ? [{ name: 'post_service.go', depth: 3, lang: 'go', tabId: 'post_service.go', path: 'internal/services/post_service.go', highlight: s === 1 }] : []),
    { name: 'tests/', isFolder: true, depth: 2, path: 'internal/tests' },
    { name: 'infra/', isFolder: true, depth: 1, path: 'infra' },
    { name: 'docker-compose.yml', depth: 2, lang: 'yaml', path: 'infra/docker-compose.yml' },
    { name: 'atlas.hcl', depth: 1, lang: 'hcl', path: 'atlas.hcl' },
    { name: 'grove.toml', depth: 1, lang: 'toml', path: 'grove.toml' },
    { name: 'go.mod', depth: 1, lang: 'go', path: 'go.mod' },
    { name: '.env.example', depth: 1, lang: 'text', path: '.env.example' },
  ]
  return items
})

// ── Open tabs (change per step) ───────────────
const openTabs = computed(() => {
  const s = currentStep.value
  const tabs = [
    { id: 'main.go', label: 'main.go', lang: 'go' },
    { id: 'app.go', label: 'app.go', lang: 'go' },
  ]
  if (s >= 1) {
    tabs.push({ id: 'routes.go', label: 'routes.go', lang: 'go' })
    tabs.push({ id: 'post.go', label: 'post.go', lang: 'go' })
    tabs.push({ id: 'post-dto.go', label: 'post-dto.go', lang: 'go' })
    tabs.push({ id: 'post_service.go', label: 'post_service.go', lang: 'go' })
    tabs.push({ id: 'post-controller.go', label: 'post-controller.go', lang: 'go' })
  }
  if (s >= 3) {
    tabs.push({ id: 'migration.sql', label: '...sql', lang: 'sql' })
  }
  return tabs
})

const currentTabMeta = computed(() => {
  const id = activeEditorTab.value === 'terminal' ? lastEditorTab.value : activeEditorTab.value
  return fileContents.value[id] ?? fileContents.value['main.go']
})

const currentFileCode = computed(() => currentTabMeta.value?.code ?? '')

const codeLineCount = computed(() => {
  return (currentFileCode.value.match(/\n/g) || []).length + 1
})

// ── Highlight ─────────────────────────────────
function highlightCode(code, lang) {
  try {
    return hljs.highlight(code ?? '', { language: lang ?? 'go' }).value
  } catch {
    return code ?? ''
  }
}

// ── Step navigation ───────────────────────────
function goToStep(idx) {
  if (idx <= currentStep.value || (idx > 0 && stepActionsCompleted.value[idx - 1]) || idx === 0) {
    currentStep.value = idx
  }
}

function prevStep() {
  if (currentStep.value > 0) currentStep.value--
}

function nextStep() {
  if (currentStep.value < steps.length - 1) {
    currentStep.value++
  }
}

// ── Terminal toggle ───────────────────────────
function toggleTerminal() {
  if (activeEditorTab.value === 'terminal') {
    activeEditorTab.value = lastEditorTab.value
  } else {
    lastEditorTab.value = activeEditorTab.value
    activeEditorTab.value = 'terminal'
  }
}

// ── Watch step → update active file & auto-open terminal ─────
watch(currentStep, (newStep) => {
  httpResponse.value = null // clear sandbox response
  if (newStep === 0) { activeEditorTab.value = 'main.go'; lastEditorTab.value = 'main.go' }
  else if (newStep === 1) { activeEditorTab.value = 'routes.go'; lastEditorTab.value = 'routes.go' }
  else if (newStep === 2) { activeEditorTab.value = 'post.go'; lastEditorTab.value = 'post.go' }
  else if (newStep === 3) { activeEditorTab.value = 'migration.sql'; lastEditorTab.value = 'migration.sql' }
  else if (newStep === 4) { activeEditorTab.value = 'terminal'; }
  else if (newStep === 5) { activeEditorTab.value = 'main.go'; lastEditorTab.value = 'main.go' }
})

// ── Simulated command execution ───────────────
function performStepAction() {
  if (typingCommand.value) return
  typingCommand.value = true
  lastEditorTab.value = activeEditorTab.value !== 'terminal' ? activeEditorTab.value : lastEditorTab.value
  activeEditorTab.value = 'terminal'

  const commands = [
    'grove setup blog-api',
    'grove make:resource Post',
    '# Editing files inside editor...',
    'grove make:migration create_posts_table',
    'grove migrate',
    'grove up'
  ]
  const cmd = commands[currentStep.value]
  let i = 0
  typedText.value = ''

  const typeTimer = setInterval(() => {
    if (i < cmd.length) {
      typedText.value += cmd[i++]
      scrollTerminal()
    } else {
      clearInterval(typeTimer)
      setTimeout(runCommand, 250)
    }
  }, 30)
}

function runCommand() {
  typingCommand.value = false
  const step = currentStep.value
  const push = (type, text) => terminalLogs.value.push({ type, text })

  if (step === 0) {
    push('cmd', `<span class="t-prompt">❯</span> grove setup blog-api`)
    push('logo', `  <span class="t-red">█▀▀ █▀█ █▀█ █░█ █▀▀</span>`)
    push('logo', `  <span class="t-red">█▄█ █▀▄ █▄█ ▀▄▀ ██▄</span>`)
    push('dim',  `  ────────────────────────────────────`)
    push('log',  `  <b>Project</b>    blog-api`)
    push('log',  `  <b>Module</b>     blog-api`)
    push('dim',  `  Template   caiolandgraf/grove-base`)
    push('dim',  `  ────────────────────────────────────`)
    push('ok',   `  <span class="t-green">✓</span> Downloading template     <span class="t-dim">425 KB</span>`)
    push('ok',   `  <span class="t-green">✓</span> Extracting files          <span class="t-dim">42 files</span>`)
    push('ok',   `  <span class="t-green">✓</span> Configuring observability <span class="t-dim">Jaeger · Prometheus · Grafana</span>`)
    push('ok',   `  <span class="t-green">✓</span> Setting Go module name    <span class="t-dim">blog-api</span>`)
    push('ok',   `  <span class="t-green">✓</span> Installing dependencies   <span class="t-dim">1.4s</span>`)
    push('done', `  <span class="t-badge">SUCCESS</span> Project ready in <b>./blog-api</b>`)
  } else if (step === 1) {
    push('cmd',  `<span class="t-prompt">❯</span> grove make:resource Post`)
    push('dim',  `  Creating resource <b>Post</b>...`)
    push('ok',   `  <span class="t-badge">CREATED</span> Model       → <span class="t-dim">internal/models/post.go</span>`)
    push('ok',   `  <span class="t-badge">CREATED</span> DTO         → <span class="t-dim">internal/dto/post-dto.go</span>`)
    push('ok',   `  <span class="t-badge">CREATED</span> Service     → <span class="t-dim">internal/services/post_service.go</span>`)
    push('ok',   `  <span class="t-badge">CREATED</span> Controller  → <span class="t-dim">internal/controllers/post-controller.go</span>`)
    push('ok',   `  <span class="t-badge">WIRED </span> Routes      → <span class="t-dim">internal/routes/routes.go</span>`)
  } else if (step === 2) {
    push('cmd',  `<span class="t-dim"># Editing model & DTO stubs...</span>`)
    push('ok',   `  <span class="t-green">✓</span> Added Title, Content, Published fields to Post GORM Model`)
    push('ok',   `  <span class="t-green">✓</span> Set up validation rules in CreatePostRequest DTO`)
    activeEditorTab.value = 'post.go'; lastEditorTab.value = 'post.go'
  } else if (step === 3) {
    push('cmd',  `<span class="t-prompt">❯</span> grove make:migration create_posts_table`)
    push('log',  `  Running Atlas schema diff...`)
    push('ok',   `  <span class="t-badge">CREATED</span> migrations/20260610_create_posts_table.sql`)
    activeEditorTab.value = 'migration.sql'; lastEditorTab.value = 'migration.sql'
  } else if (step === 4) {
    push('cmd',  `<span class="t-prompt">❯</span> grove migrate`)
    push('log',  `  Running atlas migrate apply --env local`)
    push('log',  ``)
    push('ok',   `  <span class="t-badge">MIGRATE</span> 20260610_create_posts_table`)
    push('dim',  `     CREATE TABLE  "posts" (...)`)
    push('dim',  `     CREATE INDEX  "idx_posts_deleted_at"`)
    push('log',  `  <span class="t-green">OK</span>  18.2ms · 2 SQL statements applied`)
  } else if (step === 5) {
    push('cmd',  `<span class="t-prompt">❯</span> grove up`)
    push('log',  `  Starting Docker Compose stack...`)
    push('ok',   `  <span class="t-green">✓</span> Container <b>postgres-db</b>  Started`)
    push('ok',   `  <span class="t-green">✓</span> Container <b>jaeger-otel</b>  Started`)
    push('log',  `  Booting API server (hot-reload)...`)
    push('dim',  `  2026-06-10 16:45:00 <span class="t-green">INF</span> Booting application...`)
    push('dim',  `  2026-06-10 16:45:00 <span class="t-green">INF</span> OpenTelemetry initialized service=blog-api`)
    push('dim',  `  2026-06-10 16:45:00 <span class="t-green">INF</span> Connected to database name=blog_api_dev`)
    push('done', `  <span class="t-badge">RUNNING</span> API at <b>http://localhost:8080</b>  ·  Swagger at <b>/swagger</b>`)
  }

  stepActionsCompleted.value[step] = true
  scrollTerminal()
}

function scrollTerminal() {
  nextTick(() => {
    if (terminalRef.value) terminalRef.value.scrollTop = terminalRef.value.scrollHeight
  })
}

// ── HTTP Sandbox Actions ──────────────────────
function sendMockHttpRequest() {
  if (sendingHttp.value) return
  sendingHttp.value = true
  
  setTimeout(() => {
    sendingHttp.value = false
    if (httpMethod.value === 'GET') {
      httpResponse.value = {
        status: '200 OK',
        body: JSON.stringify({
          items: mockPostsList.value,
          total: mockPostsList.value.length
        }, null, 2)
      }
    } else {
      // POST
      try {
        const bodyObj = JSON.parse(httpRequestBody.value)
        const newPost = {
          id: "3e5a31a9-bf0c-439f-9eb0-" + Math.floor(Math.random() * 1000000000).toString(16),
          title: bodyObj.title || "Untitled Post",
          content: bodyObj.content || "",
          published: bodyObj.published ?? false,
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString()
        }
        mockPostsList.value.push(newPost)
        httpResponse.value = {
          status: '201 Created',
          body: JSON.stringify(newPost, null, 2)
        }
      } catch (err) {
        httpResponse.value = {
          status: '400 Bad Request',
          body: JSON.stringify({
            error: "Invalid JSON format: " + err.message
          }, null, 2)
        }
      }
    }
  }, 600)
}
</script>

<style scoped>
/* ═══════════════════════════════════════════════
   Redesigned Split Layout (Tutorial Left, VSCode Right)
   ═══════════════════════════════════════════════ */
.tutorial-layout {
  display: flex;
  width: 100vw;
  height: 100vh;
  background: #111113;
  color: #e3e3e6;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
  overflow: hidden;
}

/* ─── Left Pane: Tutorial ───────────────────── */
.tutorial-guide-pane {
  width: 400px;
  min-width: 350px;
  max-width: 480px;
  background: #16161a;
  border-right: 1px solid rgba(255, 255, 255, 0.08);
  display: flex;
  flex-direction: column;
  height: 100%;
  flex-shrink: 0;
  box-shadow: 10px 0 30px rgba(0, 0, 0, 0.25);
  z-index: 10;
}

.guide-header {
  height: 48px;
  padding: 0 16px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  background: rgba(0, 0, 0, 0.15);
}

.back-link {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #a1a1aa;
  text-decoration: none;
  font-size: 0.78rem;
  font-weight: 600;
  transition: color 0.15s;
}
.back-link:hover {
  color: #38bdf8;
}

.guide-brand {
  font-size: 0.72rem;
  letter-spacing: 0.15em;
  text-transform: uppercase;
  font-weight: 700;
  color: #38bdf8;
}

.guide-scroll-content {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* Progress bar styles */
.guide-progress-section {
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.progress-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.progress-title {
  font-size: 0.75rem;
  font-weight: 600;
  color: #71717a;
}
.progress-percent {
  font-size: 0.75rem;
  font-weight: 700;
  color: #38bdf8;
}
.progress-bar-wrap {
  height: 5px;
  background: rgba(255, 255, 255, 0.06);
  border-radius: 99px;
  overflow: hidden;
}
.progress-bar-fill {
  height: 100%;
  background: #38bdf8;
  border-radius: 99px;
  transition: width 0.3s ease;
}

/* Steps Pipeline list */
.steps-pipeline {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  background: rgba(0, 0, 0, 0.15);
  padding: 10px;
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.04);
}
.pipeline-node {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 8px;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
  opacity: 0.45;
}
.pipeline-node:hover {
  background: rgba(255, 255, 255, 0.04);
  opacity: 0.8;
}
.pipeline-node.is-active {
  background: rgba(56, 189, 248, 0.1);
  border: 1px solid rgba(56, 189, 248, 0.2);
  opacity: 1;
}
.pipeline-node.is-done {
  opacity: 0.9;
}
.pipeline-node.is-locked {
  opacity: 0.2;
  cursor: not-allowed;
}
.pipeline-circle {
  width: 16px;
  height: 16px;
  border-radius: 50%;
  border: 1.5px solid #52525b;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.6rem;
  font-weight: 700;
  color: #a1a1aa;
}
.pipeline-node.is-active .pipeline-circle {
  border-color: #38bdf8;
  color: #38bdf8;
}
.pipeline-node.is-done .pipeline-circle {
  border-color: #4ade80;
  color: #4ade80;
  background: rgba(74, 222, 128, 0.1);
}
.pipeline-label {
  font-size: 0.72rem;
  font-weight: 500;
}

/* Step Card details */
.step-card {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 12px;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.step-card-badge {
  font-size: 0.62rem;
  text-transform: uppercase;
  letter-spacing: 0.12em;
  font-weight: 700;
  color: #38bdf8;
}
.step-card-title {
  font-size: 1.05rem;
  font-weight: 700;
  color: #fff;
  line-height: 1.25;
  margin: 0;
}
.step-card-desc {
  font-size: 0.8rem;
  color: #a1a1aa;
  line-height: 1.55;
  margin: 0;
}
.step-card-desc :deep(code) {
  background: rgba(56, 189, 248, 0.1);
  color: #38bdf8;
  padding: 0.15em 0.35em;
  border-radius: 4px;
  font-family: var(--font-mono, 'JetBrains Mono', monospace);
  font-size: 0.9em;
}
.step-card-desc :deep(strong) {
  color: #fff;
  font-weight: 600;
}

.step-card-actions {
  margin-top: 6px;
}
.btn-trigger-action {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  width: 100%;
  padding: 10px 14px;
  background: #38bdf8;
  color: #0c0a09;
  border: none;
  border-radius: 6px;
  font-size: 0.8rem;
  font-weight: 700;
  cursor: pointer;
  transition: background 0.15s, transform 0.1s;
}
.btn-trigger-action:hover:not(:disabled) {
  background: #0ea5e9;
}
.btn-trigger-action:active:not(:disabled) {
  transform: translateY(1px);
}
.btn-trigger-action:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.action-success-badge {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 12px 14px;
  background: rgba(74, 222, 128, 0.08);
  border: 1px solid rgba(74, 222, 128, 0.2);
  border-radius: 6px;
  color: #4ade80;
}
.success-icon-wrap {
  margin-top: 2px;
  flex-shrink: 0;
}
.success-text-wrap {
  font-size: 0.78rem;
  line-height: 1.4;
}

/* HTTP Sandbox Tester style */
.http-sandbox {
  margin-top: 15px;
  padding-top: 15px;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
}
.sandbox-title {
  font-size: 0.88rem;
  font-weight: 700;
  color: #fff;
  margin: 0 0 4px;
}
.sandbox-subtitle {
  font-size: 0.72rem;
  color: #71717a;
  margin: 0 0 12px;
}
.http-controls {
  display: flex;
  gap: 6px;
  margin-bottom: 10px;
}
.http-method-select {
  background: #27272a;
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: #fff;
  border-radius: 4px;
  font-size: 0.75rem;
  font-weight: 700;
  padding: 0 6px;
  cursor: pointer;
  outline: none;
}
.http-url-input {
  flex: 1;
  background: #27272a;
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: #a1a1aa;
  border-radius: 4px;
  font-size: 0.75rem;
  padding: 6px 10px;
  font-family: var(--font-mono, monospace);
  outline: none;
}
.btn-send-http {
  background: #4ade80;
  color: #052e16;
  border: none;
  font-weight: 700;
  border-radius: 4px;
  font-size: 0.75rem;
  padding: 0 14px;
  cursor: pointer;
  transition: background 0.1s;
}
.btn-send-http:hover:not(:disabled) {
  background: #22c55e;
}
.btn-send-http:disabled {
  opacity: 0.6;
}
.http-post-body {
  display: flex;
  flex-direction: column;
  gap: 4px;
  margin-bottom: 10px;
}
.post-body-label {
  font-size: 0.68rem;
  color: #a1a1aa;
  font-weight: 600;
}
.post-body-textarea {
  background: #1e1e24;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 4px;
  color: #fff;
  font-family: var(--font-mono, monospace);
  font-size: 0.72rem;
  padding: 8px;
  resize: vertical;
  outline: none;
}
.http-response-panel {
  background: #18181b;
  border-radius: 6px;
  border: 1px solid rgba(255, 255, 255, 0.06);
  overflow: hidden;
}
.panel-tab {
  background: #27272a;
  font-size: 0.62rem;
  font-weight: 700;
  color: #a1a1aa;
  padding: 4px 8px;
  letter-spacing: 0.05em;
  border-bottom: 1px solid rgba(255, 255, 255, 0.04);
}
.response-body {
  margin: 0;
  padding: 10px;
  max-height: 180px;
  overflow-y: auto;
}
.response-body code {
  font-family: var(--font-mono, monospace);
  font-size: 0.72rem;
  line-height: 1.45;
  white-space: pre-wrap;
  color: #4ade80;
}

/* Navigation footer */
.guide-nav-footer {
  height: 54px;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
  padding: 0 20px;
  display: flex;
  align-items: center;
  gap: 10px;
  background: rgba(0, 0, 0, 0.15);
}
.btn-step-prev, .btn-step-next {
  flex: 1;
  padding: 8px;
  font-size: 0.78rem;
  font-weight: 600;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.15s;
}
.btn-step-prev {
  background: #27272a;
  border: 1px solid rgba(255, 255, 255, 0.08);
  color: #d4d4d8;
}
.btn-step-prev:hover:not(:disabled) {
  background: #3f3f46;
  color: #fff;
}
.btn-step-next {
  background: rgba(56, 189, 248, 0.15);
  border: 1px solid rgba(56, 189, 248, 0.3);
  color: #38bdf8;
}
.btn-step-next:hover:not(:disabled) {
  background: rgba(56, 189, 248, 0.25);
  color: #fff;
}
.btn-step-prev:disabled, .btn-step-next:disabled {
  opacity: 0.25;
  cursor: not-allowed;
}

/* ─── Right Pane: VS Code Window ─────────── */
.vscode-pane {
  flex: 1;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
  overflow: hidden;
  box-sizing: border-box;
}

.vscode-window {
  width: 100%;
  height: 100%;
  background: #1e1e1e;
  border-radius: 8px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  box-shadow: 0 15px 45px rgba(0, 0, 0, 0.6);
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.vscode-menubar {
  height: 34px;
  background: #323233;
  display: flex;
  align-items: center;
  padding: 0 12px;
  flex-shrink: 0;
  border-bottom: 1px solid #1e1e1e;
  position: relative;
  user-select: none;
}

.menubar-traffic {
  display: flex;
  align-items: center;
  gap: 7px;
  z-index: 1;
}
.traffic-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  display: block;
  flex-shrink: 0;
}
.dot-red    { background: #ff5f56; cursor: pointer; text-decoration: none; }
.dot-red:hover { filter: brightness(1.2); }
.dot-yellow { background: #ffbd2e; }
.dot-green  { background: #27c93f; }

.menubar-title {
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.7rem;
  color: #88888b;
  font-family: -apple-system, BlinkMacSystemFont, sans-serif;
}

.menubar-close {
  margin-left: auto;
  z-index: 1;
  background: none;
  border: none;
  color: #858585;
  cursor: pointer;
  display: flex;
  align-items: center;
  transition: color 0.15s;
}
.menubar-close:hover {
  color: #fff;
}

.vscode-body {
  display: flex;
  flex: 1;
  overflow: hidden;
}

/* Activity Bar styling */
.vscode-activity-bar {
  width: 48px;
  background: #333333;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 10px 0;
  flex-shrink: 0;
  border-right: 1px solid #191919;
}
.activity-icon {
  width: 38px;
  height: 38px;
  background: none;
  border: none;
  color: #858585;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: color 0.1s;
}
.activity-icon:hover {
  color: #d4d4d4;
}
.activity-icon.active {
  color: #fff;
  border-left: 2px solid #007acc;
}

/* VS Code Sidebar styling */
.vscode-sidebar {
  width: 230px;
  background: #252526;
  border-right: 1px solid #191919;
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
  overflow: hidden;
}
.sidebar-header {
  font-size: 0.65rem;
  font-weight: 700;
  letter-spacing: 0.08em;
  color: #a1a1aa;
  padding: 10px 14px 6px;
  text-transform: uppercase;
}
.file-tree {
  flex: 1;
  overflow-y: auto;
}
.tree-item {
  display: flex;
  align-items: center;
  gap: 6px;
  height: 22px;
  cursor: pointer;
  color: #c5c5c5;
  font-family: var(--font-mono, monospace);
  font-size: 0.72rem;
  position: relative;
}
.tree-item:hover {
  background: #2a2d2e;
}
.tree-item.is-active {
  background: #37373d;
  color: #fff;
}
.tree-item.is-highlighted .tree-name {
  color: #4ade80;
  font-weight: 600;
}
.tree-item.folder {
  color: #e3e3e6;
  cursor: default;
}
.tree-icon {
  display: flex;
  align-items: center;
  color: #858585;
  flex-shrink: 0;
}
.tree-name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.tree-name.go { color: #82aaff; }
.tree-name.sql { color: #f78c6c; }
.tree-badge {
  margin-left: auto;
  margin-right: 8px;
  font-size: 0.55rem;
  font-weight: 700;
  background: rgba(74, 222, 128, 0.12);
  color: #4ade80;
  border: 1px solid rgba(74, 222, 128, 0.2);
  padding: 0 4px;
  border-radius: 3px;
}

/* Editor Area styling */
.editor-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: #1e1e1e;
  overflow: hidden;
  min-width: 0;
}

.vscode-tabs {
  height: 35px;
  background: #2d2d2d;
  display: flex;
  overflow-x: auto;
  border-bottom: 1px solid #191919;
  flex-shrink: 0;
}
.editor-tab {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 0 12px;
  height: 100%;
  font-size: 0.72rem;
  color: #969696;
  cursor: pointer;
  border-right: 1px solid #252526;
  border-bottom: 2px solid transparent;
  white-space: nowrap;
  box-sizing: border-box;
}
.editor-tab:hover {
  background: #1e1e1e;
  color: #ccc;
}
.editor-tab.is-active {
  background: #1e1e1e;
  color: #fff;
  border-bottom-color: #007acc;
}
.tab-icon {
  display: inline-flex;
  align-items: center;
}
.tab-icon.go { color: #82aaff; }
.tab-icon.sql { color: #f78c6c; }

.editor-pane {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
.editor-breadcrumb {
  height: 22px;
  background: #1e1e1e;
  display: flex;
  align-items: center;
  padding: 0 12px;
  gap: 4px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.03);
  font-size: 0.65rem;
  color: #858585;
}
.bc-sep { color: #555; }
.bc-file { color: #cccccc; }

.editor-scroll {
  flex: 1;
  display: flex;
  overflow: auto;
  background: #1e1e1e;
}
.editor-gutter {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  padding: 12px 6px;
  min-width: 34px;
  border-right: 1px solid #2d2d2d;
  flex-shrink: 0;
  user-select: none;
}
.line-num {
  font-size: 0.72rem;
  line-height: 1.5;
  color: #4f4f4f;
  font-family: var(--font-mono, 'JetBrains Mono', monospace);
}
.editor-code-body {
  flex: 1;
  margin: 0;
  padding: 12px 14px;
  overflow: visible;
  line-height: 1.5;
}
.editor-code-body code {
  font-family: var(--font-mono, 'JetBrains Mono', monospace);
  font-size: 0.78rem;
  color: #d4d4d4;
  display: block;
}

/* syntax code highlight */
:deep(.hljs-keyword)   { color: #569cd6; }
:deep(.hljs-built_in)  { color: #4ec9b0; }
:deep(.hljs-type)      { color: #4ec9b0; }
:deep(.hljs-string)    { color: #ce9178; }
:deep(.hljs-number)    { color: #b5cea8; }
:deep(.hljs-comment)   { color: #6a9955; font-style: italic; }
:deep(.hljs-title)     { color: #dcdcaa; }
:deep(.hljs-params)    { color: #9cdcfe; }
:deep(.hljs-attr)      { color: #9cdcfe; }
:deep(.hljs-literal)   { color: #569cd6; }

/* Terminal panel styling */
.terminal-pane {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: #1e1e1e;
  overflow: hidden;
}
.terminal-topbar {
  height: 28px;
  background: #2d2d2d;
  border-bottom: 1px solid #191919;
  display: flex;
  align-items: center;
  padding: 0 12px;
  font-size: 0.7rem;
  color: #a1a1aa;
}
.tbar-label {
  display: flex;
  align-items: center;
  gap: 6px;
  flex: 1;
}
.tbar-close {
  background: none;
  border: none;
  color: #858585;
  cursor: pointer;
}
.tbar-close:hover {
  color: #fff;
}
.terminal-body {
  flex: 1;
  overflow-y: auto;
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  font-family: var(--font-mono, monospace);
  font-size: 0.72rem;
  line-height: 1.45;
}
.t-line { white-space: pre-wrap; color: #d4d4d4; }
.t-line.cmd   { color: #9cdcfe; }
.t-line.done  { color: #4ade80; }
.t-line.dim   { color: #6e6e73; }
.t-line.logo  { color: #38bdf8; }
.t-line :deep(.t-prompt) { color: #38bdf8; font-weight: bold; }
.t-line :deep(.t-green)  { color: #4ade80; }
.t-line :deep(.t-red)    { color: #f87171; }
.t-line :deep(.t-dim)    { color: #71717a; }
.t-line :deep(.t-badge)  {
  display: inline-block;
  font-size: 0.6rem;
  background: rgba(74, 222, 128, 0.1);
  color: #4ade80;
  border: 1px solid rgba(74, 222, 128, 0.25);
  padding: 1px 4px;
  border-radius: 3px;
  margin-right: 4px;
}
.t-cmd { color: #9cdcfe; display: flex; align-items: center; gap: 4px; }
.t-prompt { color: #38bdf8; font-weight: bold; }
.t-typing { color: #fff; }
.t-cursor { animation: blink 0.9s step-end infinite; color: #fff; }

@keyframes blink { 0%, 100% { opacity: 1; } 50% { opacity: 0; } }

.editor-bottom-bar {
  height: 28px;
  background: #2d2d2d;
  border-top: 1px solid #191919;
  display: flex;
  align-items: center;
  padding: 0 8px;
  flex-shrink: 0;
}
.bottom-tabs {
  display: flex;
}
.bottom-tab {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 0.68rem;
  color: #969696;
  padding: 0 10px;
  height: 28px;
  background: none;
  border: none;
  cursor: pointer;
  font-weight: 600;
}
.bottom-tab.active {
  color: #fff;
  background: #1e1e1e;
  border-top: 1px solid #007acc;
}

/* Status bar styling */
.vscode-statusbar {
  height: 22px;
  background: #007acc;
  display: flex;
  align-items: center;
  padding: 0 8px;
  flex-shrink: 0;
}
.status-left, .status-right {
  display: flex;
  align-items: center;
}
.status-right {
  margin-left: auto;
}
.status-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 0.65rem;
  color: #fff;
  padding: 0 8px;
  height: 22px;
}
.status-branch { font-weight: 600; }
.status-step { background: rgba(0, 0, 0, 0.15); font-weight: 600; }

/* Responsive styles */
@media (max-width: 800px) {
  .tutorial-layout {
    flex-direction: column;
  }
  .tutorial-guide-pane {
    width: 100%;
    max-width: 100%;
    height: 40%;
  }
  .vscode-pane {
    height: 60%;
    padding: 8px;
  }
}
</style>
