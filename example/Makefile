.PHONY: deps
deps:
	@echo "Deps"
	@brew install overmind
	@brew install golang-migrate
	@go get github.com/markbates/refresh
	@go get github.com/GeertJohan/go.rice/rice
	@go get github.com/volatiletech/sqlboiler
	@go get github.com/volatiletech/sqlboiler/drivers/sqlboiler-psql
	@go get github.com/andrewstucki/typescriptify-golang-structs/tscriptify
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.24.0
	@yarn --cwd frontend

.PHONY: dev
dev:
	@mkdir -p frontend/build
	overmind s

.PHONY: test
test:
	yarn --cwd frontend test --watchAll=false
	go test ./...

.PHONY: lint
lint:
	yarn --cwd frontend lint
	golangci-lint run

.PHONY: models
models:
	sqlboiler --struct-tag-casing camel --tag-ignore google_id --no-hooks --no-rows-affected --no-tests --wipe -c .sqlboiler.toml psql

.PHONY: typings
typings:
	tscriptify -package=example/payload -target=frontend/src/models.ts \
		APIError \
		ProfileResponse

.PHONY: migrate
migrate:
	migrate -database 'postgres://postgres:postgres@localhost:5433/example-dev?sslmode=disable' -path ./migrations up
	migrate -database 'postgres://postgres:postgres@localhost:5434/example-test?sslmode=disable' -path ./migrations up

.PHONY: build
build:
	yarn --cwd frontend build
	go generate
	go build

.DEFAULT_GOAL := dev
