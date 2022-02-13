package solver

import (
	"fmt"

	"github.com/tomaszmj/sudoku/set"
)

// fieldToFill is a helper data structure used by solver to determine what fields should be filled in
type fieldToFill struct {
	x, y           int
	possibleValues *set.Set
}

// String is used only for test
func (f fieldToFill) String() string {
	return fmt.Sprintf("(%d, %d) %s", f.x, f.y, f.possibleValues.String())
}

type fieldsToFillStack []fieldToFill

func (s *fieldsToFillStack) Push(f fieldToFill) {
	*s = append(*s, f)
}

func (s *fieldsToFillStack) Pop() fieldToFill {
	f := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return f
}

func (s fieldsToFillStack) Peek() fieldToFill {
	return s[len(s)-1]
}

func (s fieldsToFillStack) IsEmpty() bool {
	return len(s) == 0
}

type fieldsToFillHeap []fieldToFill

func (h fieldsToFillHeap) Len() int {
	return len(h)
}

func (h fieldsToFillHeap) Less(i, j int) bool {
	return h[i].possibleValues.Len() < h[j].possibleValues.Len()
}

func (h fieldsToFillHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *fieldsToFillHeap) Push(x interface{}) {
	*h = append(*h, x.(fieldToFill))
}

func (h *fieldsToFillHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
