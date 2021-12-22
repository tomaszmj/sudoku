package board_test

import (
	"testing"

	"github.com/tomaszmj/sudoku/board"
)

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
