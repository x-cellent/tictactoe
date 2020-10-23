package client

import (
	"context"
	"fmt"
	"gitlab.com/kolls/networking/grpc/proto"
	"gitlab.com/kolls/networking/grpc/tictactoe"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"io"
	"math/rand"
	"net"
	"sync"

	"time"
)

func Run(clientID int, listener *bufconn.Listener) {
	rand.Seed(time.Now().Unix())

	// dial server
	bufDialer := func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
	conn, err := grpc.Dial("bufnet", grpc.WithContextDialer(bufDialer), grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// create stream
	client := proto.NewTicTacToeClient(conn)
	stream, err := client.Game(context.Background())
	if err != nil {
		panic(err)
	}

	ctx := stream.Context()
	done := make(chan bool)

	boardMutex := new(sync.RWMutex)
	board := new(tictactoe.Board)

	// first goroutine subsequently sends random numbers within [0,12), which eventually results in valid client draws
	go func() {
		for {
			boardMutex.RLock()
			req := proto.Request{
				Board: board[:],
				Draw:  int32(rand.Intn(12)),
			}
			boardMutex.RUnlock()

			err := stream.Send(&req)
			if err != nil {
				panic(err)
			}

			boardMutex.RLock()
			if board.IsFinished() {
				boardMutex.RUnlock()
				break
			}
			boardMutex.RUnlock()

			d := time.Duration(rand.Intn(20) + 2)
			time.Sleep(d * time.Millisecond)
		}

		err := stream.CloseSend()
		if err != nil {
			fmt.Println(err)
		}
	}()

	// second goroutine receives server responses, which indicates either an invalid client draw or a valid board with state information
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				panic(err)
			}

			if resp.State != proto.Response_INVALID {
				boardMutex.Lock()
				copy(board[:], resp.Board[:])
				boardMutex.Unlock()
			}

			switch resp.State {
			case proto.Response_CLIENT_WINS:
				boardMutex.RLock()
				if board.GetWinner() != 1 {
					panic("winner mismatch, client (1) should have won")
				}
				boardMutex.RUnlock()
			case proto.Response_SERVER_WINS:
				boardMutex.RLock()
				if board.GetWinner() != 2 {
					panic("winner mismatch, server (2) should have won")
				}
				boardMutex.RUnlock()
			}
		}
	}()

	// third goroutine closes done channel if context is done
	go func() {
		<-ctx.Done()
		close(done)
	}()

	<-done

	fmt.Printf("Client %d finished game:\n", clientID)
	board.PrintWinner()
}
