package deck

import "fmt"

type Deck struct {
	cards []Card
}

func (d Deck) String() string {
	str := ""
	for _, card := range d.cards {
		str += fmt.Sprint(card) + "\n"
	}
	return str
}

func New() Deck {
	deck := Deck{[]Card{}}
	return deck
}

func (d *Deck) Len() int {
	return len(d.cards)
}

func (d *Deck) getCard(index int) *Card {
	return &d.cards[index]
}
