# App
generate:
	@echo "Generating..."
	@go generate ./...

run:
	@set -o allexport; \
	if [ -f .env ]; then source .env; fi; \
	go run cmd/server/main.go

build:
	go build -gcflags=all="-N -l" -o server cmd/server/main.go

# DB
db-up:
	@echo "Starting database..."
	@sh scripts/run-postgres.sh 

db-down:
	@echo "Stopping database..."
	@sh scripts/stop-postgres.sh

db-migrate:
	@set -o allexport; \
	if [ -f .env ]; then source .env; fi; \
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate \
		-database "$$DATABASE_URL" \
		-path db/migrations \
		up

db-queries:
	go tool github.com/sqlc-dev/sqlc/cmd/sqlc generate

wait2:
	sleep 2

db-refresh: db-down db-up wait2 db-migrate

db-admin:
	docker run --rm -d --network jam-schedule --link postgres:db -p 8080:8080 --name adminer adminer

db-admin-down:
	docker stop adminer

# Docker
docker-image:
	@echo "Building docker image..."
	@docker build -t jam-schedule-api .

docker-network:
	@echo "Creating docker network..."
	@docker network create jam-schedule

container:
	@echo "Running docker container..."
	@sh scripts/run-service.sh

container-logs:
	@echo "Fetching docker container logs..."
	@docker logs -f jam-schedule-api

container-down:
	@echo "Stopping docker container..."
	@docker stop jam-schedule-api
	@docker rm jam-schedule-api

container-restart: container-down container
	@echo "Restarting docker container..."
