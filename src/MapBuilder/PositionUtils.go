package MapBuilder
import (
"github.com/adachic/lottery"
"math"
)

//四角形のなかから座標を抽選で決定して返す
//サイズ、基準点、基準点からの距離
func createRandomPositionInMap(mapSize GameMapSize, criteria GameMapPosition, distance int) GameMapPosition{
	//中央からの距離比率
	distanceFrom := distance
	//角度
	degree := lottery.GetRandomInt(0, 360)
	radian := float64(degree) / (math.Pi * 2.0)
	//半径
	var r float64
	if (mapSize.MaxX > mapSize.MaxY) {
		r = float64(mapSize.MaxX) / 2.0 * float64(distanceFrom) / 100.0
	}else {
		r = float64(mapSize.MaxY) / 2.0 * float64(distanceFrom) / 100.0
	}
	x := r * math.Cos(radian)
	y := r * math.Sin(radian)
	x2 := criteria.X + int(x)
	y2 := criteria.Y + int(y)

	if (x2 >= mapSize.MaxX) {
		x2 = mapSize.MaxX - 1
	}
	if (x2 < 0) {
		x2 = 0
	}
	if (y2 >= mapSize.MaxY) {
		y2 = mapSize.MaxY - 1
	}
	if (y2 < 0) {
		y2 = 0
	}
	//TODO:距離が守れてなければリトライ
	return GameMapPosition{x2, y2, 0}
}

