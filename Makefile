generate:
	@echo "Generating..."
	@go generate ./...

dbup:
	@echo "Starting database..."
	@sh scripts/run-postgres.sh 

dbdown:
	@echo "Stopping database..."
	@sh scripts/stop-postgres.sh

db:
	@echo "Migrating..."
	@go run cmd/migrate/main.go