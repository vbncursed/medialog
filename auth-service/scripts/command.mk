.PHONY: generate-api
generate-api:
	@./scripts/generate.sh

.PHONY: run
run:
	@configPath=./config.yaml swaggerPath=./internal/pb/swagger/auth_api/auth.swagger.json go run ./cmd/app

.PHONY: cov
cov:
	go test -cover ./internal/services/authService

.PHONY: mock
mock:
	mockery
