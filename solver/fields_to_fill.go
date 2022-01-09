package solver

import (
	"container/heap"
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

func TestHeap() {
	possibleValues1 := set.New(4)
	possibleValues1.Add(3)
	f1 := fieldToFill{
		x:              1,
		y:              1,
		possibleValues: possibleValues1.Copy(),
	}
	possibleValues1.Add(2)
	f2 := fieldToFill{
		x:              2,
		y:              2,
		possibleValues: possibleValues1.Copy(),
	}
	possibleValues1.Add(1)
	f3 := fieldToFill{
		x:              3,
		y:              3,
		possibleValues: possibleValues1.Copy(),
	}
	h := &fieldsToFillHeap{f1, f2}
	heap.Init(h)
	heap.Push(h, f3)
	fmt.Printf("minimum: %v\n", (*h)[0].String())
	for h.Len() > 0 {
		f := heap.Pop(h).(fieldToFill)
		fmt.Printf("%s; ", f.String())
	}
}
