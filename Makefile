.DEFAULT_GOAL := build

.PHONY: test
test: build
	@docker run --rm -v $(PWD):/go/src -w /go/src golang:1.15.3-buster go test ./...

.PHONY: build
build: gofmt protoc
	@docker run --rm -v $(PWD):/go/src -w /go/src golang:1.15.3-buster go build -tags netgo -o bin/tictactoe

.PHONY: gofmt
gofmt:
	@docker run --rm --user $$(id -u):$$(id -g) -v $(PWD):/go/src -w /go/src -e GO111MODULE=off golang:1.15.3-buster go fmt ./...

.PHONY: protoc
protoc:
	@docker build -t protoc .
	@docker run --rm --user $$(id -u):$$(id -g) -v $(PWD):/work -w /work protoc protoc -I pkg --go_out=plugins=grpc:pkg pkg/v1/proto/*.proto
