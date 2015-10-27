package main

import (
	"math"

	"github.com/adachic/lottery"
)

type Range struct {
	Min int
	Max int
}

//centerPointを中心として、正方形に道を描画する(面積はplazaArea)
func putLoadSquareFromCenter(
	plazaArea int,
	xyMap [][]MacroMapType,
	mapSize GameMapSize,
	centerPoint GameMapPosition) {

	sideLength := int(math.Sqrt(float64(plazaArea)))
	if sideLength <= 0 {
		sideLength = 1
	}
	offsStart := -sideLength / 2
	offsEnd := sideLength / 2

	//広場生成(正方形)
	for offsX := offsStart; offsX <= offsEnd; offsX++ {
		for offsY := offsStart; offsY <= offsEnd; offsY++ {
			x2 := centerPoint.X + offsX
			y2 := centerPoint.Y + offsY
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
			xyMap[y2][x2] = MacroMapTypeLoad
		}
	}
}

//広場生成
func createPlaza(xyMap [][]MacroMapType,
	difficult Difficult,
	mapSize GameMapSize,
	allyStartPoint GameMapPosition,
	enemyStartPoints []GameMapPosition) {

	area := mapSize.area()

	//余計に広場を追加？
	additionalPlazaCount := lottery.GetRandomInt(0, 5)
	plazaCount := (1 + len(enemyStartPoints) + additionalPlazaCount)

	//出撃ポイントの数に対するマップの大きさ
	maxAriaPerPoint := area / plazaCount
	plazaSizeRange := Range{3, maxAriaPerPoint}

	//味方用広場生成
	{
		plazaArea := lottery.GetRandomInt(plazaSizeRange.Min, plazaSizeRange.Max)
		centerPoint := allyStartPoint
		putLoadSquareFromCenter(plazaArea, xyMap, mapSize, centerPoint)
	}
	//敵用広場生成
	for i := 0; i < len(enemyStartPoints); i++ {
		plazaArea := lottery.GetRandomInt(plazaSizeRange.Min, plazaSizeRange.Max)
		centerPoint := enemyStartPoints[i]
		putLoadSquareFromCenter(plazaArea, xyMap, mapSize, centerPoint)
	}
	//余計な分生成
	/*
		for i := 0; i < additionalPlazaCount; i++ {
			plazaArea := lottery.GetRandomInt(plazaSizeRange.Min, plazaSizeRange.Max)
			centerPoint :=
			putLoadSquareFromCenter(plazaArea, xyMap, mapSize, centerPoint)
		}
	*/
}

//不可侵領域で埋める
func fillCantEnter(xyMap [][]MacroMapType, mapSize GameMapSize) {
	for x := 0; x < mapSize.MaxX; x++ {
		for y := 0; y < mapSize.MaxY; y++ {
			xyMap[y][x] = MacroMapTypeCantEnter
		}
	}
}

//見下ろしマップを返す
func createXYMap(difficult Difficult,
	mapSize GameMapSize,
	geographical Geographical,
	allyStartPoint GameMapPosition,
	enemyStartPoints []GameMapPosition) [][]MacroMapType {

	//	var xyMap [mapSize.MaxX][mapSize.MaxY]MacroMapType
	//	xyMap := new([mapSize.MaxX][mapSize.MaxY]MacroMapType)
	//	var xyMap [10][10]MacroMapType

	xyMap := make([][]MacroMapType, mapSize.MaxY, mapSize.MaxY)
	for y := 0; y < mapSize.MaxY; y++ {
		xyMap[y] = make([]MacroMapType, mapSize.MaxX, mapSize.MaxX)
	}

	//不可侵領域で埋める
	fillCantEnter(xyMap, mapSize)

	//広場生成
	createPlaza(xyMap, difficult, mapSize, allyStartPoint, enemyStartPoints)

	xyMap[allyStartPoint.Y][allyStartPoint.X] = MacroMapTypeAllyPoint
	for i := 0; i < len(enemyStartPoints); i++ {
		xyMap[enemyStartPoints[i].Y][enemyStartPoints[i].X] = MacroMapTypeEnemyPoint
	}

	//道生成

	//壁生成

	//ラフ生成

	return xyMap
}
