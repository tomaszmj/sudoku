package solver

import (
	"container/heap"
	"fmt"

	"github.com/tomaszmj/sudoku/board"
	"github.com/tomaszmj/sudoku/set"
)

type smartBacktrack struct {
	board        *board.Board
	fieldsToFill fieldsToFillHeap
	solvable     bool
	// Let's say we take from fieldsToFill
	// field (1,1) with possible numbers {1,2}.
	// We put 1 on board and push remaining options to the stack, i.e. (1,1) {2}.
	// If there are no options left, we still push field to the stack to be able
	// to trace back our choices. For example if the next field to fill
	// was (2,3) with {3}, then we put (2,3) {} to the stack.
	// NextSolution ends with success if fieldsToFill if empty,
	// with no solution if choice stack is empty.
	leftoverChoices []fieldChoice
	choicesMade     []fieldChoice
}

func NewSmartBarcktrack() Solver {
	return &smartBacktrack{}
}

func (s *smartBacktrack) Reset(board *board.Board) {
	s.solvable = true
	s.board = board.Copy()
	s.fieldsToFill = fieldsToFillHeap{}
	board.ForEachUntilError(func(x, y int) error {
		if board.Get(x, y) == 0 {
			availableNumbers, err := s.findPossibleNumbers(x, y)
			if err != nil {
				s.solvable = false
			}
			s.fieldsToFill = append(s.fieldsToFill, fieldToFill{x: x, y: y, possibleValues: availableNumbers})
		}
		return nil
	})
	heap.Init(&s.fieldsToFill)
	s.leftoverChoices = make([]fieldChoice, 0)
	s.choicesMade = make([]fieldChoice, 0, len(s.fieldsToFill))
}

func (s *smartBacktrack) NextSolution() *board.Board {
	defer func() {
		r := recover()
		if r != nil {
			fmt.Printf("panic encountered: %v\nboard:\n%s\n", r, s.board.String())
			s.solvable = false
		}
	}()
	if !s.solvable {
		return nil
	}
	for len(s.fieldsToFill) > 0 {
		f := s.fieldsToFill[0]
		if f.possibleValues.Len() == 0 {
			if s.backtrack() {
				continue
			} else {
				s.solvable = false
				return nil
			}
		}
		heap.Remove(&s.fieldsToFill, 0)
		numberToSet := s.pickFirstAvailableNumber(&f)
		s.setNumber(f.x, f.y, numberToSet)
	}
	solution := s.board.Copy()
	if !s.backtrack() {
		s.solvable = false // there will be no more solutions
	}
	return solution
}

func (s *smartBacktrack) pickFirstAvailableNumber(f *fieldToFill) uint16 {
	var numberToSet uint16
	f.possibleValues.ForEach(func(n int) bool {
		if numberToSet == 0 {
			numberToSet = uint16(n)
		} else {
			s.leftoverChoices = append(s.leftoverChoices, fieldChoice{f.x, f.y, uint16(n)})
		}
		return false
	})
	return numberToSet
}

func (s *smartBacktrack) setNumber(x, y int, n uint16) {
	s.board.Set(x, y, n)
	s.choicesMade = append(s.choicesMade, fieldChoice{x, y, n})
	sortNeeded := false
	for i := range s.fieldsToFill {
		f := &s.fieldsToFill[i]
		if f.x == x && f.y == y {
			// this check will be removed, for now just a temporary brutal panic for tests
			panic("assertion failed - setNumber while number is still in fieldsToFill")
		}
		// if field is in the same row / column / subgrid as changed field,
		// set of possibleVelues must be updated
		if f.x == x || f.y == y || s.board.HaveCommonSubgrid(x, y, f.x, f.y) {
			removed := f.possibleValues.Remove(int(n))
			sortNeeded = sortNeeded || removed
		}
	}
	// TODO we can use heap.Fix only for changed fields if each field "knows" its queue index
	// (but it is just an optimization that can be done later if needed)
	if sortNeeded {
		heap.Init(&s.fieldsToFill)
	}
}

func (s *smartBacktrack) backtrack() bool {
	if len(s.leftoverChoices) == 0 {
		return false
	}
	// pop the last leftover choice to backtrack to previous decision option
	leftoverChoice := s.leftoverChoices[len(s.leftoverChoices)-1]
	s.leftoverChoices = s.leftoverChoices[:len(s.leftoverChoices)-1]
	// revert all choices made after setting something on leftoverChoice
	for i := len(s.choicesMade) - 1; i >= 0; i-- {
		f := s.choicesMade[i]
		// check if current f is the field that we want to backtrack to (use leftoverChoice for that)
		if f.x == leftoverChoice.x && f.y == leftoverChoice.y {
			// sanity check
			possibleNumbers, err := s.findPossibleNumbers(f.x, f.y)
			if err != nil {
				panic(fmt.Sprintf("backtrack possible numbers assertion failed: %s", err))
			}
			if !possibleNumbers.Get(int(leftoverChoice.n)) {
				panic(fmt.Sprintf("backtrack possible numbers assertion failed: %d is not in possibleNumbers", leftoverChoice.n))
			}

			// just board.Set, not setNumber, because we are going to rebuild fieldsToFill from scratch anyway
			s.board.Set(f.x, f.y, leftoverChoice.n)
			// note that leftoverChoice must not be restired into fieldsToFill, because
			// restoring it might cause the algorithm infinitely process the same subtree of choices
			s.restoreFieldsToFill(s.choicesMade[(i + 1):])
			// change choicesMade - remove all that were after field to which we backtracked (f.x, f.y)
			// and change number in the choice to which we backtracked
			s.choicesMade = append(s.choicesMade[:i], fieldChoice{f.x, f.y, leftoverChoice.n})
			return true
		}
		// else - just set 0 on the board, fieldsToFill will be updated after reverting all choices
		s.board.Set(f.x, f.y, 0)
	}
	panic("assertion failed in backtrack - restoredChoice coordinates were not in choicesMade")
}

// restoreFieldsToFill is a helper function for backtrack
// It recreates fieldsToFill heap after choices from revertedChoices list
// have been removed from the board.
func (s *smartBacktrack) restoreFieldsToFill(revertedChoices []fieldChoice) {
	for i := range s.fieldsToFill {
		f := &s.fieldsToFill[i]
		possibleValues, err := s.findPossibleNumbers(f.x, f.y)
		if err != nil {
			panic(fmt.Sprintf("restoreFieldsToFill assertion failed: %s", err))
		}
		s.fieldsToFill[i].possibleValues = possibleValues
	}
	for _, f := range revertedChoices {
		possibleValues, _ := s.findPossibleNumbers(f.x, f.y)
		s.fieldsToFill = append(s.fieldsToFill, fieldToFill{f.x, f.y, possibleValues})
	}
	heap.Init(&s.fieldsToFill)
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
