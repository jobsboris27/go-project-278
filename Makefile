.PHONY: run test lint

run:
	go run main.go

test:
	go test -v -race ./...

lint:
	golangci-lint run