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
					b.CheckWin(clusterRow, clusterCol, cellRow, cellCol)
				}
			}
		}
	}
}

func (b *Board) Mistakes() int {
	return b.mistakes
}

func (b *Board) CheckWin(clusterRow, clusterCol, cellRow, cellCol int) {
	for i := range len(b.clusters[clusterRow]) {
		for j := range len(b.clusters[clusterRow][i][cellRow]) {
			if b.clusters[clusterRow][clusterCol][cellRow][cellCol].value == b.clusters[clusterRow][i][cellRow][j].value {
				return
			}
			if b.clusters[clusterRow][clusterCol][cellRow][cellCol].value == b.clusters[i][clusterCol][j][cellCol].value {
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
	b.badCells(input)

	var mx, my int = ebiten.CursorPosition()
	if !b.background.CollidePoint(attributes.Vector{X: float64(mx), Y: float64(my)}) {
		return
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		b.mistakeCounted = false
	}

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
