package board_test

import (
	"strings"
	"testing"

	"github.com/tomaszmj/sudoku/board"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBoardNewAndString(t *testing.T) {
	t.Run("subgrid size is limited", func(t *testing.T) {
		_, err := board.New(board.MaxSize+1, 1)
		assert.Error(t, err)
	})
	t.Run("board size is limited", func(t *testing.T) {
		_, err := board.New(board.MaxSize-1, board.MaxSize-1)
		assert.Error(t, err)
	})
	t.Run("board is created empty", func(t *testing.T) {
		board, err := board.New(3, 2)
		require.NoError(t, err)
		assert.Equal(t, board3x2zeros, board.String())
	})
}

func TestBoardGetAndSet(t *testing.T) {
	t.Run("board is created empty", func(t *testing.T) {
		board, err := board.New(3, 2)
		require.NoError(t, err)
		for y := 0; y < 2; y++ {
			for x := 0; x < 3; x++ {
				assert.Equal(t, uint16(0), board.Get(x, y))
			}
		}
	})
	t.Run("Get returns value previously Set", func(t *testing.T) {
		board, err := board.New(3, 2)
		require.NoError(t, err)
		board.Set(2, 1, 6)
		assert.Equal(t, uint16(6), board.Get(2, 1))
	})
	t.Run("Set panics when value set is larger than board size", func(t *testing.T) {
		board, err := board.New(3, 2)
		require.NoError(t, err)
		assert.Panics(t, func() {
			board.Set(2, 1, 7)
		})
	})
}

func TestBoardForEachInSubgrid(t *testing.T) {
	board, err := board.New(3, 2)
	require.NoError(t, err)
	board.ForEachInSubgrid(1, 3, func(x, y int) {
		board.Set(x, y, 1)
	})
	board.ForEachInSubgrid(3, 0, func(x, y int) {
		board.Set(x, y, 2)
	})
	assert.Equal(t, board3x2partiallyFilled, board.String())
}

func TestBoardForEachNeighbour(t *testing.T) {
	board, err := board.New(3, 2)
	require.NoError(t, err)
	board.ForEachNeighbour(2, 3, func(x, y int) {
		if board.Get(x, y) == 0 {
			board.Set(x, y, 1)
		} else {
			board.Set(x, y, 2) // should not happen
		}
	})
	assert.Equal(t, board3x2NeigbourFilled, board.String())
}

func TestBoardHaveCommonSubgrid(t *testing.T) {
	board, err := board.New(3, 3)
	require.NoError(t, err)
	assert.True(t, board.HaveCommonSubgrid(0, 1, 2, 2))
	assert.False(t, board.HaveCommonSubgrid(0, 0, 3, 0))
}

func TestBoardNewFromSerializedFormat(t *testing.T) {
	t.Run("board can be recreated from string", func(t *testing.T) {
		board1, err := board.New(3, 2)
		require.NoError(t, err)
		board1.Set(2, 1, 3)
		board1.Set(1, 2, 4)
		var serizalizeOutput strings.Builder
		board1.Serialize(&serizalizeOutput)
		serializedStr := serizalizeOutput.String()
		board2, err := board.NewFromSerializedFormat(strings.NewReader(serializedStr))
		require.NoError(t, err)
		assert.Equal(t, board1, board2)
	})
	t.Run("smallest possible board", func(t *testing.T) {
		board1, err := board.NewFromSerializedFormat(strings.NewReader("1 1\n1"))
		require.NoError(t, err)
		assert.Equal(t, uint16(1), board1.Get(0, 0))
	})
	t.Run("invalid first line format", func(t *testing.T) {
		_, err := board.NewFromSerializedFormat(strings.NewReader("1 \n"))
		assert.Error(t, err)
	})
	t.Run("invalid board data format", func(t *testing.T) {
		_, err := board.NewFromSerializedFormat(strings.NewReader("2 2\n0\n"))
		assert.Error(t, err)
	})
	t.Run("too many rows", func(t *testing.T) {
		_, err := board.NewFromSerializedFormat(strings.NewReader("2 1\n0 0\n0 0\n0 0\n"))
		assert.Error(t, err)
	})
	t.Run("too little rows", func(t *testing.T) {
		_, err := board.NewFromSerializedFormat(strings.NewReader("2 1\n0 0\n"))
		assert.Error(t, err)
	})
	t.Run("too long row", func(t *testing.T) {
		_, err := board.NewFromSerializedFormat(strings.NewReader("1 1\n0 0 0"))
		assert.Error(t, err)
	})
	t.Run("number parsing error despite matching regex", func(t *testing.T) {
		_, err := board.NewFromSerializedFormat(strings.NewReader("1 1\n99999999999999999999999999999999"))
		assert.Error(t, err)
	})
	t.Run("invalid number (greater than board size)", func(t *testing.T) {
		_, err := board.NewFromSerializedFormat(strings.NewReader("1 1\n2"))
		assert.Error(t, err)
	})
	t.Run("invalid board size", func(t *testing.T) {
		_, err := board.NewFromSerializedFormat(strings.NewReader("1000000 1000000\n"))
		assert.Error(t, err)
	})
}

// BenchmarkBoardString is just for fun. I checked if using strings.Builder improves performance - it does
func BenchmarkBoardString(b *testing.B) {
	board, err := board.New(50, 50)
	if err != nil {
		b.Fatal(err)
	}
	for n := 0; n < b.N; n++ {
		_ = board.String()
	}
}

const board3x2zeros = `+-------+-------+
| 0 0 0 | 0 0 0 |
| 0 0 0 | 0 0 0 |
+-------+-------+
| 0 0 0 | 0 0 0 |
| 0 0 0 | 0 0 0 |
+-------+-------+
| 0 0 0 | 0 0 0 |
| 0 0 0 | 0 0 0 |
+-------+-------+
`

const board3x2partiallyFilled = `+-------+-------+
| 0 0 0 | 2 2 2 |
| 0 0 0 | 2 2 2 |
+-------+-------+
| 1 1 1 | 0 0 0 |
| 1 1 1 | 0 0 0 |
+-------+-------+
| 0 0 0 | 0 0 0 |
| 0 0 0 | 0 0 0 |
+-------+-------+
`

const board3x2NeigbourFilled = `+-------+-------+
| 0 0 1 | 0 0 0 |
| 0 0 1 | 0 0 0 |
+-------+-------+
| 1 1 1 | 0 0 0 |
| 1 1 0 | 1 1 1 |
+-------+-------+
| 0 0 1 | 0 0 0 |
| 0 0 1 | 0 0 0 |
+-------+-------+
`
