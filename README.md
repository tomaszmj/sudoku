# Command line Sudoku solver

## Running

Example how to run it:
```
cd cmd
go run main.go boards/very_difficult_9x9.txt
```

You can also submit your own board in format similar to the example ones.


## Solver algorithm

See description, tests and benchmark in `solver/smart_backtrack.go`.
```
// smartBacktrack solver searches space of possible sudoku solutions in a "smart" way.
// When selecting where to fill in the next number, it uses heuristics -
// in each iteration we pick field, for which count of possible numbers
// (numbers which are not in the same row/column/subgrid) is smallest.
// Thanks to that, we limit number of solution space "subtrees" to be explored.
// Moreover, backtracking is implemented with stack without recursion. Thanks
// to that, in case of grids with multiple solutiions, we can generate them one-by-one on demand.
```