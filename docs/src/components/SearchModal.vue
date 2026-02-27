<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="open" class="overlay" @click.self="close">
        <div class="modal" role="dialog" aria-modal="true" aria-label="Search">

          <!-- Input -->
          <div class="modal__search">
            <svg class="modal__icon" width="16" height="16" viewBox="0 0 24 24" fill="none"
              stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/>
            </svg>
            <input
              ref="inputRef"
              v-model="query"
              type="text"
              placeholder="Search docs..."
              class="modal__input"
              autocomplete="off"
              spellcheck="false"
              @keydown.down.prevent="moveDown"
              @keydown.up.prevent="moveUp"
              @keydown.enter.prevent="confirm"
              @keydown.esc="close"
            />
            <kbd class="modal__esc" @click="close">esc</kbd>
          </div>

          <!-- Results -->
          <div class="modal__body" ref="listRef">
            <!-- Empty query → show categories -->
            <template v-if="!query.trim()">
              <div class="modal__section-label">Quick navigation</div>
              <div
                v-for="(link, i) in quickLinks"
                :key="link.url"
                class="modal__item"
                :class="{ 'modal__item--active': cursor === i }"
                @mouseenter="cursor = i"
                @click="go(link.url)"
              >
                <span class="modal__item-icon">{{ link.icon }}</span>
                <span class="modal__item-text">{{ link.title }}</span>
                <span class="modal__item-section">{{ link.section }}</span>
              </div>
            </template>

            <!-- Query → show fuzzy results -->
            <template v-else-if="results.length">
              <div
                v-for="(r, i) in results"
                :key="r.item.id"
                class="modal__item"
                :class="{ 'modal__item--active': cursor === i }"
                @mouseenter="cursor = i"
                @click="go(r.item.url)"
              >
                <span class="modal__item-icon">⬡</span>
                <div class="modal__item-body">
                  <span class="modal__item-text">{{ r.item.title }}</span>
                  <span class="modal__item-section">{{ r.item.section }}</span>
                </div>
                <svg class="modal__item-arrow" width="12" height="12" viewBox="0 0 24 24"
                  fill="none" stroke="currentColor" stroke-width="2.5"
                  stroke-linecap="round" stroke-linejoin="round">
                  <line x1="5" y1="12" x2="19" y2="12"/>
                  <polyline points="12 5 19 12 12 19"/>
                </svg>
              </div>
            </template>

            <!-- No results -->
            <div v-else class="modal__empty">
              <svg width="32" height="32" viewBox="0 0 24 24" fill="none"
                stroke="currentColor" stroke-width="1.5"
                stroke-linecap="round" stroke-linejoin="round">
                <circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/>
              </svg>
              <p>No results for <strong>"{{ query }}"</strong></p>
              <span>Try searching for a command, concept or configuration key.</span>
            </div>
          </div>

          <!-- Footer -->
          <div class="modal__footer">
            <span class="modal__hint"><kbd>↑</kbd><kbd>↓</kbd> navigate</span>
            <span class="modal__hint"><kbd>↵</kbd> open</span>
            <span class="modal__hint"><kbd>esc</kbd> close</span>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { ref, computed, watch, nextTick, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import Fuse from 'fuse.js'
import { searchIndex } from '@/data/docs.js'

const router = useRouter()

const open  = ref(false)
const query = ref('')
const cursor = ref(0)
const inputRef = ref(null)
const listRef  = ref(null)

// ── Fuse instance ──────────────────────────────
const fuse = new Fuse(searchIndex, {
  keys: [
    { name: 'title',   weight: 0.6 },
    { name: 'section', weight: 0.2 },
    { name: 'text',    weight: 0.2 },
  ],
  threshold: 0.35,
  includeScore: true,
  minMatchCharLength: 2,
})

const results = computed(() =>
  query.value.trim() ? fuse.search(query.value.trim()).slice(0, 9) : []
)

// ── Quick links (shown when query is empty) ────
const quickLinks = [
  { title: 'Getting Started',     section: 'Docs',          icon: '🚀', url: '/docs#getting-started' },
  { title: 'Quick Start',         section: 'Docs',          icon: '⚡', url: '/docs#quick-start' },
  { title: 'Project Structure',   section: 'Architecture',  icon: '🗂',  url: '/docs#project-structure' },
  { title: 'grove make:resource', section: 'Commands',      icon: '⬡',  url: '/docs#cmd-make-resource' },
  { title: 'grove migrate',       section: 'Commands',      icon: '🔄', url: '/docs#cmd-migrate' },
  { title: 'grove serve',         section: 'Commands',      icon: '🌐', url: '/docs#cmd-serve' },
  { title: 'Shell Completion',    section: 'Commands',      icon: '✦',  url: '/docs#cmd-completion' },
  { title: 'Atlas Configuration', section: 'Configuration', icon: '⚙️', url: '/docs#atlas-config' },
]

const activeList = computed(() =>
  query.value.trim() ? results.value : quickLinks
)

// ── Keyboard navigation ────────────────────────
function moveDown() {
  cursor.value = (cursor.value + 1) % activeList.value.length
  scrollActivIntoView()
}

function moveUp() {
  cursor.value = (cursor.value - 1 + activeList.value.length) % activeList.value.length
  scrollActivIntoView()
}

function scrollActivIntoView() {
  nextTick(() => {
    const el = listRef.value?.querySelector('.modal__item--active')
    el?.scrollIntoView({ block: 'nearest' })
  })
}

function confirm() {
  const item = activeList.value[cursor.value]
  if (!item) return
  const url = item.url ?? item.item?.url
  if (url) go(url)
}

function go(url) {
  close()
  if (url.startsWith('/')) {
    const [path, hash] = url.split('#')
    router.push({ path, hash: hash ? `#${hash}` : undefined })
  } else {
    window.open(url, '_blank')
  }
}

// ── Open / close ──────────────────────────────
function openModal() {
  open.value = true
  query.value = ''
  cursor.value = 0
  nextTick(() => inputRef.value?.focus())
}

function close() {
  open.value = false
  query.value = ''
}

// Reset cursor when results change
watch(results, () => { cursor.value = 0 })
watch(query,   () => { cursor.value = 0 })

// ── Global keyboard shortcut ──────────────────
function onKeydown(e) {
  // Ctrl+K or Cmd+K
  if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
    e.preventDefault()
    open.value ? close() : openModal()
  }
  if (e.key === 'Escape' && open.value) {
    close()
  }
}

onMounted(() => window.addEventListener('keydown', onKeydown))
onUnmounted(() => window.removeEventListener('keydown', onKeydown))

// expose openModal so AppNav can call it
defineExpose({ openModal })
</script>

<style scoped>
/* ── Overlay ── */
.overlay {
  position: fixed;
  inset: 0;
  z-index: 1000;
  background: rgba(6, 6, 10, 0.75);
  backdrop-filter: blur(6px);
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding-top: min(12vh, 100px);
}

/* ── Modal box ── */
.modal {
  width: 100%;
  max-width: 620px;
  background: #0f0f1a;
  border: 1px solid var(--border-md);
  border-radius: var(--radius-lg);
  box-shadow:
    0 0 0 1px rgba(200, 40, 56, 0.08),
    0 24px 80px rgba(0, 0, 0, 0.7),
    0 0 60px rgba(200, 40, 56, 0.06);
  overflow: hidden;
  display: flex;
  flex-direction: column;
  max-height: calc(100vh - 200px);
}

/* ── Search row ── */
.modal__search {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0 1.1rem;
  border-bottom: 1px solid var(--border);
  flex-shrink: 0;
}

.modal__icon {
  color: var(--text-muted);
  flex-shrink: 0;
}

.modal__input {
  flex: 1;
  background: transparent;
  border: none;
  outline: none;
  color: var(--text);
  font-size: 1rem;
  padding: 1rem 0;
  caret-color: var(--red);
}

.modal__input::placeholder {
  color: var(--text-dim);
}

.modal__esc {
  flex-shrink: 0;
  font-family: var(--font-mono);
  font-size: 0.68rem;
  color: var(--text-dim);
  background: rgba(255,255,255,0.05);
  border: 1px solid var(--border-md);
  border-radius: 4px;
  padding: 0.15em 0.45em;
  cursor: pointer;
  transition: color 0.15s, background 0.15s;
}

.modal__esc:hover {
  color: var(--text-muted);
  background: rgba(255,255,255,0.08);
}

/* ── Body / results ── */
.modal__body {
  overflow-y: auto;
  flex: 1;
  padding: 0.5rem;
}

.modal__section-label {
  font-size: 0.68rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  color: var(--text-dim);
  padding: 0.6rem 0.75rem 0.35rem;
}

/* ── Result item ── */
.modal__item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.65rem 0.75rem;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: background 0.1s;
}

.modal__item:hover,
.modal__item--active {
  background: var(--bg-hover);
}

.modal__item--active {
  background: var(--red-dim);
  outline: 1px solid var(--red-border);
}

.modal__item-icon {
  font-size: 0.9rem;
  flex-shrink: 0;
  width: 20px;
  text-align: center;
  color: var(--text-muted);
}

.modal__item-body {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 0.75rem;
  min-width: 0;
}

.modal__item-text {
  font-size: 0.9rem;
  color: var(--text);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  flex: 1;
}

.modal__item--active .modal__item-text {
  color: #fff;
}

.modal__item-section {
  font-size: 0.72rem;
  color: var(--text-dim);
  background: rgba(255,255,255,0.04);
  border: 1px solid var(--border);
  border-radius: 4px;
  padding: 0.1em 0.5em;
  white-space: nowrap;
  flex-shrink: 0;
}

.modal__item--active .modal__item-section {
  background: rgba(200, 40, 56, 0.15);
  border-color: var(--red-border);
  color: #f07080;
}

.modal__item-arrow {
  color: var(--text-dim);
  flex-shrink: 0;
  opacity: 0;
  transform: translateX(-4px);
  transition: opacity 0.15s, transform 0.15s;
}

.modal__item--active .modal__item-arrow {
  opacity: 1;
  transform: translateX(0);
  color: var(--red);
}

/* ── Empty ── */
.modal__empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
  padding: 3rem 1rem;
  color: var(--text-dim);
  text-align: center;
}

.modal__empty svg {
  opacity: 0.3;
  margin-bottom: 0.25rem;
}

.modal__empty p {
  font-size: 0.9rem;
  color: var(--text-muted);
  margin: 0;
}

.modal__empty strong {
  color: var(--text);
}

.modal__empty span {
  font-size: 0.78rem;
}

/* ── Footer ── */
.modal__footer {
  display: flex;
  align-items: center;
  gap: 1.25rem;
  padding: 0.6rem 1.1rem;
  border-top: 1px solid var(--border);
  background: rgba(255,255,255,0.02);
  flex-shrink: 0;
}

.modal__hint {
  display: flex;
  align-items: center;
  gap: 0.3rem;
  font-size: 0.7rem;
  color: var(--text-dim);
}

.modal__hint kbd {
  font-family: var(--font-mono);
  font-size: 0.65rem;
  background: rgba(255,255,255,0.06);
  border: 1px solid var(--border-md);
  border-radius: 3px;
  padding: 0.1em 0.4em;
  color: var(--text-muted);
}

/* ── Transition ── */
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.18s var(--ease);
}

.modal-enter-active .modal,
.modal-leave-active .modal {
  transition: transform 0.2s var(--ease), opacity 0.18s var(--ease);
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .modal,
.modal-leave-to .modal {
  transform: translateY(-12px) scale(0.97);
  opacity: 0;
}

/* ── Mobile ── */
@media (max-width: 660px) {
  .overlay {
    align-items: flex-end;
    padding-top: 0;
  }

  .modal {
    border-radius: var(--radius-lg) var(--radius-lg) 0 0;
    max-height: 80vh;
  }

  .modal-enter-from .modal,
  .modal-leave-to .modal {
    transform: translateY(20px);
  }
}
</style>
