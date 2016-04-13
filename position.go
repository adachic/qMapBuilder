package main

import (
	"math"

	"github.com/adachic/lottery"
"math/rand"
	"time"
)

//criteriaを基準点とし、radiansの中間の角度へ線を伸ばし、四角形との接点の配列を返す
//radiansはソート済みの角度の一覧
func CreateEdgePositionInMap(mapSize GameMapSize, criteria GameMapPosition, radians []float64) (edges []GameMapPosition){
	radNum := len(radians)
	for i := 0 ; i < radNum; i++{
		var middleRad float64
		if i == 0 {
			tmp := (radians[0] + 2.0 * math.Pi * radians[radNum - 1]) / 2.0
			middleRad = radians[0] - tmp
		} else {
			middleRad = (radians[i] - radians[i-1]) / 2.0 + radians[i - 1]
		}
		//100あれば接点とぶつかるだろう
		x := 100 * math.Cos(middleRad)
		y := 100 * math.Sin(middleRad)
		x2 := criteria.X + int(x)
		y2 := criteria.Y + int(y)

		if x2 >= mapSize.MaxX {
			x2 = mapSize.MaxX - 1
		}
		if x2 < 0 {
			x2 = 0
		}
		if y2 >= mapSize.MaxY {
			y2 = mapSize.MaxY - 1
		}
		if y2 < 0 {
			y2 = 0
		}
		edges = append(edges, GameMapPosition{x2, y2, 0})
	}
	return edges
}

//四角形のなかから座標を抽選で決定して返す
//そして、criteriaからのradianも返す
//サイズ、基準点、基準点からの距離
func CreateRandomPositionInMap(mapSize GameMapSize, criteria GameMapPosition, distance int) (GameMapPosition ,float64){
	var x2, y2 int
	var radian float64
	for {
		//中央からの距離比率
		distanceFrom := distance
		//角度
		degree := lottery.GetRandomInt(0, 360)
		radian = float64(degree) / (math.Pi * 2.0)
		//半径
		var r float64
		if mapSize.MaxX > mapSize.MaxY {
			r = float64(mapSize.MaxX) / 2.0 * float64(distanceFrom) / 100.0
		} else {
			r = float64(mapSize.MaxY) / 2.0 * float64(distanceFrom) / 100.0
		}
		x := r * math.Cos(radian)
		y := r * math.Sin(radian)
		x2 = criteria.X + int(x)
		y2 = criteria.Y + int(y)

		if x2 >= mapSize.MaxX {
			x2 = mapSize.MaxX - 1
		}
		if x2 < 0 {
			x2 = 0
		}
		if y2 >= mapSize.MaxY {
			y2 = mapSize.MaxY - 1
		}
		if y2 < 0 {
			y2 = 0
		}

		distanceFromCriteria := math.Sqrt(math.Pow((float64(x2-criteria.X)), 2) +
			math.Pow((float64(y2-criteria.Y)), 2))

		DDDlogln("distanceFromCreteria:", distanceFromCriteria,
			" ?", float64(r-1), "\n")

		if distanceFromCriteria < float64(r-1) {
			continue
		}
		DDDlogln("distanceFrom:", distanceFrom,
			" degree:", degree,
			" radian:", radian,
			" r:", r,
			" x:", x,
			" y:", y,
			" x2:", x2,
			" y2:", y2,
			" criteriaX:", criteria.X,
			" criteriaY:", criteria.Y,
			"\n",
		)
		break
	}
	//TODO:限界リトライ数を定める
	return GameMapPosition{x2, y2, 0}, radian
}

//長方形の形式の抽選結果を返す
func CreateRectForm() RectForm {
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
func CreateAspectOfRectFrom(rectForm RectForm) float32 {
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

	DDDlog("longer: %+v\n", longer)
	DDDlog("shorter: %+v\n", shorter)
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
func CreateArea(difficult Difficult) int {
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

type Range struct {
	Min int
	Max int
}
