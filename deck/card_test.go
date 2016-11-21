package deck

import "testing"

var c1 = Card{Name: "Rinoa", Level: 10, Up: 9, Right: 10, Down: 10, Left: 10}
var c2 = Card{Name: "Squall", Level: 10, Up: 8, Right: 8, Down: 8, Left: 8}

var compareTest = []struct {
	card1    Card
	card2    Card
	dir      Direction
	expected int
}{
	{c1, c2, Up, c1.Up - c2.Down},
	{c1, c2, Left, c1.Left - c2.Right},
	{c1, c2, Right, c1.Right - c2.Left},
	{c1, c2, Down, c1.Down - c2.Up},
}

func TestCompare(t *testing.T) {
	for _, tt := range compareTest {
		actual := tt.card1.Compare(tt.card2, tt.dir)
		if actual <= 0 {
			t.Errorf("Test Compare(%v), expected %v, actual %v", tt.dir, tt.expected, actual)
		}
	}
}
