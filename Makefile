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

migration-up:
	goose -dir $(MIGRATIONS_DIR) postgres "user=auth_service_user dbname=auth_service_db sslmode=disable" up

migration-down:
	goose -dir $(MIGRATIONS_DIR) postgres "user=auth_service_user dbname=auth_service_db sslmode=disable" down

