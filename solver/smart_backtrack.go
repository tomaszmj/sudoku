package solver

import (
	"github.com/tomaszmj/sudoku/board"
	"github.com/tomaszmj/sudoku/set"
)

type smartBacktrack struct {
	board                *board.Board
	originalFieldsToFill []fieldToFill
}

func NewSmartBarcktrack() Solver {
	return &bruteforce{}
}

func (s *smartBacktrack) Reset(board *board.Board) {
	fieldsToFill := make([]fieldToFill, 0)
	board.ForEachUntilError(func(x, y int) error {
		if board.Get(x, y) == 0 {
			availableNumbers := s.findPossibleNumbers(x, y)
			fieldsToFill = append(fieldsToFill, fieldToFill{x: x, y: y, possibleValues: availableNumbers})
		}
		return nil
	})
	s.board = board.Copy()
	s.originalFieldsToFill = fieldsToFill
}

func (s *smartBacktrack) NextSolution() *board.Board {
	return nil
}

func (s *smartBacktrack) findPossibleNumbers(x, y int) *set.Set {
	size := s.board.Size()
	row := set.New(size)
	col := set.New(size)
	subgrid := set.New(size)
	s.board.ForEachInRow(y, func(x, y int) {
		if n := s.board.Get(x, y); n != 0 {
			row.Add(int(n))
		}
	})
	s.board.ForEachInColumn(x, func(x, y int) {
		if n := s.board.Get(x, y); n != 0 {
			col.Add(int(n))
		}
	})
	s.board.ForEachInSubgrid(x, y, func(x, y int) {
		if n := s.board.Get(x, y); n != 0 {
			subgrid.Add(int(n))
		}
	})
	allForbiddenNumbers := set.Union(row, col, subgrid)
	return allForbiddenNumbers.Complement()
}
