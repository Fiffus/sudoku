package components

import (
	"sudoku/attributes"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Board struct {
	background *attributes.Rect

	clusters [][]Cluster

	mistakes       int
	mistakeCounted bool
	won            bool
}

func (b *Board) Construct(boardSize, clusterSize, screenWidth, screenHeight, cellSize int) {
	totalCellWidth := boardSize * clusterSize * cellSize
	totalGapWidth := boardSize*clusterSize*attributes.DISTANCE_BETWEEN_CELLS_PX + boardSize*attributes.DISTANCE_BETWEEN_CLUSTERS_PX + attributes.DISTANCE_BETWEEN_CELLS_PX + attributes.DISTANCE_BETWEEN_CLUSTERS_PX

	b.background = &attributes.Rect{
		Size: attributes.Vector{
			X: float64(totalCellWidth + totalGapWidth),
			Y: float64(totalCellWidth + totalGapWidth),
		},
	}
	var offsetY float64 = 10
	if screenWidth < screenHeight {
		if float32(screenWidth)/float32(screenHeight) < 0.65 {
			offsetY = float64(screenWidth%int(b.background.Size.X)) / 2
		} else {
			offsetY = 10
		}
	}
	b.background.Position = attributes.Vector{
		X: float64(screenWidth)/2 - b.background.Size.X/2,
		Y: offsetY,
	}

	b.mistakes = 0
	b.mistakeCounted = false
	b.won = false

	var mappedInitialValues [][][][]uint16 = attributes.GenerateBoard(boardSize, clusterSize)

	b.clusters = make([][]Cluster, boardSize)
	for clusterRow := range boardSize {
		b.clusters[clusterRow] = make([]Cluster, boardSize)

		for clusterCol := range boardSize {
			b.clusters[clusterRow][clusterCol].Construct(
				clusterRow,
				clusterCol,
				cellSize,
				b.background.Position,
				mappedInitialValues[clusterRow][clusterCol],
			)
		}
	}
}

func (b *Board) BoardOffsetY() float64 {
	return b.background.Position.Y
}

func (b *Board) badRow(clusterRow, cellRow int) {
	for clC := range len(b.clusters[clusterRow]) {
		for ceC := range len(b.clusters[clusterRow][clC][cellRow]) {
			b.clusters[clusterRow][clC][cellRow][ceC].SetBadChoice(true)
		}
	}
}

func (b *Board) badCol(clusterCol, cellCol int) {
	for clR := range len(b.clusters) {
		for ceR := range len(b.clusters[clR][clusterCol]) {
			b.clusters[clR][clusterCol][ceR][cellCol].SetBadChoice(true)
		}
	}
}

func (b *Board) badCluster(clusterRow, clusterCol int) {
	for cellRow := range len(b.clusters[clusterRow][clusterCol]) {
		for cellCol := range len(b.clusters[clusterRow][clusterCol][cellRow]) {
			b.clusters[clusterRow][clusterCol][cellRow][cellCol].SetBadChoice(true)
		}
	}
}

func (b *Board) badCells(input uint16) {
	if input == 0 {
		return
	}
	for clusterRow := range len(b.clusters) {
		for clusterCol := range len(b.clusters[clusterRow]) {
			for cellRow := range len(b.clusters[clusterRow][clusterCol]) {
				for cellCol := range len(b.clusters[clusterRow][clusterCol][cellRow]) {
					b.clusters[clusterRow][clusterCol][cellRow][cellCol].SetBadChoice(false)
					b.clusters[clusterRow][clusterCol][cellRow][cellCol].SetNormal()
				}
			}
		}
	}
	for clusterRow := range len(b.clusters) {
		for clusterCol := range len(b.clusters[clusterRow]) {
			for cellRow := range len(b.clusters[clusterRow][clusterCol]) {
				for cellCol := range len(b.clusters[clusterRow][clusterCol][cellRow]) {
					if b.clusters[clusterRow][clusterCol][cellRow][cellCol].value == input {
						b.badRow(clusterRow, cellRow)
						b.badCol(clusterCol, cellCol)
						b.badCluster(clusterRow, clusterCol)
					}
				}
			}
		}
	}
	for clusterRow := range len(b.clusters) {
		for clusterCol := range len(b.clusters[clusterRow]) {
			for cellRow := range len(b.clusters[clusterRow][clusterCol]) {
				for cellCol := range len(b.clusters[clusterRow][clusterCol][cellRow]) {
					if b.clusters[clusterRow][clusterCol][cellRow][cellCol].value == input {
						b.clusters[clusterRow][clusterCol][cellRow][cellCol].SetHighlight()
					}
				}
			}
		}
	}
}

func (b *Board) Mistakes() int {
	return b.mistakes
}

func (b *Board) checkRow(globalRow, clusterSize int) bool {
	found := make(map[uint16]bool)

	for clusterCol := 0; clusterCol < len(b.clusters); clusterCol++ {
		rowInCluster := globalRow % clusterSize
		clusterRow := globalRow / clusterSize

		for col := 0; col < clusterSize; col++ {
			val := b.clusters[clusterRow][clusterCol][rowInCluster][col].value
			if val == 0 || found[val] {
				return false
			}
			found[val] = true
		}
	}
	return true
}

func (b *Board) checkCol(globalCol, clusterSize int) bool {
	found := make(map[uint16]bool)

	for clusterRow := 0; clusterRow < len(b.clusters); clusterRow++ {
		colInCluster := globalCol % clusterSize
		clusterCol := globalCol / clusterSize

		for row := 0; row < clusterSize; row++ {
			val := b.clusters[clusterRow][clusterCol][row][colInCluster].value
			if val == 0 || found[val] {
				return false
			}
			found[val] = true
		}
	}
	return true
}

func (b *Board) checkCluster(clusterRow, clusterCol, clusterSize int) bool {
	found := make(map[uint16]bool)

	for row := 0; row < clusterSize; row++ {
		for col := 0; col < clusterSize; col++ {
			val := b.clusters[clusterRow][clusterCol][row][col].value
			if val == 0 || found[val] {
				return false
			}
			found[val] = true
		}
	}
	return true
}

func (b *Board) CheckWin() {
	boardSize := len(b.clusters)
	clusterSize := len(b.clusters[0][0])

	for i := 0; i < boardSize*clusterSize; i++ {
		if !b.checkRow(i, clusterSize) || !b.checkCol(i, clusterSize) {
			b.won = false
			return
		}
	}

	for clusterRow := 0; clusterRow < boardSize; clusterRow++ {
		for clusterCol := 0; clusterCol < boardSize; clusterCol++ {
			if !b.checkCluster(clusterRow, clusterCol, clusterSize) {
				b.won = false
				return
			}
		}
	}

	b.won = true
}

func (b *Board) Won() bool {
	return b.won
}

func (b *Board) Update(input uint16) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		b.mistakeCounted = false
	}

	b.badCells(input)

	var mx, my int = ebiten.CursorPosition()
	if !b.background.CollidePoint(attributes.Vector{X: float64(mx), Y: float64(my)}) {
		return
	}

	b.CheckWin()

	for clusterRow := range len(b.clusters) {
		for clusterCol := range len(b.clusters[clusterRow]) {
			cell, cellRow, cellCol := b.clusters[clusterRow][clusterCol].TouchedCell()

			if cell == nil {
				continue
			}

			if input == 0 {
				continue
			}

			if !cell.mutable {
				continue
			}

			if cell.value == input {
				cell.value = 0
				continue
			}

			if b.clusters[clusterRow][clusterCol][cellRow][cellCol].IsBadChoice() && !b.mistakeCounted {
				if cell.value != input {
					b.mistakes++
					b.mistakeCounted = true
				}
			}

			cell.value = input
		}
	}
}

func (b *Board) Draw(surface *ebiten.Image) {
	b.background.Draw(surface, attributes.BACKGROUND_COLOR)

	for row := range len(b.clusters) {
		for _, cluster := range b.clusters[row] {
			cluster.Draw(surface)
		}
	}
}
