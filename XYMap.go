package main

import (
	"math"
	"github.com/adachic/lottery"
	"fmt"
)

type xymap struct {
	mapSize GameMapSize
	matrix [][]MacroMapType
}

func NewXYMap(mapSize GameMapSize) *xymap{
	xy := &xymap{}
	return xy.init(mapSize)
}


func (xy *xymap) init(mapSize GameMapSize) *xymap{
	xy.mapSize = mapSize
	xy.matrix = make([][]MacroMapType, mapSize.MaxY)
	for y := 0; y < mapSize.MaxY; y++ {
		xy.matrix[y] = make([]MacroMapType, mapSize.MaxX)
	}
	xy.fillCantEnter()
	return xy
}

//不可侵領域で埋める
func (xy *xymap) fillCantEnter() {
	for x := 0; x < xy.mapSize.MaxX; x++ {
		for y := 0; y < xy.mapSize.MaxY; y++ {
			xy.matrix[y][x] = MacroMapTypeCantEnter
		}
	}
}

//広場生成
func (xy *xymap) putPlazas(
difficult Difficult,
allyStartPoint GameMapPosition,
enemyStartPoints []GameMapPosition) {

	area := xy.mapSize.area()

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
		xy.putPlaza(plazaArea, centerPoint)
	}

	//敵用広場生成
	for i := 0; i < len(enemyStartPoints); i++ {
		plazaArea := lottery.GetRandomInt(plazaSizeRange.Min, plazaSizeRange.Max)
		centerPoint := enemyStartPoints[i]
		xy.putPlaza(plazaArea, centerPoint)
	}
}

//広場生成
//centerPointを中心として、正方形に道を描画する(面積はplazaArea)
func (xy *xymap) putPlaza(
plazaArea int,
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
			if x2 >= xy.mapSize.MaxX {
				x2 = xy.mapSize.MaxX - 1
			}
			if x2 < 0 {
				x2 = 0
			}
			if y2 >= xy.mapSize.MaxY {
				y2 = xy.mapSize.MaxY - 1
			}
			if y2 < 0 {
				y2 = 0
			}
			xy.matrix[y2][x2] = MacroMapTypeLoad
		}
	}
}

//pointをmacroMapTypeにする
func (xy *xymap) putPoint(point GameMapPosition, macroMapType MacroMapType){
	xy.matrix[point.Y][point.X] = macroMapType
}

func (xy *xymap) printMapForDebug() {
	for y := 0; y < xy.mapSize.MaxY; y++ {
		for x := 0; x < xy.mapSize.MaxX; x++ {
			switch xy.matrix[y][x] {
			case MacroMapTypeCantEnter:
				fmt.Print("#")
			case MacroMapTypeLoad:
				fmt.Print(".")
			case MacroMapTypeAllyPoint:
				fmt.Print("A")
			case MacroMapTypeEnemyPoint:
				fmt.Print("E")
			}
		}
		fmt.Print("\n")
	}
}
