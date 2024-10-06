godoc: ## godocサーバーを起動する。
	@godoc -http=:6060 & sleep 2 && open http://localhost:6060/pkg/github.com/ucho456job/pocgo/

unit_test: ## ユニットテストを実行する。
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

