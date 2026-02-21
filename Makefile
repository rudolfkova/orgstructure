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

.PHONY: mock
mock:
	mockery --name=DepartmentRepository --dir=./internal/repository --output=./mocks/repository --outpkg=mocks
	mockery --name=EmployeeRepository --dir=./internal/repository --output=./mocks/repository --outpkg=mocks