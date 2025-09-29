run:
	go run .\cmd\sso\main.go --config="./config/local.yaml"

migrate:
	go run .\cmd\migrator --storage_path=./storage/sso.db --migrations_path=./migrations

migrate-test:
	go run ./cmd/migrator --storage_path=./storage/sso.db --migrations_path=./tests/migrations --migrations_table=migrations_test
