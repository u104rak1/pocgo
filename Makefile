.PHONY: develop_start dependencies_start dependencies_stop migrate_refresh \
				migrate_up migrate_down migrate_reset drop_tables seed run clean \
				unit_test godoc swagger help

develop_start: ## Dockerを動かしDBのデータを準備する
	@docker compose -f ./docker/docker-compose.yml up -d
	sleep 5
	make migrate_refresh
	make migrate_up
	make seed

dependencies_start: ## Dockerコンテナ起動（PostgreSQL）
	@docker compose -f ./docker/docker-compose.yml up -d

dependencies_stop: ## Dockerコンテナ停止
	@docker compose -f ./docker/docker-compose.yml down

migrate_refresh: ## schema.sql更新 & migration file生成（migration file生成の為にDBを一度リセットするので本番環境では使わないで下さい）
	@go run ./cmd/postgres/main.go migrate refresh

migrate_up: ## migration実行
	@go run ./cmd/postgres/main.go migrate up

migrate_down: ## migrationを1つ戻す
	@go run ./cmd/postgres/main.go migrate down

migrate_reset: ## migrationを全て戻す
	@go run ./cmd/postgres/main.go migrate reset

drop_tables: ## 全テーブル削除
	@go run ./cmd/postgres/main.go drop tables

seed: ## seedデータ投入
	make drop_tables
	make migrate_up
	@go run ./cmd/postgres/main.go insert seed

run: ## サーバー起動
	@go run ./cmd/pocgo/main.go

clean: ## キャッシュ削除
	@go clean -cache -modcache

unit_test: ## ユニットテスト実行
	@mkdir -p tmp
	@echo "Running tests and generating log and coverage reports..."
	@go test ./internal/... ./pkg/... -v -coverprofile=tmp/coverage.out 2>&1 | tee tmp/unit_test.log
	@go tool cover -html=tmp/coverage.out -o tmp/unit_test.cover.html
	@go tool cover -func=tmp/coverage.out

unit_coverage: ## ユニットテストのカバレッジレポートをログに表示
	@go tool cover -func=tmp/coverage.out

swagger: ## Swaggerドキュメント生成
	@swag init -g ./cmd/pocgo/main.go

help: ## ヘルプ表示
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
