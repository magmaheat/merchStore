up:
	docker-compose up --build -d && docker-compose logs -f
.PHONY: compose-up

down:
	docker-compose down --remove-orphans
.PHONY: compose-down

unit-tests:
	go test -coverprofile=coverage.out ./internal/service
	go tool cover -func=coverage.out
.PHONY: unit-tests

tests:
	go test ./tests
.PHONY: tests

linter-golangci:
	golangci-lint run
.PHONY: linter-golangci
