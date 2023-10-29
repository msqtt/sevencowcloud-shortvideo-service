URL_DB=mysql://root:password@tcp(localhost:3306)/shortvideo
TEST_URL_DB=mysql://root:password@tcp(api.kexie.space:3306)/shortvideo
MOD_NAME=github.com/msqtt/sevencowcloud-shortvideo-service
EXE_NAME=cowserver

dev:
	go run cmd/server/main.go
gen:
	protoc --proto_path=api/protos/v1 \
		--go_out=api/pb/v1 --go_opt=paths=source_relative \
		--go-grpc_out=api/pb/v1 --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=api/pb/v1 --grpc-gateway_opt paths=source_relative \
		--openapiv2_out=api/openapi/v1 \
		api/protos/v1/user/*.proto \
		api/protos/v1/profile/*.proto
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
cowserver:
	CGO_ENABLED=0 GOARCH="amd64" go build -o $(EXE_NAME) cmd/server/main.go
migrate-test-up:
	migrate -path internal/db/migration/ -database "$(TEST_URL_DB)" up
migrate-test-down:
	migrate -path internal/db/migration/ -database "$(TEST_URL_DB)" down
deploy: $(EXE_NAME)
	rsync -av $(EXE_NAME) ubuntu@kexieserver:cow/
	rsync -av configs/app.env ubuntu@kexieserver:cow/configs/
	ssh ubuntu@kexieserver "cd cow/ ; ./start.sh"
	rm cowserver

.PHONY: dev gen sqlc dev-up dev-down mock migrate-up migrate-down deploy
