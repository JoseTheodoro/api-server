MIGRATIONS_DIR=./internal/repository/postgres/migrations

build:
	@CGO_ENABLED=0 GOOS=linux go build -o ./bin/api ./cmd/api
create-migration:
	migrate create -ext sql -dir $(MIGRATIONS_DIR)  -seq $(NAME)
run-migration-up:
	migrate -database ${DATABASE_DSN} -path $(MIGRATIONS_DIR) up
run-migration-down:
	migrate -database ${DATABASE_DSN} -path $(MIGRATIONS_DIR) down