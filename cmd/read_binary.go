package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type FileJson struct {
	Json struct {
		Image struct {
			Data struct {
				Type string `json:"type"`
				Data []byte `json:"data"`
			} `json:"data"`
		} `json:"image"`
	} `json:"json"`
}

type CodeFile = []FileJson

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	jsonBytes, err := os.ReadFile("Code.json")
	check(err)

	var fileJson CodeFile
	err = json.Unmarshal(jsonBytes, &fileJson)
	check(err)

	file, err := os.Create("image.png")
	check(err)

	n, err := file.Write(fileJson[0].Json.Image.Data.Data)
	check(err)
	fmt.Println(n)
}
