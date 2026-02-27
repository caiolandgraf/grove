<template>
  <header class="nav" :class="{ scrolled }">
    <div class="nav-inner">
      <!-- Logo -->
      <router-link to="/" class="nav-logo" aria-label="Grove home">
        <span class="nav-logo-top">█▀▀ █▀█ █▀█ █░█ █▀▀</span>
        <span class="nav-logo-bot">█▄█ █▀▄ █▄█ ▀▄▀ ██▄</span>
      </router-link>

      <!-- Links -->
      <nav class="nav-links" aria-label="Main navigation">
        <router-link to="/" class="nav-link" exact-active-class="active">Home</router-link>
        <router-link to="/docs" class="nav-link" active-class="active">Docs</router-link>
        <router-link to="/contributors" class="nav-link" active-class="active">Contributors</router-link>
      </nav>

      <!-- Right actions -->
      <div class="nav-actions">
        <!-- Search trigger -->
        <button class="nav-search" @click="openSearch" aria-label="Search (Ctrl+K)">
          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="11" cy="11" r="8"/>
            <line x1="21" y1="21" x2="16.65" y2="16.65"/>
          </svg>
          <span class="nav-search-text">Search</span>
          <kbd class="nav-search-kbd">⌘K</kbd>
        </button>

        <!-- GitHub -->
        <a
          href="https://github.com/caiolandgraf/grove"
          target="_blank"
          rel="noopener"
          class="nav-github"
          aria-label="GitHub repository"
        >
          <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0 1 12 6.844a9.59 9.59 0 0 1 2.504.337c1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.02 10.02 0 0 0 22 12.017C22 6.484 17.522 2 12 2z"/>
          </svg>
        </a>

        <!-- Mobile menu toggle -->
        <button class="nav-mobile-toggle" @click="mobileOpen = !mobileOpen" :aria-expanded="mobileOpen" aria-label="Toggle menu">
          <span class="bar" :class="{ open: mobileOpen }"></span>
          <span class="bar" :class="{ open: mobileOpen }"></span>
          <span class="bar" :class="{ open: mobileOpen }"></span>
        </button>
      </div>
    </div>

    <!-- Mobile drawer -->
    <transition name="drawer">
      <div v-if="mobileOpen" class="nav-mobile" @click="mobileOpen = false">
        <router-link to="/" class="nav-mobile-link" exact-active-class="active">Home</router-link>
        <router-link to="/docs" class="nav-mobile-link" active-class="active">Docs</router-link>
        <router-link to="/contributors" class="nav-mobile-link" active-class="active">Contributors</router-link>
        <a href="https://github.com/caiolandgraf/grove" target="_blank" rel="noopener" class="nav-mobile-link">
          GitHub ↗
        </a>
      </div>
    </transition>
  </header>
</template>

<script setup>
import { ref, onMounted, onUnmounted, inject } from 'vue'

const scrolled    = ref(false)
const mobileOpen  = ref(false)
const openSearch  = inject('openSearch', () => {})

function onScroll() {
  scrolled.value = window.scrollY > 10
}

function onKeydown(e) {
  if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
    e.preventDefault()
    openSearch()
  }
}

onMounted(() => {
  window.addEventListener('scroll', onScroll, { passive: true })
  window.addEventListener('keydown', onKeydown)
})

onUnmounted(() => {
  window.removeEventListener('scroll', onScroll)
  window.removeEventListener('keydown', onKeydown)
})
</script>

<style scoped>
/* ── Base ── */
.nav {
  position: fixed;
  inset: 0 0 auto;
  z-index: 100;
  height: var(--nav-h);
  border-bottom: 1px solid transparent;
  transition: background 0.2s, border-color 0.2s, backdrop-filter 0.2s;
}

.nav.scrolled {
  background: rgba(6, 6, 10, 0.85);
  border-color: var(--border);
  backdrop-filter: blur(16px) saturate(1.4);
  -webkit-backdrop-filter: blur(16px) saturate(1.4);
}

.nav-inner {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 2rem;
  height: 100%;
  display: flex;
  align-items: center;
  gap: 2rem;
}

/* ── Logo ── */
.nav-logo {
  display: flex;
  flex-direction: column;
  font-family: var(--font-mono);
  font-size: 0.52rem;
  line-height: 1.35;
  letter-spacing: 0.04em;
  flex-shrink: 0;
  transition: opacity 0.15s;
}

.nav-logo:hover { opacity: 0.8; }

.nav-logo-top { color: var(--red); }
.nav-logo-bot { color: #961428; }

/* ── Links ── */
.nav-links {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  margin-right: auto;
}

.nav-link {
  font-size: 0.875rem;
  font-weight: 450;
  color: var(--text-muted);
  padding: 0.35em 0.75em;
  border-radius: var(--radius-sm);
  transition: color 0.15s, background 0.15s;
}

.nav-link:hover {
  color: var(--text);
  background: rgba(255, 255, 255, 0.05);
}

.nav-link.active {
  color: var(--text);
  background: rgba(255, 255, 255, 0.06);
}

/* ── Actions ── */
.nav-actions {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

/* Search button */
.nav-search {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  height: 34px;
  padding: 0 0.75rem;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border-md);
  background: var(--bg-input);
  color: var(--text-muted);
  font-size: 0.8rem;
  min-width: 160px;
  transition: border-color 0.15s, background 0.15s, color 0.15s;
}

.nav-search:hover {
  border-color: var(--red-border);
  background: var(--red-dim);
  color: var(--text);
}

.nav-search-text {
  flex: 1;
  text-align: left;
}

.nav-search-kbd {
  font-family: var(--font-mono);
  font-size: 0.68rem;
  color: var(--text-dim);
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.09);
  border-radius: 4px;
  padding: 0.1em 0.4em;
  line-height: 1.4;
}

/* GitHub */
.nav-github {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 34px;
  height: 34px;
  border-radius: var(--radius-sm);
  color: var(--text-muted);
  border: 1px solid var(--border);
  transition: color 0.15s, border-color 0.15s, background 0.15s;
}

.nav-github:hover {
  color: var(--text);
  border-color: var(--border-md);
  background: rgba(255, 255, 255, 0.05);
}

/* Mobile toggle */
.nav-mobile-toggle {
  display: none;
  flex-direction: column;
  justify-content: center;
  gap: 5px;
  width: 34px;
  height: 34px;
  padding: 6px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border);
}

.bar {
  display: block;
  height: 1.5px;
  background: var(--text-muted);
  border-radius: 2px;
  transition: transform 0.22s var(--ease), opacity 0.22s;
}

.bar:nth-child(1).open { transform: translateY(6.5px) rotate(45deg); }
.bar:nth-child(2).open { opacity: 0; }
.bar:nth-child(3).open { transform: translateY(-6.5px) rotate(-45deg); }

/* ── Mobile drawer ── */
.nav-mobile {
  display: flex;
  flex-direction: column;
  border-top: 1px solid var(--border);
  background: rgba(6, 6, 10, 0.97);
  backdrop-filter: blur(20px);
  padding: 0.75rem 0;
}

.nav-mobile-link {
  padding: 0.75rem 2rem;
  font-size: 0.9rem;
  color: var(--text-muted);
  transition: color 0.15s, background 0.15s;
}

.nav-mobile-link:hover,
.nav-mobile-link.active {
  color: var(--text);
  background: rgba(255, 255, 255, 0.04);
}

/* Drawer animation */
.drawer-enter-active,
.drawer-leave-active { transition: opacity 0.18s, transform 0.18s var(--ease); }
.drawer-enter-from,
.drawer-leave-to     { opacity: 0; transform: translateY(-6px); }

/* ── Responsive ── */
@media (max-width: 768px) {
  .nav-links       { display: none; }
  .nav-search      { display: none; }
  .nav-mobile-toggle { display: flex; }
  .nav-inner       { gap: 1rem; }
}

@media (max-width: 480px) {
  .nav-inner { padding: 0 1.25rem; }
}
</style>
