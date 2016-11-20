package main

import (
	"fmt"

	"github.com/tshih/tripletriad/deck"
)

func main() {
	c := deck.Card{Name: "Rinoa", Level: 10, Up: 9, Right: 10, Down: 10, Left: 9}
	fmt.Println(c)
}
