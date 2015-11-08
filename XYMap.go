package main

import (
	"math"
	"github.com/adachic/lottery"
	"fmt"
)

type xymap struct {
	mapSize GameMapSize
	matrix [][]MacroMapType //種別
	high [][]int //高さ
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

	xy.high = make([][]int, mapSize.MaxY)
	for y := 0; y < mapSize.MaxY; y++ {
		xy.high[y] = make([]int, mapSize.MaxX)
	}
	xy.fillHeightZero()

	return xy
}

//高さ0で埋める
func (xy *xymap) fillHeightZero() {
	for x := 0; x < xy.mapSize.MaxX; x++ {
		for y := 0; y < xy.mapSize.MaxY; y++ {
			xy.high[y][x] = 1
		}
	}
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

//道を引く
func (xy *xymap) putRoads(
difficult Difficult,
allyStartPoint GameMapPosition,
enemyStartPoints []GameMapPosition) {
	//最も近いポイントに単純に線を引く方式

	mapPositions := make([]GameMapPosition, 0)
	{
		mapPositions = append(mapPositions, allyStartPoint)
		for i := 0 ; i < len(enemyStartPoints) ; i++ {
			mapPositions = append(mapPositions, enemyStartPoints[i])
		}
	}


	alreadyPutPathPositions := make([]PathPosition,0)

	for i := 0 ; i< len(mapPositions); i++{
		src := mapPositions[i]

		//最も近いポイントを検索
		dst ,err := src.searchNearPositionWithOutMe(mapPositions)
		if(err){
			//近いポイントなかったマン(ありえない)
		}

		path := PathPosition{src, dst}
		if (containsPath(alreadyPutPathPositions, path)){
			//すでに道はひかれている
			continue
		}

		//直線道路を引く
		xy.putRoadStraight(path)

		alreadyPutPathPositions = append(alreadyPutPathPositions, path)
	}

	/*
	//何をつなげるか
	for i := 0 ; i< len(mapPositions); i++{
		alreadyPos := make([]GameMapPosition, 0)

		src := pathPositions[i]
		append(alreadyPos, src)

		//いくつのポイントとつなげるか
		maxPosIdx := len(mapPositions - 1)
		if (maxPosIdx < 1){
			maxPosIdx = 1
		}
		pathNum := lottery.GetRandomInt(1, maxPosIdx)

		//何とつなげるか
		for j := 0 ; j< pathNum ; j++ {
			aho:
			dstPosIdx := lottery.GetRandomInt(0, maxPosIdx)
			for k := 0; k < len(alreadyPos); k++ {
				if alreadyPos[k] == mapPositions[dstPosIdx] {
					break aho
				}
			}
			append(alreadyPos, mapPositions[dstPosIdx])
			append(pathPositions, PathPosition{src, mapPositions[dstPosIdx]})
		}
	}

	//それぞれの道生成
	for i := 0 ; i< len(pathPositions) ; i++ {
		xy.putRoad(pathPositions[i])
	}
	*/
}

//道(直線)
type PathPosition struct {
	src GameMapPosition
	dst GameMapPosition
}

//pathPositionsにtargetを含むならtrue
func containsPath(pathPositions []PathPosition, target PathPosition) bool {
	for i := 0 ; i< len(pathPositions); i++{
		if(pathPositions[i].equalXYTo(target)){
			return true
		}
	}
	return false
}

//同じパスか,逆の組み合わせも判定
func (path PathPosition) equalXYTo(another PathPosition) bool{
	if(path.src.equalXYTo(another.src) && path.dst.equalXYTo(another.dst)){
		return true
	}
	if(path.src.equalXYTo(another.dst) && path.dst.equalXYTo(another.src)){
		//逆の組み合わせ
		return true
	}
	return false
}

//道を引く
func (xy *xymap) putRoadStraight(path PathPosition){
	// y = ax + b
	offsY := float64(path.dst.Y) - float64(path.src.Y)
	offsX := float64(path.dst.X) - float64(path.src.X)
	a := offsY / offsX
	b := float64(path.src.Y) + 0.5 - a * (float64(path.src.X) + 0.5)

	minX, maxX, minY, maxY := getMinMaxXY(path)

	for y := minY ; y <= maxY ; y++ {
		for x := minX; x <= maxX; x++ {
			if(!isHitToSquare(a, b, float64(x), float64(y))){
				continue
			}
			xy.putPoint(GameMapPosition{x,y,0}, MacroMapTypeLoad)
		}
	}
}

//x,yを起点とする1x1の四角形にy=ax+bの直線が重なるならtrue
func isHitToSquare(a float64, b float64, x float64, y float64) bool{
	//下の辺
	{
		yy1 := x * a + b
		yy2 := (x+1.0) * a + b
		if(yy1 <= y && yy2 > y){
			return true
		}
		if(yy1 > y && yy2 <= y){
			return true
		}
	}
	//左の辺
	{
		xx1 := (y - b) / a
		xx2 := (y + 1.0 - b) / a
		if(xx1 <= x && xx2 > x){
			return true
		}
		if(xx1 > x && xx2 <= x){
			return true
		}
	}
	//上の辺
	{
		yy1 := x * a + b
		yy2 := (x+1.0) * a + b
		if(yy1 <= (y + 1.0) && yy2 > (y + 1.0)){
			return true
		}
		if(yy1 > (y + 1.0) && yy2 <= (y + 1.0)){
			return true
		}
	}
	//右の辺
	{
		xx1 := (y - b) / a
		xx2 := (y + 1.0 - b) / a
		if(xx1 <= (x + 1.0) && xx2 > (x + 1.0)){
			return true
		}
		if(xx1 > (x + 1.0) && xx2 <= (x + 1.0)){
			return true
		}
	}
	return false
}

func getMinMaxXY(path PathPosition) (minX int, maxX int, minY int, maxY int){
	if(path.src.X < path.dst.X){
		minX = path.src.X
		maxX = path.dst.X
	}else{
		minX = path.dst.X
		maxX = path.src.X
	}
	if(path.src.Y < path.dst.Y){
		minY = path.src.Y
		maxY = path.dst.Y
	}else{
		minY = path.dst.Y
		maxY = path.src.Y
	}
	return minX, maxX, minY, maxY
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
		fmt.Print("   ")
		for x := 0; x < xy.mapSize.MaxX; x++ {
			fmt.Print(xy.high[y][x])
		}
		fmt.Print("\n")
	}
}

//勾配を生成
//ルール: x,yが大きいほど高い
func (xy *xymap) makeGradient(geo Geographical) {
	rowestHigh := 0
	//まず地形でだいたいの高さ
	switch geo {
	case GeographicalStep    :
		rowestHigh = 3
		break
	case GeographicalMountain:
		rowestHigh = 7
		break
	case GeographicalCave    :
		rowestHigh = 7
		break
	case GeographicalFort    :
		rowestHigh = 4
		break
	case GeographicalShrine  :
		rowestHigh = 4
		break
	case GeographicalTown    :
		rowestHigh = 4
		break
	case GeographicalCastle  :
		rowestHigh = 4
		break
	}

	//手前から道の高さを調整していく（歩けるように高さ調整していく）
	xy.makeGradientRoad(rowestHigh)

	//だんだんと段差になっていく


	//それ以外のところは2個あがったり下がったりさせる
}

//手前から道の高さを調整していく（歩けるように高さ調整していく）
func (xy *xymap) makeGradientRoad(rowestHigh int) {
	currentHeight := rowestHigh
	for y := 0; y < xy.mapSize.MaxY; y++ {
		for x := 0; x < xy.mapSize.MaxX; x++ {
			switch xy.matrix[y][x] {
			case MacroMapTypeLoad:
				fallthrough
			case MacroMapTypeAllyPoint:
				fallthrough
			case MacroMapTypeEnemyPoint:
				xy.high[y][x] = currentHeight
			}
		}
	}
}
