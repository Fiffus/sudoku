package components

import (
	"sudoku/attributes"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Cluster [][]Cell

func (c *Cluster) Construct(clusterRow, clusterCol, cellSize int, boardPosition attributes.Vector, values [][]uint16) {
	*c = make(Cluster, len(values))
	for row := range len(*c) {
		(*c)[row] = make([]Cell, len(values[row]))
		for col := range len((*c)[row]) {
			(*c)[row][col].Construct(
				attributes.Vector{
					X: boardPosition.X + float64(attributes.DISTANCE_BETWEEN_CLUSTERS_PX+attributes.DISTANCE_BETWEEN_CELLS_PX) + float64(col*attributes.DISTANCE_BETWEEN_CELLS_PX+col*cellSize) + float64(clusterCol*(len(values[0])*(cellSize+attributes.DISTANCE_BETWEEN_CELLS_PX)+attributes.DISTANCE_BETWEEN_CLUSTERS_PX)),
					Y: boardPosition.Y + float64(attributes.DISTANCE_BETWEEN_CLUSTERS_PX+attributes.DISTANCE_BETWEEN_CELLS_PX) + float64(row*(attributes.DISTANCE_BETWEEN_CELLS_PX+cellSize)) + float64(clusterRow*(len(values)*(cellSize+attributes.DISTANCE_BETWEEN_CELLS_PX)+attributes.DISTANCE_BETWEEN_CLUSTERS_PX)),
				},
				cellSize,
				values[row][col],
			)
		}
	}
}

func (c *Cluster) TouchedCell(touches []ebiten.TouchID) (cell *Cell, row, col int) {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return nil, 0, 0
	}

	var mx, my int = ebiten.CursorPosition()

	for row := range len(*c) {
		for col := range (*c)[row] {
			if (*c)[row][col].rect.CollidePoint(attributes.Vector{X: float64(mx), Y: float64(my)}) {
				return &(*c)[row][col], row, col
			}
		}
	}
	/*

		if len(touches) < 1 {
			return
		}
		tx, ty := ebiten.TouchPosition(touches[0])

		for row := range len(*c) {
			for col := range (*c)[row] {
				if (*c)[row][col].rect.CollidePoint(attributes.Vector{X: float64(tx), Y: float64(ty)}) {
					return &(*c)[row][col], row, col
				}
			}
		}*/

	return nil, 0, 0
}

func (c *Cluster) Draw(surface *ebiten.Image) {
	for row := range len(*c) {
		for _, cell := range (*c)[row] {
			cell.Draw(surface)
		}
	}
}
