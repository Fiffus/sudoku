package attributes

import "image/color"

type Vector struct {
	X float64
	Y float64
}

type Color struct {
	Current     color.RGBA
	White       color.RGBA
	Gray        color.RGBA
	PurpleLight color.RGBA
	PurpleDark  color.RGBA
}
