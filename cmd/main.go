package main

import (
	"fmt"
	"os"

	"github.com/tomaszmj/sudoku/board"
	"github.com/tomaszmj/sudoku/solver"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("exactly 1 argument required (path to board)")
		return
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("error opening file: %s\n", err)
		return
	}
	board1, err := board.NewFromSerializedFormat(file)
	if err2 := file.Close(); err2 != nil {
		fmt.Printf("error closing file: %s\n", err2)
	}
	if err != nil {
		fmt.Printf("error creating board from file %s: %s\n", os.Args[1], err)
		return
	}
	fmt.Printf("input:\n%s\n", board1)
	solver := solver.NewSmartBarcktrack()
	solver.Reset(board1)
	solution := solver.NextSolution()
	if solution == nil {
		fmt.Println("no solution")
	} else {
		fmt.Printf("solution:\n%s\n", solution)
		solution2 := solver.NextSolution()
		if solution2 != nil {
			fmt.Println("there are more solutions to this board (only 1 has been shown)")
		}
	}
}
