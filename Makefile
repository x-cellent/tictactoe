.DEFAULT_GOAL := help

.PHONY: help
help: test
	@echo
	@bin/tictactoe -h

.PHONY: test
test: build
	@docker run --rm -v $(PWD):/go/src -w /go/src golang:1.15.6-buster go test ./...

.PHONY: build
build: gofmt protoc
	@docker run --rm -v $(PWD):/go/src -w /go/src golang:1.15.6-buster go build -tags netgo -o bin/tictactoe

.PHONY: gofmt
gofmt:
	@docker run --rm --user $$(id -u):$$(id -g) -v $(PWD):/go/src -w /go/src -e GO111MODULE=off golang:1.15.6-buster go fmt ./...

.PHONY: protoc
protoc:
	@docker build -t protoc .
	@docker run --rm --user $$(id -u):$$(id -g) -v $(PWD):/work -w /work protoc protoc -I pkg \
	    --go_out=pkg --go_opt=paths=source_relative \
        --go-grpc_out=require_unimplemented_servers=true:pkg --go-grpc_opt=paths=source_relative \
        pkg/v1/tictactoe/*.proto
