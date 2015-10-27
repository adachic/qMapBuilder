package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/adachic/lottery"
)

//マップの大きさ
type GameMapSize struct {
	MaxX int
	MaxY int
	MaxZ int
}

//面積
func (s GameMapSize) area() int {
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

//マップメタ(ここから詳細なパーツを決定)
type MacroMapType int
const (
	MacroMapTypeLoad = iota
	MacroMapTypeRough
	MacroMapTypeWall
	MacroMapTypeCantEnter //進入不可地形
	MacroMapTypeAllyPoint
	MacroMapTypeEnemyPoint
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

func NewGameMap(condition GameMapCondition) *GameMap{
	game_map := &GameMap{}
	return game_map.init(condition)
}

func (game_map *GameMap) init(condition GameMapCondition) *GameMap{
	//難易度を初期化
	game_map.initMapDifficult()
	fmt.Printf("difficult: %+v\n", game_map.difficult)

	//マップのサイズを初期化
	game_map.initMapSize()
	fmt.Printf("mapSize: %+v\n", game_map.Size)

	//大まかな地形を初期化
	game_map.initMapGeographical()

	//味方開始ポイントを初期化
	game_map.initAllyStartPoint()
	fmt.Printf("allyStartPoint: %+v\n", game_map.AllyStartPoint)

	//敵開始ポイントを決定
	game_map.initEnemyStartPoints()
	for _, enemyStartPoint := range game_map.EnemyStartPoints { // キーは使われません
		fmt.Printf("enemyStartPoint: %+v\n", enemyStartPoint)
	}

	//2次元マップの生成
	{
		xymap := NewXYMap(game_map.Size)
		//広場生成
		xymap.putPlazas(game_map.difficult, game_map.AllyStartPoint, game_map.EnemyStartPoints)

		//味方、敵ポイント描画
		xymap.putPoint(game_map.AllyStartPoint, MacroMapTypeAllyPoint)
		for i := 0; i < len(game_map.EnemyStartPoints); i++ {
			xymap.putPoint(game_map.EnemyStartPoints[i], MacroMapTypeEnemyPoint)
		}

		//道生成

		//壁生成

		//ラフ生成
		xymap.printMapForDebug()
	}

	return game_map
}

//マップ難易度の抽選結果を返す
func (game_map *GameMap) initMapDifficult() {
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
	game_map.difficult = difficults[result].(Difficult)
}


//マップサイズの抽選結果を返す
func (game_map *GameMap) initMapSize() {
	//横長か縦長か
	rectForm := CreateRectForm()
	//x/yアスペクト比
	aspect := CreateAspectOfRectFrom(rectForm)
	//面積
	area := CreateArea(game_map.difficult)

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

	game_map.Size = GameMapSize{x, y, 30}
}

//地形の抽選結果を返す
func (game_map *GameMap) initMapGeographical() {
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
	game_map.Geographical = geographicals[result].(Geographical)
}

//味方の出撃座標を返す
func (game_map *GameMap) initAllyStartPoint() {
	var seed Range
	//難度が高いと中央寄りになる
	switch game_map.difficult{
	case easy:
		seed = Range{70, 100}
		break
	case normal:
		seed = Range{30, 80}
		break
	case hard:
		seed = Range{0, 50}
		break
	case exhard:
		seed = Range{0, 20}
		break
	}
	fmt.Print("seed",seed.Min," ",seed.Max)
	distanceFrom := lottery.GetRandomInt(seed.Min, seed.Max)
	game_map.AllyStartPoint = CreateRandomPositionInMap(
		game_map.Size,
		GameMapPosition{game_map.Size.MaxX / 2, game_map.Size.MaxY / 2, 0},
		distanceFrom)
}

//敵出現座標の一覧を返す
func (game_map *GameMap) initEnemyStartPoints(){
	type rangeFromAlly struct {
		Min int
		Max int
	}
	//味方と敵との大体の距離感
	var rangeFrom rangeFromAlly
	//敵出撃座標の数
	var sattyPointNum int

	switch game_map.difficult{
	case easy:
		sattyPointNum = lottery.GetRandomInt(1, 3)
		rangeFrom = rangeFromAlly{50, 100}
		break
	case normal:
		sattyPointNum = lottery.GetRandomInt(3, 6)
		rangeFrom = rangeFromAlly{30, 100}
		break
	case hard:
		sattyPointNum = lottery.GetRandomInt(5, 10)
		rangeFrom = rangeFromAlly{13, 100}
		break
	case exhard:
		sattyPointNum = lottery.GetRandomInt(10, 20)
		rangeFrom = rangeFromAlly{10, 100}
		break
	}
	var sattyPoints []GameMapPosition
	for sattyPointNum > 0 {
		sattyPointNum--
		//味方ポイントからの距離
		distance := lottery.GetRandomInt(rangeFrom.Min, rangeFrom.Max)
		sattyPoint := CreateRandomPositionInMap(game_map.Size, game_map.AllyStartPoint, distance)
		sattyPoints = append(sattyPoints, sattyPoint)
	}
	game_map.EnemyStartPoints = sattyPoints
}

//マップ
type GameMap struct {
	JungleGym [][][]GameParts
	Size      GameMapSize
	difficult Difficult
	Geographical Geographical
	AllyStartPoint GameMapPosition
	EnemyStartPoints []GameMapPosition
}
