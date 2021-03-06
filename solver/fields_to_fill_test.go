package solver

import (
	"container/heap"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tomaszmj/sudoku/set"
)

func TestHeap(t *testing.T) {
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
	h := &fieldsToFillHeap{f1, f3}
	heap.Init(h)
	heap.Push(h, f2)
	gotHeapValues := make([]fieldToFill, 0, 3)
	expectedHeapValues := []fieldToFill{f1, f2, f3}
	for h.Len() > 0 {
		f := heap.Pop(h).(fieldToFill)
		gotHeapValues = append(gotHeapValues, f)
	}
	assert.Equal(t, expectedHeapValues, gotHeapValues)
}
