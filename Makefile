.PHONY: lint lint-fix test ci

lint:
	golangci-lint run ./...

lint-fix:
	golangci-lint run --fix ./...

test:
	go test ./...
