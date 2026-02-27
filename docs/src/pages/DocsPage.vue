<template>
  <div class="docs-layout">

    <!-- ─── Sidebar ──────────────────────────── -->
    <aside class="sidebar" :class="{ open: sidebarOpen }">
      <div class="sidebar-inner">
        <nav class="sidebar-nav" aria-label="Documentation">
          <div v-for="section in sections" :key="section.id" class="sidebar-section">
            <span class="sidebar-section-title">{{ section.title }}</span>
            <ul>
              <li v-for="item in section.items" :key="item.id">
                <a
                  :href="`#${item.id}`"
                  class="sidebar-link"
                  :class="{ active: activeId === item.id }"
                  @click.prevent="scrollTo(item.id)"
                >
                  {{ item.title }}
                </a>
              </li>
            </ul>
          </div>
        </nav>
      </div>
    </aside>

    <!-- ─── Mobile sidebar toggle ────────────── -->
    <button class="sidebar-toggle" @click="sidebarOpen = !sidebarOpen" :aria-expanded="sidebarOpen">
      <svg v-if="!sidebarOpen" width="18" height="18" viewBox="0 0 24 24" fill="none"
        stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <line x1="3" y1="6" x2="21" y2="6"/>
        <line x1="3" y1="12" x2="15" y2="12"/>
        <line x1="3" y1="18" x2="18" y2="18"/>
      </svg>
      <svg v-else width="18" height="18" viewBox="0 0 24 24" fill="none"
        stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/>
      </svg>
      <span>{{ sidebarOpen ? 'Close' : 'Menu' }}</span>
    </button>

    <!-- ─── Overlay (mobile) ──────────────────── -->
    <div v-if="sidebarOpen" class="sidebar-overlay" @click="sidebarOpen = false" />

    <!-- ─── Main content ─────────────────────── -->
    <main class="docs-main" ref="mainRef">
      <div class="docs-content prose">
        <template v-for="section in sections" :key="section.id">

          <!-- Section anchor -->
          <div :id="section.id" class="section-anchor" />
          <h1 class="section-heading">{{ section.title }}</h1>

          <!-- Items -->
          <template v-for="item in section.items" :key="item.id">
            <div :id="item.id" class="item-anchor" />

            <h2 class="item-heading">
              <a :href="`#${item.id}`" class="anchor-link" @click.prevent="scrollTo(item.id)" aria-hidden="true">#</a>
              {{ item.title }}
            </h2>

            <!-- Blocks -->
            <template v-for="(block, bi) in item.blocks" :key="bi">

              <!-- Paragraph -->
              <p v-if="block.type === 'paragraph'" v-html="block.text" />

              <!-- Code block -->
              <CodeBlock
                v-else-if="block.type === 'code'"
                :code="block.code"
                :lang="block.lang"
                :label="block.label"
              />

              <!-- Note / callout -->
              <div v-else-if="block.type === 'note'"
                class="callout"
                :class="`callout--${block.kind}`"
              >
                <span class="callout-icon">
                  <span v-if="block.kind === 'tip'">💡</span>
                  <span v-else-if="block.kind === 'warning'">⚠️</span>
                  <span v-else>ℹ️</span>
                </span>
                <span v-html="block.text" />
              </div>

              <!-- Steps list -->
              <ol v-else-if="block.type === 'steps'" class="steps">
                <li v-for="(step, si) in block.items" :key="si" class="step">
                  <span class="step-num">{{ si + 1 }}</span>
                  <div class="step-body">
                    <strong class="step-title">{{ step.title }}</strong>
                    <span v-html="step.text" />
                  </div>
                </li>
              </ol>

              <!-- Table -->
              <div v-else-if="block.type === 'table'" class="table-wrap">
                <table>
                  <thead>
                    <tr>
                      <th v-for="h in block.head" :key="h" v-html="h" />
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="(row, ri) in block.rows" :key="ri">
                      <td v-for="(cell, ci) in row" :key="ci" v-html="cell" />
                    </tr>
                  </tbody>
                </table>
              </div>

            </template>

            <div class="item-divider" />
          </template>

        </template>
      </div>
    </main>

    <!-- ─── Right TOC ────────────────────────── -->
    <aside class="toc" aria-label="On this page">
      <div class="toc-inner">
        <span class="toc-title">On this page</span>
        <nav>
          <template v-for="section in sections" :key="section.id">
            <div class="toc-section">
              <a
                :href="`#${section.id}`"
                class="toc-section-link"
                :class="{ active: activeId === section.id }"
                @click.prevent="scrollTo(section.id)"
              >
                {{ section.title }}
              </a>
              <ul>
                <li v-for="item in section.items" :key="item.id">
                  <a
                    :href="`#${item.id}`"
                    class="toc-link"
                    :class="{ active: activeId === item.id }"
                    @click.prevent="scrollTo(item.id)"
                  >
                    {{ item.title }}
                  </a>
                </li>
              </ul>
            </div>
          </template>
        </nav>
      </div>
    </aside>

  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { sections } from '@/data/docs.js'
import CodeBlock from '@/components/CodeBlock.vue'

const route      = useRoute()
const mainRef    = ref(null)
const activeId   = ref('')
const sidebarOpen = ref(false)

// ── Scroll to anchor ──────────────────────────
function scrollTo(id) {
  sidebarOpen.value = false
  const el = document.getElementById(id)
  if (!el) return
  const top = el.getBoundingClientRect().top + window.scrollY - 80
  window.scrollTo({ top, behavior: 'smooth' })
  activeId.value = id
  history.replaceState(null, '', `#${id}`)
}

// ── Intersection observer for active item ─────
let observer = null

function buildObserver() {
  const allIds = sections.flatMap(s => [s.id, ...s.items.map(i => i.id)])

  observer = new IntersectionObserver(
    entries => {
      for (const entry of entries) {
        if (entry.isIntersecting) {
          activeId.value = entry.target.id
        }
      }
    },
    { rootMargin: '-80px 0px -60% 0px', threshold: 0 }
  )

  allIds.forEach(id => {
    const el = document.getElementById(id)
    if (el) observer.observe(el)
  })
}

// ── Handle hash on mount ──────────────────────
onMounted(async () => {
  await nextTick()
  buildObserver()

  const hash = route.hash?.replace('#', '') || window.location.hash?.replace('#', '')
  if (hash) {
    scrollTo(hash)
  } else if (sections[0]) {
    activeId.value = sections[0].id
  }
})

onUnmounted(() => {
  observer?.disconnect()
})
</script>

<style scoped>
/* ─────────────────────────────────────────────
   Layout
───────────────────────────────────────────── */
.docs-layout {
  display: grid;
  grid-template-columns: 260px minmax(0, 1fr) 220px;
  grid-template-areas: "sidebar main toc";
  min-height: calc(100vh - var(--nav-h));
  padding-top: var(--nav-h);
  max-width: 1400px;
  margin: 0 auto;
  align-items: start;
}

/* ─────────────────────────────────────────────
   Sidebar
───────────────────────────────────────────── */
.sidebar {
  grid-area: sidebar;
  position: sticky;
  top: var(--nav-h);
  height: calc(100vh - var(--nav-h));
  overflow-y: auto;
  border-right: 1px solid var(--border);
  padding: 2rem 0;
}

.sidebar-inner {
  padding: 0 1rem;
}

.sidebar-section {
  margin-bottom: 1.75rem;
}

.sidebar-section-title {
  display: block;
  font-size: 0.68rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  color: var(--text-muted);
  padding: 0 0.6rem;
  margin-bottom: 0.4rem;
}

.sidebar-section ul {
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.sidebar-link {
  display: block;
  font-size: 0.85rem;
  color: var(--text-muted);
  padding: 0.38em 0.6em;
  border-radius: var(--radius-sm);
  transition: color 0.15s, background 0.15s;
  border-left: 2px solid transparent;
}

.sidebar-link:hover {
  color: var(--text);
  background: rgba(255, 255, 255, 0.04);
}

.sidebar-link.active {
  color: var(--red-hover);
  background: var(--red-dim);
  border-left-color: var(--red);
}

/* Mobile toggle */
.sidebar-toggle {
  display: none;
  position: fixed;
  bottom: 1.5rem;
  left: 1.5rem;
  z-index: 90;
  align-items: center;
  gap: 0.4rem;
  font-size: 0.8rem;
  font-weight: 500;
  color: var(--text);
  background: var(--bg-card);
  border: 1px solid var(--border-md);
  border-radius: 99px;
  padding: 0.5em 1em 0.5em 0.75em;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.4);
  transition: background 0.15s, border-color 0.15s;
}

.sidebar-toggle:hover {
  background: var(--bg-hover);
  border-color: var(--red-border);
}

.sidebar-overlay {
  display: none;
  position: fixed;
  inset: 0;
  z-index: 79;
  background: rgba(6, 6, 10, 0.7);
  backdrop-filter: blur(4px);
}

/* ─────────────────────────────────────────────
   Main content
───────────────────────────────────────────── */
.docs-main {
  grid-area: main;
  min-width: 0;
  padding: 3rem 3rem 6rem;
  border-right: 1px solid var(--border);
}

.docs-content {
  max-width: 720px;
}

/* Section heading */
.section-heading {
  font-size: 1.75rem;
  font-weight: 700;
  letter-spacing: -0.03em;
  margin: 3rem 0 0.5rem;
  padding-bottom: 0.5rem;
  border-bottom: 1px solid var(--border);
  color: var(--text);
}

.docs-content > .section-anchor:first-child + .section-heading {
  margin-top: 0;
}

/* Item heading */
.item-heading {
  position: relative;
  font-size: 1.2rem;
  font-weight: 600;
  letter-spacing: -0.02em;
  margin: 2.5rem 0 0.75rem;
  color: var(--text);
}

.anchor-link {
  position: absolute;
  left: -1.4em;
  top: 50%;
  transform: translateY(-50%);
  color: var(--text-dim);
  font-weight: 400;
  font-size: 1rem;
  opacity: 0;
  transition: opacity 0.15s, color 0.15s;
  text-decoration: none;
}

.item-heading:hover .anchor-link {
  opacity: 1;
  color: var(--red);
}

/* Anchors (invisible scroll targets) */
.section-anchor,
.item-anchor {
  display: block;
  height: 0;
  margin-top: -88px;
  padding-top: 88px;
  pointer-events: none;
  visibility: hidden;
}

/* Item divider */
.item-divider {
  height: 1px;
  background: var(--border);
  margin: 2.5rem 0;
}

/* ─────────────────────────────────────────────
   Callouts
───────────────────────────────────────────── */
.callout {
  display: flex;
  align-items: flex-start;
  gap: 0.75rem;
  padding: 0.85rem 1rem;
  border-radius: var(--radius-sm);
  border-left: 3px solid;
  margin: 1.2em 0;
  font-size: 0.875rem;
  line-height: 1.65;
}

.callout-icon {
  font-size: 0.9rem;
  flex-shrink: 0;
  margin-top: 0.1em;
}

.callout--info {
  background: rgba(80, 160, 255, 0.07);
  border-left-color: #5090ff;
  color: rgba(226, 228, 239, 0.78);
}

.callout--tip {
  background: rgba(74, 222, 128, 0.07);
  border-left-color: #4ade80;
  color: rgba(226, 228, 239, 0.78);
}

.callout--warning {
  background: rgba(250, 200, 40, 0.07);
  border-left-color: #fac828;
  color: rgba(226, 228, 239, 0.78);
}

.callout :deep(a) {
  color: var(--red-hover);
  text-decoration: underline;
  text-underline-offset: 3px;
}

.callout :deep(code) {
  font-family: var(--font-mono);
  font-size: 0.82em;
  background: rgba(255, 255, 255, 0.08);
  padding: 0.1em 0.35em;
  border-radius: 3px;
}

.callout :deep(strong) {
  color: var(--text);
  font-weight: 600;
}

/* ─────────────────────────────────────────────
   Steps
───────────────────────────────────────────── */
.steps {
  display: flex;
  flex-direction: column;
  gap: 0;
  margin: 1.2em 0;
  padding-left: 0;
  list-style: none;
}

.step {
  display: flex;
  gap: 1rem;
  position: relative;
  padding-bottom: 1.4rem;
}

.step:last-child {
  padding-bottom: 0;
}

.step:not(:last-child)::before {
  content: '';
  position: absolute;
  left: 15px;
  top: 32px;
  bottom: 0;
  width: 1px;
  background: var(--border-md);
}

.step-num {
  flex-shrink: 0;
  width: 30px;
  height: 30px;
  border-radius: 50%;
  background: var(--red-dim);
  border: 1px solid var(--red-border);
  color: var(--red-hover);
  font-size: 0.78rem;
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: var(--font-mono);
}

.step-body {
  padding-top: 0.3rem;
  font-size: 0.875rem;
  color: rgba(226, 228, 239, 0.78);
  line-height: 1.65;
}

.step-title {
  display: block;
  color: var(--text);
  font-weight: 600;
  margin-bottom: 0.2em;
}

.step-body :deep(code) {
  font-family: var(--font-mono);
  font-size: 0.82em;
  background: rgba(200, 40, 56, 0.10);
  color: #f07080;
  padding: 0.15em 0.45em;
  border-radius: 4px;
  border: 1px solid rgba(200, 40, 56, 0.18);
}

/* ─────────────────────────────────────────────
   Table
───────────────────────────────────────────── */
.table-wrap {
  overflow-x: auto;
  margin: 1.2em 0;
  border: 1px solid var(--border);
  border-radius: var(--radius);
}

.table-wrap table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.875rem;
  margin: 0;
}

.table-wrap th {
  text-align: left;
  padding: 0.65em 1em;
  border-bottom: 1px solid var(--border-md);
  color: var(--text-muted);
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  background: rgba(255, 255, 255, 0.02);
  white-space: nowrap;
}

.table-wrap td {
  padding: 0.65em 1em;
  border-bottom: 1px solid var(--border);
  color: rgba(226, 228, 239, 0.78);
  vertical-align: middle;
}

.table-wrap tr:last-child td {
  border-bottom: none;
}

.table-wrap tr:hover td {
  background: rgba(255, 255, 255, 0.02);
}

.table-wrap :deep(code) {
  font-family: var(--font-mono);
  font-size: 0.82em;
  background: rgba(200, 40, 56, 0.10);
  color: #f07080;
  padding: 0.12em 0.4em;
  border-radius: 4px;
  border: 1px solid rgba(200, 40, 56, 0.18);
  white-space: nowrap;
}

.table-wrap :deep(a) {
  color: var(--red-hover);
  text-decoration: underline;
  text-underline-offset: 3px;
}

/* ─────────────────────────────────────────────
   TOC
───────────────────────────────────────────── */
.toc {
  grid-area: toc;
  position: sticky;
  top: var(--nav-h);
  height: calc(100vh - var(--nav-h));
  overflow-y: auto;
  padding: 2.5rem 1.25rem;
}

.toc-title {
  display: block;
  font-size: 0.68rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  color: var(--text-muted);
  margin-bottom: 0.75rem;
  padding: 0 0.4rem;
}

.toc-section {
  margin-bottom: 0.75rem;
}

.toc-section-link {
  display: block;
  font-size: 0.78rem;
  font-weight: 600;
  color: var(--text-muted);
  padding: 0.25em 0.4em;
  border-radius: 4px;
  transition: color 0.15s;
  margin-bottom: 0.1rem;
}

.toc-section-link:hover {
  color: var(--text);
}

.toc-section-link.active {
  color: var(--red);
}

.toc-section ul {
  padding-left: 0.6rem;
  border-left: 1px solid var(--border);
  display: flex;
  flex-direction: column;
}

.toc-link {
  display: block;
  font-size: 0.78rem;
  color: var(--text-dim);
  padding: 0.25em 0.6em;
  border-radius: 4px;
  transition: color 0.15s, background 0.15s;
  border-left: 2px solid transparent;
  margin-left: -1px;
}

.toc-link:hover {
  color: var(--text-muted);
}

.toc-link.active {
  color: var(--red-hover);
  border-left-color: var(--red);
}

/* ─────────────────────────────────────────────
   Prose overrides (inside docs-content)
───────────────────────────────────────────── */
.docs-content :deep(p) {
  font-size: 0.9375rem;
  color: rgba(226, 228, 239, 0.78);
  margin-bottom: 0.9em;
  line-height: 1.7;
}

.docs-content :deep(code:not(pre code)) {
  font-family: var(--font-mono);
  font-size: 0.82em;
  background: rgba(200, 40, 56, 0.10);
  color: #f07080;
  padding: 0.15em 0.45em;
  border-radius: 4px;
  border: 1px solid rgba(200, 40, 56, 0.18);
}

.docs-content :deep(strong) {
  color: var(--text);
  font-weight: 600;
}

.docs-content :deep(a:not(.anchor-link)) {
  color: var(--red-hover);
  text-decoration: underline;
  text-decoration-color: var(--red-border);
  text-underline-offset: 3px;
  transition: color 0.15s;
}

.docs-content :deep(a:not(.anchor-link):hover) {
  color: #ff5060;
}

.docs-content :deep(kbd) {
  font-family: var(--font-mono);
  font-size: 0.8em;
  background: var(--bg-hover);
  border: 1px solid var(--border-md);
  border-radius: 4px;
  padding: 0.1em 0.4em;
  color: var(--text-muted);
}

/* ─────────────────────────────────────────────
   Responsive
───────────────────────────────────────────── */
@media (max-width: 1200px) {
  .docs-layout {
    grid-template-columns: 240px minmax(0, 1fr);
    grid-template-areas:
      "sidebar main";
  }

  .toc { display: none; }
}

@media (max-width: 768px) {
  .docs-layout {
    grid-template-columns: 1fr;
    grid-template-areas: "main";
  }

  .sidebar {
    position: fixed;
    top: var(--nav-h);
    left: 0;
    width: 280px;
    height: calc(100vh - var(--nav-h));
    z-index: 80;
    background: var(--bg-elevated);
    transform: translateX(-100%);
    transition: transform 0.26s var(--ease);
    border-right: 1px solid var(--border-md);
  }

  .sidebar.open {
    transform: translateX(0);
  }

  .sidebar-toggle { display: flex; }
  .sidebar-overlay { display: block; }

  .docs-main {
    padding: 2rem 1.25rem 5rem;
    border-right: none;
  }

  .item-heading {
    font-size: 1.1rem;
  }

  .anchor-link { display: none; }
}
</style>
