package deck

import "fmt"

type Card struct {
	cardID int
	Name   string
	Level  int
	Up     int
	Right  int
	Down   int
	Left   int
}

func (c Card) String() string {
	return fmt.Sprintf("Name: %v, Level: %v, Up: %v, Right: %v, Down: %v, Left: %v", c.Name, c.Level, c.Up, c.Right, c.Down, c.Left)
}

type Direction int

type DirectionError error

const (
	Up Direction = iota
	Right
	Left
	Down
	MaxDirection = iota - 1
)

func (c *Card) Compare(k Card, d Direction) int {
	return c.GetValue(d) - c.GetOppDirValue(d)
}

func (c *Card) GetValue(d Direction) int {
	if d == Up {
		return c.Up
	}
	if d == Right {
		return c.Right
	}
	if d == Down {
		return c.Down
	}
	if d == Left {
		return c.Left
	}
	return -1
}

func (c *Card) GetOppDirValue(d Direction) int {
	return c.GetValue(MaxDirection - d)
}
