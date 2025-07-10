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

	touchTracker TouchTracker
}

func (ns *NumberSelector) Construct(clusterSize, cellSize, screenWidth, screenHeight int, mobileOffset, boardOffsetY float64) {
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
		Y: float64(screenHeight) - ns.background.Size.Y - boardOffsetY - mobileOffset,
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

	ns.selectedCell = nil

	ns.numbers.Construct(0, 0, cellSize, ns.background.Position, values)
	ns.touchTracker.Construct()
}

func (ns *NumberSelector) CurrentValue() uint16 {
	if ns.selectedCell == nil {
		return 0
	}
	return ns.selectedCell.value
}

func (ns *NumberSelector) UsedUp() {
	for i := 1; i < len(ns.numbers)*len(ns.numbers[0])+1; i++ {
		if ns.CurrentValue() == uint16(i) {
			ns.selectedCell.MarkAsUsed()
		}
	}
}

func (ns *NumberSelector) NotUsedUp() {
	for i := 1; i < len(ns.numbers)*len(ns.numbers[0])+1; i++ {
		if ns.CurrentValue() == uint16(i) {
			ns.selectedCell.MarkAsUnUsed()
		}
	}
}

func (ns *NumberSelector) Update() {
	if ns.selectedCell != nil {
		for row := range len(ns.numbers) {
			for col := range len(ns.numbers[row]) {
				if ns.numbers[row][col].PlayerUsedAll() {
					continue
				}
				if ns.CurrentValue() == ns.numbers[row][col].value {
					ns.numbers[row][col].SetHighlight()
				} else {
					ns.numbers[row][col].SetNormal()
				}
			}
		}
	}

	var prevCell *Cell = ns.selectedCell
	var touches []ebiten.TouchID = ns.touchTracker.JustPressedTouchIDs()

	cell, row, col := ns.numbers.TouchedCell(touches)

	if cell == nil {
		return
	}

	if prevCell != nil {
		if !prevCell.PlayerUsedAll() {
			prevCell.SetNormal()
		}
	}

	ns.selectedCell = cell
	ns.selectedCellRow = row
	ns.selectedCellCol = col
}

func (ns *NumberSelector) Draw(surface *ebiten.Image) {
	ns.background.Draw(surface, attributes.BACKGROUND_COLOR)
	ns.numbers.Draw(surface)
}
