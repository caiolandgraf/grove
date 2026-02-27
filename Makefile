# ──────────────────────────────────────────────────────────────────────────────
# Grove CLI
# ──────────────────────────────────────────────────────────────────────────────

grove-build: ## Build the grove CLI binary → bin/grove
	@mkdir -p bin
	go build -o bin/grove .
	@echo "\033[32m  ✔ Binary built at bin/grove\033[0m"

grove-install: ## Install grove globally (go install → adds to PATH)
	go install .
	@echo "\033[32m  ✔ grove installed — run: grove --help\033[0m"

docs-dev: ## Start the docs site dev server (http://localhost:5173)
	cd docs && npm install --silent && npm run dev

docs-build: ## Build the docs site for production → docs/dist/
	cd docs && npm install --silent && npm run build
	@echo "\033[32m  ✔ Docs built at docs/dist/\033[0m"

docs-preview: ## Preview the production docs build (http://localhost:4173)
	cd docs && npm run preview
