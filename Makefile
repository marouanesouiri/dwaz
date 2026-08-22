.PHONY: all fmt lint lint-fix test

all: fmt lint test

## fmt: Format Go source code
fmt:
	go fmt ./...

## lint: Run linter checks
lint:
	golangci-lint run ./...

## lint-fix: Run linter and apply automated fixes
lint-fix:
	golangci-lint run --fix ./...

## test: Run unit tests with race detection
test:
	go test -v -race ./...
