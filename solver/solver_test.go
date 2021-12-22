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
}

func TestBrutefoce(t *testing.T) {
	solver := solver.NewBruteforce()
	genericTestSolver(t, solver)
}
