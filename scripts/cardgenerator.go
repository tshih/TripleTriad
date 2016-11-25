package scripts

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/tshih/tripletriad/deck"
)

var DataPath = "./../data/"
var imgPath = "./../img/"
var jsonName = "cardData.json"

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

	if val, _ := exists(DataPath + jsonName); val {
		os.RemoveAll(DataPath + jsonName)
	}
	//os.MkdirAll(DataPath, os.ModeDir)

	levels := GenerateLevels(levelDirs)

	cards := make([]deck.Card, 0, 110)

	statMap := ParseStats(DataPath + "stats.txt")

	for _, level := range levels {
		for _, cardName := range level.cards {
			card := deck.Card{ImgName: cardName, Level: level.level, Name: cardName[2 : len(cardName)-4]}
			AddStats(&card, statMap)
			cards = append(cards, card)
		}
	}

	bytes, err := json.Marshal(cards)
	err = ioutil.WriteFile(DataPath+jsonName, bytes, os.ModePerm)

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

func ParseStats(filePath string) map[string][]int {

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	statMap := make(map[string][]int)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Name") {
			vals := strings.Split(line, ":")
			stats := make([]int, 0, 4)

			for scanner.Scan() {
				// line 1
				line2 := scanner.Text()
				if strings.HasPrefix(line2, "Statistics") {
					line = strings.Split(line2, ":")[1]
					top, err1 := ConvertValue(line)
					if err1 != nil {
						log.Fatal(err1.Error())
					}
					stats = append(stats, top)

					//line 2
					scanner.Scan()
					line2 = scanner.Text()
					leftRight := strings.Split(strings.TrimSpace(line2), " ")
					for _, str := range leftRight {
						val, err2 := ConvertValue(str)
						if err2 != nil {
							continue
						}
						stats = append(stats, val)
					}

					scanner.Scan()
					line2 = scanner.Text()
					down, err3 := ConvertValue(line2)
					if err3 != nil {
						log.Fatal(err3)
					}
					stats = append(stats, down)
					break
				}

			}
			//fmt.Println(stats)
			//fmt.Println(StripSpace(vals[1]))
			statMap[StripSpace(vals[1])] = stats
		}
	}

	return statMap

}

func AddStats(card *deck.Card, stats map[string][]int) {
	vals, found := stats[card.Name]
	if !found {
		fmt.Printf("%v\n", card.Name)
	} else {
		fmt.Println(vals)
	}
}

func ConvertValue(str string) (int, error) {
	if len(str) == 0 {
		return -1, errors.New("str length is 0")
	}
	str = strings.TrimSpace(str)
	var val int
	var err error
	if str == "A" {
		val = 10
	} else {
		val, err = strconv.Atoi(str)
	}

	if err != nil {
		log.Fatal(err.Error())
	}
	return val, nil
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

func StripSpace(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) || unicode.IsSymbol(r) {
			return -1
		}
		return r
	}, str)
}
