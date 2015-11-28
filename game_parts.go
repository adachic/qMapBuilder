package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"github.com/adachic/lottery"
	"strconv"
	"github.com/ojrac/opensimplex-go"
	"sort"
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

	IsEmpty       bool
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

	for _, value := range partsDict {
		value.IsEmpty = false;
	}

	return partsDict
}

//gamePartsがmacroMapTypeを含むならtrue
func IsIncludedMacroType(gameParts GameParts, tgt MacroMapType) bool {
	for _, value := range gameParts.MacroTypes {
		//		fmt.Printf("unko2 %+v,\n", value)
		if (value == tgt) {
			return true
		}
	}
	return false
}

//ソート用
type Roads []GameParts
func (p Roads) Len() int {
	return len(p)
}

func (p Roads) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p Roads) Less(i, j int) bool {
	return p[i].PavementType < p[j].PavementType
}

//主幹となる道の種類を決定し、id群を返す
func GetIdsRoad(game_map *GameMap, gamePartsDict map[string]GameParts) []int {
	var partsRoads Roads
	var idsRoad []int

	//カテゴリ集合体
	for _, parts := range gamePartsDict {
		if (parts.Category != game_map.Category) {
			//			fmt.Printf("unko%+v, %+v\n", parts.Category ,game_map.Category)
			continue
		}
		if (!IsIncludedMacroType(parts, MacroMapTypeRoad)) {
			//			fmt.Printf("unko2%+v, %+v\n", parts.,game_map.Category)
			continue
		}
		//		fmt.Printf("unko3\n")
		partsRoads = append(partsRoads, parts)
	}

	//舗装度でソートする
	sort.Sort(partsRoads)

	for _, parts := range partsRoads {
		i, _ := strconv.Atoi(parts.Id)
		idsRoad = append(idsRoad, i)
	}

	fmt.Printf("idsRoad:%+v\n", idsRoad)
	return idsRoad
}

//主幹となるラフの種類を決定し、id群を返す
func GetIdsRough(game_map *GameMap, gamePartsDict map[string]GameParts) []int {
	var idsRough []int

	//カテゴリ集合体
	for id, parts := range gamePartsDict {
		if (parts.Category != game_map.Category) {
			continue
		}
		if (!IsIncludedMacroType(parts, MacroMapTypeRough)) {
			continue
		}
		i, _ := strconv.Atoi(id)
		idsRough = append(idsRough, i)
	}
	fmt.Printf("idsRough:%+v\n", idsRough)
	return idsRough
}

//主幹となる壁の種類を決定し、id群を返す
func GetIdsWall(game_map *GameMap, gamePartsDict map[string]GameParts) []int {
	var idsWall []int

	//カテゴリ集合体
	for id, parts := range gamePartsDict {
		if (parts.Category != game_map.Category) {
			continue
		}
		if (!IsIncludedMacroType(parts, MacroMapTypeWall)) {
			continue
		}
		i, _ := strconv.Atoi(id)
		idsWall = append(idsWall, i)
	}
	fmt.Printf("idsWall:%+v\n", idsWall)
	return idsWall
}

//表層(道,ラフ,壁)
func GetGamePartsSurface(idsWall []int, idsRough []int, idsRoad []int,
gamePartsDict map[string]GameParts, macro MacroMapType, x int, y int, z int) GameParts {
	switch(macro){
	case MacroMapTypeRoad:
		return GetGamePartsRoad(idsWall, idsRough, idsRoad, gamePartsDict, macro, x, y, z);
	case MacroMapTypeRough:
		return GetGamePartsRough(idsWall, idsRough, idsRoad, gamePartsDict, macro, z);
	case MacroMapTypeWall:
		return GetGamePartsWall(idsWall, idsRough, idsRoad, gamePartsDict, macro, z);
	}
	return GetGamePartsRough(idsWall, idsRough, idsRoad, gamePartsDict, macro, z);
}

//道を返す
func GetGamePartsRoad(idsWall []int, idsRough []int, idsRoad []int,
gamePartsDict map[string]GameParts, macro MacroMapType, x int, y int, z int, ) GameParts {

	idx := GetIdxWithEval2(x,y,idsRoad)
	id := idsRoad[idx]
	return gamePartsDict[strconv.Itoa(id)]
}

//ラフを返す
func GetGamePartsRough(idsWall []int, idsRough []int, idsRoad []int,
gamePartsDict map[string]GameParts, macro MacroMapType, z int) GameParts {
	opensimplex.NewWithSeed(0);


	id := idsRough[0]
	return gamePartsDict[strconv.Itoa(id)]
}

//壁を返す
func GetGamePartsWall(idsWall []int, idsRough []int, idsRoad []int,
gamePartsDict map[string]GameParts, macro MacroMapType, z int) GameParts {
	id := idsWall[0]
	return gamePartsDict[strconv.Itoa(id)]
}



//土を返す
func GetGamePartsFoundation(idsWall []int, idsRough []int, idsRoad []int, gamePartsDict map[string]GameParts) GameParts {
	wallIdsCount := len(idsWall)
	if (wallIdsCount > 0) {
		//		id := lottery.GetRandomInt(0, wallIdsCount)
		fmt.Printf("wall id: %2d", idsWall[0])
		return gamePartsDict[strconv.Itoa(idsWall[0])]
	}

	roughIdsCount := len(idsRough)
	if (roughIdsCount > 0) {
		id := lottery.GetRandomInt(0, roughIdsCount)
		fmt.Printf("rough id: %2d", idsRough[0])
		return gamePartsDict[strconv.Itoa(idsRough[id])]
	}

	roadIdsCount := len(idsRoad)
	if (roadIdsCount > 0) {
		id := lottery.GetRandomInt(0, roadIdsCount)
		fmt.Printf("road id: %2d", idsRoad[0])
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

//パーリンノイズに従って[]からidxを得る
func GetIdxWithEval2(x int, y int, ids []int) int{
	coefficient := 0.1 //パーリンノイズのサンプリング粒度小さいほどなだらか

	val := opensimplex.NewWithSeed(0).Eval2(float64(x) * coefficient, float64(y) * coefficient)
	num := len(ids)
	floatId := float64(num) * (val + 1.0) / 2.0
	idx := int(floatId)
	fmt.Printf("x:%d,y:%d,val:%f,num:%f,idx:%d\n", x, y, val, num, idx);

	return idx
}

