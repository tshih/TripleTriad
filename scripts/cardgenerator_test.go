package scripts

import "testing"

func TestCardGenerator(t *testing.T) {
	CreateCardJSON()
	ParseStats(DataPath + "stats.txt")
	// test := gosseract.Must(gosseract.Params{
	// 	Src:       "./../img/Level1/TTBiteBug.png",
	// 	Languages: "eng",
	// })
	//
	// fmt.Println(test)
}
