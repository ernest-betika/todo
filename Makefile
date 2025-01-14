include .env

# make migration name=CreateUsers
migration:
	goose -dir internal/db/migrations create $(name) sql

migrate:
	goose -dir 'internal/db/migrations' postgres ${DATABASE_URL} up

rollback:
	goose -dir 'internal/db/migrations' postgres ${DATABASE_URL} down

up:
	docker-compose up --remove-orphans

ps:
	docker-compose ps

down:
	docker-compose down --remove-orphans

server:
	gow run cmd/todo-api/main.go -e .env
