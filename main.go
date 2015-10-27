package main

import (
	"flag"
	"fmt"
)

func main() {
	fmt.Printf("Hello, world.\n")

	var filePath string
	flag.StringVar(&filePath, "file", "せつめい", "APP_PARTS_FILE_PATH")
	flag.Parse()

	//pre;
	//- [] パーツ情報のロード
	gamePartsDict := CreateGamePartsDict(filePath)

	//loop;
	//- [] アルゴリズムで自動生成
	CreateGameMap(gamePartsDict)

	//post;

	/*
		//loop
		//- [] アルゴリズムで自動生成
		createGameMap()
		//- [] png生成->アップロード
		createPngFromMap()
		//- [] jsonの生成->アップロード
		createJsonFromMap()

		//- [] エディタでjsonロード
	*/
	fmt.Printf("Hello, world2.\n")
}
