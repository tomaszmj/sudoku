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
}
