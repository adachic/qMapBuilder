package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"github.com/adachic/lottery"
	"strconv"
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
	PavementType  int

	WaterType     WaterType `json:"waterType"`
	Category      Category
	StructureType StructureType

	IsEmpty 	  bool
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

	for _, value := range partsDict{
		value.IsEmpty = false;
	}

	return partsDict
}

//gamePartsがmacroMapTypeを含むならtrue
func IsIncludedMacroType(gameParts GameParts,tgt MacroMapType) bool{
	for _, value := range gameParts.MacroTypes{
//		fmt.Printf("unko2 %+v,\n", value)
		if(value == tgt){
			return true
		}
	}
	return false
}

//主幹となる道の種類を決定し、id群を返す
func GetIdsRoad(game_map *GameMap, gamePartsDict map[string]GameParts) []int {
	var idsRoad []int

	//カテゴリ集合体
	for id, parts:= range gamePartsDict{
		if (parts.Category != game_map.Category){
//			fmt.Printf("unko%+v, %+v\n", parts.Category ,game_map.Category)
			continue
		}
		if (!IsIncludedMacroType(parts, MacroMapTypeRoad)){
//			fmt.Printf("unko2%+v, %+v\n", parts.,game_map.Category)
			continue
		}
//		fmt.Printf("unko3\n")
		i , _ := strconv.Atoi(id)
		idsRoad = append(idsRoad, i)
	}

	fmt.Printf("idsRoad:%+v\n", idsRoad)
	return idsRoad
}

//主幹となるラフの種類を決定し、id群を返す
func GetIdsRough(game_map *GameMap, gamePartsDict map[string]GameParts) []int {
	var idsRough []int

	//カテゴリ集合体
	for id, parts:= range gamePartsDict{
		if (parts.Category != game_map.Category){
			continue
		}
		if (!IsIncludedMacroType(parts, MacroMapTypeRough)){
			continue
		}
		i , _ := strconv.Atoi(id)
		idsRough = append(idsRough, i)
	}
	fmt.Printf("idsRough:%+v\n", idsRough)
	return idsRough
}

//主幹となる壁の種類を決定し、id群を返す
func GetIdsWall(game_map *GameMap, gamePartsDict map[string]GameParts) []int {
	var idsWall []int

	//カテゴリ集合体
	for id, parts:= range gamePartsDict{
		if (parts.Category != game_map.Category){
			continue
		}
		if (!IsIncludedMacroType(parts, MacroMapTypeWall)){
			continue
		}
		i , _ := strconv.Atoi(id)
		idsWall= append(idsWall, i)
	}
	fmt.Printf("idsWall:%+v\n", idsWall)
	return idsWall
}

//表層(道,ラフ,壁)
func GetGamePartsSurface(idsWall []int, idsRough []int, idsRoad []int,
gamePartsDict map[string]GameParts,macro MacroMapType, z int) GameParts{
	switch(macro){
	case MacroMapTypeRoad:
		return GetGamePartsRoad(idsWall, idsRough, idsRoad, gamePartsDict, macro, z);
	case MacroMapTypeRough:
		return GetGamePartsRough(idsWall, idsRough, idsRoad, gamePartsDict, macro, z);
	case MacroMapTypeWall:
		return GetGamePartsWall(idsWall, idsRough, idsRoad, gamePartsDict, macro, z);
	}
	return GetGamePartsRough(idsWall, idsRough, idsRoad, gamePartsDict, macro, z);
}

//道を返す
func GetGamePartsRoad(idsWall []int, idsRough []int, idsRoad []int,
gamePartsDict map[string]GameParts,macro MacroMapType, z int) GameParts{
	id := idsRoad[0]
	return gamePartsDict[strconv.Itoa(id)]
}

//ラフを返す
func GetGamePartsRough(idsWall []int, idsRough []int, idsRoad []int,
gamePartsDict map[string]GameParts,macro MacroMapType, z int) GameParts{
	id := idsRough[0]
	return gamePartsDict[strconv.Itoa(id)]
}

//壁を返す
func GetGamePartsWall(idsWall []int, idsRough []int, idsRoad []int,
gamePartsDict map[string]GameParts,macro MacroMapType, z int) GameParts{
	id := idsWall[0]
	return gamePartsDict[strconv.Itoa(id)]
}


//土を返す
func GetGamePartsFoundation(idsWall []int, idsRough []int, idsRoad []int, gamePartsDict map[string]GameParts) GameParts{
	wallIdsCount := len(idsWall)
	if(wallIdsCount > 0){
//		id := lottery.GetRandomInt(0, wallIdsCount)
		return gamePartsDict[strconv.Itoa(idsWall[0])]
	}

	roughIdsCount := len(idsRough)
	if(roughIdsCount > 0){
		id := lottery.GetRandomInt(0, roughIdsCount)
		return gamePartsDict[strconv.Itoa(idsRough[id])]
	}

	roadIdsCount := len(idsRoad)
	if(roadIdsCount > 0){
		id := lottery.GetRandomInt(0,roadIdsCount)
		return gamePartsDict[strconv.Itoa(idsRoad[id])]
	}
	return GameParts{}
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

