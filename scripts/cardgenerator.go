package scripts

import (
	"bufio"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/tshih/tripletriad/deck"
)

//DataPath path to data folder
var DataPath = "./../data/"

//ImgPath path to img folder
var ImgPath = "./../img/"

var jsonName = "cardData.json"

//Level contains cards for a given level
type Level struct {
	level int
	cards []string
}

type statData struct {
	name  string
	stats []int
}

//Adds a card to the given Level
func (l *Level) addCard(cardName string) {
	l.cards = append(l.cards, cardName)
}

//CreateCardJSON creates Triple Triad card data using the images in img and
//a given text file with card stats
func CreateCardJSON() int {
	levelDirs, err := ioutil.ReadDir(ImgPath)

	if err != nil {
		log.Fatal(err)
	}

	if val, _ := exists(DataPath + jsonName); val {
		os.RemoveAll(DataPath + jsonName)
	}

	levels := GenerateLevels(levelDirs)
	statMap := ParseStats(DataPath + "stats.txt")

	cards := make([]deck.Card, 0, 110)
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
}

//GenerateLevels Generates the Level structs from a slice of img directories
func GenerateLevels(dirs []os.FileInfo) []Level {

	levels := make([]Level, 0, 10)

	for _, level := range dirs {

		files, err := ioutil.ReadDir(ImgPath + level.Name())
		levelInt, err := strconv.Atoi(level.Name()[5:])

		if err != nil {
			log.Fatal(err)
		}

		currLevel := Level{level: levelInt, cards: []string{}}

		for _, file := range files {
			currLevel.addCard(file.Name())
		}

		levels = append(levels, currLevel)
	}
	return levels
}

//ParseStats parses data file for card stats
func ParseStats(filePath string) map[string]statData {

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	statMap := make(map[string]statData)

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
					stats = append(stats, extractStat(line))

					//line 2
					scanner.Scan()
					line2 = scanner.Text()
					leftRight := strings.Split(strings.TrimSpace(line2), " ")
					for _, str := range leftRight {
						if len(str) > 0 {
							val := extractStat(str)
							if val <= 0 {
								log.Fatal("Attempted to convert bad value")
							}
							stats = append(stats, val)
						}
					}
					scanner.Scan()
					stats = append(stats, extractStat(scanner.Text()))
					break
				}
			}

			statMap[StripSpace(vals[1])] = statData{name: strings.TrimSpace(vals[1]), stats: stats}
		}
	}

	return statMap

}

//AddStats Add stats from a map to the card
func AddStats(card *deck.Card, stats map[string]statData) {
	vals, found := stats[card.Name]
	if !found {
		log.Fatalf("Not found by name: %v\n", card.Name)
	}
	card.Up = vals.stats[0]
	card.Left = vals.stats[1]
	card.Right = vals.stats[2]
	card.Down = vals.stats[3]
	card.Name = vals.name
}

//ConvertValue converts value
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

func extractStat(str string) int {
	stat, err := ConvertValue(str)
	if err != nil {
		log.Println("Error converting string to card value")
		return -1
	}
	return stat
}

//StripSpace strips spaces and -
func StripSpace(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) || r == '-' {
			return -1
		}
		return r
	}, str)
}
