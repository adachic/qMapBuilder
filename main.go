package main

import (
//	"flag"
	"fmt"
	"sync"
)

var gamePartsDict map[string]GameParts

func main() {
	fmt.Printf("Hello, world.\n")

	//pre;
	//- [x] パーツ情報のロード
	gamePartsDict = CreateGamePartsDict("./assets/IntegratedPartsAll3.json") //harfId対応済み

	//- [x] アルゴリズムで自動生成
	condition := GameMapCondition{}
	bulc(condition)

	fmt.Printf("Hello, world2.\n")
}

//基本フロー
func flow(condition GameMapCondition) {
	//マップ生成
	fmt.Println("====createMap====")
	game_map := NewGameMap(condition)

	fmt.Println("====bind====")

	//実際のパーツとのひも付け
	if(!game_map.bindToGameParts(gamePartsDict)){
		//紐付けるパーツがない
		fmt.Println("====dame===")
		return;
	}
	fmt.Println("====drawMap====:geographical:",game_map.Geographical)

	//png生成
	game_map.createPng(gamePartsDict)

	//json_export
	game_map.createJson(gamePartsDict)
}

//雑に100回まわしてみる
func bulc(condition GameMapCondition) {
	x := 100
	/*
	for x > 0 {
		x--
		flow(condition)
	}
	return;
	*/
	wt := sync.WaitGroup{}
	for x > 0 {
		x--
		wt.Add(1)
		go func (){
			flow(condition)
			fmt.Printf("\n")
			wt.Done()
		}()
	}
	wt.Wait()
}
