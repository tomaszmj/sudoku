package set

import (
	"fmt"
	"strings"
)

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

// MaxNumber returns maximum number that can be stored in the set.
// Set can store integers from 1 to MaxNumber. Calling Get / Set / Remove with other value may cause panic.
func (s *Set) MaxNumber() int {
	return len(s.slice)
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

// Complement returns set that contains all valid values (integers 1..maxNumber) that are NOT in the set.
func (s *Set) Complement() *Set {
	maxNumber := s.MaxNumber()
	result := New(maxNumber)
	for n := 1; n <= maxNumber; n++ {
		if !s.Get(n) {
			result.Add(n)
		}
	}
	return result
}

func (s *Set) Copy() *Set {
	newSlice := make([]uint8, len(s.slice))
	copy(newSlice, s.slice)
	return &Set{
		slice: newSlice,
		count: s.Len(),
	}
}

// String returns human-readable representation of Set. It should be used only for tests / debugging.
func (s *Set) String() string {
	var b strings.Builder
	b.WriteString("{")
	for n := 1; n <= s.MaxNumber(); n++ {
		if s.Get(n) {
			b.WriteString(fmt.Sprintf("%d,", n))
		}
	}
	b.WriteString("}")
	return b.String()
}

func Intersection(sets ...*Set) *Set {
	if len(sets) == 0 {
		return nil
	}
	maxNumber := sets[0].MaxNumber()
	for _, s := range sets {
		if s.MaxNumber() != maxNumber {
			panic(fmt.Sprintf("set Intersection failed - incompatible maxNumber %d, %d", maxNumber, s.MaxNumber()))
		}
	}
	result := New(maxNumber)
	for n := 1; n <= maxNumber; n++ {
		ok := true
		for _, set := range sets {
			ok = ok && set.Get(n)
			if !ok {
				break
			}
		}
		if ok {
			result.Add(n)
		}
	}
	return result
}

func Union(sets ...*Set) *Set {
	if len(sets) == 0 {
		return nil
	}
	maxNumber := sets[0].MaxNumber()
	for _, s := range sets {
		if s.MaxNumber() != maxNumber {
			panic(fmt.Sprintf("set Intersection failed - incompatible maxNumber %d, %d", maxNumber, s.MaxNumber()))
		}
	}
	result := New(maxNumber)
	for n := 1; n <= maxNumber; n++ {
		ok := false
		for _, set := range sets {
			ok = ok || set.Get(n)
		}
		if ok {
			result.Add(n)
		}
	}
	return result
}
