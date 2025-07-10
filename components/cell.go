package components

import (
	"bytes"
	"fmt"
	"log"
	"sudoku/attributes"
	"sudoku/fonts"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Cell struct {
	rect  *attributes.Rect
	value uint16

	mutable   bool
	badChoice bool
	usedAll   bool

	clr      attributes.Color
	fontFace *text.GoTextFace
}

func (c *Cell) Construct(position attributes.Vector, cellSize int, initialValue uint16) {
	c.rect = &attributes.Rect{
		Position: position,
		Size:     attributes.Vector{X: float64(cellSize), Y: float64(cellSize)},
	}
	c.value = initialValue
	c.mutable = false
	if c.value == 0 {
		c.mutable = true
	}
	c.clr = attributes.CELL_COLOR
	if !c.mutable {
		c.clr.Current = c.clr.Gray
	}
	var err error
	fontFace, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.Inter))
	if err != nil {
		log.Fatal(err)
	}
	c.fontFace = &text.GoTextFace{
		Source: fontFace,
		Size:   c.rect.Size.Y * 0.9,
	}
	c.usedAll = false
}

func (c *Cell) SetBadChoice(isBad bool) {
	c.badChoice = isBad
}

func (c *Cell) IsBadChoice() bool {
	return c.badChoice
}

func (c *Cell) SetNormal() {
	if c.mutable {
		c.clr.Current = c.clr.White
		return
	}
	c.clr.Current = c.clr.Gray
}

func (c *Cell) SetHighlight() {
	if c.mutable {
		c.clr.Current = c.clr.PurpleLight
		return
	}
	c.clr.Current = c.clr.PurpleDark
}

func (c *Cell) PlayerUsedAll() bool {
	return c.usedAll
}

func (c *Cell) MarkAsUsed() {
	if c == nil {
		return
	}
	c.usedAll = true
	c.clr.Current = c.clr.Yellow
}

func (c *Cell) MarkAsUnUsed() {
	if c == nil {
		return
	}
	c.usedAll = false
}

func (c *Cell) Draw(surface *ebiten.Image) {
	c.rect.Draw(surface, c.clr.Current)

	var strVal string = fmt.Sprintf("%d", c.value)
	if c.value == 0 {
		strVal = ""
	}

	var textWidth, textHeight = text.Measure(strVal, c.fontFace, c.fontFace.Size+10)

	options := &text.DrawOptions{}
	options.GeoM.Translate(c.rect.Left()+c.rect.Size.X/2-0.5-textWidth/2, c.rect.Top()+c.rect.Size.Y/2-textHeight/2)
	options.ColorScale.Scale(0, 0, 0, 1)

	text.Draw(surface, strVal, c.fontFace, options)
}
