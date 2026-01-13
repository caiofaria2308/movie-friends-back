USER := $(shell id -u):$(shell id -g)

.PHONY: up down logs shell build

up:
	docker compose -f docker-compose.yml -f docker-compose.override.yml up -d

stop:
	docker compose -f docker-compose.yml -f docker-compose.override.yml stop

restart_app:
	docker compose -f docker-compose.yml -f docker-compose.override.yml restart app

down:
	docker compose -f docker-compose.yml -f docker-compose.override.yml down

log_db:
	docker compose -f docker-compose.yml -f docker-compose.override.yml logs -f postgres

log:
	docker compose -f docker-compose.yml -f docker-compose.override.yml logs -f app

sh_app:
	docker compose -f docker-compose.yml -f docker-compose.override.yml exec app sh

sh_db:
	docker compose -f docker-compose.yml -f docker-compose.override.yml exec postgres sh

build:
	docker compose -f docker-compose.yml -f docker-compose.override.yml build

	