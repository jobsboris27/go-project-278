build:
	go build -o bin/gendiff ./cmd/gendiff

run: build
	bin/gendiff

.PHONY: build run

test:
	go mod tidy
	go test -v -coverprofile=coverage.out ./...

install:
	go install ./cmd/gendiff

lint:
	golangci-lint run ./...