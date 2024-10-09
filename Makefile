dependencies_start: ## Dockerコンテナ起動（PostgreSQL）
	@docker compose -f ./docker/docker-compose.yml up -d

dependencies_stop: ## Dockerコンテナ停止
	@docker compose -f ./docker/docker-compose.yml down

migrate_generate: ## schema.sql更新 & migration file生成（migration file生成の為にDBを一度リセットするので本番環境では使わないで下さい）
	@go run ./cmd/postgres/main.go migrate generate

migrate_up: ## migration実行
	@go run ./cmd/postgres/main.go migrate up

migrate_down_one: ## migrationを1つ戻す
	@go run ./cmd/postgres/main.go migrate downone

migrate_down_all: ## migrationを全て戻す
	@go run ./cmd/postgres/main.go migrate downall

seed: ## seedデータ投入
	make migrate_down_all
	make migrate_up
	@go run ./cmd/postgres/main.go seed

godoc: ## godocサーバー起動
	@godoc -http=:6060 & sleep 2 && open http://localhost:6060/pkg/github.com/ucho456job/pocgo/

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

