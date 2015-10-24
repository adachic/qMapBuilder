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

//確率を返す
func (d Difficult) Prob() int {
	return int(d)
}

func (d RectForm) Prob() int {
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

func CreateGameMap(gamePartsDict map[string]GameParts) GameMap {
	x := 100
	for x > 0 {
		x--
		//難易度を決定
		difficult := createMapDifficult()
		fmt.Printf("difficult: %+v\n", difficult)

		//マップのサイズを決定
		mapSize := createMapSize(difficult)
		fmt.Printf("mapSize: %+v\n\n", mapSize)
	}
	//難易度を決定
	difficult := createMapDifficult()
	fmt.Printf("difficult: %+v\n", difficult)

	//マップのサイズを決定
	mapSize := createMapSize(difficult)
	fmt.Printf("mapSize: %+v\n", mapSize)

	//大まかな地形を決定

	//味方開始ポイントを決定

	//敵開始ポイントを決定


	//広場生成

	//道生成

	//壁生成

	//ラフ生成

	//勾配生成


	//バリデーション

	return GameMap{}
}

func createJsonFromMap() {

}

func createPngFromMap() {

}


