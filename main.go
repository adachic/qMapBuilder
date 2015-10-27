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
//	gamePartsDict := CreateGamePartsDict(filePath)

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
	NewGameMap(condition)

	//実際のパーツとのひも付け

	//勾配生成

	//バリデーション

}

//雑に100回まわしてみる
func bulc(condition GameMapCondition) {
	x := 100
	for x > 0 {
		x--
		flow(condition)
		fmt.Printf("\n")
	}
}



