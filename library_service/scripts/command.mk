.PHONY: generate-api
generate-api:
	@./scripts/generate.sh

.PHONY: run
run:
	@configPath=./config.yaml swaggerPath=./internal/pb/swagger/library_api/library.swagger.json go run ./cmd/app

.PHONY: cov
cov:
	go test -cover ./internal/services/library_service

.PHONY: mock
mock:
	mockery
