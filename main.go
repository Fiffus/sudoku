package main

import "sudoku/game"

// desktop
func main() {
	var sudoku game.Sudoku = game.Sudoku{}
	sudoku.Construct()
	sudoku.Run()
}
