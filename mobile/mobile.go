package mobile

import (
	"sudoku/game"

	"github.com/hajimehoshi/ebiten/v2/mobile"
)

// export PATH="$PATH:$HOME/go/bin"
func init() {
	var sudoku game.Sudoku = game.Sudoku{}
	sudoku.Construct()
	mobile.SetGame(&sudoku)
}

func ExportedFunctionForEbitenMobile() {}
