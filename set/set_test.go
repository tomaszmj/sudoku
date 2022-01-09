package set_test

import (
	"testing"

	"github.com/tomaszmj/sudoku/set"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	s := set.New(4)
	assert.Equal(t, 0, s.Len()) // set is created empty

	// Add some elements
	assert.True(t, s.Add(1))  // 1 was added
	assert.False(t, s.Add(1)) // 1 was already in the set - Add does nothing
	assert.True(t, s.Add(4))  // 4 was added as well

	// Check what Get and Len return
	assert.True(t, s.Get(1))
	assert.False(t, s.Get(2))
	assert.False(t, s.Get(3))
	assert.True(t, s.Get(4))
	assert.Equal(t, 2, s.Len())

	// Invalid numbers may cause panic
	assert.Panics(t, func() { s.Add(0) })
	assert.Panics(t, func() { s.Add(5) })

	// Remove some element
	assert.False(t, s.Remove(2)) // 2 was not in the set - Remove does nothing
	assert.True(t, s.Remove(1))  // 1 was removed
	assert.False(t, s.Get(1))
	assert.Equal(t, 1, s.Len())
}

func TestSetUnion(t *testing.T) {
	s1 := set.New(4)
	s1.Add(1)
	s1.Add(4)
	s2 := set.New(4)
	s2.Add(1)
	s2.Add(2)
	s := s1.Union(s2)
	assert.True(t, s.Get(1))
	assert.True(t, s.Get(2))
	assert.False(t, s.Get(3))
	assert.True(t, s.Get(4))
}

func TestSetIntersection(t *testing.T) {
	s1 := set.New(4)
	s1.Add(1)
	s1.Add(4)
	s2 := set.New(4)
	s2.Add(1)
	s2.Add(2)
	s := s1.Intersection(s2)
	assert.True(t, s.Get(1))
	assert.False(t, s.Get(2))
	assert.False(t, s.Get(3))
	assert.False(t, s.Get(4))
}

func BenchmarkSet(b *testing.B) {
	var size int = 16
	s := set.New(size)
	for n := 0; n < b.N; n++ {
		s.Add(1)
		s.Add(4)
		s.Add(6)
		s.Get(4)
		s.Add(size - 1)
		s.Add(size / 2)
		s.Get(size)
		s.Remove(6)
		s.Remove(1)
		_ = s.Len()
	}
}
