PROTO_DIR=.
GEN_DIR=pkg/user-service
MIGRATIONS_DIR = ./migrations

.PHONY: proto-gen new-migration migration-up migration-down

proto-gen:
	buf generate

# Make new migration sql
ifndef NAME
	$(error Usage: make new-migration NAME=your_migration_name)
endif

new-migration:
	goose -dir $(MIGRATIONS_DIR) create $(NAME) sql

# Применить миграции
migration-up:
	goose -dir $(MIGRATIONS_DIR) postgres "user=$(DB_USER) dbname=$(DB_NAME) sslmode=$(DB_SSLMODE)" up

# Откатить миграции
migration-down:
	goose -dir $(MIGRATIONS_DIR) postg