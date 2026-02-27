<template>
  <footer class="footer">
    <div class="footer-inner">
      <div class="footer-brand">
        <span class="footer-logo">
          <span class="logo-top">█▀▀ █▀█ █▀█ █░█ █▀▀</span>
          <span class="logo-bot">█▄█ █▀▄ █▄█ ▀▄▀ ██▄</span>
        </span>
        <p class="footer-tagline">
          An opinionated Go foundation for building structured,<br />
          observable, production-ready applications.
        </p>
      </div>

      <div class="footer-links">
        <div class="footer-col">
          <span class="footer-col-title">Docs</span>
          <router-link to="/docs#getting-started">Getting Started</router-link>
          <router-link to="/docs#commands">Commands</router-link>
          <router-link to="/docs#architecture">Architecture</router-link>
          <router-link to="/docs#configuration">Configuration</router-link>
        </div>
        <div class="footer-col">
          <span class="footer-col-title">Project</span>
          <a href="https://github.com/caiolandgraf/grove" target="_blank" rel="noopener">GitHub</a>
          <a href="https://github.com/caiolandgraf/grove/releases" target="_blank" rel="noopener">Releases</a>
          <a href="https://github.com/caiolandgraf/grove/issues" target="_blank" rel="noopener">Issues</a>
          <router-link to="/contributors">Contributors</router-link>
        </div>
        <div class="footer-col">
          <span class="footer-col-title">Stack</span>
          <a href="https://gorm.io" target="_blank" rel="noopener">GORM</a>
          <a href="https://github.com/go-fuego/fuego" target="_blank" rel="noopener">fuego</a>
          <a href="https://atlasgo.io" target="_blank" rel="noopener">Atlas</a>
          <a href="https://github.com/air-verse/air" target="_blank" rel="noopener">air</a>
        </div>
      </div>
    </div>

    <div class="footer-bottom">
      <span class="footer-copy">
        © {{ year }} <a href="https://github.com/caiolandgraf" target="_blank" rel="noopener">Caio Landgraf</a>.
        Released under the MIT License.
      </span>
      <span class="footer-built">
        Built with
        <a href="https://vuejs.org" target="_blank" rel="noopener">Vue 3</a>
        &
        <a href="https://vitejs.dev" target="_blank" rel="noopener">Vite</a>
      </span>
    </div>
  </footer>
</template>

<script setup>
const year = new Date().getFullYear()
</script>

<style scoped>
.footer {
  margin-top: auto;
  border-top: 1px solid var(--border);
  background: var(--bg-elevated);
}

.footer-inner {
  max-width: 1200px;
  margin: 0 auto;
  padding: 3rem 2rem 2rem;
  display: flex;
  gap: 4rem;
  align-items: flex-start;
}

/* ── Brand ── */
.footer-brand {
  flex: 1;
  min-width: 200px;
}

.footer-logo {
  display: flex;
  flex-direction: column;
  font-family: var(--font-mono);
  font-size: 0.6rem;
  line-height: 1.4;
  letter-spacing: 0.05em;
  margin-bottom: 1rem;
}

.logo-top { color: var(--red); }
.logo-bot { color: var(--red-dark, #961428); }

.footer-tagline {
  font-size: 0.8rem;
  color: var(--text-muted);
  line-height: 1.7;
}

/* ── Links grid ── */
.footer-links {
  display: flex;
  gap: 3.5rem;
}

.footer-col {
  display: flex;
  flex-direction: column;
  gap: 0.6rem;
}

.footer-col-title {
  font-size: 0.7rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  color: var(--text-muted);
  margin-bottom: 0.25rem;
}

.footer-col a {
  font-size: 0.85rem;
  color: var(--text-muted);
  transition: color 0.15s;
}

.footer-col a:hover {
  color: var(--text);
}

.footer-col a.router-link-active {
  color: var(--red);
}

/* ── Bottom bar ── */
.footer-bottom {
  max-width: 1200px;
  margin: 0 auto;
  padding: 1rem 2rem;
  border-top: 1px solid var(--border);
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  font-size: 0.78rem;
  color: var(--text-dim);
}

.footer-bottom a {
  color: var(--text-muted);
  transition: color 0.15s;
}

.footer-bottom a:hover {
  color: var(--red);
}

/* ── Responsive ── */
@media (max-width: 768px) {
  .footer-inner {
    flex-direction: column;
    gap: 2rem;
    padding: 2rem 1.25rem 1.5rem;
  }

  .footer-links {
    gap: 2rem;
    flex-wrap: wrap;
  }

  .footer-bottom {
    flex-direction: column;
    align-items: flex-start;
    padding: 1rem 1.25rem;
  }
}
</style>
