.PHONY: develop_start dependencies_start dependencies_stop migrate_refresh \
				migrate_up migrate_down migrate_reset drop_tables seed run clean \
				unit_test unit_coverage integration_test integration_coverage \
				swagger mockgen help

develop_start: ## Start docker container, migrate, seed
	@docker compose -f ./docker/docker-compose.yml up -d
	sleep 5
	make migrate_refresh
	make migrate_up
	make seed

dependencies_start: ## Start docker container
	@docker compose -f ./docker/docker-compose.yml up -d

dependencies_stop: ## Stop docker container
	@docker compose -f ./docker/docker-compose.yml down

migrate_refresh: ## Update schema.sql & generate migration file
	@go run ./cmd/postgres/main.go migrate refresh

migrate_up: ## Run migration
	@go run ./cmd/postgres/main.go migrate up

migrate_down: ## Revert one migration
	@go run ./cmd/postgres/main.go migrate down

migrate_reset: ## Revert all migrations
	@go run ./cmd/postgres/main.go migrate reset

drop_tables: ## Drop all tables
	@go run ./cmd/postgres/main.go drop tables

seed: ## Seed data input
	make drop_tables
	make migrate_up
	@go run ./cmd/postgres/main.go insert seed

run: ## Run application
	@go run ./cmd/pocgo/main.go

clean: ## Delete cache
	@go clean -cache -modcache

unit_test: ## Run unit test (Specify CASE to run only a specific test, e.g. CASE=TestSignup)
	@mkdir -p tmp
	@if [ -z "$(CASE)" ]; then \
		go test ./internal/... ./pkg/... -v -coverprofile=tmp/unit_coverage.out 2>&1 | tee tmp/unit_test.log; \
	else \
		go test ./internal/... ./pkg/... -v -coverprofile=tmp/unit_coverage.out -run ^$(CASE)$$ 2>&1 | tee tmp/unit_test.log; \
	fi
	@go tool cover -html=tmp/unit_coverage.out -o tmp/unit_test.cover.html

unit_coverage: ## Show unit test coverage report in terminal
	@go tool cover -func=tmp/unit_coverage.out

integration_test: ## Run integration tests (Specify CASE to run only a specific test, e.g., CASE=TestSignup | Use UPDATE=-update to refresh golden files)
	@mkdir -p tmp
	@if [ -z "$(CASE)" ]; then \
		go test ./test/... -v -coverprofile=tmp/integration_coverage.out -coverpkg=./internal/... $(UPDATE) 2>&1 | tee tmp/integration_test.log; \
	else \
		go test ./test/... -v -coverprofile=tmp/integration_coverage.out -coverpkg=./internal/... -run ^$(CASE)$$ $(UPDATE) 2>&1 | tee tmp/integration_test.log; \
	fi
	@go tool cover -html=tmp/integration_coverage.out -o tmp/integration_test.cover.html

integration_coverage: ## Show integration test coverage report in terminal
	@go tool cover -func=tmp/integration_coverage.out

swagger: ## Generate swagger document
	@swag init -g ./cmd/pocgo/main.go

mockgen: ## Generate mock (e.g. make mockgen path=internal/domain/user/user_repository.go)
	@if echo "$(path)" | grep -q "internal/application"; then \
		mockgen -source=$(path) \
								-destination=internal/application/mock/mock_$(notdir $(basename $(path))).go \
								-package=mock; \
	elif echo "$(path)" | grep -q "internal/domain"; then \
		mockgen -source=$(path) \
								-destination=internal/domain/mock/mock_$(notdir $(basename $(path))).go \
								-package=mock; \
	else \
		echo "Unsupported path: $(path)"; \
		exit 1; \
	fi

help: ## Show help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
