.PHONY: run build test clean docker-build docker-run docker-dev docker-down docker-logs deps dev install-air

run:
	go run cmd/server/main.go

build:
	go build -o bin/server cmd/server/main.go

test:
	go test -v ./...

clean:
	rm -rf bin/
	rm -f data/*.db
	rm -rf tmp/

docker-build:
	docker build -t consultant-tracker:latest -f docker/Dockerfile .

docker-dev:
	docker compose -f docker/docker-compose.yml up --build

docker-down:
	docker compose -f docker/docker-compose.yml down

docker-logs:
	docker compose -f docker/docker-compose.yml logs -f

deps:
	go mod download
	go mod tidy

dev:
	air -c .air.toml

install-air:
	go install github.com/air-verse/air@latest