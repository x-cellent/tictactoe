package tictactoe

import "fmt"

type Field = int32

const (
	Free Field = iota
	Client
	Server
)

type Board [9]Field

func Parse(fields []Field) (*Board, bool) {
	if len(fields) != 9 {
		return nil, false
	}

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

	b := new(Board)
	copy(b[:], fields)
	return b, true
}

func (b *Board) Draw(index int, field Field) bool {
	if !b.IsDrawValid(index, field) {
		return false
	}

	b[index] = field

	return true
}

func (b *Board) IsFinished() bool {
	return b.GetWinner() > 0
}

func (b *Board) PrintWinner() {
	fmt.Printf("%d %d %d\n%d %d %d\n%d %d %d\n", b[0], b[1], b[2], b[3], b[4], b[5], b[6], b[7], b[8])
	switch b.GetWinner() {
	case 1:
		println("Client (1) wins")
	case 2:
		println("Server (2) wins")
	case 3:
		println("Drawn")
	default:
		println("Game has not been decided yet")
	}
}

func (b *Board) IsDrawValid(index int, field Field) bool {
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

func (b *Board) GetWinner() int {
	// diagonals
	s := b[0] + b[4] + b[8]
	if s == 3 && b[0] == 1 && b[4] == 1 {
		return 1
	} else if s == 6 {
		return 2
	}

	// rows
	for i := 0; i < 9; i += 3 {
		s = b[i] + b[i+1] + b[i+2]
		if s == 3 && b[i] == 1 && b[i+1] == 1 {
			return 1
		} else if s == 6 {
			return 2
		}
	}

	// columns
	for i := 0; i < 3; i++ {
		s = b[i] + b[i+3] + b[i+6]
		if s == 3 && b[i] == 1 && b[i+3] == 1 {
			return 1
		} else if s == 6 {
			return 2
		}
	}

	for _, f := range b {
		if f == Free {
			return 0
		}
	}

	return 3
}
