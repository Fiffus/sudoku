package attributes

import "image/color"

const (
	DISTANCE_BETWEEN_CELLS_PX    int = 2
	DISTANCE_BETWEEN_CLUSTERS_PX int = 4
)

var (
	CELL_COLOR Color = Color{
		Current:     color.RGBA{255, 255, 255, 255},
		White:       color.RGBA{255, 255, 255, 255},
		Gray:        color.RGBA{220, 220, 220, 255},
		PurpleLight: color.RGBA{255, 40, 190, 255},
		PurpleDark:  color.RGBA{200, 22, 157, 255},
	}
	BACKGROUND_COLOR color.RGBA = color.RGBA{100, 30, 80, 255}
)
