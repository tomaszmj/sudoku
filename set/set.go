package set

// Set is holds some integers, from 1 to provided maxNumber.
// We could just use map[int]bool, but this implementation provides much better performance
// for dense dataset with small number of possible values.
// Thanks to holding all data in consistent slice, it should use cache well.
// In Sudoku solver, we are going to check existence of some number
// from a very small set (usually 1-9) very often.
type Set struct {
	slice []uint8
	count int
}

// New creates new Set, which can hold integers from 1 to maxNumber.
func New(maxNumber int) *Set {
	return &Set{
		slice: make([]uint8, maxNumber),
		count: 0,
	}
}

// Add argument must be integer from 1 to maxNumber, otherwise it will panic.
func (s *Set) Add(n int) {
	if s.slice[n-1] == 0 {
		s.count += 1
	}
	s.slice[n-1] = 1
}

// Get argument must be integer from 1 to maxNumber, otherwise it will panic.
func (s *Set) Get(n int) bool {
	return s.slice[n-1] != 0
}

// Remove argument must be integer from 1 to maxNumber, otherwise it will panic.
func (s *Set) Remove(n int) {
	if s.slice[n-1] != 0 {
		s.count -= 1
	}
	s.slice[n-1] = 0
}

// Len returns number of elements actually stored in the set (not to be confused with maxNumber).
func (s *Set) Len() int {
	return s.count
}

// RawData is provided as convenience if iterating for each set element is needed.
// Returned slice should be used read-only.
func (s *Set) RawData() []uint8 {
	return s.slice
}
