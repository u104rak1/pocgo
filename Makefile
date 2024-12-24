.PHONY: develop_start dependencies_start dependencies_stop migrate_refresh \
				migrate_up migrate_down migrate_reset drop_tables seed run clean \
				unit_test unit_coverage integration_test integration_coverage \
				swagger mockgen help

develop_start: ## Dockerコンテナを起動し、マイグレーションを実行してシードデータを挿入
	@docker compose -f ./docker/docker-compose.yml up -d
	sleep 5
	make migrate_refresh
	make migrate_up
	make seed

dependencies_start: ## Dockerコンテナを起動
	@docker compose -f ./docker/docker-compose.yml up -d

dependencies_stop: ## Dockerコンテナを停止
	@docker compose -f ./docker/docker-compose.yml down

migrate_refresh: ## schema.sqlを更新してマイグレーションファイルを生成
	@go run ./cmd/postgres/main.go migrate refresh

migrate_up: ## マイグレーションを実行
	@go run ./cmd/postgres/main.go migrate up

migrate_down: ## 1つのマイグレーションを取り消す
	@go run ./cmd/postgres/main.go migrate down

migrate_reset: ## すべてのマイグレーションを取り消す
	@go run ./cmd/postgres/main.go migrate reset

drop_tables: ## すべてのテーブルを削除
	@go run ./cmd/postgres/main.go drop tables

seed: ## テーブルを削除してからすべてのマイグレーションを実行した後にシードデータを挿入
	make drop_tables
	make migrate_up
	@go run ./cmd/postgres/main.go insert seed

run: ## アプリケーションを実行
	@go run ./cmd/pocgo/main.go

clean: ## キャッシュを削除
	@go clean -cache -modcache

unit_test: ## 単体テストを実行
	@mkdir -p tmp
	@go test -v ./internal/... ./pkg/... -coverprofile=tmp/unit_coverage.out 2>&1 | tee tmp/unit_test.log
	@go tool cover -html=tmp/unit_coverage.out -o tmp/unit_test.cover.html

integration_test: ## 統合テストを実行 (特定のテストのみを実行するには CASE を指定します。例: CASE=TestSignup | ゴールデンファイルを更新するには UPDATE=-update を使用します)
	@mkdir -p tmp
	@if [ -z "$(CASE)" ]; then \
		go test -v ./test/integration -coverprofile=tmp/integration_coverage.out -coverpkg=./internal/... $(UPDATE) 2>&1 | tee tmp/integration_test.log; \
	else \
		go test -v ./test/integration -coverprofile=tmp/integration_coverage.out -coverpkg=./internal/... -run ^$(CASE)$$ $(UPDATE) 2>&1 | tee tmp/integration_test.log; \
	fi
	@go tool cover -html=tmp/integration_coverage.out -o tmp/integration_test.cover.html

integration_coverage: ## ターミナルに統合テストカバレッジレポートを表示
	@go tool cover -func=tmp/integration_coverage.out

swagger: ## Swaggerドキュメントを生成
	@swag init -g ./cmd/pocgo/main.go

mockgen: ## Mockを生成 (例: make mockgen path=internal/domain/user/user_repository.go)
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

help: ## ヘルプを表示
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
