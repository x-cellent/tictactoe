.DEFAULT_GOAL := run

.PHONY: run
run: build
	bin/tictactoe

.PHONY: build
build: gofmt protoc
	go mod download
	go mod tidy
	go build -tags netgo -o bin/tictactoe

.PHONY: gofmt
gofmt:
	GO111MODULE=off go fmt ./...

.PHONY: protoc
protoc:
	@protoc --go_out=plugins=grpc:. proto/*.proto || true
