package solver

import (
	"container/heap"
	"fmt"

	"github.com/tomaszmj/sudoku/board"
	"github.com/tomaszmj/sudoku/set"
)

type smartBacktrack struct {
	board               *board.Board
	currentFieldsToFill fieldsToFillHeap
	solvable            bool
	// Let's say we take from currentFieldsToFill
	// field (1,1) with possible numbers {1,2}.
	// We put 1 on board and push remaining options to the stack, i.e. (1,1) {2}.
	// If there are no options left, we still push field to the stack to be able
	// to trace back our choices. For example if the next field to fill
	// was (2,3) with {3}, then we put (2,3) {} to the stack.
	// NextSolution ends with success if currentFieldsToFill if empty,
	// with no solution if choice stack is empty.
	choiceStack  fieldChoicesLeftStack
	chosenFields []fieldCoordinates
}

func NewSmartBarcktrack() Solver {
	return &smartBacktrack{}
}

func (s *smartBacktrack) Reset(board *board.Board) {
	s.solvable = true
	s.choiceStack = fieldChoicesLeftStack{}
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
			if s.backtrack() {
				continue
			} else {
				s.solvable = false
				return nil
			}
		}
		f.possibleValues.ForEach(func(n int) bool {
			s.choiceStack.Push(fieldChoiceLeft{f.x, f.y, uint16(n)})
			return false
		})
		choice := s.choiceStack.Pop()
		s.setNumber(choice.x, choice.y, uint16(choice.n))
	}
	solution := s.board.Copy()
	if !s.backtrack() {
		s.solvable = false // there will be no more solutions
	}
	return solution
}

func (s *smartBacktrack) backtrack() bool {
	if s.choiceStack.IsEmpty() {
		return false
	}
	restoredChoice := s.choiceStack.Pop()
	for i := len(s.chosenFields) - 1; i >= 0; i-- {
		f := s.chosenFields[i]
		if f.x == restoredChoice.x && f.y == restoredChoice.y {
			s.resetNumber(f.x, f.y, true) //TODO remove hack
			s.chosenFields = s.chosenFields[:i]
			s.setNumber(f.x, f.y, restoredChoice.n)
			return true
		}
		s.resetNumber(f.x, f.y, false)
	}
	panic("assertion failed in backtrack - restoredChoice coordinates were not in chosenFields")
}

func (s *smartBacktrack) setNumber(x, y int, n uint16) {
	s.board.Set(x, y, n)
	s.chosenFields = append(s.chosenFields, fieldCoordinates{x, y})
	sortNeeded := false
	for _, f := range s.currentFieldsToFill {
		if f.x == x && f.y == y {
			// this check will be removed, for now just a temporary brutal panic for tests
			panic("assertion failed - setNumber while number is still in currentFieldsToFill")
		}
		// if field is in the same row / column / subgrid as changed field,
		// set of possibleVelues must be updated
		if f.x == x || f.y == y || s.board.HaveCommonSubgrid(x, y, f.x, f.y) {
			sortNeeded = sortNeeded || f.possibleValues.Remove(int(n))
		}
	}
	// TODO we can use heap.Fix only for changed fields if each field "knows" its queue index
	// (but it is just an optimization that can be done later if needed)
	if sortNeeded {
		heap.Init(&s.currentFieldsToFill)
	}
}

func (s *smartBacktrack) resetNumber(x, y int, hack bool) {
	n := s.board.Get(x, y)
	s.board.Set(x, y, 0)
	sortNeeded := false
	for i, f := range s.currentFieldsToFill {
		if f.x == x && f.y == y {
			// this check will be removed, for now just a temporary brutal panic for tests
			panic("assertion failed - resetNumber while number is still in currentFieldsToFill")
		}
		if f.x == x || f.y == y || s.board.HaveCommonSubgrid(x, y, f.x, f.y) {
			if s.fieldCanHaveNumber(f.x, f.y, n) {
				s.currentFieldsToFill[i].possibleValues.Add(int(n))
				sortNeeded = true
			}
		}
	}
	if sortNeeded {
		if !hack {
			possibleNumbers, _ := s.findPossibleNumbers(x, y)
			newField := fieldToFill{
				x:              x,
				y:              y,
				possibleValues: possibleNumbers,
			}
			s.currentFieldsToFill.Push(newField)
		}
		heap.Init(&s.currentFieldsToFill)
	} else {
		if !hack {
			possibleNumbers, _ := s.findPossibleNumbers(x, y)
			newField := fieldToFill{
				x:              x,
				y:              y,
				possibleValues: possibleNumbers,
			}
			heap.Push(&s.currentFieldsToFill, newField)
		}
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

func (s *smartBacktrack) fieldCanHaveNumber(x, y int, n uint16) bool {
	ok := true
	s.board.ForEachInRow(y, func(x, y int) {
		if s.board.Get(x, y) == n {
			ok = false
		}
	})
	if !ok {
		return false
	}
	s.board.ForEachInColumn(x, func(x, y int) {
		if s.board.Get(x, y) == n {
			ok = false
		}
	})
	if !ok {
		return false
	}
	s.board.ForEachInSubgrid(x, y, func(x, y int) {
		if s.board.Get(x, y) == n {
			ok = false
		}
	})
	return ok
}
