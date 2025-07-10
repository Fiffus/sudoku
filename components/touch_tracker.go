package components

import "github.com/hajimehoshi/ebiten/v2"

type TouchTracker struct {
	prev map[ebiten.TouchID]struct{}
}

func (tt *TouchTracker) Construct() {
	tt.prev = make(map[ebiten.TouchID]struct{})
}

func (t *TouchTracker) JustPressedTouchIDs() []ebiten.TouchID {
	current := ebiten.AppendTouchIDs(nil)
	justPressed := []ebiten.TouchID{}

	currentSet := make(map[ebiten.TouchID]struct{})
	for _, id := range current {
		currentSet[id] = struct{}{}
		if _, alreadyTouched := t.prev[id]; !alreadyTouched {
			justPressed = append(justPressed, id)
		}
	}

	t.prev = currentSet
	return justPressed
}
