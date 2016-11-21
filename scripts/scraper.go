package scripts

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

func Scrape() {
	pageReader, _ := getPage("http://finalfantasy.wikia.com/wiki/List_of_Final_Fantasy_VIII_Triple_Triad_cards")
	wd, _ := os.Getwd()
	os.Mkdir(wd+"/img", os.ModePerm)
	ParsePage(pageReader)
}

func getPage(url string) (io.Reader, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func ParsePage(pageReader io.Reader) error {
	tk := html.NewTokenizer(pageReader)
	findCardTable(tk)
	return nil
}

func SearchPageForAnchor(tk *html.Tokenizer, anchFunc AnchorActor) error {
loop:
	for {
		tt := tk.Next()

		switch tt {
		case html.ErrorToken:
			return tk.Err()

		case html.StartTagToken:
			t := tk.Token()

			isAnchor := anchFunc.Evaluator(t)

			if isAnchor {
				anchFunc.Action(tk, t)
				if anchFunc.stopOnFound(t) {
					break loop
				}
			}
		}
	}
	return nil
}

func checkAttribute(attrs []html.Attribute, key string, attrVal string) (bool, string) {
	for _, val := range attrs {
		if val.Key == key && strings.Contains(val.Val, attrVal) {
			return true, val.Val
		}
	}
	return false, ""
}

type AnchorEval func(html.Token) bool
type AnchorAction func(*html.Tokenizer, html.Token)
type AnchorActor struct {
	Evaluator   AnchorEval
	Action      AnchorAction
	stopOnFound AnchorEval
}

func isTable() AnchorEval {
	return func(t html.Token) bool {
		foundTable, _ := checkAttribute(t.Attr, "class", "FVIII table")
		return t.Data == "table" && foundTable
	}
}

func isImage() AnchorEval {
	return func(t html.Token) bool {
		foundImg, _ := checkAttribute(t.Attr, "data-image-name", ".png")
		if !foundImg {
			foundImg, _ = checkAttribute(t.Attr, "src", ".png")
		}
		return t.Data == "img" && foundImg
	}
}

func stopOnImageFound() AnchorEval {
	i := 0
	return func(t html.Token) bool {
		i++
		if i < 11 {
			return false
		} else {
			i = 0
			return true
		}
	}
}

func falseFunc() AnchorEval {
	return func(t html.Token) bool {
		return false
	}
}

func findCardTable(tk *html.Tokenizer) {
	aa := AnchorActor{Evaluator: isTable(), Action: findCardImage(), stopOnFound: falseFunc()}
	SearchPageForAnchor(tk, aa)
}

func findCardImage() AnchorAction {
	return func(tk *html.Tokenizer, t html.Token) {
		fmt.Printf("Called findCardImage, token: %v\n", t)
		aa := AnchorActor{Evaluator: isImage(), Action: DownloadImage(), stopOnFound: stopOnImageFound()}
		SearchPageForAnchor(tk, aa)
	}
}

var i = 0

func DownloadImage() AnchorAction {
	return func(tk *html.Tokenizer, t html.Token) {
		found, url := checkAttribute(t.Attr, "data-src", "http")
		if !found {
			found, url = checkAttribute(t.Attr, "src", ".png")
		}

		if found {
			_, name := checkAttribute(t.Attr, "data-image-name", ".png")
			fmt.Println(strconv.Itoa((i/11)+1) + ":" + name)
			response, err := http.Get(url)
			if err != nil {
				log.Fatal(err)
			}
			defer response.Body.Close()

			//open a file for writing
			wd, _ := os.Getwd()
			if i%11 == 0 {
				os.Mkdir(wd+"/img/"+"Level"+strconv.Itoa((i/11)+1), os.ModePerm)
			}

			file, err := os.Create(wd + "/img/" + "Level" + strconv.Itoa((i/11)+1) + "/" + name)
			if err != nil {
				log.Fatal(err)
			}

			_, err = io.Copy(file, response.Body)
			if err != nil {
				log.Fatal(err)
			}
			file.Close()
			i++
		}
	}
}
