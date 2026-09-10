
DB_PASSWORD ?= default_pass	# Set this to your password

DSN = 'postgres://shodruz:$(DB_PASSWORD)@localhost:5432/myshop?sslmode=disable'

.PHONY: migrate-up migrate-down server-run

migrate-up:
	migrate -path migrations -database $(DSN) up

migrate-down:
	migrate -path migrations -database $(DSN) down 1

server-run:
	go run ./cmd/api -db-password=$(DB_PASSWORD)