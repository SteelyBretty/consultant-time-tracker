.PHONY: run build test clean docker-build docker-run

run:
	go run cmd/server/main.go

build:
	go build -o bin/server cmd/server/main.go

test:
	go test -v ./...

clean:
	rm -rf bin/
	rm -f data/*.db

docker-build:
	docker build -t consultant-tracker:latest -f docker/Dockerfile .

docker-run:
	docker-compose -f docker/docker-compose.yml up

deps:
	go mod download
	go mod tidy

dev:
	air -c .air.toml