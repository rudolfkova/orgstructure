.PHONY: build
build:
	go build -v ./cmd/app

.PHONY: start
start:
	./app.exe -config-path=config.toml

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: gofmt
gofmt:
	gofmt -w -s .

.PHONY: mock
mock:
	mockery --name=DepartmentRepository --dir=./internal/repository --output=./mocks/repository --outpkg=mocks
	mockery --name=EmployeeRepository --dir=./internal/repository --output=./mocks/repository --outpkg=mocks

.PHONY: test
test:
	go test -v -race -timeout 30s ./...