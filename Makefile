.PHONY: develop_start dependencies_start dependencies_stop migrate_refresh \
				migrate_up migrate_down migrate_reset drop_tables seed \
				unit_test godoc help

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
	@-export UNIT_TEST_CMD="test ./pkg/... -v -coverprofile=tmp/unit_test_cover.out";\
	mkdir -p tmp; \
	if [ -z "$(TEST_TARGET)" ]; \
		then UNIT_TEST_CMD="$$UNIT_TEST_CMD ./..." ;\
		else UNIT_TEST_CMD="$$UNIT_TEST_CMD -run $(TEST_TARGET)" ;\
	fi;\
	if [ -n "$(TEST_UPDATE)" ]; \
		then UNIT_TEST_CMD="$$UNIT_TEST_CMD -update" ;\
	fi;\
	echo $$UNIT_TEST_CMD;\
	go $$UNIT_TEST_CMD | tee tmp/unit_test.log;\
	go $$UNIT_TEST_CMD;\
	EXECUTE_CODE=$$?;\
	go tool cover -html=tmp/unit_test_cover.out -o tmp/unit_test_cover.html;\
	if [ $$EXECUTE_CODE -eq 1 ]; then false; fi;

help: ## ヘルプ表示
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
