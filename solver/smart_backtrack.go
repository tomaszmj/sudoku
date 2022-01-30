package solver

import (
	"container/heap"
	"fmt"

	"github.com/tomaszmj/sudoku/board"
	"github.com/tomaszmj/sudoku/set"
)

type smartBacktrack struct {
	board                *board.Board
	originalFieldsToFill []fieldToFill
	currentFieldsToFill fieldsToFillHeap
	solvable bool
}

func NewSmartBarcktrack() Solver {
	return &smartBacktrack{}
}

func (s *smartBacktrack) Reset(board *board.Board) {
	s.solvable = true
	s.board = board.Copy()
	fieldsToFill := make([]fieldToFill, 0)
	board.ForEachUntilError(func(x, y int) error {
		if board.Get(x, y) == 0 {
			availableNumbers, err := s.findPossibleNumbers(x, y)
			if err != nil {
				s.solvable = false
			}
			fieldsToFill = append(fieldsToFill, fieldToFill{x: x, y: y, possibleValues: availableNumbers})
		}
		return nil
	})
	s.originalFieldsToFill = fieldsToFill
	currentFieldsToFill := make([]fieldToFill, len(fieldsToFill))
	copy(currentFieldsToFill, fieldsToFill)
	s.currentFieldsToFill = fieldsToFillHeap(currentFieldsToFill)
	heap.Init(&s.currentFieldsToFill)
}

func (s *smartBacktrack) NextSolution() *board.Board {
	if !s.solvable {
		return nil
	}
	for len(s.currentFieldsToFill) > 0 {
		f := heap.Pop(&s.currentFieldsToFill).(fieldToFill)
		if f.possibleValues.Len() == 0 {
			return nil
		}
		if f.possibleValues.Len() == 1 {
			n := f.possibleValues.ForEach(func(n int) bool {
				return true
			})
			s.setNumber(f.x, f.y, uint16(n))
		} else {
			fmt.Printf("got %d choices for (%d, %d) - " +
				"backtracking not implemented yet :(\n", f.possibleValues.Len(), f.x, f.y)
			return nil
		}
		//TODO proper backtracking
		//var solution *board.Board
		//n := f.possibleValues.ForEach(func(n int) bool {
		//	s.setNumber(f.x, f.y, uint16(n))
		//	ns := s.NextSolution() // recursive call
		//	if ns != nil { // ok - return from recursive call
		//		solution = ns
		//		return true
		//	}
		//	// else - backtrack - revert selecting number for given field
		//	s.resetNumber(f.x, f.y)
		//	return false
		//})
		//if solution != nil {
		//	return solution // ok - solution found
		//}
		//// else - backtrack - revert selecting field
		//heap.Push(&s.currentFieldsToFill, f)
	}
	s.solvable = false // ensure solution is returned only once
	return s.board.Copy()
}


func (s *smartBacktrack) setNumber(x, y int, n uint16) {
	s.board.Set(x, y, n)
	toRemove := -1
	sortNeeded := false
	for i, f := range s.currentFieldsToFill {
		// if given field is still in currentFieldsToFill list - mark to be removed
		if f.x == x && f.y == y {
			toRemove = i
			continue
		}
		// if field is in the same row / column / subgrid as changed field,
		// if set of possibleVelues must be updated
		if f.x == x || f.y == y || s.board.HaveCommonSubgrid(x, y, f.x, f.y) {
			sortNeeded = sortNeeded || f.possibleValues.Remove(int(n))
		}
	}
	if toRemove >= 0 {
		heap.Remove(&s.currentFieldsToFill, toRemove)
	}
	// TODO we can use heap.Fix only for changed fields if each field "knows" its queue index
	// (but it is just an optimization that can be done later if needed)
	if sortNeeded {
		heap.Init(&s.currentFieldsToFill)
	}
}

func (s *smartBacktrack) resetNumber(x, y int) {
	n := s.board.Get(x, y)
	s.board.Set(x, y, 0)
	possibleNumbers, _ := s.findPossibleNumbers(x, y)
	sortNeeded := false
	for _, f := range s.currentFieldsToFill {
		if f.x == x || f.y == y || s.board.HaveCommonSubgrid(x, y, f.x, f.y) {
			f.possibleValues.Add(int(n))
			sortNeeded = true
		}
	}
	newField := fieldToFill{
		x:              x,
		y:              y,
		possibleValues: possibleNumbers,
	}
	if sortNeeded {
		s.currentFieldsToFill.Push(newField)
		heap.Init(&s.currentFieldsToFill)
	} else {
		heap.Push(&s.currentFieldsToFill, newField)
	}
}

func (s *smartBacktrack) findPossibleNumbers(x, y int) (*set.Set, error) {
	size := s.board.Size()
	row := set.New(size)
	col := set.New(size)
	subgrid := set.New(size)
	var err error
	s.board.ForEachInRow(y, func(x, y int) {
		if n := s.board.Get(x, y); n != 0 {
			if !row.Add(int(n)) {
				err = fmt.Errorf("number %d is repeated in row %d", n, y)
			}
		}
	})
	s.board.ForEachInColumn(x, func(x, y int) {
		if n := s.board.Get(x, y); n != 0 {
			if !col.Add(int(n)) {
				err = fmt.Errorf("number %d is repeated in column %d", n, x)
			}
		}
	})
	s.board.ForEachInSubgrid(x, y, func(x, y int) {
		if n := s.board.Get(x, y); n != 0 {
			if !subgrid.Add(int(n)) {
				err = fmt.Errorf("number %d is repeated in subgrid for field (%d, %d)", n, x, y)
			}
		}
	})
	if err != nil {
		return set.New(size), err
	}
	allForbiddenNumbers := set.Union(row, col, subgrid)
	return allForbiddenNumbers.Complement(), nil
}
