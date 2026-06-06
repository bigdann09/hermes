include .env

DATABASE := postgresql://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL}
MIGRATION_PATH := ./internal/infrastructure/database/migrations

run:
	go run cmd/api/main.go

migrate_create:
	migrate create -ext sql -dir ${MIGRATION_PATH} -seq ${name}

migrate_up:
	migrate -database ${DATABASE} -path ${MIGRATION_PATH} up

migrate_down:
	migrate -database ${DATABASE} -path ${MIGRATION_PATH} down

migrate_force:
	migrate -database ${DATABASE} -path ${MIGRATION_PATH} force ${version}

docker_run:
	docker compose up --build --watch

docker_down:
	docker compose down

docker_seed:
	docker compose run --rm --no-deps backend ./cli seed

docker_migrate:
	docker compose run --rm migrate

seed:
	go build -o cli cmd/cli/main.go
	./cli seed