package tictactoe

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBoard(t *testing.T) {
	b, ok := Parse([]Field{0, 0, 0, 0, 1, 0, 1, 2, 0})
	require.False(t, ok)
	require.Nil(t, b)

	b, ok = Parse([]Field{0, 0, 0, 0, 0, 0, 1, 2, 0})
	require.True(t, ok)
	require.NotNil(t, b)
	require.False(t, b.IsFinished())

	require.True(t, b.Draw(0, 1))
	require.False(t, b.IsFinished())

	require.False(t, b.Draw(1, 1))
	require.False(t, b.Draw(0, 2))
	require.True(t, b.Draw(1, 2))
	require.False(t, b.IsFinished())

	require.True(t, b.Draw(4, 1))
	require.False(t, b.IsFinished())

	require.True(t, b.Draw(3, 2))
	require.False(t, b.IsFinished())

	require.True(t, b.Draw(8, 1))
	require.True(t, b.IsFinished())

	b, ok = Parse([]Field{1, 2, 2, 2, 1, 0, 1, 1, 2})
	require.True(t, ok)
	require.NotNil(t, b)
	require.False(t, b.IsFinished())

	require.True(t, b.Draw(5, 1))
	require.Equal(t, 3, b.GetWinner())
}
