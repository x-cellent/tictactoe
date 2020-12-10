package server

import (
	"fmt"
	"github.com/x-cellent/tictactoe/pkg/v1/tictactoe"
)

type Field = int32

const (
	Free Field = iota
	Client
	Server
)

type board [9]Field

func parse(fields []Field, strict bool) (*board, bool) {
	if len(fields) != 9 {
		return nil, false
	}

	if strict {
		n1, n2 := 0, 0
		for _, f := range fields {
			switch f {
			case Client:
				n1++
			case Server:
				n2++
			}
		}

		if n1 != n2 {
			return nil, false
		}
	}

	b := new(board)
	copy(b[:], fields)
	return b, true
}

func (b *board) draw(index int, field Field) bool {
	if !b.isDrawValid(index, field) {
		return false
	}

	b[index] = field

	return true
}

func (b *board) isFinished() bool {
	w := b.getWinner()
	return w == tictactoe.DrawResponse_DRAWN || w == tictactoe.DrawResponse_SERVER_WINS || w == tictactoe.DrawResponse_CLIENT_WINS
}

func (b *board) getResult() string {
	r := fmt.Sprintf("%d %d %d\n%d %d %d\n%d %d %d\n", b[0], b[1], b[2], b[3], b[4], b[5], b[6], b[7], b[8])
	switch b.getWinner() {
	case tictactoe.DrawResponse_CLIENT_WINS:
		r += "Client wins"
	case tictactoe.DrawResponse_SERVER_WINS:
		r += "Server wins"
	case tictactoe.DrawResponse_DRAWN:
		r += "Drawn"
	default:
		r += "Game has not been decided yet"
	}
	return r
}

func (b *board) isDrawValid(index int, field Field) bool {
	if index < 0 || index >= len(b) || field == Free || b[index] != Free {
		return false
	}

	n1, n2 := 0, 0
	for _, f := range b {
		switch f {
		case Client:
			n1++
		case Server:
			n2++
		}
	}
	switch field {
	case Client:
		if n1 > n2 {
			return false
		}
	case Server:
		if n2 > n1 {
			return false
		}
	}

	return true
}

func (b *board) getWinner() tictactoe.DrawResponse_State {
	// diagonals
	s := b[0] + b[4] + b[8]
	if s == 3 && b[0] == 1 && b[4] == 1 {
		return tictactoe.DrawResponse_CLIENT_WINS
	} else if s == 6 {
		return tictactoe.DrawResponse_SERVER_WINS
	}

	// rows
	for i := 0; i < 9; i += 3 {
		s = b[i] + b[i+1] + b[i+2]
		if s == 3 && b[i] == 1 && b[i+1] == 1 {
			return tictactoe.DrawResponse_CLIENT_WINS
		} else if s == 6 {
			return tictactoe.DrawResponse_SERVER_WINS
		}
	}

	// columns
	for i := 0; i < 3; i++ {
		s = b[i] + b[i+3] + b[i+6]
		if s == 3 && b[i] == 1 && b[i+3] == 1 {
			return tictactoe.DrawResponse_CLIENT_WINS
		} else if s == 6 {
			return tictactoe.DrawResponse_SERVER_WINS
		}
	}

	for _, f := range b {
		if f == Free {
			return tictactoe.DrawResponse_NOT_FINISHED
		}
	}

	return tictactoe.DrawResponse_DRAWN
}
