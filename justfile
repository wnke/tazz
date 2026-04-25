generate-sql:
	sqlc generate -f internal/adapter/driven/persistence/sqlc/sqlc.yaml

migrate-up:
	goose -dir internal/adapter/driven/persistence sqlite3 ./devices.db up

migrate-down:
	goose -dir internal/adapter/driven/persistence sqlite3 ./devices.db down
