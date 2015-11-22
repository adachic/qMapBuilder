package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Category int
const (
	CategoryStep     Category = 14
	CategoryMountain Category = 9
	CategoryCave     Category = 8
	CategoryShrine   Category = 7
	CategoryTown     Category = 6
	CategoryCastle   Category = 5
)

type StructureType int

const (
	StructureTypeRoad StructureType = iota
	StructureTypeWall
	StructureTypeStep
	StructureTypeWatar
	StructureTypeWatarDamage1
	StructureTypeWatarDamage2
	StructureTypeWatarHeal
)

type WaterType int

const (
	WaterTypeNone WaterType = iota //個体、ソリッド
	WaterTypeWater
	WaterTypePoison
	WaterTypeFlame
	WaterTypeHeal
)

type RezoTypeRect int

const (
	RezoTypeRect32 RezoTypeRect = iota
	RezoTypeRect64
)

type Tile struct {
	FilePath string `json:"tile"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
	Width    int    `json:"w"`
	Height   int    `json:"h"`
}

type GameParts struct {
	Id            string //`json:"id"`
	Tiles         []Tile //`json:"tiles"`

	Walkable      bool   //`json:"walkable"`
	Harf          bool   //`json:"harf"`
	RezoType      RezoTypeRect `json:"rezo"`

	Snow          int
	MacroTypes    []MacroMapType
	Pavement      int

	WaterType     WaterType `json:"waterType"`
	Category      Category
	StructureType StructureType
}

//jsonから辞書作成
func CreateGamePartsDict(filePath string) map[string]GameParts {
	// Loading jsonfile
	file, err := ioutil.ReadFile(filePath)
	// 指定したDataset構造体が中身になるSliceで宣言する
	var partsDict map[string]GameParts

	json_err := json.Unmarshal(file, &partsDict)
	if err != nil {
		fmt.Println("Format Error: ", json_err)
	}

	fmt.Printf("%+v\n", partsDict)
	fmt.Printf("%+v\n", len(partsDict))
	fmt.Printf("%+v\n", partsDict["15"])

	return partsDict
}

/*
//具体的なパーツを返す
func GetGameParts(macroType MacroMapType, geographical Geographical, z int) GameParts {
	shouldHarf := ((z % 2) == 1)
	var parts GameParts
	switch macroType {
	case MacroMapTypeLoad:
		parts = getRoad(geographical, shouldHarf);
	case MacroMapTypeRough:
		parts = getRough(geographical);
	case MacroMapTypeWall:
		parts = getWall(geographical);
	case MacroMapTypeCantEnter: //進入不可地形
		parts = getCantEnter(geographical);
	case MacroMapTypeAllyPoint:
	case MacroMapTypeEnemyPoint:
	}
	return parts
}

//道
func getRoad(geographical Geographical, shouldHarf bool) GameParts {
	var ids []int
	if (!shouldHarf) {
		switch {
		case GeographicalStep:
			idsHosou := []int{847, 462, 845, 31, 166, 465}
			idsMichi := []int{848}
			idsBoro := []int{462, 463, 464}
		//847舗装//848木
		case GeographicalMountain:
		case GeographicalCave:
		case GeographicalFort:
		case GeographicalShrine:
		case GeographicalTown:
		case GeographicalCastle:
		}
		return;
	}
	switch {
	case GeographicalStep:
		ids = []int{848, 847, 462, 845}
	//847舗装//848木
	case GeographicalMountain:
	case GeographicalCave:
	case GeographicalFort:
	case GeographicalShrine:
	case GeographicalTown:
	case GeographicalCastle:
	}
}

//ラフ
func getRough(geographical Geographical) GameParts {
	var ids []int
	switch {
	case GeographicalStep:
		ids = []int{995, 994}
	case GeographicalMountain:
	case GeographicalCave:
	case GeographicalFort:
	case GeographicalShrine:
	case GeographicalTown:
	case GeographicalCastle:
	}
}

//壁
func getWall(geographical Geographical) GameParts {
}

//進入不可
func getCantEnter(geographical Geographical) GameParts {
}

*/

