package main

import "sudoku/game"

func main() {
	var sudoku game.Sudoku = game.Sudoku{}
	sudoku.Construct()
	sudoku.Run()
}
