package main

import (
	"fmt"
	_ "MapBuilder"
	"MapBuilder"
)

func main() {
	fmt.Printf("Hello, world.\n")

	//pre;
	//- [] パーツ情報のロード
	gamePartsDict := MapBuilder.CreateGamePartsDict()

	//loop;
	//- [] アルゴリズムで自動生成
	MapBuilder.CreateGameMap(gamePartsDict)

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
