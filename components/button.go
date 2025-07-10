package components

import (
	"bytes"
	"log"
	"sudoku/attributes"
	"sudoku/fonts"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Button struct {
	rect         attributes.Rect
	clr          attributes.Color
	text         string
	fontFace     *text.GoTextFace
	touchTracker TouchTracker
}

func (b *Button) Construct(position, size attributes.Vector, buttonText string) {
	b.rect = attributes.Rect{
		Position: position,
		Size:     size,
	}
	b.text = buttonText
	var err error
	fontFace, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.Inter))
	if err != nil {
		log.Fatal(err)
	}
	b.fontFace = &text.GoTextFace{
		Source: fontFace,
		Size:   size.Y * 0.7,
	}
	b.clr = attributes.CELL_COLOR
	b.touchTracker = TouchTracker{}
	b.touchTracker.Construct()
}

func (b *Button) HighLight() {
	x, y := ebiten.CursorPosition()
	if b.rect.CollidePoint(attributes.Vector{X: float64(x), Y: float64(y)}) {
		b.clr.Current = b.clr.PurpleLight
		return
	}
	b.clr.Current = b.clr.PurpleDark
}

func (b *Button) Pressed() bool {
	touches := b.touchTracker.JustPressedTouchIDs()

	if len(touches) > 0 {
		tx, ty := ebiten.TouchPosition(touches[0])

		if b.rect.CollidePoint(attributes.Vector{X: float64(tx), Y: float64(ty)}) {
			return true
		}
	}

	var x, y int = ebiten.CursorPosition()
	if b.rect.CollidePoint(attributes.Vector{X: float64(x), Y: float64(y)}) {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			return true
		}
	}
	return false
}

func (b *Button) Draw(surface *ebiten.Image) {
	b.rect.Draw(surface, b.clr.Current)

	var textWidth, textHeight = text.Measure(b.text, b.fontFace, b.fontFace.Size+10)

	options := &text.DrawOptions{}
	options.GeoM.Translate(b.rect.Left()+b.rect.Size.X/2-float64(len(b.text))/2-textWidth/2, b.rect.Top()+b.rect.Size.Y/2-textHeight/2)
	options.ColorScale.Scale(0, 0, 0, 1)
	text.Draw(surface, b.text, b.fontFace, options)
}
