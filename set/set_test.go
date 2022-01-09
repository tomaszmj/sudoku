package set_test

import (
	"testing"

	"github.com/tomaszmj/sudoku/set"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	s := set.New(4)
	assert.Equal(t, 0, s.Len())
	s.Add(1)
	s.Add(4)
	assert.True(t, s.Get(1))
	assert.False(t, s.Get(2))
	assert.False(t, s.Get(3))
	assert.True(t, s.Get(4))
	assert.Equal(t, 2, s.Len())
	assert.Panics(t, func() { s.Add(0) })
	assert.Panics(t, func() { s.Add(5) })
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
