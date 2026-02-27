<template>
  <main class="home">
    <!-- ─── Hero ─────────────────────────────── -->
    <section class="hero">
      <div class="hero-glow" aria-hidden="true"></div>

      <div class="hero-inner">
        <div class="hero-badge">
          <span class="hero-badge-dot"></span>
          Open Source · MIT License
        </div>

        <div class="hero-logo" aria-label="grove">
          <span class="hero-logo-top">█▀▀ █▀█ █▀█ █░█ █▀▀</span>
          <span class="hero-logo-bot">█▄█ █▀▄ █▄█ ▀▄▀ ██▄</span>
        </div>

        <h1 class="hero-title">
          An opinionated Go foundation for building
          <span class="hero-title-accent"
            >structured, observable,<br />production-ready</span
          >
          applications.
        </h1>

        <p class="hero-sub">
          Grove wires together GORM, fuego and Atlas into a single CLI so you
          can scaffold models, controllers, DTOs and migrations in seconds — and
          focus entirely on your business logic.
        </p>

        <div class="hero-actions">
          <router-link to="/docs" class="btn btn-primary">
            Get Started
            <svg
              width="14"
              height="14"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2.5"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <line x1="5" y1="12" x2="19" y2="12" />
              <polyline points="12 5 19 12 12 19" />
            </svg>
          </router-link>
          <a
            href="https://github.com/caiolandgraf/grove"
            target="_blank"
            rel="noopener"
            class="btn btn-ghost"
          >
            <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
              <path
                d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0 1 12 6.844a9.59 9.59 0 0 1 2.504.337c1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.02 10.02 0 0 0 22 12.017C22 6.484 17.522 2 12 2z"
              />
            </svg>
            GitHub
          </a>
        </div>

        <!-- Install snippet -->
        <div class="hero-install">
          <span class="hero-install-prompt">$</span>
          <span class="hero-install-cmd"
            >go install github.com/caiolandgraf/grove@latest</span
          >
          <button
            class="hero-install-copy"
            @click="copyInstall"
            :class="{ copied: installCopied }"
            :aria-label="installCopied ? 'Copied!' : 'Copy install command'"
          >
            <svg
              v-if="!installCopied"
              width="13"
              height="13"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <rect x="9" y="9" width="13" height="13" rx="2" />
              <path
                d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"
              />
            </svg>
            <svg
              v-else
              width="13"
              height="13"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2.5"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <polyline points="20 6 9 17 4 12" />
            </svg>
          </button>
        </div>
      </div>

      <!-- Scroll indicator -->
      <div class="hero-scroll" aria-hidden="true">
        <div class="hero-scroll-line"></div>
      </div>
    </section>

    <!-- ─── Features ─────────────────────────── -->
    <section class="features">
      <div class="section-inner">
        <div class="section-label">Why Grove</div>
        <h2 class="section-title">Everything you need, nothing you don't</h2>

        <div class="features-grid">
          <div v-for="f in features" :key="f.title" class="feature-card">
            <div class="feature-icon" v-html="f.icon"></div>
            <h3 class="feature-title">{{ f.title }}</h3>
            <p class="feature-desc">{{ f.desc }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- ─── Stack ─────────────────────────────── -->
    <section class="stack">
      <div class="section-inner">
        <div class="section-label">The Stack</div>
        <h2 class="section-title">
          Battle-tested libraries, intelligently connected
        </h2>

        <div class="stack-grid">
          <a
            v-for="s in stack"
            :key="s.name"
            :href="s.url"
            target="_blank"
            rel="noopener"
            class="stack-card"
          >
            <span class="stack-name">{{ s.name }}</span>
            <span class="stack-role">{{ s.role }}</span>
            <svg
              class="stack-ext"
              width="11"
              height="11"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <line x1="7" y1="17" x2="17" y2="7" />
              <polyline points="7 7 17 7 17 17" />
            </svg>
          </a>
        </div>
      </div>
    </section>

    <!-- ─── Quick Start ──────────────────────── -->
    <section class="quickstart">
      <div class="section-inner qs-inner">
        <div class="qs-content">
          <div class="section-label">Quick Start</div>
          <h2 class="section-title">
            From zero to running API<br />in under a minute
          </h2>
          <p class="qs-desc">
            Scaffold a full project, configure your environment and start the
            server. Grove takes care of the boilerplate so your first commit can
            be meaningful business logic.
          </p>
          <router-link to="/docs#quick-start" class="btn btn-primary">
            Read the Docs
            <svg
              width="14"
              height="14"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2.5"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <line x1="5" y1="12" x2="19" y2="12" />
              <polyline points="12 5 19 12 12 19" />
            </svg>
          </router-link>
        </div>

        <div class="qs-terminal">
          <div class="terminal">
            <div class="terminal-bar">
              <span class="terminal-dot terminal-dot--red"></span>
              <span class="terminal-dot terminal-dot--yellow"></span>
              <span class="terminal-dot terminal-dot--green"></span>
              <span class="terminal-title">zsh</span>
            </div>
            <div class="terminal-body">
              <div
                v-for="(line, i) in terminalLines"
                :key="i"
                class="terminal-line"
                :class="`terminal-line--${line.type}`"
                :style="{ animationDelay: `${i * 0.08}s` }"
              >
                <span v-if="line.type === 'cmd'" class="terminal-prompt"
                  >❯</span
                >
                <span v-html="line.text"></span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- ─── CTA ───────────────────────────────── -->
    <section class="cta">
      <div class="cta-glow" aria-hidden="true"></div>
      <div class="section-inner cta-inner">
        <h2 class="cta-title">Ready to build?</h2>
        <p class="cta-sub">
          Install grove and scaffold your first resource in seconds.
        </p>
        <div class="cta-actions">
          <router-link to="/docs" class="btn btn-primary btn-lg"
            >Get Started</router-link
          >
          <router-link to="/contributors" class="btn btn-ghost btn-lg"
            >Contributors</router-link
          >
        </div>
      </div>
    </section>
  </main>
</template>

<script setup>
import { ref } from 'vue'

// ── Copy install command ───────────────────────
const installCopied = ref(false)
async function copyInstall() {
  try {
    await navigator.clipboard.writeText(
      'go install github.com/caiolandgraf/grove@latest'
    )
    installCopied.value = true
    setTimeout(() => {
      installCopied.value = false
    }, 2000)
  } catch {
    /* noop */
  }
}

// ── Features data ──────────────────────────────
const features = [
  {
    icon: `<svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round"><polyline points="16 18 22 12 16 6"/><polyline points="8 6 2 12 8 18"/></svg>`,
    title: 'One command, full resource',
    desc: 'grove make:resource generates a typed GORM model, fuego controller and DTO structs in a single command — fully wired together.'
  },
  {
    icon: `<svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round"><ellipse cx="12" cy="5" rx="9" ry="3"/><path d="M21 12c0 1.66-4 3-9 3s-9-1.34-9-3"/><path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5"/></svg>`,
    title: 'Schema migrations via Atlas',
    desc: 'No manual SQL. grove make:migration diffs your Go structs against the live schema and generates precise, versioned migration files.'
  },
  {
    icon: `<svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/><polyline points="10 9 9 9 8 9"/></svg>`,
    title: 'OpenAPI 3.1 for free',
    desc: 'fuego infers request/response types from your handler signatures and serves a live Swagger UI — zero annotation boilerplate.'
  },
  {
    icon: `<svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="3"/><path d="M19.07 4.93a10 10 0 0 1 0 14.14M4.93 4.93a10 10 0 0 0 0 14.14"/><path d="M15.54 8.46a5 5 0 0 1 0 7.07M8.46 8.46a5 5 0 0 0 0 7.07"/></svg>`,
    title: 'Hot-reload out of the box',
    desc: 'If air is installed, grove serve uses it automatically. Instant feedback loop with no configuration required.'
  },
  {
    icon: `<svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="18" height="18" rx="2"/><path d="M3 9h18"/><path d="M9 21V9"/></svg>`,
    title: 'Typed repository layer',
    desc: 'Every model gets a generic Repository[T] — Find, All, Create, Update, Delete — with GORM available for custom queries when needed.'
  },
  {
    icon: `<svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round"><path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/></svg>`,
    title: 'Opinionated, not rigid',
    desc: "The project layout enforces separation of concerns via Go's internal/ boundary. You always own the code — grove just scaffolds it."
  }
]

// ── Stack data ─────────────────────────────────
const stack = [
  { name: 'GORM', role: 'ORM & repository layer', url: 'https://gorm.io' },
  {
    name: 'fuego',
    role: 'HTTP router + OpenAPI generation',
    url: 'https://github.com/go-fuego/fuego'
  },
  { name: 'Atlas', role: 'Schema migration engine', url: 'https://atlasgo.io' },
  {
    name: 'Cobra',
    role: 'CLI framework',
    url: 'https://github.com/spf13/cobra'
  },
  {
    name: 'air',
    role: 'Hot-reload for development',
    url: 'https://github.com/air-verse/air'
  },
  { name: 'Go', role: '1.22+ required', url: 'https://go.dev' }
]

// ── Terminal animation data ────────────────────
const terminalLines = [
  { type: 'cmd', text: 'grove setup my-api' },
  {
    type: 'out',
    text: '<span style="color:#4ade80">  ✓</span>  Downloading template      <span style="color:#3a3c55">86 KB</span>'
  },
  {
    type: 'out',
    text: '<span style="color:#4ade80">  ✓</span>  Extracting files          <span style="color:#3a3c55">34 files</span>'
  },
  {
    type: 'out',
    text: '<span style="color:#4ade80">  ✓</span>  Configuring module        <span style="color:#3a3c55">github.com/you/my-api</span>'
  },
  {
    type: 'out',
    text: '<span style="color:#4ade80">  ✓</span>  Installing dependencies   <span style="color:#3a3c55">1.4s</span>'
  },
  { type: 'spacer', text: '' },
  { type: 'cmd', text: 'cd my-api && grove serve' },
  {
    type: 'out',
    text: '<span style="color:#c82838">  ▸ AIR </span>  Starting server with <b>air</b>'
  },
  {
    type: 'out',
    text: '<span style="color:#3a3c55">        http://localhost:8080</span>'
  },
  {
    type: 'out',
    text: '<span style="color:#3a3c55">        http://localhost:8080/swagger</span>'
  },
  { type: 'spacer', text: '' },
  { type: 'cmd', text: 'grove make:resource Post' },
  {
    type: 'out',
    text: '  <span style="color:#4ade80;background:rgba(74,222,128,.1);padding:1px 6px;border-radius:3px"> CREATED </span>  Model      <b>Post</b>  <span style="color:#3a3c55">→ internal/models/post.go</span>'
  },
  {
    type: 'out',
    text: '  <span style="color:#4ade80;background:rgba(74,222,128,.1);padding:1px 6px;border-radius:3px"> CREATED </span>  Controller <b>Post</b>  <span style="color:#3a3c55">→ internal/controllers/post-controller.go</span>'
  },
  {
    type: 'out',
    text: '  <span style="color:#4ade80;background:rgba(74,222,128,.1);padding:1px 6px;border-radius:3px"> CREATED </span>  Request    <b>Post</b>  <span style="color:#3a3c55">→ internal/dto/post-dto.go</span>'
  }
]
</script>

<style scoped>
/* ─────────────────────────────────────────────
   Shared
───────────────────────────────────────────── */
.home {
  padding-top: var(--nav-h);
  overflow: hidden;
}

.section-inner {
  max-width: 1160px;
  margin: 0 auto;
  padding: 0 2rem;
}

.section-label {
  font-size: 0.72rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.12em;
  color: var(--red);
  margin-bottom: 0.75rem;
}

.section-title {
  font-size: clamp(1.6rem, 3vw, 2.2rem);
  font-weight: 700;
  letter-spacing: -0.03em;
  line-height: 1.2;
  margin-bottom: 0.5em;
  color: var(--text);
}

/* ─────────────────────────────────────────────
   Buttons
───────────────────────────────────────────── */
.btn {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.9rem;
  font-weight: 500;
  padding: 0.6em 1.3em;
  border-radius: var(--radius-sm);
  transition: all 0.18s var(--ease);
  border: 1px solid transparent;
  white-space: nowrap;
}

.btn-primary {
  background: var(--red);
  color: #fff;
  border-color: var(--red);
}
.btn-primary:hover {
  background: var(--red-hover);
  border-color: var(--red-hover);
  transform: translateY(-1px);
  box-shadow: 0 4px 20px rgba(200, 40, 56, 0.35);
}

.btn-ghost {
  color: var(--text-muted);
  border-color: var(--border-md);
  background: transparent;
}
.btn-ghost:hover {
  color: var(--text);
  border-color: var(--border-md);
  background: rgba(255, 255, 255, 0.05);
  transform: translateY(-1px);
}

.btn-lg {
  font-size: 1rem;
  padding: 0.7em 1.7em;
}

/* ─────────────────────────────────────────────
   Hero
───────────────────────────────────────────── */
.hero {
  position: relative;
  min-height: calc(100vh - var(--nav-h));
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 5rem 2rem 4rem;
  overflow: hidden;
}

.hero-glow {
  position: absolute;
  inset: -20% -20%;
  background:
    radial-gradient(
      ellipse 60% 50% at 50% 0%,
      rgba(200, 40, 56, 0.13) 0%,
      transparent 70%
    ),
    radial-gradient(
      ellipse 40% 30% at 50% 100%,
      rgba(200, 40, 56, 0.05) 0%,
      transparent 60%
    );
  pointer-events: none;
}

.hero-inner {
  position: relative;
  z-index: 1;
  max-width: 780px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0;
}

/* Badge */
.hero-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.75rem;
  color: var(--text-muted);
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid var(--border-md);
  border-radius: 99px;
  padding: 0.3em 0.9em;
  margin-bottom: 2rem;
  letter-spacing: 0.02em;
}

.hero-badge-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--red);
  box-shadow: 0 0 6px var(--red);
  animation: pulse 2s ease-in-out infinite;
}

@keyframes pulse {
  0%,
  100% {
    opacity: 1;
    box-shadow: 0 0 6px var(--red);
  }
  50% {
    opacity: 0.6;
    box-shadow: 0 0 12px var(--red);
  }
}

/* Block-char logo */
.hero-logo {
  display: flex;
  flex-direction: column;
  font-family: var(--font-mono);
  font-size: clamp(0.9rem, 2.5vw, 1.35rem);
  line-height: 1.4;
  letter-spacing: 0.06em;
  margin-bottom: 2.2rem;
}

.hero-logo-top {
  color: var(--red);
  text-shadow: 0 0 40px rgba(200, 40, 56, 0.5);
}

.hero-logo-bot {
  color: #8b1c28;
  text-shadow: 0 0 30px rgba(139, 28, 40, 0.4);
}

/* Title */
.hero-title {
  font-size: clamp(1.55rem, 3.5vw, 2.4rem);
  font-weight: 700;
  letter-spacing: -0.03em;
  line-height: 1.25;
  color: var(--text);
  margin-bottom: 1.25rem;
}

.hero-title-accent {
  color: transparent;
  background: linear-gradient(
    135deg,
    var(--red) 0%,
    #e85060 50%,
    var(--red-hover) 100%
  );
  -webkit-background-clip: text;
  background-clip: text;
}

/* Subtitle */
.hero-sub {
  font-size: 1rem;
  color: var(--text-muted);
  line-height: 1.7;
  max-width: 560px;
  margin-bottom: 2rem;
}

/* Actions */
.hero-actions {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 2rem;
  flex-wrap: wrap;
  justify-content: center;
}

/* Install snippet */
.hero-install {
  display: inline-flex;
  align-items: center;
  gap: 0.65rem;
  background: var(--bg-card);
  border: 1px solid var(--border-md);
  border-radius: var(--radius-sm);
  padding: 0.6em 0.9em 0.6em 1.1em;
  font-family: var(--font-mono);
  font-size: 0.82rem;
  color: var(--text-muted);
}

.hero-install-prompt {
  color: var(--red);
  font-weight: 600;
  user-select: none;
}

.hero-install-cmd {
  color: var(--text);
}

.hero-install-copy {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 26px;
  height: 26px;
  border-radius: 4px;
  color: var(--text-dim);
  transition:
    color 0.15s,
    background 0.15s;
}

.hero-install-copy:hover {
  color: var(--text-muted);
  background: rgba(255, 255, 255, 0.06);
}

.hero-install-copy.copied {
  color: #4ade80;
}

/* Scroll indicator */
.hero-scroll {
  position: absolute;
  bottom: 2rem;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  flex-direction: column;
  align-items: center;
}

.hero-scroll-line {
  width: 1px;
  height: 48px;
  background: linear-gradient(to bottom, var(--red), transparent);
  animation: scrollLine 1.8s ease-in-out infinite;
}

@keyframes scrollLine {
  0% {
    opacity: 0;
    transform: scaleY(0);
    transform-origin: top;
  }
  40% {
    opacity: 1;
  }
  100% {
    opacity: 0;
    transform: scaleY(1);
    transform-origin: top;
  }
}

/* ─────────────────────────────────────────────
   Features
───────────────────────────────────────────── */
.features {
  padding: 7rem 0 6rem;
  border-top: 1px solid var(--border);
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1px;
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  overflow: hidden;
  margin-top: 3rem;
  background: var(--border);
}

.feature-card {
  background: var(--bg-elevated);
  padding: 1.8rem 1.6rem;
  transition: background 0.2s;
}

.feature-card:hover {
  background: var(--bg-card);
}

.feature-icon {
  color: var(--red);
  margin-bottom: 1rem;
  opacity: 0.9;
}

.feature-title {
  font-size: 0.975rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
  color: var(--text);
}

.feature-desc {
  font-size: 0.84rem;
  color: var(--text-muted);
  line-height: 1.65;
  margin: 0;
}

/* ─────────────────────────────────────────────
   Stack
───────────────────────────────────────────── */
.stack {
  padding: 6rem 0;
  border-top: 1px solid var(--border);
}

.stack-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 0.75rem;
  margin-top: 2.5rem;
}

.stack-card {
  position: relative;
  display: flex;
  flex-direction: column;
  gap: 0.3rem;
  padding: 1.2rem 1.3rem;
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  transition:
    border-color 0.2s,
    background 0.2s,
    transform 0.2s;
}

.stack-card:hover {
  border-color: var(--red-border);
  background: var(--bg-hover);
  transform: translateY(-2px);
}

.stack-name {
  font-family: var(--font-mono);
  font-size: 0.9rem;
  font-weight: 600;
  color: var(--text);
}

.stack-role {
  font-size: 0.78rem;
  color: var(--text-muted);
}

.stack-ext {
  position: absolute;
  top: 1rem;
  right: 1rem;
  color: var(--text-dim);
  transition: color 0.15s;
}

.stack-card:hover .stack-ext {
  color: var(--red);
}

/* ─────────────────────────────────────────────
   Quick Start
───────────────────────────────────────────── */
.quickstart {
  padding: 6rem 0;
  border-top: 1px solid var(--border);
}

.qs-inner {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 5rem;
  align-items: center;
}

.qs-desc {
  color: var(--text-muted);
  line-height: 1.7;
  margin-bottom: 1.75rem;
  font-size: 0.95rem;
}

/* Terminal */
.qs-terminal {
  perspective: 1000px;
}

.terminal {
  background: #080810;
  border: 1px solid var(--border-md);
  border-radius: var(--radius);
  overflow: hidden;
  box-shadow:
    0 0 0 1px rgba(200, 40, 56, 0.06),
    0 24px 60px rgba(0, 0, 0, 0.6);
}

.terminal-bar {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  padding: 0.65rem 1rem;
  background: rgba(255, 255, 255, 0.03);
  border-bottom: 1px solid var(--border);
}

.terminal-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
}

.terminal-dot--red {
  background: #ff5f57;
}
.terminal-dot--yellow {
  background: #febc2e;
}
.terminal-dot--green {
  background: #28c840;
}

.terminal-title {
  flex: 1;
  text-align: center;
  font-family: var(--font-mono);
  font-size: 0.7rem;
  color: var(--text-dim);
  margin-right: 1.4rem;
}

.terminal-body {
  padding: 1.1rem 1.3rem;
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.terminal-line {
  display: flex;
  align-items: baseline;
  gap: 0.6rem;
  font-family: var(--font-mono);
  font-size: 0.78rem;
  line-height: 1.6;
  opacity: 0;
  animation: fadeInLine 0.3s var(--ease) forwards;
}

@keyframes fadeInLine {
  from {
    opacity: 0;
    transform: translateY(4px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.terminal-line--cmd {
  color: var(--text);
}
.terminal-line--out {
  color: var(--text-muted);
  padding-left: 0.3rem;
}
.terminal-line--spacer {
  height: 0.4rem;
}

.terminal-prompt {
  color: var(--red);
  user-select: none;
  flex-shrink: 0;
}

/* ─────────────────────────────────────────────
   CTA
───────────────────────────────────────────── */
.cta {
  position: relative;
  padding: 7rem 0;
  border-top: 1px solid var(--border);
  text-align: center;
  overflow: hidden;
}

.cta-glow {
  position: absolute;
  inset: 0;
  background: radial-gradient(
    ellipse 70% 60% at 50% 100%,
    rgba(200, 40, 56, 0.1) 0%,
    transparent 70%
  );
  pointer-events: none;
}

.cta-inner {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.cta-title {
  font-size: clamp(1.8rem, 4vw, 2.8rem);
  font-weight: 700;
  letter-spacing: -0.035em;
  color: var(--text);
  margin-bottom: 0.6rem;
}

.cta-sub {
  color: var(--text-muted);
  font-size: 1rem;
  margin-bottom: 2rem;
}

.cta-actions {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  flex-wrap: wrap;
  justify-content: center;
}

/* ─────────────────────────────────────────────
   Responsive
───────────────────────────────────────────── */
@media (max-width: 900px) {
  .features-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .stack-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .qs-inner {
    grid-template-columns: 1fr;
    gap: 3rem;
  }
}

@media (max-width: 600px) {
  .features-grid {
    grid-template-columns: 1fr;
  }

  .stack-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .hero {
    padding: 3rem 1.25rem 4rem;
  }

  .section-inner {
    padding: 0 1.25rem;
  }

  .hero-install {
    font-size: 0.72rem;
    padding: 0.55em 0.75em 0.55em 0.9em;
    max-width: 100%;
    overflow-x: auto;
  }
}
</style>
