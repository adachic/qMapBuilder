package MapBuilder

import (
	"fmt"
	"github.com/adachic/lottery"
	"math"
	"math/rand"
	"time"
)

//マップ
type GameMap struct {
	Size      GameMapSize
	JungleGym [][][]GameParts
}

//マップの大きさ
type GameMapSize struct {
	MaxX int
	MaxY int
	MaxZ int
}

//面積
func (s GameMapSize) area() int{
	return s.MaxY * s.MaxX
}

//マップの難度
type Difficult int

const (
	easy   Difficult = 10
	normal Difficult = 3
	hard   Difficult = 2
	exhard Difficult = 1
)

//長方形の形式
type RectForm int

const (
	horizontalLong RectForm = 6 //横長
	verticalLong   RectForm = 5
	square         RectForm = 2
)

//地形
type Geographical int

const (
	GeographicalStep     Geographical = 14 + 10
	GeographicalMountain Geographical = 9 + 10
	GeographicalCave     Geographical = 8 + 10
	GeographicalFort     Geographical = 7 + 10
	GeographicalShrine   Geographical = 6 + 10
	GeographicalTown     Geographical = 5 + 10
	GeographicalCastle   Geographical = 4 + 10
)

type MacroMapType int

const (
	MacroMapTypeLoad = iota
	MacroMapTypeRough
	MacroMapTypeWall
	MacroMapTypeCantEnter //進入不可地形
)

//座標
type GameMapPosition struct {
	X int
	Y int
	Z int
}

//確率を返す
func (d Difficult) Prob() int {
	return int(d)
}

func (d RectForm) Prob() int {
	return int(d)
}

func (d Geographical) Prob() int {
	return int(d)
}

//マップ難易度の抽選結果を返す
func createMapDifficult() Difficult {
	lot := lottery.New(rand.New(rand.NewSource(time.Now().UnixNano())))
	difficults := []lottery.Interface{
		easy,
		normal,
		hard,
		exhard,
	}
	result := lot.Lots(
		difficults...,
	)
	return difficults[result].(Difficult)
}

//長方形の形式の抽選結果を返す
func createRectForm() RectForm {
	lot := lottery.New(rand.New(rand.NewSource(time.Now().UnixNano())))
	forms := []lottery.Interface{
		horizontalLong,
		verticalLong,
		square,
	}
	result := lot.Lots(
		forms...,
	)
	return forms[result].(RectForm)
}

// x/yのRatioを返す
func createAspectOfRectFrom(rectForm RectForm) float32 {
	var ret float32
	var longer int
	var shorter int

	const minLength = 3
	const maxLength = 10

	x := lottery.GetRandomNormInt(minLength, maxLength)
	y := (minLength + maxLength) - x
	if x > y {
		longer = x
		shorter = y
	} else {
		longer = y
		shorter = x
	}

	fmt.Printf("longer: %+v\n", longer)
	fmt.Printf("shorter: %+v\n", shorter)
	switch rectForm {
	case horizontalLong:
		ret = float32(longer) / float32(shorter)
		break
	case verticalLong:
		ret = float32(shorter) / float32(longer)
		break
	case square:
		ret = 1.0
		break
	}
	return ret
}

//マップ面積を返す
func createArea(difficult Difficult) int {
	var ret int
	//10x10を最小とし、30x30を最大とする
	switch difficult {
	case easy:
		ret = lottery.GetRandomInt(10*10, 15*15)
		break
	case normal:
		ret = lottery.GetRandomInt(10*10, 25*25)
		break
	case hard:
		ret = lottery.GetRandomInt(15*15, 30*30)
		break
	case exhard:
		ret = lottery.GetRandomInt(15*15, 30*30)
		break
	}
	return ret
}

//マップサイズの抽選結果を返す
func createMapSize(difficult Difficult) GameMapSize {
	//横長か縦長か
	rectForm := createRectForm()
	//x/yアスペクト比
	aspect := createAspectOfRectFrom(rectForm)
	//面積
	area := createArea(difficult)

	fmt.Printf("form: %+v\n", rectForm)
	fmt.Printf("area: %+v\n", area)
	fmt.Printf("aspect: %f\n", aspect)
	yy := float32(area) / float32(aspect)
	fmt.Printf("yy: %f\n", yy)
	y := int(math.Sqrt(float64(yy)))
	if y < 1 {
		y = 1
	}
	x := int(area / y)

	return GameMapSize{x, y, 30}
}

//地形の抽選結果を返す
func createMapGeographical() Geographical {
	lot := lottery.New(rand.New(rand.NewSource(time.Now().UnixNano())))
	geographicals := []lottery.Interface{
		GeographicalStep,
		GeographicalMountain,
		GeographicalCave,
		GeographicalFort,
		GeographicalShrine,
		GeographicalTown,
		GeographicalCastle,
	}
	result := lot.Lots(
		geographicals...,
	)
	return geographicals[result].(Geographical)
}

//味方の出撃座標を返す
func createAllyStartPoint(difficult Difficult, mapSize GameMapSize) GameMapPosition {
	type distanceFromCenter struct {
		distanceFromCenterMin int
		distanceFromCenterMax int //これを下げると難しいのが作られやすい
	}
	var seed distanceFromCenter

	//難度が高いと中央寄りになる
	switch difficult {
	case easy:
		seed = distanceFromCenter{70, 100}
		break
	case normal:
		seed = distanceFromCenter{30, 80}
		break
	case hard:
		seed = distanceFromCenter{0, 50}
		break
	case exhard:
		seed = distanceFromCenter{0, 20}
		break
	}
	distanceFrom := lottery.GetRandomInt(seed.distanceFromCenterMin, seed.distanceFromCenterMax)
	return createRandomPositionInMap(mapSize, GameMapPosition{mapSize.MaxX / 2, mapSize.MaxY / 2, 0}, distanceFrom)
}

//敵出現座標の一覧を返す
func createEnemyStartPoints(difficult Difficult,
	mapSize GameMapSize,
	allyStartPoint GameMapPosition) []GameMapPosition {
	type rangeFromAlly struct {
		Min int
		Max int
	}
	//味方と敵との大体の距離感
	var rangeFrom rangeFromAlly
	//敵出撃座標の数
	var sattyPointNum int

	switch difficult {
	case easy:
		sattyPointNum = lottery.GetRandomInt(1, 3)
		rangeFrom = rangeFromAlly{15, 30}
		break
	case normal:
		sattyPointNum = lottery.GetRandomInt(3, 6)
		rangeFrom = rangeFromAlly{8, 30}
		break
	case hard:
		sattyPointNum = lottery.GetRandomInt(5, 10)
		rangeFrom = rangeFromAlly{5, 30}
		break
	case exhard:
		sattyPointNum = lottery.GetRandomInt(10, 20)
		rangeFrom = rangeFromAlly{5, 30}
		break
	}
	var sattyPoints []GameMapPosition
	for sattyPointNum > 0 {
		sattyPointNum--
		//味方ポイントからの距離
		distance := lottery.GetRandomInt(rangeFrom.Min, rangeFrom.Max)
		sattyPoint := createRandomPositionInMap(mapSize, allyStartPoint, distance)
		sattyPoints = append(sattyPoints, sattyPoint)
	}
	return sattyPoints
}


//雑に100回まわしてみる
func bulc() {
	x := 100
	for x > 0 {
		x--
		flow()
		fmt.Printf("\n")
	}
}

//基本フロー
func flow() {
	//難易度を決定
	difficult := createMapDifficult()
	fmt.Printf("difficult: %+v\n", difficult)

	//マップのサイズを決定
	mapSize := createMapSize(difficult)
	fmt.Printf("mapSize: %+v\n", mapSize)

	//大まかな地形を決定
	geographical := createMapGeographical()

	//味方開始ポイントを決定
	allyStartPoint := createAllyStartPoint(difficult, mapSize)
	fmt.Printf("allyStartPoint: %+v\n", allyStartPoint)

	//敵開始ポイントを決定
	enemyStartPoints := createEnemyStartPoints(difficult, mapSize, allyStartPoint)
	for _, enemyStartPoint := range enemyStartPoints { // キーは使われません
		fmt.Printf("enemyStartPoint: %+v\n", enemyStartPoint)
	}

	//見下ろしのマップを作成
	xyMap := createXYMap(difficult, mapSize, geographical, allyStartPoint, enemyStartPoints)
	//dump
	printGameMap(xyMap, mapSize)

	//実際のパーツとのひも付け

	//勾配生成

	//バリデーション

}

func printGameMap(xyMap [][]MacroMapType, mapSize GameMapSize){
	for y:=0 ; y < mapSize.MaxY ; y++ {
		for x:=0 ; x < mapSize.MaxX ; x++ {
			switch(xyMap[y][x]){
			case MacroMapTypeCantEnter:
				fmt.Print("#")
			case MacroMapTypeLoad:
				fmt.Print(".")
			}
		}
		fmt.Print("\n")
	}
}

func CreateGameMap(gamePartsDict map[string]GameParts) GameMap {
	bulc()
	//flow()

	return GameMap{}
}

func createJsonFromMap() {

}

func createPngFromMap() {

}
