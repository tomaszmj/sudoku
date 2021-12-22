package solver

import "github.com/tomaszmj/sudoku/board"

type Solver interface {
	Reset(b *board.Board)
	NextSolution() *board.Board
}
