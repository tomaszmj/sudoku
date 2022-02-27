package solver

import (
	"fmt"

	"github.com/tomaszmj/sudoku/board"
	"github.com/tomaszmj/sudoku/set"
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
	board.ForEach(func(x, y int, n uint16) {
		if n == 0 {
			fieldsToFill[field{x, y}] = struct{}{}
		}
	})
	b.board = board.Copy()
	b.solutions = nil
	b.fieldsToFill = fieldsToFill
}

// NextSolution in buruteforce case performs very brutal
// backtracking: for each field that was initially empty
// it tries to fill ANY number and recursively calls itself
// until all number are filled (without any initial validation!).
// When all numbers are filled, it validates board and
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
	numbersFound := set.New(b.board.Size())
	validateFunc := func(x, y int, n uint16) error {
		if n == 0 {
			return fmt.Errorf("field %d, %d is empty", x, y)
		}
		if !numbersFound.Add(int(n)) {
			return fmt.Errorf("number %d is repeated in row/column/subgrid at %d, %d", n, x, y)
		}
		return nil
	}
	err := b.board.Validate(validateFunc, numbersFound.Clear)
	return err == nil
}
