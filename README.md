# Card Rest API

## Requirement

- Go 1.18
- Docker
- Docker-compose

## Dev requirement
- cosmtrek/air
- swaggo/swag
- momaek/formattag
- golang-migrate/migrate
- GNU Make

## Instalation

- githook
```
git config core.hooksPath .githooks
```
- cosmtrek/air
```
go install github.com/cosmtrek/air@latest
```

- swaggo/swag
```
go install github.com/swaggo/swag/cmd/swag@latest
```

- momaek/formattag
```
go install github.com/momaek/formattag@latest
```

- golang-migrate/migrate
```
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

- GNU Make (on windows need chocolatey)
Make is on linux by default
```
choco install make
```

- Install dependencies
```
make install
```

## Database
- Create Migration
```
migrate create -seq -ext sql -dir migrations [name] or
make migration [name]
```

- Migration Up
```
make migrate
```

- Migration Down
```
make migrate-down
```

## Running
- Runing
```
make updb
make migrate
make down
make up
```
the open localhost:3000/swagger

