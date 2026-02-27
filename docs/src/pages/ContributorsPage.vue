<template>
  <main class="contributors-page">

    <!-- ─── Header ───────────────────────────── -->
    <section class="contrib-header">
      <div class="contrib-glow" aria-hidden="true" />
      <div class="inner">
        <div class="section-label">Contributors</div>
        <h1 class="contrib-title">Built by the community</h1>
        <p class="contrib-sub">
          Grove is open source and grows with every contribution.
          From bug fixes to new features, every bit counts.
        </p>
      </div>
    </section>

    <!-- ─── Stats ────────────────────────────── -->
    <section class="stats">
      <div class="inner stats-grid">
        <div v-for="s in stats" :key="s.label" class="stat-card">
          <span class="stat-value">{{ s.value }}</span>
          <span class="stat-label">{{ s.label }}</span>
        </div>
      </div>
    </section>

    <!-- ─── Contributors grid ────────────────── -->
    <section class="contrib-list">
      <div class="inner">
        <h2 class="list-title">Core contributors</h2>

        <div class="cards-grid">
          <a
            v-for="c in contributors"
            :key="c.login"
            :href="c.url"
            target="_blank"
            rel="noopener"
            class="contrib-card"
          >
            <div class="contrib-avatar-wrap">
              <img
                :src="c.avatar"
                :alt="c.name"
                class="contrib-avatar"
                loading="lazy"
              />
              <span v-if="c.role === 'Author'" class="contrib-badge">Author</span>
            </div>

            <div class="contrib-info">
              <span class="contrib-name">{{ c.name }}</span>
              <span class="contrib-login">@{{ c.login }}</span>
            </div>

            <div class="contrib-tags">
              <span
                v-for="tag in c.contributions"
                :key="tag"
                class="contrib-tag"
              >
                {{ contributionTypes[tag]?.icon }}
                {{ contributionTypes[tag]?.label }}
              </span>
            </div>

            <svg class="contrib-ext" width="12" height="12" viewBox="0 0 24 24"
              fill="none" stroke="currentColor" stroke-width="2"
              stroke-linecap="round" stroke-linejoin="round">
              <line x1="7" y1="17" x2="17" y2="7"/>
              <polyline points="7 7 17 7 17 17"/>
            </svg>
          </a>

          <!-- Placeholder slots -->
          <div v-for="i in placeholders" :key="`ph-${i}`" class="contrib-card contrib-card--ghost">
            <div class="ghost-avatar" />
            <div class="ghost-lines">
              <div class="ghost-line ghost-line--name" />
              <div class="ghost-line ghost-line--login" />
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- ─── CTA ───────────────────────────────── -->
    <section class="join-section">
      <div class="join-glow" aria-hidden="true" />
      <div class="inner join-inner">

        <div class="join-icon" aria-hidden="true">
          <svg width="36" height="36" viewBox="0 0 24 24" fill="none"
            stroke="currentColor" stroke-width="1.4"
            stroke-linecap="round" stroke-linejoin="round">
            <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
          </svg>
        </div>

        <h2 class="join-title">Become a contributor</h2>
        <p class="join-sub">
          Whether it's a typo fix, a new command, or a performance improvement —
          all contributions are welcome. Check the open issues or propose something new.
        </p>

        <div class="join-actions">
          <a
            href="https://github.com/caiolandgraf/grove/issues"
            target="_blank"
            rel="noopener"
            class="btn btn-primary"
          >
            <svg width="15" height="15" viewBox="0 0 24 24" fill="none"
              stroke="currentColor" stroke-width="2"
              stroke-linecap="round" stroke-linejoin="round">
              <circle cx="12" cy="12" r="10"/>
              <line x1="12" y1="8" x2="12" y2="12"/>
              <line x1="12" y1="16" x2="12.01" y2="16"/>
            </svg>
            Browse Issues
          </a>
          <a
            href="https://github.com/caiolandgraf/grove/fork"
            target="_blank"
            rel="noopener"
            class="btn btn-ghost"
          >
            <svg width="15" height="15" viewBox="0 0 24 24" fill="none"
              stroke="currentColor" stroke-width="2"
              stroke-linecap="round" stroke-linejoin="round">
              <circle cx="6"  cy="18" r="3"/>
              <circle cx="6"  cy="6"  r="3"/>
              <circle cx="18" cy="6"  r="3"/>
              <path d="M6 9v6"/>
              <path d="M18 9a9 9 0 0 1-9 9"/>
            </svg>
            Fork on GitHub
          </a>
        </div>

        <!-- How to contribute steps -->
        <div class="how-grid">
          <div v-for="step in howSteps" :key="step.title" class="how-card">
            <span class="how-num">{{ step.num }}</span>
            <strong class="how-title">{{ step.title }}</strong>
            <p class="how-desc">{{ step.desc }}</p>
          </div>
        </div>

      </div>
    </section>

  </main>
</template>

<script setup>
import { computed } from 'vue'
import { contributors, contributionTypes } from '@/data/contributors.js'

// Fill up the grid with ghost placeholders
const TARGET_CARDS = 6
const placeholders = computed(() =>
  Math.max(0, TARGET_CARDS - contributors.length)
)

const stats = [
  { value: '1',    label: 'Contributor' },
  { value: '0.1',  label: 'Current version' },
  { value: 'MIT',  label: 'License' },
  { value: '∞',    label: 'Contributions welcome' },
]

const howSteps = [
  {
    num: '01',
    title: 'Fork & clone',
    desc: 'Fork the repository on GitHub and clone it locally to start working.',
  },
  {
    num: '02',
    title: 'Make your change',
    desc: 'Fix a bug, add a feature or improve the docs. Build with make grove-build.',
  },
  {
    num: '03',
    title: 'Open a pull request',
    desc: 'Push your branch and open a PR with a clear description of what changed and why.',
  },
]
</script>

<style scoped>
/* ─────────────────────────────────────────────
   Layout
───────────────────────────────────────────── */
.contributors-page {
  padding-top: var(--nav-h);
  overflow: hidden;
}

.inner {
  max-width: 1100px;
  margin: 0 auto;
  padding: 0 2rem;
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
  padding: 0.62em 1.3em;
  border-radius: var(--radius-sm);
  border: 1px solid transparent;
  transition: all 0.18s var(--ease);
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
}
.btn-ghost:hover {
  color: var(--text);
  background: rgba(255, 255, 255, 0.05);
  transform: translateY(-1px);
}

/* ─────────────────────────────────────────────
   Header
───────────────────────────────────────────── */
.contrib-header {
  position: relative;
  padding: 6rem 0 4rem;
  text-align: center;
  border-bottom: 1px solid var(--border);
  overflow: hidden;
}

.contrib-glow {
  position: absolute;
  inset: -10% -20%;
  background: radial-gradient(ellipse 70% 60% at 50% 0%, rgba(200, 40, 56, 0.11) 0%, transparent 70%);
  pointer-events: none;
}

.contrib-header .inner {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.section-label {
  font-size: 0.72rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.12em;
  color: var(--red);
  margin-bottom: 0.75rem;
}

.contrib-title {
  font-size: clamp(2rem, 5vw, 3.2rem);
  font-weight: 700;
  letter-spacing: -0.04em;
  color: var(--text);
  margin-bottom: 0.75rem;
}

.contrib-sub {
  color: var(--text-muted);
  font-size: 1rem;
  line-height: 1.7;
  max-width: 500px;
}

/* ─────────────────────────────────────────────
   Stats
───────────────────────────────────────────── */
.stats {
  padding: 3rem 0;
  border-bottom: 1px solid var(--border);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 1px;
  background: var(--border);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  overflow: hidden;
}

.stat-card {
  background: var(--bg-elevated);
  padding: 1.6rem 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 0.3rem;
  transition: background 0.2s;
}

.stat-card:hover {
  background: var(--bg-card);
}

.stat-value {
  font-family: var(--font-mono);
  font-size: 1.65rem;
  font-weight: 700;
  color: var(--text);
  letter-spacing: -0.03em;
}

.stat-label {
  font-size: 0.78rem;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.06em;
  font-weight: 500;
}

/* ─────────────────────────────────────────────
   Contributors grid
───────────────────────────────────────────── */
.contrib-list {
  padding: 4rem 0 5rem;
  border-bottom: 1px solid var(--border);
}

.list-title {
  font-size: 1.1rem;
  font-weight: 600;
  color: var(--text-muted);
  letter-spacing: -0.01em;
  margin-bottom: 2rem;
}

.cards-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1rem;
}

/* ── Contributor card ── */
.contrib-card {
  position: relative;
  display: flex;
  flex-direction: column;
  gap: 1rem;
  padding: 1.5rem;
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  transition: border-color 0.2s, background 0.2s, transform 0.2s, box-shadow 0.2s;
  text-decoration: none;
  overflow: hidden;
}

.contrib-card::before {
  content: '';
  position: absolute;
  inset: 0;
  background: radial-gradient(ellipse 80% 60% at 50% 0%, var(--red-glow) 0%, transparent 70%);
  opacity: 0;
  transition: opacity 0.3s;
}

.contrib-card:hover {
  border-color: var(--red-border);
  background: var(--bg-hover);
  transform: translateY(-3px);
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.3), 0 0 0 1px rgba(200, 40, 56, 0.1);
}

.contrib-card:hover::before {
  opacity: 1;
}

/* Avatar */
.contrib-avatar-wrap {
  position: relative;
  align-self: flex-start;
}

.contrib-avatar {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  border: 2px solid var(--border-md);
  object-fit: cover;
  transition: border-color 0.2s;
}

.contrib-card:hover .contrib-avatar {
  border-color: var(--red-border);
}

.contrib-badge {
  position: absolute;
  bottom: -4px;
  right: -4px;
  font-size: 0.6rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  background: var(--red);
  color: #fff;
  border-radius: 99px;
  padding: 0.15em 0.5em;
  white-space: nowrap;
  border: 2px solid var(--bg-card);
}

/* Info */
.contrib-info {
  display: flex;
  flex-direction: column;
  gap: 0.15rem;
  position: relative;
  z-index: 1;
}

.contrib-name {
  font-size: 1rem;
  font-weight: 600;
  color: var(--text);
}

.contrib-login {
  font-size: 0.8rem;
  color: var(--text-muted);
  font-family: var(--font-mono);
}

/* Tags */
.contrib-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 0.4rem;
  position: relative;
  z-index: 1;
}

.contrib-tag {
  display: inline-flex;
  align-items: center;
  gap: 0.3rem;
  font-size: 0.72rem;
  color: var(--text-muted);
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid var(--border);
  border-radius: 99px;
  padding: 0.2em 0.65em;
  transition: color 0.15s, border-color 0.15s, background 0.15s;
}

.contrib-card:hover .contrib-tag {
  border-color: var(--red-border);
  background: var(--red-dim);
  color: #f07080;
}

/* External arrow */
.contrib-ext {
  position: absolute;
  top: 1.25rem;
  right: 1.25rem;
  color: var(--text-dim);
  opacity: 0;
  transform: translate(-3px, 3px);
  transition: opacity 0.2s, transform 0.2s, color 0.2s;
}

.contrib-card:hover .contrib-ext {
  opacity: 1;
  transform: translate(0, 0);
  color: var(--red);
}

/* ── Ghost / placeholder ── */
.contrib-card--ghost {
  pointer-events: none;
  gap: 1.25rem;
  border-style: dashed;
  background: transparent;
}

.ghost-avatar {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.04);
  border: 2px dashed rgba(255, 255, 255, 0.07);
}

.ghost-lines {
  display: flex;
  flex-direction: column;
  gap: 0.55rem;
}

.ghost-line {
  height: 10px;
  border-radius: 99px;
  background: rgba(255, 255, 255, 0.04);
}

.ghost-line--name  { width: 120px; }
.ghost-line--login { width: 80px; }

/* ─────────────────────────────────────────────
   Join section
───────────────────────────────────────────── */
.join-section {
  position: relative;
  padding: 6rem 0 7rem;
  overflow: hidden;
}

.join-glow {
  position: absolute;
  inset: 0;
  background: radial-gradient(ellipse 70% 80% at 50% 100%, rgba(200, 40, 56, 0.09) 0%, transparent 70%);
  pointer-events: none;
}

.join-inner {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
}

.join-icon {
  width: 72px;
  height: 72px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--red-dim);
  border: 1px solid var(--red-border);
  border-radius: var(--radius-lg);
  color: var(--red);
  margin-bottom: 1.5rem;
}

.join-title {
  font-size: clamp(1.6rem, 3.5vw, 2.4rem);
  font-weight: 700;
  letter-spacing: -0.035em;
  color: var(--text);
  margin-bottom: 0.75rem;
}

.join-sub {
  color: var(--text-muted);
  font-size: 1rem;
  line-height: 1.7;
  max-width: 520px;
  margin-bottom: 2rem;
}

.join-actions {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  flex-wrap: wrap;
  justify-content: center;
  margin-bottom: 4rem;
}

/* How-to steps */
.how-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1px;
  background: var(--border);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  overflow: hidden;
  width: 100%;
  max-width: 780px;
  text-align: left;
}

.how-card {
  background: var(--bg-elevated);
  padding: 1.5rem 1.4rem;
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
  transition: background 0.2s;
}

.how-card:hover {
  background: var(--bg-card);
}

.how-num {
  font-family: var(--font-mono);
  font-size: 0.75rem;
  font-weight: 700;
  color: var(--red);
  letter-spacing: 0.04em;
}

.how-title {
  font-size: 0.9rem;
  font-weight: 600;
  color: var(--text);
}

.how-desc {
  font-size: 0.82rem;
  color: var(--text-muted);
  line-height: 1.6;
  margin: 0;
}

/* ─────────────────────────────────────────────
   Responsive
───────────────────────────────────────────── */
@media (max-width: 900px) {
  .cards-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .how-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 600px) {
  .inner {
    padding: 0 1.25rem;
  }

  .cards-grid {
    grid-template-columns: 1fr;
  }

  .contrib-header {
    padding: 4rem 0 3rem;
  }

  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .join-section {
    padding: 4rem 0 5rem;
  }
}
</style>
