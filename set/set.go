package set

import "fmt"

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

// Add tries to add given number (n) to the set.
// If n was already in the set, Add returns false.
// If n was not in the set, Add adds it and returns true.
// Argument n must be integer from 1 to maxNumber, otherwise it will panic.
func (s *Set) Add(n int) bool {
	if s.slice[n-1] != 0 {
		return false
	}
	s.count += 1
	s.slice[n-1] = 1
	return true
}

// Get argument must be integer from 1 to maxNumber, otherwise it will panic.
func (s *Set) Get(n int) bool {
	return s.slice[n-1] != 0
}

// Remove tries to remove given number (n) from the set.
// If n was not in the set, Remove returns false.
// If n was in the set, Remove removes it and returns true.
// Argument n must be integer from 1 to maxNumber, otherwise it will panic.
func (s *Set) Remove(n int) bool {
	if s.slice[n-1] == 0 {
		return false
	}
	s.count -= 1
	s.slice[n-1] = 0
	return true
}

// Len returns number of elements actually stored in the set (not to be confused with maxNumber).
func (s *Set) Len() int {
	return s.count
}

// Clear removes all elements from set (underlying storage remains allocated).
func (s *Set) Clear() {
	for i := range s.slice {
		s.slice[i] = 0
	}
	s.count = 0
}

// RawData is provided as convenience if iterating for each set element is needed.
// Returned slice should be used read-only.
func (s *Set) RawData() []uint8 {
	return s.slice
}

func (s *Set) Union(s1 *Set) *Set {
	if len(s.slice) != len(s1.slice) {
		panic(fmt.Sprintf("Sets cannot be unioned to different maxNumber: %d vs %d", len(s.slice), len(s1.slice)))
	}
	newSet := New(len(s.slice))
	for i := range s.slice {
		if s.slice[i] != 0 || s1.slice[i] != 0 {
			newSet.Add(i + 1)
		}
	}
	return newSet
}

func (s *Set) Intersection(s1 *Set) *Set {
	if len(s.slice) != len(s1.slice) {
		panic(fmt.Sprintf("Sets cannot be intersectioned to different maxNumber: %d vs %d", len(s.slice), len(s1.slice)))
	}
	newSet := New(len(s.slice))
	for i := range s.slice {
		if s.slice[i] != 0 && s1.slice[i] != 0 {
			newSet.Add(i + 1)
		}
	}
	return newSet
}
