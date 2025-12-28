DOCKER_COMPOSE ?= docker compose

.PHONY: up
up:
	$(DOCKER_COMPOSE) up -d

.PHONY: down
down:
	$(DOCKER_COMPOSE) down

.PHONY: restart
restart: down up

.PHONY: ps
ps:
	$(DOCKER_COMPOSE) ps

.PHONY: logs
logs:
	$(DOCKER_COMPOSE) logs -f --tail=200

.PHONY: pull
pull:
	$(DOCKER_COMPOSE) pull

.PHONY: reset
reset:
	$(DOCKER_COMPOSE) down -v --remove-orphans

# Удобные команды для auth-service (запуск сервиса — отдельно от docker-compose).
.PHONY: auth-generate-api
auth-generate-api:
	$(MAKE) -C auth_service generate-api

.PHONY: auth-run
auth-run: auth-generate-api auth-test
	$(MAKE) -C auth_service run

.PHONY: auth-test
auth-test:
	$(MAKE) -C auth_service cov

