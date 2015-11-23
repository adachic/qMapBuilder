package main

import (
//	"flag"
	"fmt"
)

var gamePartsDict map[string]GameParts

func main() {
	fmt.Printf("Hello, world.\n")
	/*
	var filePath string
	flag.StringVar(&filePath, "file", "せつめい", "APP_PARTS_FILE_PATH")
	flag.Parse()
	*/

	//pre;
	//- [] パーツ情報のロード
//	gamePartsDict = CreateGamePartsDict("./AppParts.json")
	gamePartsDict = CreateGamePartsDict("./assets/IntegratedPartsAll.json")

	//loop;
	//- [] アルゴリズムで自動生成
	condition := GameMapCondition{}
	bulc(condition)

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

//基本フロー
func flow(condition GameMapCondition) {
	//マップ生成
	fmt.Println("====createMap====")
	game_map := NewGameMap(condition)

	fmt.Println("====bind====")

	//実際のパーツとのひも付け
	game_map.bindToGameParts(gamePartsDict)

	//バリデーション

	fmt.Println("====drawMap====")
	//png生成
	game_map.createPng(gamePartsDict)

	//json_export
}

//雑に100回まわしてみる
func bulc(condition GameMapCondition) {
	//x := 100
	x := 1
	for x > 0 {
		x--
		flow(condition)
		fmt.Printf("\n")
	}
}



