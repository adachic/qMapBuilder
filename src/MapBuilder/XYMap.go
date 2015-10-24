package MapBuilder
import (
	"github.com/adachic/lottery"
	"math"
)

type Range struct {
	Min int
	Max int
}

//広場生成
func createPlaza(xyMap *[][]MacroMapType,
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

		i := 0
		r := 0
		for i < plazaArea {

			*xyMap[x][]

		}

		sideLength := int(math.Sqrt(float64(plazaArea)))
		if(sideLength <=0){
			sideLength = 1
		}
		offsStart := -sideLength/2
		offsEnd := sideLength/2

		//円形に広場生成
		for offsX := offsStart; offsX <= offsEnd ; offsX++{

		}

	}
	for i := 0; i < plazaCount; i++ {




	}
}

//不可侵領域で埋める
func fillCantEnter(xyMap *[][]MacroMapType, mapSize GameMapSize) {
	for x := 0; x := mapSize.MaxX; x++ {
		for y := 0; y := mapSize.MaxY; y++ {
			*xyMap[y][x] = MacroMapTypeCantEnter
		}
	}
}

//見下ろしマップを返す
func createXYMap(difficult Difficult,
mapSize GameMapSize,
geographical Geographical,
allyStartPoint GameMapPosition,
enemyStartPoints []GameMapPosition) [][]MacroMapType {

	var xyMap [mapSize.MaxX][mapSize.MaxY]MacroMapType

	//不可侵領域で埋める
	fillCantEnter(&xyMap, mapSize)

	//広場生成
	createPlaza(&xyMap, difficult, mapSize, allyStartPoint, enemyStartPoints)

	//道生成

	//壁生成

	//ラフ生成

	return xyMap
}

