.PHONY: build test lint

build:
	go build -o willys-cli .

test:
	go test ./...

lint:
	golangci-lint run
