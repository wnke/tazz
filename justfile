set shell := ["bash", "-cu"]

gen-api:
	oapi-codegen -config internal/driver/http/oapi-codegen.yaml internal/driver/http/openapi.yaml

test:
	go test ./...

lint:
	golangci-lint run

check:
	just lint
	just test
generate-sql:
	sqlc generate -f internal/adapter/driven/persistence/sqlc/sqlc.yaml

migrate-up:
	goose -dir internal/adapter/driven/persistence/migrations sqlite3 ./devices.db up

migrate-down:
	goose -dir internal/adapter/driven/persistence/migrations sqlite3 ./devices.db down
