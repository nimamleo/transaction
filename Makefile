.PHONY: build run test integration-test clean docker-up docker-down

build:
	go build -o bin/api ./cmd/api

run:
	go run ./cmd/api

test:
	go test -v ./...

integration-test:
	docker-compose up --build -d
	sleep 10
	go test -v -tags=integration ./...
	docker-compose down

clean:
	rm -rf bin/

docker-up:
	docker-compose up --build -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

migrate-up:
	docker-compose exec api ./main

migrate-down:
	MIGRATION_DIRECTION=down docker-compose exec api ./main
