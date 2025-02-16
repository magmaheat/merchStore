compose-up:
	docker-compose up --build -d && docker-compose logs -f
.PHONY: compose-up

compose-down:
	docker-compose down --remove-orphans
.PHONY: compose-down

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

linter-golangci:
	golangci-lint run
.PHONY: linter-golangci
