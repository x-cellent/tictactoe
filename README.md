# TicTacToe via gRPC

This project serves as a demonstration of using gRPC with Go.
It implements the TicTacToe game between a grand master and several opponents, each of them playing a distinct game against the grand master in parallel.
They communicate via gRPC, whereas the grand master is implemented by a gRPC server and the opponents as gRPC clients.

For simplicity, there is no AI at all, i.e. the server as well as the clients draw randomly. Furthermore, neither the server nor the clients store the state of each game, instead the game board is exchanged and updated within each draw request-response cycle. The board validation and game rules are ensured by the server.

As an example, the client begins a new game by sending a draw request with an empty board alongside the intended client draw to the server.
The server then validates that the board is either empty or has as much client draws as server draws. It then validates that the intended client draw only affects an empty field. If either of these checks fails it responds with an invalid state indication alongside the given board. The client then tries another draw request.
Eventually there will be a valid client draw request and the server updates the board according to the intended client draw.
If the client wins the server responds accordingly. Otherwise it updates the board again with its own draw and responds with the updated board alongside an indication that the server has won or the game is drawn, if either happened.

As soon as a game has finished the client finally requests and prints the result.

## Goals

This scenario demonstrates the following gRPC features:
- classic request/response services like REST
- bidirectional request and response streams like WebSockets
- high parallelism (multi-threading) using simple goroutines without any synchronization overhead regarding gRPC
- gRPC definition files written in ProtoBuf (`*.proto`) as well as the corresponding Go client generation by `protoc` through a Docker container

This project also provides a Dockerfile that downloads and installs `protoc` as well as `protoc-gen-go`.

## Prerequisites

Docker

## Build

```shell
make
```

## Usage

```shell
bin/tictactoe -h
```

### Example with 1000 clients

```shell
bin/tictactoe -c 1000
```

### Example using in-memory gRPC connections

```shell
bin/tictactoe -m
```

## Blog article

This project is related to [this](https://www.x-cellent.com/blog/using-grpc-with-go/) blog article.
