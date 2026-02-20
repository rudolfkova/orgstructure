.PHONY: build
build:
	go build -v ./cmd/app

.PHONY: lint
lint:
	golangci-lint run ./...