.PHONY: generate-api
generate-api:
	@./scripts/generate.sh

.PHONY: run
run:
	@configPath=./config.yaml swaggerPath=./internal/pb/swagger/metadata_api/metadata.swagger.json go run ./cmd/app

.PHONY: cov
cov:
	go test -cover ./internal/services/metadata_service

.PHONY: mock
mock:
	mockery
