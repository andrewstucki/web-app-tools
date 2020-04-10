.PHONY: deps
deps:
	@echo "Deps"
	@brew install overmind
	@go get github.com/markbates/refresh
	@go get github.com/GeertJohan/go.rice/rice
	@go get github.com/volatiletech/sqlboiler
	@go get github.com/volatiletech/sqlboiler/drivers/sqlboiler-psql
	@yarn --cwd frontend

.PHONY: dev
dev:
	@mkdir -p frontend/build
	overmind s

.PHONY: test
test:
	yarn --cwd frontend test --watchAll=false
	go test ./...

.PHONY: build
build:
	yarn --cwd frontend build
	go generate
	go build

.DEFAULT_GOAL := dev
