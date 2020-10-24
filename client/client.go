package client

import (
	"context"
	"fmt"
	"github.com/x-cellent/tictactoe/pkg/v1/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"io"
	"math/rand"
	"net"
	"sync"

	"time"
)

func Connect(clientID int, listener net.Listener, inMemory bool) {
	rand.Seed(time.Now().Unix())

	// dial server
	var conn *grpc.ClientConn
	var err error

	if inMemory {
		bufDialer := func(context.Context, string) (net.Conn, error) {
			return listener.(*bufconn.Listener).Dial()
		}
		conn, err = grpc.Dial("bufnet", grpc.WithContextDialer(bufDialer), grpc.WithBlock(), grpc.WithInsecure())
	} else {
		conn, err = grpc.Dial(":50005", grpc.WithBlock(), grpc.WithInsecure())
	}

	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// create stream
	client := proto.NewTicTacToeClient(conn)
	stream, err := client.Play(context.Background())
	if err != nil {
		panic(err)
	}

	// play game...
	boardMutex := new(sync.RWMutex)
	keepDrawing := true
	board := make([]int32, 9)
	wg := new(sync.WaitGroup)
	wg.Add(2)

	// first goroutine subsequently sends random numbers within [0,12), which eventually results in valid client draws
	go func() {
		defer wg.Done()
		for keepDrawing {
			boardMutex.RLock()
			req := proto.DrawRequest{
				Board: &proto.Board{Fields: board},
				Draw:  int32(rand.Intn(12)),
			}
			boardMutex.RUnlock()

			err := stream.Send(&req)
			if err != nil {
				panic(err)
			}

			// sleep some random time
			d := time.Duration(rand.Intn(20) + 2)
			time.Sleep(d * time.Millisecond)
		}

		err := stream.CloseSend()
		if err != nil {
			fmt.Println(err)
		}
	}()

	// second goroutine receives server responses, which indicates either an invalid client draw or the actual board including state information
	go func() {
		defer wg.Done()
		for keepDrawing {
			resp, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				panic(err)
			}

			if resp.State != proto.DrawResponse_INVALID {
				keepDrawing = resp.State == proto.DrawResponse_NOT_FINISHED
				boardMutex.Lock()
				copy(board, resp.Board.Fields[:])
				boardMutex.Unlock()
			}
		}
	}()

	wg.Wait()

	result := fmt.Sprintf("Client %d finished game:\n", clientID)
	boardMutex.RLock()
	defer boardMutex.RUnlock()
	resp, err := client.Result(context.Background(), &proto.Board{Fields: board})
	if err != nil {
		panic(err)
	}
	result += resp.Text
	println(result)
}
