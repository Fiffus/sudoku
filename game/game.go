package game

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"os"
	"sudoku/attributes"
	"sudoku/components"
	"sudoku/fonts"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Sudoku struct {
	board          *components.Board
	numberSelector *components.NumberSelector

	pauseButton *components.Button
	resetButton *components.Button

	fontFace *text.GoTextFace

	secondsElapsed int
	delay          int64
	cellSize       int
	boardSize      int
	clusterSize    int

	screenWidth, screenHeight int

	mobileOffset float64

	gamePaused bool
}

func (s *Sudoku) Construct() {
	// desktop
	s.screenWidth, s.screenHeight = ebiten.Monitor().Size()

	// mobile
	if s.screenWidth <= 0 || s.screenHeight <= 0 {
		s.screenWidth = 1080
		s.screenHeight = 2400
		s.mobileOffset = 160
	} else {
		s.mobileOffset = 0
		ebiten.SetWindowTitle("Sudoku")
		ebiten.SetFullscreen(true)
	}

	ebiten.SetWindowSize(s.screenWidth, s.screenHeight)

	s.boardSize = 3
	s.clusterSize = 3
	if s.screenWidth < s.screenHeight {
		if float32(s.screenWidth)/float32(s.screenHeight) < 0.65 {
			s.cellSize = int(float32(s.screenWidth/(s.boardSize*s.clusterSize)) * 0.9)
		} else {
			s.cellSize = int(float32(s.screenWidth/(s.boardSize*s.clusterSize)) * 0.6)
		}
	}
	if s.screenWidth >= s.screenHeight {
		s.cellSize = int(float32(s.screenHeight/(s.boardSize*s.clusterSize)) * 0.64)
	}

	var err error
	fontFace, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.Inter))
	if err != nil {
		log.Fatal(err)
	}
	s.fontFace = &text.GoTextFace{
		Source: fontFace,
		Size:   float64(s.cellSize) * 0.6,
	}

	s.secondsElapsed = 0
	s.delay = time.Now().UnixMilli()

	s.board = &components.Board{}
	s.board.Construct(s.boardSize, s.clusterSize, s.screenWidth, s.screenHeight, s.cellSize)

	s.numberSelector = &components.NumberSelector{}
	s.numberSelector.Construct(s.clusterSize, s.cellSize, s.screenWidth, s.screenHeight, s.mobileOffset, s.board.BoardOffsetY())

	s.pauseButton = &components.Button{}
	s.pauseButton.Construct(
		attributes.Vector{
			X: float64(s.screenWidth) - float64(s.cellSize)*2.6 - s.board.BoardOffsetY(),
			Y: float64(s.screenHeight) - float64(s.cellSize)*0.85 - float64(s.cellSize)*2.3 - s.mobileOffset,
		},
		attributes.Vector{
			X: float64(s.cellSize) * 2.6,
			Y: float64(s.cellSize) * 0.85,
		},
		"Pauza",
	)
	s.resetButton = &components.Button{}
	s.resetButton.Construct(
		attributes.Vector{
			X: float64(s.screenWidth) - float64(s.cellSize)*2.6 - s.board.BoardOffsetY(),
			Y: float64(s.screenHeight) - float64(s.cellSize)*0.85 - float64(s.cellSize)*1.4 - s.mobileOffset,
		},
		attributes.Vector{
			X: float64(s.cellSize) * 2.6,
			Y: float64(s.cellSize) * 0.85,
		},
		"Reset",
	)

	s.gamePaused = false
}

func (s *Sudoku) DrawMistakes(surface *ebiten.Image) {
	var strVal string = fmt.Sprintf("Chyby: %d", s.board.Mistakes())

	var _, textHeight float64 = text.Measure(strVal, s.fontFace, s.fontFace.Size+10)

	options := &text.DrawOptions{}
	options.GeoM.Translate(s.board.BoardOffsetY(), float64(s.screenHeight)-textHeight-float64(s.cellSize)*2.3-s.mobileOffset)
	options.ColorScale.Scale(0, 0, 0, 1)

	text.Draw(surface, strVal, s.fontFace, options)
}

func (s *Sudoku) DrawTime(surface *ebiten.Image) {
	var strVal string = fmt.Sprintf("Čas: %d:%d", s.secondsElapsed/60, s.secondsElapsed%60)

	var _, textHeight float64 = text.Measure(strVal, s.fontFace, s.fontFace.Size+10)

	options := &text.DrawOptions{}
	options.GeoM.Translate(s.board.BoardOffsetY(), float64(s.screenHeight)-textHeight-float64(s.cellSize)*1.2-s.mobileOffset)
	options.ColorScale.Scale(0, 0, 0, 1)

	text.Draw(surface, strVal, s.fontFace, options)
}

func (s *Sudoku) DrawWin(surface *ebiten.Image) {
	if !s.board.Won() {
		return
	}

	var strVal string = "Skvělá práce!"

	var textWidth, textHeight float64 = text.Measure(strVal, s.fontFace, s.fontFace.Size+10)

	options := &text.DrawOptions{}
	options.GeoM.Translate(float64(s.screenWidth/2)-textWidth/2, float64(s.cellSize*s.boardSize*s.boardSize)*1.3-s.board.BoardOffsetY()-textHeight)
	options.ColorScale.Scale(0.2, 0.6, 0.45, 1)

	text.Draw(surface, strVal, s.fontFace, options)
}

// BUG: board kontroluje počet čísel ale nechechkne jestli poslední číslo je valid a uzamkne ho i když se jedná o chybu
func (s *Sudoku) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}

	s.pauseButton.HighLight()
	s.resetButton.HighLight()

	if s.pauseButton.Pressed() {
		s.gamePaused = !s.gamePaused
	}

	if s.resetButton.Pressed() {
		s.board.Construct(s.boardSize, s.clusterSize, s.screenWidth, s.screenHeight, s.cellSize)
		s.numberSelector.Construct(s.clusterSize, s.cellSize, s.screenWidth, s.screenHeight, s.mobileOffset, s.board.BoardOffsetY())
		s.secondsElapsed = 0
	}

	if s.gamePaused {
		return nil
	}

	if s.board.Won() {
		return nil
	}

	if s.delay < time.Now().UnixMilli() {
		s.delay = time.Now().UnixMilli() + time.Second.Milliseconds()
		s.secondsElapsed++
	}

	if s.board.FinishedPlacing(s.numberSelector.CurrentValue()) {
		s.numberSelector.UsedUp()
	} else {
		s.numberSelector.NotUsedUp()
	}

	s.numberSelector.Update()
	s.board.Update(s.numberSelector.CurrentValue())

	return nil
}

func (s *Sudoku) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)

	s.board.Draw(screen)
	s.numberSelector.Draw(screen)

	s.pauseButton.Draw(screen)
	s.resetButton.Draw(screen)

	s.DrawMistakes(screen)
	s.DrawTime(screen)
	s.DrawWin(screen)
}

func (s *Sudoku) Run() {
	if err := ebiten.RunGame(s); err != nil {
		log.Fatal(err)
	}
}

func (s *Sudoku) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return s.screenWidth, s.screenHeight
}
