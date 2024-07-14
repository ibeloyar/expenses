include config/docker/.env.pg

MAIN_FILE = cmd/expenses.go
DATABASE_STRING = "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable"
DATABASE_MIGRATIONS_PATH = internal/storage/postgres/migrations

.PHONY: install-tools
install-tools:
	go install github.com/swaggo/swag/cmd/swag@latest # swagger
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest # golang-migrate CLI

.PHONY: run
run:
	go run $(MAIN_FILE)

.PHONY: swagger-gen
swagger-gen:
	$(GOPATH)/bin/swag init -o ./api -g $(MAIN_FILE)

.PHONY: migrate-up
migrate-up:
	migrate \
	-path $(DATABASE_MIGRATIONS_PATH) \
	-database $(DATABASE_STRING) up

.PHONY: migrate-down
migrate-down:
	migrate \
	-path $(DATABASE_MIGRATIONS_PATH) \
	-database $(DATABASE_STRING) down

.PHONY: migrate-create
migrate-create:
	migrate create \
	-ext sql \
	-dir $(DATABASE_MIGRATIONS_PATH) \
	-seq $(NAME)