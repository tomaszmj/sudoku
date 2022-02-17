package solver_test

import (
	"strings"

	"github.com/tomaszmj/sudoku/board"
)

func mustCreateBoard(s string) *board.Board {
	board, err := board.NewFromSerializedFormat(strings.NewReader(s))
	if err != nil {
		panic(err)
	}
	return board
}

var (
	// simple boards for tests:

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

	// boards for benchmarks:

	board6x6 = mustCreateBoard(`3 2
+-------+-------+
| 0 5 6 | 3 2 0 |
| 3 0 0 | 6 4 5 |
+-------+-------+
| 6 1 5 | 0 0 4 |
| 2 0 3 | 0 0 6 |
+-------+-------+
| 1 0 0 | 4 0 0 |
| 0 0 4 | 0 6 0 |
+-------+-------+
`)

	board9x9Easy = mustCreateBoard(`3 3
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
+-------+-------+-------+
`)

	board9x9Difficult = mustCreateBoard(`3 3
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
+-------+-------+-------+
`)

	board12x12 = mustCreateBoard(`3 4
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
+----------+----------+----------+----------+
`)

	board25x25 = mustCreateBoard(`5 5
+----------------+----------------+----------------+----------------+----------------+
|  0  0  0  0  0 |  0  0  0  0  0 |  0  0  0  0  0 |  0  0  0  0  0 |  0  0  0  0  0 |
|  0  0  0  0  0 |  0  0  0  0  0 |  0  0  0  0  0 |  0  0  0  0  0 |  0  0  0  0  8 |
|  0  0  0  0  0 |  0  0  0  0  0 |  0  0  0  0  0 |  0  0  0  0  0 |  0  0  0  1  0 |
|  0  0  0  0  0 |  0  0  0  0  0 |  0  0  0  0  0 |  0  0  0  0  0 |  0  0  3  0  0 |
|  0  0  0  0  0 |  0  0  0  0  0 |  0  0  0  0  0 |  0  0  0  0  0 |  0 20  0  0  0 |
+----------------+----------------+----------------+----------------+----------------+
|  0  0  0  0  0 |  0  0  0  0  0 |  0  0  0  0  0 |  0  0  0  0  0 | 16  0  0  0  0 |
|  0  0  0  0  0 |  0  0  0  0  0 |  0  0  0  0  0 |  0  0  0  0 24 |  0  0  0  0  0 |
|  0  0  0  0  0 |  0  0  0  0  0 |  0  0  0  0  0 |  0  0  0 14  0 |  0  0  0  0  0 |
|  0  0  0  0  0 |  0  0  0  0  0 |  0  0  0  0  0 |  0  0 11  0  0 |  0  0  0  0  0 |
|  0  0  0  0  0 |  0  0  0  0  0 |  0  0  0  0  0 |  0  2  0  0  0 |  0  0  0  0  0 |
+----------------+----------------+----------------+----------------+----------------+
|  0  0  0  0  0 |  0  0  0  0  0 |  0  0  0  0  0 | 21  0  0  0  0 |  0  0  0  0  0 |
|  0  0  0  0  0 |  0  0  0  0  0 |  0  0  0  0 22 |  0  0  0  0  0 |  0  0  0  0  0 |
|  0  0  0  0  0 |  0  0  0  0  0 |  0  0  0 12  0 |  0  0  0  0  0 |  0  0  0  0  0 |
|  0  0  0  0  0 |  0  0  0  0  0 |  0  0 23  0  0 |  0  0  0  0  0 |  0  0  0  0  0 |
|  0  0  0  0  0 |  0  0  0  0  0 |  0  4  0  0  0 |  0  0  0  0  0 |  0  0  0  0  0 |
+----------------+----------------+----------------+----------------+----------------+
|  0  0  0  0  0 |  0  0  0  0  0 |  6  0  0  0  0 |  0  0  0  0  0 |  0  0  0  0  0 |
|  0  0  0  0  0 |  0  0  0  0 25 |  0  0  0  0  0 |  0  0  0  0  0 |  0  0  0  0  0 |
|  0  0  0  0  0 |  0  0  0 10  0 |  0  0  0  0  0 |  0  0  0  0  0 |  0  0  0  0  0 |
|  0  0  0  0  0 |  0  0 18  0  0 |  0  0  0  0  0 |  0  0  0  0  0 |  0  0  0  0  0 |
|  0  0  0  0  0 |  0 17  0  0  0 |  0  0  0  0  0 |  0  0  0  0  0 |  0  0  0  0  0 |
+----------------+----------------+----------------+----------------+----------------+
|  0  0  0  0  0 | 15  0  0  0  0 |  0  0  0  0  0 |  0  0  0  0  0 |  0  0  0  0  0 |
|  0  0  0  0 19 |  0  0  0  0  0 |  0  0  0  0  0 |  0  0  0  0  0 |  0  0  0  0  0 |
|  0  0  0 13  0 |  0  0  0  0  0 |  0  0  0  0  0 |  0  0  0  0  0 |  0  0  0  0  0 |
|  0  0  7  0  0 |  0  0  0  0  0 |  0  0  0  0  0 |  0  0  0  0  0 |  0  0  0  0  0 |
|  0  5  0  0  0 |  0  0  0  0  0 |  0  0  0  0  0 |  0  0  0  0  0 |  0  0  0  0  0 |
+----------------+----------------+----------------+----------------+----------------+
`)
)
