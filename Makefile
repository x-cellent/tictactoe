.DEFAULT_GOAL := run

.PHONY: run
run: test
	bin/tictactoe

.PHONY: test
test: build
	go test ./...

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
	@protoc --go_out=plugins=grpc:. proto/*.proto
