.PHONY: all install run-backend run-frontend run-all sync-matches build-all

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
