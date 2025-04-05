generate:
	@echo "Generating..."
	@go generate ./...

run:
	@echo "Running..."
	@go run cmd/server/main.go

dbup:
	@echo "Starting database..."
	@sh scripts/run-postgres.sh 

dbdown:
	@echo "Stopping database..."
	@sh scripts/stop-postgres.sh

db:
	@echo "Migrating..."
	@go run cmd/migrate/main.go

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
