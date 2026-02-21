.PHONY: build
build:
	go build -v ./cmd/app
docker-build:
	docker-compose up --build

.PHONY: start
start:
	./app.exe

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