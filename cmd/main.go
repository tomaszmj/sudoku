package main

import (
	"fmt"

	"github.com/tomaszmj/sudoku/board"
)

func main() {
	board, err := board.New(4, 3)
	if err != nil {
		fmt.Println(err)
		return
	}
	board.FillExampleData()
	fmt.Println(board.String())
}
