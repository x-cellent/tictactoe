package server

import (
	"context"
	"fmt"
	"gitlab.com/kolls/networking/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"io"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"
)

func Run(inMemory bool) net.Listener {
	rand.Seed(time.Now().Unix())

	s := grpc.NewServer()
	proto.RegisterTicTacToeServer(s, &server{})

	var listener net.Listener
	if inMemory {
		listener = bufconn.Listen(1024 * 1024)
	} else {
		var err error
		listener, err = net.Listen("tcp", ":50005")
		if err != nil {
			panic(err)
		}
	}

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
}

func (s *server) Play(srv proto.TicTacToe_PlayServer) error {
	ctx := srv.Context()

	for {
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

		resp := proto.DrawResponse{}

		board, ok := parse(req.Board.Fields, true)
		if !ok {
			println("not ok")
			resp.State = proto.DrawResponse_INVALID
		} else if board.isFinished() {
			println("already finished")
			resp.State = board.getWinner()
		} else if !board.draw(clientDraw, Client) {
			resp.State = proto.DrawResponse_INVALID
		} else if board.isFinished() {
			resp.State = board.getWinner() // may also be drawn here
		} else {
			// make server draw
			for {
				serverDraw := rand.Intn(9)
				if board.draw(serverDraw, Server) {
					break
				}
			}
			if board.isFinished() {
				resp.State = proto.DrawResponse_SERVER_WINS
			} else {
				resp.State = proto.DrawResponse_NOT_FINISHED
			}
		}

		if resp.State != proto.DrawResponse_INVALID {
			resp.Board = &proto.Board{Fields: board[:]}
		}

		err = srv.Send(&resp)
		if err != nil {
			fmt.Printf("send error %v\n", err)
		}
	}
}

func (s *server) Result(ctx context.Context, board *proto.Board) (*proto.ResultResponse, error) {
	resp := &proto.ResultResponse{}
	b, ok := parse(board.Fields, false)
	if !ok {
		resp.Text = "invalid board"
	} else {
		resp.Text = b.getResult()
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return resp, nil
	}
}
