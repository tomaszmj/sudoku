package solver_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tomaszmj/sudoku/solver"
)

func genericTestSolver(t *testing.T, solver solver.Solver) {
	t.Run("puzzle already solved", func(t *testing.T) {
		solver.Reset(solvedBoard)
		solution := solver.NextSolution()
		require.NotNil(t, solution)
		assert.Equal(t, solvedBoard.String(), solution.String())
	})

	t.Run("unsolveable puzzle", func(t *testing.T) {
		solver.Reset(unsolveableBoard)
		require.Nil(t, solver.NextSolution())
	})

	t.Run("solve puzzle", func(t *testing.T) {
		solver.Reset(boardToSolve)
		solution := solver.NextSolution()
		require.NotNil(t, solution)
		assert.Equal(t, solvedBoard.String(), solution.String())
		require.Nil(t, solver.NextSolution())
	})

	t.Run("puzzle with many solutions", func(t *testing.T) {
		solver.Reset(boardWithManySoltions)
		i := 0
		solution := solver.NextSolution()
		for solution != nil {
			i++
			solution = solver.NextSolution()
		}
		assert.Equal(t, 2, i)
	})
}

func TestBrutefoce(t *testing.T) {
	solver := solver.NewBruteforce()
	genericTestSolver(t, solver)
}

func TestSmartBacktrack(t *testing.T) {
	solver := solver.NewSmartBarcktrack()
	// TODO for now this fails for "puzzle_with_many_solutions" and
	// normal solution is also hacked (it would not work for more complex
	// puzzles with backtracking needed)
	genericTestSolver(t, solver)

	// this is not tested for Bruteforce because of computational complexity
	t.Run("difficult puzzle with 1 solution", func(t *testing.T) {
		solver.Reset(difficultBoard)
		solution := solver.NextSolution()
		require.NotNil(t, solution)
		assert.Equal(t, difficultBoardSolution.String(), solution.String())
		require.Nil(t, solver.NextSolution())
	})
}

func BenchmarkSmartBacktrack(b *testing.B) {
	solver := solver.NewSmartBarcktrack()

	b.Run("6x6", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			solver.Reset(board6x6)
			require.NotNil(b, solver.NextSolution())
			require.Nil(b, solver.NextSolution())
		}
	})

	b.Run("9x9 easy", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			solver.Reset(board9x9Easy)
			require.NotNil(b, solver.NextSolution())
			require.Nil(b, solver.NextSolution())
		}
	})

	b.Run("9x9 difficult", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			solver.Reset(board9x9Difficult)
			require.NotNil(b, solver.NextSolution())
			require.Nil(b, solver.NextSolution())
		}
	})

	b.Run("12x12", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			solver.Reset(board12x12)
			require.NotNil(b, solver.NextSolution())
			require.Nil(b, solver.NextSolution())
		}
	})

	b.Run("25x25", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			solver.Reset(board25x25)
			require.NotNil(b, solver.NextSolution())
			// there are more solutions - ignore that
		}
	})
}
