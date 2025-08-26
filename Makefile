.PHONY: db-up db-down db-logs test build run clean

# Database operations
db-up:
	docker compose up -d
	@echo "Waiting for database to be ready..."
	@timeout 30 bash -c 'until docker compose exec postgres pg_isready -U developer -d gamevault_test; do sleep 1; done'
	@echo "Database is ready!"

db-down:
	docker compose down

db-logs:
	docker compose logs postgres

db-reset:
	docker compose down -v
	docker compose up -d

# Application operations  
test:
	cd tests && go test -v

test-race:
	cd tests && go test -race -v

bench:
	cd tests && go test -bench=.

build:
	cd starter && go build -o orders-system .

run:
	cd starter && go run . -db="$$(grep DATABASE_URL ../.env | cut -d= -f2)"

run-enrich:
	cd starter && go run . -enrich -db="$$(grep DATABASE_URL ../.env | cut -d= -f2)"

clean:
	cd starter && rm -f orders-system

# Development workflow
dev: db-up test
	@echo "Development environment ready!"