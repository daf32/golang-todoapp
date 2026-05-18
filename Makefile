include .env
export

export PROJECT_ROOT=$(shell pwd)

# ============
# Enviroments
# ============

env-up:
	@docker compose up -d todoapp-postgres

env-down:
	@docker compose down todoapp-postgres

env-cleanup:
	@read -p "Clear all volume files in the environment? Risk of data loss. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down todoapp-postgres port-forwarder && \
		rm -rf ${PROJECT_ROOT}/out/pgdata && \
		echo "Environment files have been cleared"; \
	else \
		echo "Environment cleanup cancelled"; \
	fi

env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder

logs-cleanup:
	@read -p "Clear all log files? Risk of losing logs. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		rm -rf ${PROJECT_ROOT}/out/logs && \
		echo "Logs files have been cleared"; \
	else \
		echo "Logs files cleanup cancelled"; \
	fi

# ============
# Migrations
# ============

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "The required seq parameter is missing. Example: make migrate-create seq=init"; \
		exit 1; \
	fi; \
	docker compose run --rm todoapp-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "The required action parameter is missing. Example: make migrate-action action=up"; \
		exit 1; \
	fi; \
	docker compose run --rm todoapp-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todoapp-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

admin-promote:
	@if [ -z "$(email)" ]; then \
		echo "The required email parameter is missing. Example: make admin-promote email=ivan@example.com"; \
		exit 1; \
	fi; \
	result=$$(docker exec \
		-e PGPASSWORD=${POSTGRES_PASSWORD} \
		todoapp-env-postgres \
		psql \
			-U ${POSTGRES_USER} \
			-d ${POSTGRES_DB} \
			-tA \
			-v ON_ERROR_STOP=1 \
			-c "UPDATE todoapp.users SET role = 'admin' WHERE email = '$(email)' RETURNING email;"); \
	if [ -z "$$result" ]; then \
		echo "User with email $(email) not found"; \
		exit 1; \
	fi; \
	echo "Promoted $$result to admin";

# ============
# Utils
# ============

ps:
	@docker compose ps

# ============
# Startap app
# ============

todoapp-run:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	export POSTGRES_HOST=localhost && \
	go mod tidy && \
	go run ${PROJECT_ROOT}/cmd/todoapp/main.go

# ============
# Deploy
# ============

todoapp-deploy:
	@docker compose up -d --build todoapp

todoapp-undeploy:
	@docker compose down todoapp

# ============
# Swagger
# ============

swagger-gen:
	@docker compose run --rm swagger \
		init \
		-g cmd/todoapp/main.go \
		-o docs \
		--parseInternal \
		--parseDependency

# ============
# Tests
# ============

test:
	@go test -v ./... 2>&1 | grep -v '\[no test files\]'

generate:
	@mockery
