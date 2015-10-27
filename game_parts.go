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
	Id    string //`json:"id"`
	Tiles []Tile //`json:"tiles"`

	Walkable bool         //`json:"walkable"`
	Harf     bool         //`json:"harf"`
	RezoType RezoTypeRect `json:"rezo"`

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
