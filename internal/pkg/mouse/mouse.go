package mouse

import (
	"github.com/go-vgo/robotgo"
)

//Position is a struct for defining mouse position
type Position struct {
	MouseX int
	MouseY int
}

//GetPosition gets the mouse position at this instant
func GetPosition() *Position {
	x, y := robotgo.GetMousePos()
	return &Position{
		MouseX: x,
		MouseY: y,
	}
}
