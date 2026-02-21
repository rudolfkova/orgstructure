.PHONY: build
build:
	go build -v ./cmd/app

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: migrate
migrate-up:
	migrate -path migrations -database "$(DB_DSN)" up

.PHONY: gofmt
gofmt:
	gofmt -w -s .