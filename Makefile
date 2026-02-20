.PHONY: run test lint migrate-up migrate-down sqlc-generate

run:
	go run cmd/server/main.go

test:
	go test -v -race ./...

lint:
	golangci-lint run

migrate-up:
	goose -dir db/migrations up

migrate-down:
	goose -dir db/migrations down

sqlc-generate:
	sqlc generate