URL_DB=mysql://root:password@tcp(localhost:3306)/shortvideo
MOD_NAME=github.com/msqtt/sevencowcloud-shortvideo-service

gen:
	protoc --proto_path=api/protos/v1 \
		--go_out=api/pb/v1 --go_opt=paths=source_relative \
		api/protos/v1/user/*.proto
sqlc:
	sqlc generate -f configs/sqlc.yaml
dev-up:
	docker compose -f ./deploy/dev-compose.yml up -d
dev-down:
	docker compose -f ./deploy/dev-compose.yml down
mock:
	mockgen -package mockdb -destination internal/repo/mock/store.go $(MOD_NAME)/internal/repo/sqlc Store
migrate-up:
	migrate -path internal/db/migration/ -database "$(URL_DB)" up
migrate-down:
	migrate -path internal/db/migration/ -database "$(URL_DB)" down
