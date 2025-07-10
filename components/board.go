package components

import (
	"sudoku/attributes"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Board struct {
	background *attributes.Rect

	clusters [][]Cluster

	correctBoard [][][][]uint16

	mistakes       int
	mistakeCounted bool
	won            bool

	touchTracker TouchTracker
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

	b.correctBoard = attributes.GenerateBoard(boardSize, clusterSize)
	var finalValues [][][][]uint16 = attributes.CreateEmptyCellsForPlayer(b.correctBoard, boardSize, clusterSize)

	b.clusters = make([][]Cluster, boardSize)
	for clusterRow := range boardSize {
		b.clusters[clusterRow] = make([]Cluster, boardSize)

		for clusterCol := range boardSize {
			b.clusters[clusterRow][clusterCol].Construct(
				clusterRow,
				clusterCol,
				cellSize,
				b.background.Position,
				finalValues[clusterRow][clusterCol],
			)
		}
	}

	b.touchTracker.Construct()
}

func (b *Board) BoardOffsetY() float64 {
	return b.background.Position.Y
}

func (b *Board) badCells(input uint16) {
	if input == 0 {
		return
	}
	for clusterRow := range len(b.clusters) {
		for clusterCol := range len(b.clusters[clusterRow]) {
			for cellRow := range len(b.clusters[clusterRow][clusterCol]) {
				for cellCol := range len(b.clusters[clusterRow][clusterCol][cellRow]) {
					b.clusters[clusterRow][clusterCol][cellRow][cellCol].SetBadChoice(
						input != b.correctBoard[clusterRow][clusterCol][cellRow][cellCol],
					)

					if b.clusters[clusterRow][clusterCol][cellRow][cellCol].value == input {
						b.clusters[clusterRow][clusterCol][cellRow][cellCol].SetHighlight()
						continue
					}
					b.clusters[clusterRow][clusterCol][cellRow][cellCol].SetNormal()
				}
			}
		}
	}
}

func (b *Board) Mistakes() int {
	return b.mistakes
}

func (b *Board) CheckWin() {
	for clusterRow := range len(b.clusters) {
		for clusterCol := range len(b.clusters[clusterRow]) {
			for cellRow := range len(b.clusters[clusterRow][clusterCol]) {
				for cellCol := range len(b.clusters[clusterRow][clusterCol][cellRow]) {
					if b.clusters[clusterRow][clusterCol][cellRow][cellCol].value != b.correctBoard[clusterRow][clusterCol][cellRow][cellCol] {
						return
					}
				}
			}
		}
	}
	b.won = true
}

func (b *Board) Won() bool {
	return b.won
}

func (b *Board) FinishedPlacing(number uint16) bool {
	var count int = 0
	for clusterRow := range len(b.clusters) {
		for clusterCol := range len(b.clusters[clusterRow]) {
			for cellRow := range len(b.clusters[clusterRow][clusterCol]) {
				for cellCol := range len(b.clusters[clusterRow][clusterCol][cellRow]) {
					if b.clusters[clusterRow][clusterCol][cellRow][cellCol].value == number {
						count++
						if count == len(b.clusters)*len(b.clusters[0]) {
							return true
						}
					}
				}
			}
		}
	}
	return false
}

func (b *Board) Update(input uint16) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		b.mistakeCounted = false
	}

	b.CheckWin()

	for clusterRow := range len(b.clusters) {
		for clusterCol := range len(b.clusters[clusterRow]) {
			var touches []ebiten.TouchID = b.touchTracker.JustPressedTouchIDs()

			cell, cellRow, cellCol := b.clusters[clusterRow][clusterCol].TouchedCell(touches)

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

	b.badCells(input)
}

func (b *Board) Draw(surface *ebiten.Image) {
	b.background.Draw(surface, attributes.BACKGROUND_COLOR)

	for row := range len(b.clusters) {
		for _, cluster := range b.clusters[row] {
			cluster.Draw(surface)
		}
	}
}
