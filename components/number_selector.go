package components

import (
	"sudoku/attributes"

	"github.com/hajimehoshi/ebiten/v2"
)

type NumberSelector struct {
	background *attributes.Rect

	numbers Cluster

	selectedCell    *Cell
	selectedCellRow int
	selectedCellCol int
}

func (ns *NumberSelector) Construct(clusterSize, cellSize, screenWidth, screenHeight int, boardOffsetY float64) {
	totalCellWidth := clusterSize * cellSize
	totalGapWidth := clusterSize*attributes.DISTANCE_BETWEEN_CELLS_PX + attributes.DISTANCE_BETWEEN_CLUSTERS_PX + attributes.DISTANCE_BETWEEN_CELLS_PX + attributes.DISTANCE_BETWEEN_CLUSTERS_PX

	ns.background = &attributes.Rect{
		Size: attributes.Vector{
			X: float64(totalCellWidth + totalGapWidth),
			Y: float64(totalCellWidth + totalGapWidth),
		},
	}
	ns.background.Position = attributes.Vector{
		X: float64(screenWidth)/2 - ns.background.Size.X/2,
		Y: float64(screenHeight) - ns.background.Size.Y - boardOffsetY,
	}

	var values [][]uint16 = make([][]uint16, clusterSize)
	var val uint16 = 1
	for row := range clusterSize {
		values[row] = make([]uint16, clusterSize)
		for col := range clusterSize {
			values[row][col] = val
			val++
		}
	}

	ns.numbers.Construct(0, 0, cellSize, ns.background.Position, values)
}

func (ns *NumberSelector) CurrentValue() uint16 {
	if ns.selectedCell == nil {
		return 0
	}
	return ns.selectedCell.value
}

func (ns *NumberSelector) Update() {
	var mx, my int = ebiten.CursorPosition()
	if !ns.background.CollidePoint(attributes.Vector{X: float64(mx), Y: float64(my)}) {
		return
	}

	var prevCell *Cell = ns.selectedCell
	cell, row, col := ns.numbers.TouchedCell()

	if cell == nil {
		return
	}

	if prevCell != nil {
		prevCell.SetNormal()
	}

	ns.selectedCell = cell
	ns.selectedCellRow = row
	ns.selectedCellCol = col
	ns.selectedCell.SetHighlight()
}

func (ns *NumberSelector) Draw(surface *ebiten.Image) {
	ns.background.Draw(surface, attributes.BACKGROUND_COLOR)
	ns.numbers.Draw(surface)
}
