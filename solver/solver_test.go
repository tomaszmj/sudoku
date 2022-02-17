package solver_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tomaszmj/sudoku/board"
	"github.com/tomaszmj/sudoku/solver"
)

func mustCreateBoard(s string) *board.Board {
	board, err := board.NewFromSerializedFormat(strings.NewReader(s))
	if err != nil {
		panic(err)
	}
	return board
}

var (
	boardToSolve = mustCreateBoard(`2 2
+-----+-----+
| 0 0 | 0 3 |
| 0 1 | 0 4 |
+-----+-----+
| 4 2 | 3 1 |
| 1 3 | 4 2 |
+-----+-----+
`)

	boardWithManySoltions = mustCreateBoard(`2 1
+-----+
| 0 0 |
+-----+
| 0 0 |
+-----+
`)

	solvedBoard = mustCreateBoard(`2 2
+-----+-----+
| 2 4 | 1 3 |
| 3 1 | 2 4 |
+-----+-----+
| 4 2 | 3 1 |
| 1 3 | 4 2 |
+-----+-----+
`)

	unsolveableBoard = mustCreateBoard(`2 2
+-----+-----+
| 2 4 | 1 3 |
| 1 0 | 2 4 |
+-----+-----+
| 4 2 | 3 1 |
| 1 3 | 4 2 |
+-----+-----+
`)

	// difficultBoard - taken from https://sandiway.arizona.edu/sudoku/examples.html
	difficultBoard = mustCreateBoard(`3 3
+-------+-------+-------+
| 0 0 0 | 6 0 0 | 4 0 0 |
| 7 0 0 | 0 0 3 | 6 0 0 |
| 0 0 0 | 0 9 1 | 0 8 0 |
+-------+-------+-------+
| 0 0 0 | 0 0 0 | 0 0 0 |
| 0 5 0 | 1 8 0 | 0 0 3 |
| 0 0 0 | 3 0 6 | 0 4 5 |
+-------+-------+-------+
| 0 4 0 | 2 0 0 | 0 6 0 |
| 9 0 3 | 0 0 0 | 0 0 0 |
| 0 2 0 | 0 0 0 | 1 0 0 |
+-------+-------+-------+
`)

	difficultBoardSolution = mustCreateBoard(`3 3
+-------+-------+-------+
| 5 8 1 | 6 7 2 | 4 3 9 |
| 7 9 2 | 8 4 3 | 6 5 1 |
| 3 6 4 | 5 9 1 | 7 8 2 |
+-------+-------+-------+
| 4 3 8 | 9 5 7 | 2 1 6 |
| 2 5 6 | 1 8 4 | 9 7 3 |
| 1 7 9 | 3 2 6 | 8 4 5 |
+-------+-------+-------+
| 8 4 5 | 2 1 9 | 3 6 7 |
| 9 1 3 | 7 6 8 | 5 2 4 |
| 6 2 7 | 4 3 5 | 1 9 8 |
+-------+-------+-------+
`)
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

	board6x6 := mustCreateBoard(`3 2
+-------+-------+
| 0 5 6 | 3 2 0 |
| 3 0 0 | 6 4 5 |
+-------+-------+
| 6 1 5 | 0 0 4 |
| 2 0 3 | 0 0 6 |
+-------+-------+
| 1 0 0 | 4 0 0 |
| 0 0 4 | 0 6 0 |
+-------+-------+`)

	board9x9Easy := mustCreateBoard(`3 3
+-------+-------+-------+
| 0 0 0 | 2 6 0 | 7 0 1 |
| 6 8 0 | 0 7 0 | 0 9 0 |
| 1 9 0 | 0 0 4 | 5 0 0 |
+-------+-------+-------+
| 8 2 0 | 1 0 0 | 0 4 0 |
| 0 0 4 | 6 0 2 | 9 0 0 |
| 0 5 0 | 0 0 3 | 0 2 8 |
+-------+-------+-------+
| 0 0 9 | 3 0 0 | 0 7 4 |
| 0 4 0 | 0 5 0 | 0 3 6 |
| 7 0 3 | 0 1 8 | 0 0 0 |
+-------+-------+-------+`)

	board9x9Difficult := mustCreateBoard(`3 3
+-------+-------+-------+
| 0 2 0 | 0 0 0 | 0 0 0 |
| 0 0 0 | 6 0 0 | 0 0 3 |
| 0 7 4 | 0 8 0 | 0 0 0 |
+-------+-------+-------+
| 0 0 0 | 0 0 3 | 0 0 2 |
| 0 8 0 | 0 4 0 | 0 1 0 |
| 6 0 0 | 5 0 0 | 0 0 0 |
+-------+-------+-------+
| 0 0 0 | 0 1 0 | 7 8 0 |
| 5 0 0 | 0 0 9 | 0 0 0 |
| 0 0 0 | 0 0 0 | 0 4 0 |
+-------+-------+-------+`)

	board12x12 := mustCreateBoard(`3 4
+----------+----------+----------+----------+
|  0  7  0 | 10  0  0 |  9  0 11 | 12  5  0 |
|  0  0  0 | 11  2  0 |  0  0  0 |  0  4  0 |
|  6  0  0 |  0  4  0 | 10  1  2 |  0  0  9 |
|  5  0  0 |  0  6  0 |  8  0  0 |  1  0  0 |
+----------+----------+----------+----------+
| 11  0  4 |  0  0  9 |  0  0  0 |  3  8  0 |
|  0  1  3 |  6  0  0 |  5  0  0 |  0  0  0 |
|  0  0  0 |  0  0  0 |  0  0  0 |  0  0  0 |
|  7  9  5 |  0  0 11 |  3  0 12 |  2  0  0 |
+----------+----------+----------+----------+
| 12  4  0 |  3  0  0 |  7  0  0 |  0  0  0 |
|  0  0  0 |  0  0  0 |  0  0  0 |  0  0  2 |
|  0  8  0 |  0  0 10 |  0  5  4 |  0  7  0 |
|  0  0  0 |  5  0  0 |  0 10  0 |  4  0  6 |
+----------+----------+----------+----------+`)

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
}
