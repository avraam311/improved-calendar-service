.PHONY: up, down, lint, test

up:
	docker-compose up -d --build

down:
	docker-compose down -v

lint:
	go vet ./...
	golangci-lint run ./...

test:
	go clean --testcache
	go test ./...