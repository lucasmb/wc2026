.PHONY: all install run-backend run-frontend run-all sync-matches build-all
ENV_DIR = .
# Load it (use absolute path to be safe)
-include $(shell pwd)/$(ENV_DIR)/.env

BINARY_NAME = wc2026-api

all: run-all

# Install dependencies for both frontend and backend

# Run Go Pocketbase backend server
run-backend:
	@echo "Starting Pocketbase backend serve..."
	cd backend && go run main.go serve --http="127.0.0.1:8090"

# Run Quasar frontend dev server
run-frontend:
	@echo "Starting Quasar frontend dev server..."
	pnpm --filter frontend dev || pnpm -C frontend dev

# Run both frontend and backend concurrently in local terminal
run-all:
	@echo "Starting both services concurrently..."
	@zsh -c " \
		trap 'kill 0' INT; \
		(cd backend && go run main.go serve --http='127.0.0.1:8090') & \
		(pnpm --filter frontend dev || pnpm -C frontend dev) & \
		wait \
	"

# Trigger backend match sync endpoint
sync-matches:
	@echo "Triggering World Cup 2026 match synchronizer..."
	curl -X POST http://127.0.0.1:8090/api/wc2026/sync

# Build all applications for production
build-all:
	@echo "Building Go backend..."
	cd backend && go build -o wc2026-backend main.go
	@echo "Building Quasar frontend..."
	pnpm --filter frontend build || pnpm -C frontend build

build:
	cd backend && GOOS=linux GOARCH=amd64 go build -o ./bin/$(BINARY_NAME) .

.PHONY: restart-pb-remote
restart-remote:
	@echo $(ENV_DIR)
	@echo "Restarting PocketBase service on $(REMOTE_HOST)..."
	@ssh -t ""$$REMOTE_USER@$$REMOTE_HOST"" "sudo systemctl restart $$PB_SERVICE_NAME"
	@echo "PocketBase service restart command sent to $(REMOTE_HOST)."


.PHONY: upload
upload:
	@rsync -avz \
		--chown=:www \
		-e "ssh -p $(REMOTE_PORT)" \
		backend/bin/$(BINARY_NAME) \
		$(REMOTE_USER)@$(REMOTE_HOST):$(REMOTE_PATH)
	@ssh -p $(REMOTE_PORT) -t $(REMOTE_USER)@$(REMOTE_HOST) "sudo systemctl restart calometrics-api.service"
