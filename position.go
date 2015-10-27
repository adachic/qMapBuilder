package main

import (
	"fmt"
	"math"

	"github.com/adachic/lottery"
)

//四角形のなかから座標を抽選で決定して返す
//サイズ、基準点、基準点からの距離
func createRandomPositionInMap(mapSize GameMapSize, criteria GameMapPosition, distance int) GameMapPosition {

	var x2, y2 int
	for {
		//中央からの距離比率
		distanceFrom := distance
		//角度
		degree := lottery.GetRandomInt(0, 360)
		radian := float64(degree) / (math.Pi * 2.0)
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

		fmt.Print("distanceFromCreteria:", distanceFromCriteria,
			" ?", float64(r-1), "\n")

		if distanceFromCriteria < float64(r-1) {
			continue
		}
		fmt.Print("distanceFrom:", distanceFrom,
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
	return GameMapPosition{x2, y2, 0}
}
