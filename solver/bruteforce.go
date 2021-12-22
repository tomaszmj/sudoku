package solver

import (
	"fmt"

	"github.com/tomaszmj/sudoku/board"
)

type bruteforce struct {
	board        *board.Board
	lastSolution *board.Board
	fieldsToFill map[field]struct{}
}

type field struct {
	x int
	y int
}

func NewBruteforce() Solver {
	return &bruteforce{}
}

func (b *bruteforce) Reset(board *board.Board) {
	fieldsToFill := make(map[field]struct{})
	board.ForEachUntilError(func(x, y int) error {
		if board.Get(x, y) == 0 {
			fieldsToFill[field{x, y}] = struct{}{}
		}
		return nil
	})
	b.board = board.Copy()
	b.lastSolution = nil
	b.fieldsToFill = fieldsToFill
}

func (b *bruteforce) NextSolution() *board.Board {
	if BoardValid(b.board) {
		return b.board.Copy()
	}
	return nil
}

func BoardValid(b *board.Board) bool {
	err := b.ForEachUntilError(func(x, y int) error {
		number := b.Get(x, y)
		if number == 0 {
			return fmt.Errorf("unfilled number")
		}

		rowValid := true
		b.ForEachInRow(y, func(x2, y2 int) {
			if x2 == x && y2 == y {
				return
			}
			if b.Get(x2, y2) == number {
				rowValid = false
			}
		})
		if !rowValid {
			return fmt.Errorf("row invalid")
		}

		columnValid := true
		b.ForEachInColumn(x, func(x2, y2 int) {
			if x2 == x && y2 == y {
				return
			}
			if b.Get(x2, y2) == number {
				columnValid = false
			}
		})
		if !columnValid {
			return fmt.Errorf("column invalid")
		}

		subgridValid := true
		b.ForEachInSubgrid(x, y, func(x2, y2 int) {
			if x2 == x && y2 == y {
				return
			}
			if b.Get(x2, y2) == number {
				subgridValid = false
			}
		})
		if !subgridValid {
			return fmt.Errorf("subgrid invalid")
		}

		return nil
	})
	return err == nil
}
