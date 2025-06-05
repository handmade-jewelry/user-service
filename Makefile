PROTO_DIR=.
GEN_DIR=pkg/user-service
MIGRATIONS_DIR = ./migrations
DB_USER=user_service_user
DB_NAME=user_service_db
DB_SSLMODE=disable

.PHONY: proto-gen new-migration migration-up migration-down

proto-gen:
	buf generate

# Make new migration sql
new-migration:
ifndef NAME
	$(error Usage: make new-migration NAME=your_migration_name)
endif
	goose -dir $(MIGRATIONS_DIR) create $(NAME) sql

# Apply migrations
migration-up:
	goose -dir $(MIGRATIONS_DIR) postgres "user=$(DB_USER) dbname=$(DB_NAME) sslmode=$(DB_SSLMODE)" up

# Rollback migrations
migration-down:
	goose -dir $(MIGRATIONS_DIR) postgres "user=$(DB_USER) dbname=$(DB_NAME) sslmode=$(DB_SSLMODE)" down