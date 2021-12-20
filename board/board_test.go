package board_test

import (
	"testing"

	"github.com/tomaszmj/sudoku/board"
)

func BenchmarkBoardString(b *testing.B) {
	board, err := board.New(6, 8)
	if err != nil {
		b.Fatal(err)
	}
	for n := 0; n < b.N; n++ {
		_ = board.String()
	}
}
