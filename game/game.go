package game

import (
	"bytes"
	"fmt"

	"github.com/tshih/tripletriad/deck"
)

type BaseGame struct {
	playerOneId int
	playerTwoId int
	board       Board
	moveHistory []GameMove
}

type GameMove struct {
	card   deck.Card
	row    int
	column int
}

type Board struct {
	field [][]deck.Card
}

func NewBoardStandard() Board {
	return NewBoard(3, 3)
}

func NewBoard(rows int, cols int) Board {
	field := make([][]deck.Card, cols)

	for i := range field {
		field[i] = make([]deck.Card, rows)
	}

	board := Board{field: field}
	return board
}

func (b Board) String() string {
	var buffer bytes.Buffer

	for _, i := range b.field {
		for _, j := range i {
			buffer.WriteString(j.String())
			buffer.WriteString("\n")
		}
	}

	return fmt.Sprintf("%v", buffer.String())
}

func (b *Board) playCard(row int, col int, card *deck.Card) {
	if &b.field[row][col] == nil {
		b.field[row][col] = *card
	}

}
