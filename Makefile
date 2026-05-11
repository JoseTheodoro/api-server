MIGRATIONS_DIR=./internal/repository/postgres/migrations

create-migration:
	migrate create -ext sql -dir $(MIGRATIONS_DIR)  -seq $(NAME)
run-migration-up:
	migrate -database ${DATABASE_DSN} -path $(MIGRATIONS_DIR) up
run-migration-down:
	migrate -database ${DATABASE_DSN} -path $(MIGRATIONS_DIR) down