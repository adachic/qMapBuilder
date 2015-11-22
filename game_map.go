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
	MacroMapTypeOther //他
	MacroMapTypeAllyPoint
	MacroMapTypeEnemyPoint
)


//alreadyに登録されてなくて自分を除く最も近いポイントを返す
//errがtrueなら見つからなかった
func (pos GameMapPosition) searchNearPositionWithOutMe(positions []GameMapPosition, alreadys []PathPosition) (nearPos GameMapPosition, err bool) {
	err = true;
	minDistance := 10000
	for i := 0; i < len(positions); i++ {
		tgtPos := positions[i]
		if (pos.equalXYTo(tgtPos)) {
			continue
		}
		if (containsPath(alreadys, PathPosition{pos, tgtPos})) {
			continue
		}
		distance := pos.distance(tgtPos)
		if (distance < minDistance) {
			minDistance = distance
			nearPos = tgtPos
			err = false
		}
	}
	return nearPos, err
}

func (pos GameMapPosition) equalXYTo(another GameMapPosition) bool {
	return (pos.X == another.X) && (pos.Y == another.Y)
}

func (pos GameMapPosition) distance(another GameMapPosition) int {
	absx := pos.X - another.X
	absy := pos.Y - another.Y
	if (absx < 0) {
		absx *= -1
	}
	if (absy < 0) {
		absy *= -1
	}
	return absx + absy
}

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

func NewGameMap(condition GameMapCondition) *GameMap {
	game_map := &GameMap{}
	return game_map.init(condition)
}

func (game_map *GameMap) init(condition GameMapCondition) *GameMap {
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

	//2次元マップ生成
	{
		xymap := NewXYMap(game_map.Size)

		//広場配置
		xymap.putPlazas(game_map.difficult, game_map.AllyStartPoint, game_map.EnemyStartPoints)

		//道配置
		xymap.putRoads(game_map.difficult, game_map.AllyStartPoint, game_map.EnemyStartPoints)

		//壁配置

		//ラフ配置

		//味方、敵ポイント
		xymap.putPoint(game_map.AllyStartPoint, MacroMapTypeAllyPoint)
		for i := 0; i < len(game_map.EnemyStartPoints); i++ {
			xymap.putPoint(game_map.EnemyStartPoints[i], MacroMapTypeEnemyPoint)
		}

		//勾配を生成
		xymap.makeGradient(game_map.Geographical)

		//水、毒沼配置

		//alloc/init
		game_map.allocToJungleGym()

		//xymap情報をコピー
		game_map.copyFromXY(xymap)

		//dump
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
	fmt.Print("seed", seed.Min, " ", seed.Max)
	distanceFrom := lottery.GetRandomInt(seed.Min, seed.Max)
	game_map.AllyStartPoint = CreateRandomPositionInMap(
		game_map.Size,
		GameMapPosition{game_map.Size.MaxX / 2, game_map.Size.MaxY / 2, 0},
		distanceFrom)
}

//敵出現座標の一覧を返す
func (game_map *GameMap) initEnemyStartPoints() {
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
	JungleGym        [][][]GameParts
	MacroMapTypes    [][][]MacroMapType
	High 			 [][]int
	Size             GameMapSize
	difficult        Difficult
	Geographical     Geographical
	AllyStartPoint   GameMapPosition
	EnemyStartPoints []GameMapPosition
}

func (game_map *GameMap) allocToJungleGym() {
	game_map.JungleGym = make([][][]GameParts, game_map.Size.MaxZ)
	for z := 0; z < game_map.Size.MaxZ; z++ {
		game_map.JungleGym[z] = make([][]GameParts, game_map.Size.MaxY)
		for y := 0; y < game_map.Size.MaxY; y++ {
			game_map.JungleGym[z][y] = make([]GameParts, game_map.Size.MaxX)
		}
	}

	game_map.MacroMapTypes = make([][][]MacroMapType, game_map.Size.MaxZ)
	for z := 0; z < game_map.Size.MaxZ; z++ {
		game_map.MacroMapTypes[z] = make([][]MacroMapType, game_map.Size.MaxY)
		for y := 0; y < game_map.Size.MaxY; y++ {
			game_map.MacroMapTypes[z][y] = make([]MacroMapType, game_map.Size.MaxX)
		}
	}

	game_map.High = make([][]int,game_map.Size.MaxY)
	for y := 0; y < game_map.Size.MaxY; y++ {
		game_map.High[y] = make([]int, game_map.Size.MaxX)
	}
}

//ここでxyをjungleGymへ移行
func (game_map *GameMap) copyFromXY(xy *xymap) {
	for x := 0; x < xy.mapSize.MaxX; x++ {
		for y := 0; y < xy.mapSize.MaxY; y++ {
			macro := xy.getMatrix(x, y);
			high := xy.getHigh(x, y);
			game_map.High[y][x] = high;
			for z := 0; z <= high; z++ {
				game_map.MacroMapTypes[z][y][x] = macro;
			}
		}
	}
}

//パーツとのひも付け
func (game_map *GameMap) bindToGameParts(gamePartsDict map[string]GameParts) {



	/*
	for x := 0; x < game_map.Size.MaxX; x++ {
		for y := 0; y < game_map.Size.MaxY; y++ {
			for z := 0; z < game_map.Size.MaxZ; z++ {
//				macro := game_map.MacroMapTypes[z][y][x]
//				high := xymap.high[y][x];
				for z := 0; z <= high; z++ {
//					game_map.JungleGym[z][y][x] = GetGameParts(macro, game_map.Geographical, z);
				}
			}
		}
	}
	*/
}
