.PHONY: build test clean run-standard run-custom lint

build:
	go build -o bin/standard-server examples/standard_server/main.go
	go build -o bin/custom-server examples/custom_server/main.go

test:
	go test -v ./...

clean:
	rm -rf bin/
	go clean

run-standard:
	go run examples/std_server/main.go

run-custom:
	go run examples/custom/main.go

lint:
	golangci-lint run

fmt:
	go fmt ./...

deps:
	go mod download

check: fmt lint test
