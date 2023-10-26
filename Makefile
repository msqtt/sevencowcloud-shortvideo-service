gen:

sqlc:

dev-up:
	docker compose -f ./deploy/dev-compose.yml up -d

dev-down:
	docker compose -f ./deploy/dev-compose.yml down

migrate-up:
	migrate -path internal/db/migration/ -database "mysql://root:password@tcp(localhost:3306)/shortvideo" up

migrate-down:
	migrate -path internal/db/migration/ -database "mysql://root:password@tcp(localhost:3306)/shortvideo" down
