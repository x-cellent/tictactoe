package server

import (
	"github.com/stretchr/testify/require"
	"github.com/x-cellent/tictactoe/pkg/v1/proto"
	"testing"
)

func TestBoard(t *testing.T) {
	b, ok := parse([]Field{0, 0, 0, 0, 1, 0, 1, 2, 0}, true)
	require.False(t, ok)
	require.Nil(t, b)

	b, ok = parse([]Field{0, 0, 0, 0, 0, 0, 1, 2, 0}, true)
	require.True(t, ok)
	require.NotNil(t, b)
	require.False(t, b.isFinished())

	require.True(t, b.draw(0, 1))
	require.False(t, b.isFinished())

	require.False(t, b.draw(1, 1))
	require.False(t, b.draw(0, 2))
	require.True(t, b.draw(1, 2))
	require.False(t, b.isFinished())

	require.True(t, b.draw(4, 1))
	require.False(t, b.isFinished())

	require.True(t, b.draw(3, 2))
	require.False(t, b.isFinished())

	require.True(t, b.draw(8, 1))
	require.True(t, b.isFinished())

	b, ok = parse([]Field{1, 2, 2, 2, 1, 0, 1, 1, 2}, true)
	require.True(t, ok)
	require.NotNil(t, b)
	require.False(t, b.isFinished())

	require.True(t, b.draw(5, 1))
	require.Equal(t, proto.DrawResponse_DRAWN, b.getWinner())
}
