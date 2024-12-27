# Include variables from the .envrc file
# include .envrc

DB_DSN = postgres://digilletsadmin:passwordllets@localhost/digillets?sslmode=disable
current_time = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
linker_flags = '-s -X main.buildTime=${current_time}'

# all:
# 	@echo ${linker_flags} # prints Hello Make

#==================================================================================== #
# HELPERS
#==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

# Create the new confirm target.
.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

#==================================================================================== #
# DEVELOPMENT
#==================================================================================== #

## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	go run ./cmd/

## db/psql: connect to the database using psql, not setup yet
.PHONY: db/psql
db/psql:
	psql ${DB_DSN}

## db/migrations/version: check migrations versions of DB
.PHONY: db/migrations/version
db/migrations/version:
	migrate -path=./migrations -database=${DB_DSN} version

## db/migrations/new name=$1: create a new database migration eg make migration name=testing
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up: apply all up database migrations. confirm included as prerequisite.
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path=./migrations -database ${DB_DSN} up

## db/migrations/goto position=$1: takes database to particular migrations. confirm included as prerequisite.
.PHONY: db/migrations/goto
db/migrations/goto: confirm
	@echo 'Running up migrations...'
	migrate -path=./migrations -database ${DB_DSN} goto ${position}

## db/migrations/goto position=$1: takes database to particular migrations. confirm included as prerequisite.
.PHONY: db/migrations/force
db/migrations/force: confirm
	@echo 'Running up migrations...'
	migrate -path=./migrations -database ${DB_DSN} force ${position}

#==================================================================================== #
# QUALITY CONTROL
#==================================================================================== #

## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Vendoring dependencies...'
	go mod vendor

#==================================================================================== #
# BUILD
#==================================================================================== #


# current_time = $(shell date --iso-8601=seconds)
# linker_flags = '-s -X main.buildTime=${current_time}'
# go build -ldflags=${linker_flags} -o=./bin/api ./cmd/api
# GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flags} -o=./bin/linux_amd64/api ./cmd/api
# GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64/api ./cmd/api

## build/api: build the cmd/api application
.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	go build -ldflags='-s' -o=./bin/api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/application ./cmd/api
	@echo 'Building completed...'

## build/zip: zip the application for aws deployment
.PHONY: build/zip
build/zip:
	rm api-app.zip
	rm -rf bin
	@echo 'Building cmd/api...'
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/application ./cmd/api
	@echo 'Building completed...'
	@echo 'zipping app...'
	zip -r api-app.zip bin Procfile go.mod
	@echo 'zipping completed...'
