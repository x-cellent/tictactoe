package server

import (
	"fmt"
	"gitlab.com/kolls/networking/grpc/proto"
	"gitlab.com/kolls/networking/grpc/tictactoe"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"io"
	"log"
	"math/rand"
	"sync"
	"time"
)

func Run() *bufconn.Listener {
	rand.Seed(time.Now().Unix())

	// create in-memory listener
	listener := bufconn.Listen(1024 * 1024)

	// create gRPC server
	s := grpc.NewServer()
	proto.RegisterTicTacToeServer(s, &server{
		board: new(tictactoe.Board),
	})

	// and start...
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		wg.Done()
		log.Fatal(s.Serve(listener))
	}()
	wg.Wait()
	return listener
}

type server struct {
	board *tictactoe.Board
}

func (s server) Game(srv proto.TicTacToe_GameServer) error {
	ctx := srv.Context()

	for {
		// exit if context is done or continue
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// receive (potential invalid) draw from stream
		req, err := srv.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			fmt.Printf("receive error %v\n", err)
			continue
		}
		clientDraw := int(req.Draw)

		resp := proto.Response{}

		// parse client board and check if client draw is valid
		board, ok := tictactoe.Parse(req.Board)
		if !ok || board.IsFinished() {
			resp.State = proto.Response_INVALID
		} else if !board.Draw(clientDraw, tictactoe.Client) {
			resp.State = proto.Response_INVALID
		} else if board.IsFinished() {
			resp.State = proto.Response_State(board.GetWinner()) // may also be drawn
		} else {
			// make server draw
			for {
				serverDraw := rand.Intn(9)
				if board.Draw(serverDraw, tictactoe.Server) {
					break
				}
			}
			// update state
			if board.IsFinished() {
				resp.State = proto.Response_State(board.GetWinner()) // may also be drawn
			} else {
				resp.State = proto.Response_NOT_FINISHED
			}
		}

		if ok {
			resp.Board = board[:]
		}

		err = srv.Send(&resp)
		if err != nil {
			fmt.Printf("send error %v\n", err)
		}
	}
}
