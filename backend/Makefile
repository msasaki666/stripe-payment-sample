.PHONY: migrate-reset
migrate-reset:
	migrate -database postgres://postgres:postgres@db:5432/app_database?sslmode=disable -path db/migrations drop -f
	migrate -database postgres://postgres:postgres@db:5432/app_database?sslmode=disable -path db/migrations up
	sqlboiler psql

.PHONY: migrate-reset-seed
migrate-reset-seed:
	$(MAKE) migrate-reset
	go run ./cmd/seed
