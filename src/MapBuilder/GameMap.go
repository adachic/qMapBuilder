package MapBuilder

import (
	"github.com/adachic/lottery"
	"time"
	"math/rand"
	"math"
	"fmt"
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

//マップの難度
type Difficult int
const (
	easy Difficult = 10
	normal Difficult = 3
	hard Difficult = 2
	exhard Difficult = 1
)

//長方形の形式
type RectForm int
const (
	horizontalLong RectForm = 6//横長
	verticalLong RectForm = 5
	square RectForm = 2
)

//地形
type Geographical int
const (
	GeographicalStep Geographical = 14 + 10
	GeographicalMountain Geographical = 9 + 10
	GeographicalCave Geographical = 8 + 10
	GeographicalFort Geographical = 7 + 10
	GeographicalShrine Geographical = 6 + 10
	GeographicalTown Geographical = 5 + 10
	GeographicalCastle Geographical = 4 + 10
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

	x := lottery.GetRandomNormInt(minLength, maxLength);
	y := (minLength + maxLength) - x
	if x > y {
		longer = x
		shorter = y
	}else {
		longer = y
		shorter = x
	}

	fmt.Printf("longer: %+v\n", longer)
	fmt.Printf("shorter: %+v\n", shorter)
	switch rectForm {
	case horizontalLong:
		ret = float32(longer) / float32(shorter)
		break;
	case verticalLong:
		ret = float32(shorter) / float32(longer)
		break;
	case square:
		ret = 1.0
		break;
	}
	return ret
}

//マップ面積を返す
func createArea(difficult Difficult) int {
	var ret int
	//10x10を最小とし、30x30を最大とする
	switch difficult {
	case easy:
		ret = lottery.GetRandomInt(10 * 10, 15 * 15)
		break;
	case normal:
		ret = lottery.GetRandomInt(10 * 10, 25 * 25)
		break;
	case hard:
		ret = lottery.GetRandomInt(15 * 15, 30 * 30)
		break;
	case exhard:
		ret = lottery.GetRandomInt(15 * 15, 30 * 30)
		break;
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
	if (y < 1) {
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
		break;
	case normal:
		seed = distanceFromCenter{30, 80}
		break;
	case hard:
		seed = distanceFromCenter{0, 50}
		break;
	case exhard:
		seed = distanceFromCenter{0, 20}
		break;
	}
	//中央からの距離比率
	distanceFrom := lottery.GetRandomInt(seed.distanceFromCenterMin, seed.distanceFromCenterMax)
	//角度
	degree := lottery.GetRandomInt(0, 360)
	radian := float64(degree) / (math.Pi * 2.0)
	//半径
	var r float64
	if(mapSize.MaxX > mapSize.MaxY){
		r = float64(mapSize.MaxX) / 2.0 * float64(distanceFrom) / 100.0
	}else{
		r = float64(mapSize.MaxY) / 2.0 * float64(distanceFrom) / 100.0
	}
	x := r * math.Cos(radian)
	y := r * math.Sin(radian)
	x2 := int(float64(mapSize.MaxX) / 2.0 + x)
	y2 := int(float64(mapSize.MaxY) / 2.0 + y)

	if(x2 >= mapSize.MaxX){
		x2 = mapSize.MaxX - 1
	}
	if(x2 < 0){
		x2 = 0
	}
	if(y2 >= mapSize.MaxY){
		y2 = mapSize.MaxY - 1
	}
	if(y2 < 0){
		y2 = 0
	}
	return GameMapPosition{x2, y2, 0}
}

//敵の出現座標の一覧を返す
func createEnemyStartPoint() []GameMapPosition {

	return nil
}

//雑に100回まわしてみる
func bulc(){
	x := 100
	for x > 0 {
		x--
		flow()
		fmt.Printf("\n")
	}
}

func flow(){
	//難易度を決定
	difficult := createMapDifficult()
	fmt.Printf("difficult: %+v\n", difficult)

	//マップのサイズを決定
	mapSize := createMapSize(difficult)
	fmt.Printf("mapSize: %+v\n", mapSize)

	//大まかな地形を決定
	//geographical := createMapGeographical()
	createMapGeographical()

	//味方開始ポイントを決定
	allyStartPoint := createAllyStartPoint(difficult, mapSize)
	fmt.Printf("allyStartPoint: %+v\n", allyStartPoint)

	//敵開始ポイントを決定
	//	createEnemyStartPoints()


	//広場生成

	//道生成

	//壁生成

	//ラフ生成

	//勾配生成


	//バリデーション
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


