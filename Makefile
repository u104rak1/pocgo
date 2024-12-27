.PHONY: run_local dependencies_start dependencies_stop migrate_refresh \
				migrate_up migrate_down migrate_reset drop_tables seed run clean \
				unit_test unit_coverage integration_test integration_coverage \
				swagger mockgen build_docker run_docker_local help

##
## Dockerを使用したコマンド群です。`docker container exec`を使用してホストマシンからコンテナにアクセスして実行します。
## コンテナにはアプリケーション用、PostgreSQL用、デバッグ用の3つのコンテナがあります。
## Atlasなどの依存ツールを全てコンテナ内でインストールする為、開発環境を簡単に構築できます。
##
develop_start: ## アプリケーションの準備と実行
	@make dependencies_start
	make migrate_refresh
	make migrate_up
	make seed

dependencies_start: ## アプリケーションとPostgreSQLのコンテナを起動
	@docker compose -f ./docker/docker-compose.yml up -d pocgo postgres

dependencies_start_delve: ## delveでデバッグに対応したアプリケーションコンテナを追加で起動 (デバッガをアタッチした後、port 8081でアクセスできる)
	@docker compose -f ./docker/docker-compose.yml up -d delve

dependencies_stop: ## 全てのDockerコンテナを停止
	@docker compose -f ./docker/docker-compose.yml down

migrate_refresh: ## schema.sqlを更新してマイグレーションファイルを生成
	@docker container exec pocgo_app go run ./cmd/postgres/main.go migrate refresh

migrate_up: ## マイグレーションを実行
	@docker container exec pocgo_app go run ./cmd/postgres/main.go migrate up

migrate_down: ## 最新のマイグレーションを取り消す
	@docker container exec pocgo_app go run ./cmd/postgres/main.go migrate down

migrate_reset: ## すべてのマイグレーションを取り消す
	@docker container exec pocgo_app go run ./cmd/postgres/main.go migrate reset

drop_tables: ## すべてのテーブルを削除
	@docker container exec pocgo_app go run ./cmd/postgres/main.go drop tables

seed: ## シードデータを挿入 (一度データをリセットしてからシードデータを挿入)
	make drop_tables
	make migrate_up
	@docker container exec pocgo_app go run ./cmd/postgres/main.go insert seed

unit_test: ## 単体テストを実行 (詳細表示するには SHOW=-v を使用)
	@mkdir -p tmp
	@docker container exec pocgo_app go test $(SHOW) ./internal/... ./pkg/... -coverprofile=tmp/unit_coverage.out 2>&1 | tee tmp/unit_test.log
	@docker container exec pocgo_app go tool cover -html=tmp/unit_coverage.out -o tmp/unit_test.cover.html

integration_test: ## 統合テストを実行 (特定のテストのみを実行するには CASE を指定。例: CASE=TestSignup | ゴールデンファイルを更新するには UPDATE=-update を使用 | 詳細表示するには SHOW=-v を使用)
	@mkdir -p tmp
	@if [ -z "$(CASE)" ]; then \
		docker container exec pocgo_app go test $(SHOW) ./test/integration -coverprofile=tmp/integration_coverage.out -coverpkg=./internal/... $(UPDATE) 2>&1 | tee tmp/integration_test.log; \
	else \
		docker container exec pocgo_app go test $(SHOW) ./test/integration -coverprofile=tmp/integration_coverage.out -coverpkg=./internal/... -run ^$(CASE)$$ $(UPDATE) 2>&1 | tee tmp/integration_test.log; \
	fi
	@docker container exec pocgo_app go tool cover -html=tmp/integration_coverage.out -o tmp/integration_test.cover.html

swagger: ## Swaggerドキュメントを生成
	@docker container exec pocgo_app swag init -g ./cmd/pocgo/main.go

mockgen: ## Mockを生成 (例: make mockgen path=internal/domain/user/user_repository.go)
	@if echo "$(path)" | grep -q "internal/application"; then \
		docker container exec pocgo_app mockgen -source=$(path) \
								-destination=internal/application/mock/mock_$(notdir $(basename $(path))).go \
								-package=mock; \
	elif echo "$(path)" | grep -q "internal/domain"; then \
		docker container exec pocgo_app mockgen -source=$(path) \
								-destination=internal/domain/mock/mock_$(notdir $(basename $(path))).go \
								-package=mock; \
	else \
		echo "Unsupported path: $(path)"; \
		exit 1; \
	fi

attach: ## コンテナにbashをアタッチ
	docker exec -it pocgo_app /bin/bash

##
## 
## ホストマシンでアプリケーションを実行するコマンド群。
## インメモリをDBとして活用する為、依存関係がなく、Goを実行する環境があれば簡単にアプリケーションを実行できます。
## ホストマシンでさくっとアプリケーションを動かしたい場合はこちらを使用します。
##
run_host: ## アプリケーションを実行
	USE_INMEMORY=true go run ./cmd/pocgo/main.go

clean: ## キャッシュを削除
	@go clean -cache -modcache

unit_test_host: ## 単体テストを実行
	@mkdir -p tmp
	@go test $(SHOW) ./internal/... ./pkg/... -coverprofile=tmp/unit_coverage.out 2>&1 | tee tmp/unit_test.log
	@go tool cover -html=tmp/unit_coverage.out -o tmp/unit_test.cover.html

integration_test_host: ## 統合テストを実行 (DBコンテナが起動している必要があります)
	@mkdir -p tmp
	@if [ -z "$(CASE)" ]; then \
		POSTGRES_HOST=localhost go test $(SHOW) ./test/integration -coverprofile=tmp/integration_coverage.out -coverpkg=./internal/... $(UPDATE) 2>&1 | tee tmp/integration_test.log; \
	else \
		POSTGRES_HOST=localhost go test $(SHOW) ./test/integration -coverprofile=tmp/integration_coverage.out -coverpkg=./internal/... -run ^$(CASE)$$ $(UPDATE) 2>&1 | tee tmp/integration_test.log; \
	fi
	@go tool cover -html=tmp/integration_coverage.out -o tmp/integration_test.cover.html

swagger_host: ## Swaggerドキュメントを生成 （go install github.com/swaggo/swag/cmd/swag@latestが必要です）
	@swag init -g ./cmd/pocgo/main.go

mockgen_host: ## Mockを生成 (go install github.com/golang/mock/mockgen@latestが必要です)
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

##
##
## 本番環境用のコマンド群
##
build_prod: ## Dockerコンテナをビルド
	@docker build -t pocgo -f ./docker/Dockerfile .

run_prod: ## Dockerコンテナを起動し、インメモリモードでアプリケーションを実行
	@docker run -e USE_INMEMORY=true -p 8080:8080 pocgo

##
##
## その他
##
start_dependencies_for_integration_test_in_ci: ## CIでインテグレーションテストを実行する為に必要な依存関係まわりの実行コマンド 
	@docker compose -f ./docker/docker-compose.yml up -d postgres
	sleep 5s


help: ## ヘルプを表示
	@grep -E '(^##|^[a-zA-Z_-]+:.*?##)' $(MAKEFILE_LIST) | \
		awk '/^##/ {print substr($$0, 4)} /^[a-zA-Z_-]+:/ {split($$0, a, ":.*?## "); printf "\033[36m%-30s\033[0m %s\n", a[1], a[2]}'
