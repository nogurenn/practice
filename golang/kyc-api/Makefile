producer:
	docker-compose up --build --force-recreate --no-deps producer

up:
	docker-compose up -d db
.PHONY: up

down:
	docker-compose down
.PHONY: down

migrate:
	docker-compose run --rm flyway
.PHONY: migrate

mocks:
	docker-compose run --rm mockery --all --dir=./internal
.PHONY: mocks

test:
	docker-compose up --build --force-recreate --no-deps test
.PHONY: test
