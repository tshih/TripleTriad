package scripts

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/tshih/tripletriad/deck"
)

var dataPath = "./../data/"
var imgPath = "./../img/"

type Level struct {
	level int
	cards []string
}

func (l *Level) AddCard(cardName string) {
	l.cards = append(l.cards, cardName)
}

func CreateCardJSON() int {
	levelDirs, err := ioutil.ReadDir(imgPath)

	if err != nil {
		log.Fatal(err)
	}

	if val, _ := exists(dataPath); val {
		os.RemoveAll(dataPath)
	}
	os.MkdirAll(dataPath, os.ModeDir)

	levels := GenerateLevels(levelDirs)

	cards := make([]deck.Card, 0, 110)

	for _, level := range levels {
		for _, cardName := range level.cards {
			card := deck.Card{ImgName: cardName, Level: level.level, Name: cardName[2 : len(cardName)-4]}
			cards = append(cards, card)
		}
	}

	bytes, err := json.Marshal(cards)
	err = ioutil.WriteFile("./../data/cardData.json", bytes, os.ModePerm)

	return len(cards)
	//names, err := os.Readdirnames(-1)
}

func GenerateLevels(dirs []os.FileInfo) []Level {

	levels := make([]Level, 0, 10)

	for _, level := range dirs {

		files, err := ioutil.ReadDir(imgPath + level.Name())
		levelInt, err := strconv.Atoi(level.Name()[5:])

		if err != nil {
			log.Fatal(err)
		}

		currLevel := Level{level: levelInt, cards: []string{}}

		for _, file := range files {
			currLevel.AddCard(file.Name())
		}

		levels = append(levels, currLevel)
	}
	return levels
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
