<template>
  <div class="code-block">
    <div class="code-block__header">
      <span class="code-block__label">{{ label }}</span>
      <button
        class="code-block__copy"
        :class="{ copied }"
        @click="copy"
        :aria-label="copied ? 'Copied!' : 'Copy code'"
      >
        <svg v-if="!copied" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
          <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
        </svg>
        <svg v-else width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
          <polyline points="20 6 9 17 4 12"/>
        </svg>
        {{ copied ? 'Copied' : 'Copy' }}
      </button>
    </div>
    <div class="code-block__body">
      <pre><code v-html="highlighted" :class="`language-${lang}`" /></pre>
    </div>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import hljs from 'highlight.js/lib/core'
import go from 'highlight.js/lib/languages/go'
import bash from 'highlight.js/lib/languages/bash'
import sql from 'highlight.js/lib/languages/sql'
import json from 'highlight.js/lib/languages/json'
import yaml from 'highlight.js/lib/languages/yaml'

hljs.registerLanguage('go', go)
hljs.registerLanguage('bash', bash)
hljs.registerLanguage('shell', bash)
hljs.registerLanguage('sql', sql)
hljs.registerLanguage('json', json)
hljs.registerLanguage('yaml', yaml)

const props = defineProps({
  code: { type: String, required: true },
  lang: { type: String, default: 'bash' },
  label: { type: String, default: '' },
})

const copied = ref(false)

const highlighted = computed(() => {
  try {
    const lang = hljs.getLanguage(props.lang) ? props.lang : 'plaintext'
    return hljs.highlight(props.code, { language: lang }).value
  } catch {
    return props.code
  }
})

async function copy() {
  try {
    await navigator.clipboard.writeText(props.code)
    copied.value = true
    setTimeout(() => { copied.value = false }, 2000)
  } catch {
    // fallback
    const ta = document.createElement('textarea')
    ta.value = props.code
    document.body.appendChild(ta)
    ta.select()
    document.execCommand('copy')
    document.body.removeChild(ta)
    copied.value = true
    setTimeout(() => { copied.value = false }, 2000)
  }
}
</script>

<style scoped>
.code-block {
  border: 1px solid var(--border-md);
  border-radius: var(--radius);
  overflow: hidden;
  margin: 1.2em 0;
  background: #0a0a12;
}

.code-block__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.45em 1em;
  background: rgba(255, 255, 255, 0.03);
  border-bottom: 1px solid var(--border);
}

.code-block__label {
  font-family: var(--font-mono);
  font-size: 0.72rem;
  color: var(--text-muted);
  letter-spacing: 0.03em;
}

.code-block__copy {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 0.72rem;
  color: var(--text-muted);
  padding: 0.2em 0.55em;
  border-radius: 4px;
  transition: color 0.15s, background 0.15s;
}

.code-block__copy:hover {
  color: var(--text);
  background: rgba(255, 255, 255, 0.06);
}

.code-block__copy.copied {
  color: #4ade80;
}

.code-block__body {
  overflow-x: auto;
}

.code-block__body pre {
  margin: 0;
  padding: 1.1em 1.3em;
  background: transparent;
}

.code-block__body code {
  font-family: var(--font-mono);
  font-size: 0.845rem;
  line-height: 1.7;
  background: transparent;
  color: #c8cce8;
  border: none;
  padding: 0;
}

/* ── highlight.js token overrides ── */
:deep(.hljs-keyword)   { color: #c82838; font-style: italic; }
:deep(.hljs-built_in)  { color: #e06c75; }
:deep(.hljs-type)      { color: #e5c07b; }
:deep(.hljs-string)    { color: #98c379; }
:deep(.hljs-number)    { color: #d19a66; }
:deep(.hljs-comment)   { color: #4a4e6a; font-style: italic; }
:deep(.hljs-function)  { color: #61afef; }
:deep(.hljs-title)     { color: #61afef; }
:deep(.hljs-params)    { color: #c8cce8; }
:deep(.hljs-attr)      { color: #e06c75; }
:deep(.hljs-variable)  { color: #c8cce8; }
:deep(.hljs-operator)  { color: #56b6c2; }
:deep(.hljs-punctuation){ color: #7a7ea8; }
:deep(.hljs-meta)      { color: #c82838; }
:deep(.hljs-literal)   { color: #d19a66; }
:deep(.hljs-symbol)    { color: #56b6c2; }
:deep(.hljs-deletion)  { color: #e06c75; background: rgba(224,108,117,0.12); }
:deep(.hljs-addition)  { color: #98c379; background: rgba(152,195,121,0.10); }
:deep(.hljs-section)   { color: #e5c07b; }
:deep(.hljs-name)      { color: #e06c75; }
:deep(.hljs-selector-tag)  { color: #c82838; }
:deep(.hljs-selector-id)   { color: #61afef; }
:deep(.hljs-selector-class){ color: #e5c07b; }
</style>
