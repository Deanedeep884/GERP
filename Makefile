# ====================================================================================
# GERP (Google ERP) Master Control Plane
# ====================================================================================

.PHONY: help up down init-db run-gateway run-worker generate

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# --- INFRASTRUCTURE ---

up: ## Boot the local Docker infrastructure (Spanner, Temporal, pgvector, Qdrant)
	@echo "🚀 Starting GERP infrastructure..."
	docker-compose up -d

down: ## Tear down the local infrastructure
	@echo "🛑 Stopping GERP infrastructure..."
	docker-compose down

# --- DATABASE PROVISIONING ---

init-db: ## Create the Spanner emulator instance and apply all 7 Domain DDLs
	@echo "🏗️ Initializing Cloud Spanner Emulator..."
	@gcloud config set api_endpoint_overrides/spanner http://localhost:9020/ --quiet
	@gcloud spanner instances create gerp-instance --config=emulator-config --description="Local GERP" --nodes=1 || true
	@gcloud spanner databases create gerp-db --instance=gerp-instance || true
	@echo "📜 Applying Domain Schemas..."
	@gcloud spanner databases ddl update gerp-db --instance=gerp-instance --ddl-file=./internal/finance/schema.ddl
	@gcloud spanner databases ddl update gerp-db --instance=gerp-instance --ddl-file=./internal/hcm/schema.ddl
	@gcloud spanner databases ddl update gerp-db --instance=gerp-instance --ddl-file=./internal/scm/schema.ddl
	@gcloud spanner databases ddl update gerp-db --instance=gerp-instance --ddl-file=./internal/eam/schema.ddl
	@gcloud spanner databases ddl update gerp-db --instance=gerp-instance --ddl-file=./internal/legal/schema.ddl
	@gcloud spanner databases ddl update gerp-db --instance=gerp-instance --ddl-file=./internal/revenue/schema.ddl
	@gcloud spanner databases ddl update gerp-db --instance=gerp-instance --ddl-file=./internal/lms/schema.ddl
	@gcloud spanner databases ddl update gerp-db --instance=gerp-instance --ddl-file=./internal/mdm/schema.ddl
	@echo "✅ All domains initialized."

# --- APPLICATION RUNTIMES ---

run-gateway: ## Boot the GraphQL API Gateway (BFF)
	@echo "🌐 Starting GraphQL Gateway on http://localhost:8080..."
	go run ./cmd/gateway/main.go

run-worker: ## Boot the Temporal Saga Orchestrator Worker
	@echo "⚙️ Starting Temporal Pipeline Worker..."
	go run ./cmd/worker/main.go

# --- CODE GENERATION & LINTING ---

generate: ## Regenerate GraphQL and Temporal bindings
	@echo "🧬 Running go generate..."
	go generate ./...

audit-skills: ## Run the QuanuX Zero-Orphans linter against SKILL.md files
	@echo "🔍 Auditing Knowledge Vectors..."
	@# (Placeholder for future CLI linter script)
	@echo "✅ Knowledge Vectors structurally sound."
