package solver

import (
	"container/heap"
	"fmt"

	"github.com/tomaszmj/sudoku/board"
	"github.com/tomaszmj/sudoku/set"
)

// smartBacktrack solver searches space of possible sudoku solutions in a "smart" way.
// When selecting where to fill in the next number, it uses heuristics -
// in each iteration we pick field, for which count of possible numbers
// (numbers which are not in the same row/column/subgrid) is smallest.
// Thanks to that, we limit number of solution space "subtrees" to be explored.
// Moreover, backtracking is implemented with stack without recursion. Thanks
// to that, in case of grids with multiple solutiions, we can generate them one-by-one on demand.
//
// The algothim works on the following data structures:
// fieldsToFill fieldsToFillHeap - priority queue to pick field to be filled in in each iteration,
// leftoverChoices []fieldChoice - stack of choices that can be made when bactracking,
// choicesMade []fieldChoice - stack of choices that were made, and will have to be reversed when backtracking.
//
// Each fieldToFill (element of fieldsToFillHeap) contains information about
// its possible values (i.e. numbers that are not in the same row / column / subgrid).
// Initially, fieldsToFill are filled in with all empty fileds from the board.
// In each iteration of filling in sudoku, we pick field from fieldsToFill,
// using heuristic - number that has smallest number of possible values.
// Among possible numbers we select one and put in on the board (other choices are pushed
// to leftoverChoices stack). When number on a field is set, we remove the field from fieldsToFill,
// push inforamtion about decision made to choicesMade stack, and update remaining fieldsToFill
// (after putting new number on the board there will be less possible choices for other fields to fill
// in the same row / column / subgrid).
//
// If there are no possible values to choose from in one of the fieldsToFill, it means that given
// solution space "subtree" cannot be solved - we try to do backtracking. In this case,
// we revert all choices made (from choicesMade stack) until fields choice that was left on
// leftoverFields stack (which opens other "subtree" that we did not explore yet). Reverting
// choices puts fields back in fieldsToFill queue and updates possible values of other fieldsToFill
// (except for the last reverted choice - the one for which we have some leftover choice to be selected -
// for this field only numbers that were not picked yet are still available).
//
// Solving sudoku (NextSolution) ends if we fill in whole board (fieldsToFill queue is empty) or if we
// encounter field with no choices and there are no backtracking options left (in this case there is no solution).
// After finding solution we copy board state and try to perform backtracking, so that NextSolution can be
// called again to find more solutions (if they exist). Thanks to that, we can force searching whole
// possible solutions "tree".

type smartBacktrack struct {
	board           *board.Board
	fieldsToFill    fieldsToFillHeap
	solvable        bool
	leftoverChoices []fieldChoice
	choicesMade     []fieldChoice
}

func NewSmartBarcktrack() Solver {
	return &smartBacktrack{}
}

func (s *smartBacktrack) Reset(board *board.Board) {
	if err := s.validateInitialBoard(board); err != nil {
		s.solvable = false
		return
	}
	s.solvable = true
	s.board = board.Copy()
	s.fieldsToFill = fieldsToFillHeap{}
	board.ForEach(func(x, y int, n uint16) {
		if n == 0 {
			availableNumbers := s.findPossibleNumbers(x, y)
			s.fieldsToFill = append(s.fieldsToFill, fieldToFill{x: x, y: y, possibleValues: availableNumbers})
		}
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
			possibleNumbers := s.findPossibleNumbers(f.x, f.y)
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
		possibleValues := s.findPossibleNumbers(f.x, f.y)
		s.fieldsToFill[i].possibleValues = possibleValues
	}
	for _, f := range revertedChoices {
		possibleValues := s.findPossibleNumbers(f.x, f.y)
		s.fieldsToFill = append(s.fieldsToFill, fieldToFill{f.x, f.y, possibleValues})
	}
	heap.Init(&s.fieldsToFill)
}

func (s *smartBacktrack) findPossibleNumbers(x, y int) *set.Set {
	allForbiddenNumbers := set.New(s.board.Size())
	s.board.ForEachNeighbour(x, y, func(x, y int) {
		if n := s.board.Get(x, y); n != 0 {
			allForbiddenNumbers.Add(int(n))
		}
	})
	return allForbiddenNumbers.Complement()
}

func (s *smartBacktrack) validateInitialBoard(b *board.Board) error {
	numbersFound := set.New(b.Size())
	validateFunc := func(x, y int, n uint16) error {
		if n == 0 {
			return nil // initial validation accepts unfilled fields
		}
		if !numbersFound.Add(int(n)) {
			return fmt.Errorf("number %d is repeated in row/column/subgrid at %d, %d", n, x, y)
		}
		return nil
	}
	return b.Validate(validateFunc, numbersFound.Clear)
}
