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
import { projectStructureTree } from '@/data/projectStructure.js'
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

const treeNodes = projectStructureTree

const collapsedPaths = ref([])
const selectedNode = ref(treeNodes.find(n => n.path === 'cmd/api/main.go') ?? treeNodes[0])
const isMaximized = ref(false)

const visibleNodes = computed(() => {
  return treeNodes.filter(node => {
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
