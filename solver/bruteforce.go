package solver

import (
	"fmt"

	"github.com/tomaszmj/sudoku/board"
)

type bruteforce struct {
	board        *board.Board
	solutions    []*board.Board
	fieldsToFill map[field]struct{}
}

type field struct {
	x int
	y int
}

// NewBruteforce returns Solver which is written in a very naive way
// and has poor performance. It is just a first step to have any
// working solution to be referred to when implementing proper solver.
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
	b.solutions = nil
	b.fieldsToFill = fieldsToFill
}

// NextSolution in buruteforce case performs very brutal
// backtracking: for each field that was initially empty
// it tries to fill ANY number and recursively calls itself
// until all number are filled (without any initial validation!).
// When all numbers are filled, it checks "boardIsValid" and
// if given solution has not been returned before.
// If solution is new, it is saved in b.solutions and returned.
// Otherwise, nil is returned.
func (b *bruteforce) NextSolution() *board.Board {
	for field := range b.fieldsToFill {
		if b.board.Get(field.x, field.y) != 0 {
			continue
		}
		for x := 1; x <= b.board.Size(); x++ {
			b.board.Set(field.x, field.y, uint16(x))
			ns := b.NextSolution()
			if ns != nil {
				return ns // return from recursive call
			}
			b.board.Set(field.x, field.y, 0)
		}
	}
	// break recursion and validate board state if all fieldsToFill are already set
	if b.boardIsValid() {
		if b.recordSolution() {
			return b.solutions[len(b.solutions)-1]
		}
	}
	return nil
}

// recordSolution checks if given solution has been encountered before.
// If it has been, it returns false. If it has not been, it is appended
// to b.solutions, current board state is reset, and true is returned.
func (b *bruteforce) recordSolution() bool {
	for _, solution := range b.solutions {
		if b.board.Equal(solution) {
			return false
		}
	}
	b.solutions = append(b.solutions, b.board.Copy())
	for field := range b.fieldsToFill {
		b.board.Set(field.x, field.y, 0)
	}
	return true
}

func (b *bruteforce) boardIsValid() bool {
	err := b.board.ForEachUntilError(func(x, y int) error {
		number := b.board.Get(x, y)
		if number == 0 {
			return fmt.Errorf("unfilled number")
		}

		rowValid := true
		b.board.ForEachInRow(y, func(x2, y2 int) {
			if x2 == x && y2 == y {
				return
			}
			if b.board.Get(x2, y2) == number {
				rowValid = false
			}
		})
		if !rowValid {
			return fmt.Errorf("row invalid")
		}

		columnValid := true
		b.board.ForEachInColumn(x, func(x2, y2 int) {
			if x2 == x && y2 == y {
				return
			}
			if b.board.Get(x2, y2) == number {
				columnValid = false
			}
		})
		if !columnValid {
			return fmt.Errorf("column invalid")
		}

		subgridValid := true
		b.board.ForEachInSubgrid(x, y, func(x2, y2 int) {
			if x2 == x && y2 == y {
				return
			}
			if b.board.Get(x2, y2) == number {
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
