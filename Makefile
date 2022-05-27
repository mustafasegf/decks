install:
	go mod tidy

run:
	go run main.go

dev:
	air

build:
	go build -o ./bin/main main.go

run-build:
	./bin/main

test:
	go test -cover -v ./pkg/...

swagger:
	swag init -g main.go 

up:
	docker compose up -d
	docker compose logs -f

upb:
	docker compose up --build -d
	docker compose logs -f

down:
	docker compose down

updb:
	docker compose up -d db 
	docker compose logs -f

tool:
	sh ./tools/install.sh

migration:
	migrate create -seq -ext sql -dir migrations $(filter-out $@,$(MAKECMDGOALS))

migrate:
	cd ./migrations/script/ && \
	go run migrate.go

migrate-down:
	cd ./migrations/script/ && \
	go run migrate.go -action down

.PHONY: build
