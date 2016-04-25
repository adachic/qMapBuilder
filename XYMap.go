package main

import (
	"math"
	"github.com/adachic/lottery"
	"github.com/ojrac/opensimplex-go"
	"strconv"
	"fmt"
)

type xymap struct {
	mapSize    GameMapSize
	matrix     [][]MacroMapType //種別
	high       [][]int          //高さ
	areaId     [][]int          //A*のためのエリアID
	areaPath   [][]int
	maxAreaId  int
	zoneMarked bool
}

func (xy *xymap) getMatrix(x int, y int) MacroMapType {
	return xy.matrix[y][x]
}

func (xy *xymap) getHigh(x int, y int) int {
	return xy.high[y][x]
}

func (xy *xymap) getAreaId(x int, y int) int {
	return xy.areaId[y][x]
}

func NewXYMap(mapSize GameMapSize) *xymap {
	xy := &xymap{}
	return xy.init(mapSize)
}

func (xy *xymap) init(mapSize GameMapSize) *xymap {
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

	xy.areaId = make([][]int, mapSize.MaxY)
	for y := 0; y < mapSize.MaxY; y++ {
		xy.areaId[y] = make([]int, mapSize.MaxX)
	}
	xy.fillArea0()
	return xy
}

//areaId=-1で埋める
func (xy *xymap) fillArea0() {
	for x := 0; x < xy.mapSize.MaxX; x++ {
		for y := 0; y < xy.mapSize.MaxY; y++ {
			xy.areaId[y][x] = -1
		}
	}
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
			xy.matrix[y2][x2] = MacroMapTypeRoad
		}
	}
}

//pointをmacroMapTypeにする
func (xy *xymap) putPoint(point GameMapPosition, macroMapType MacroMapType) {
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
		for i := 0; i < len(enemyStartPoints); i++ {
			mapPositions = append(mapPositions, enemyStartPoints[i])
		}
	}

	alreadyPutPathPositions := make([]PathPosition, 0)

	for i := 0; i < len(mapPositions); i++ {
		src := mapPositions[i]

		//最も近いポイントを検索
		dst, err := src.searchNearPositionWithOutMe(mapPositions, alreadyPutPathPositions)
		if (err) {
			//近いポイントなかったマン(ありえない)
			continue
		}

		path := PathPosition{src, dst}
		/*
		if (containsPath(alreadyPutPathPositions, path)){
			//すでに道はひかれている
			continue
		}
		*/

		//直線道路を引く
		xy.putRoadStraight(path)
		alreadyPutPathPositions = append(alreadyPutPathPositions, path)
	}
}

//道(直線)
type PathPosition struct {
	src GameMapPosition
	dst GameMapPosition
}

//pathPositionsにtargetを含むならtrue
func containsPath(pathPositions []PathPosition, target PathPosition) bool {
	for i := 0; i < len(pathPositions); i++ {
		if (pathPositions[i].equalXYTo(target)) {
			return true
		}
	}
	return false
}

//同じパスか,逆の組み合わせも判定
func (path PathPosition) equalXYTo(another PathPosition) bool {
	if (path.src.equalXYTo(another.src) && path.dst.equalXYTo(another.dst)) {
		return true
	}
	if (path.src.equalXYTo(another.dst) && path.dst.equalXYTo(another.src)) {
		//逆の組み合わせ
		return true
	}
	return false
}

//道を引く
func (xy *xymap) putRoadStraight(path PathPosition) {
	// y = ax + b
	offsY := float64(path.dst.Y) - float64(path.src.Y)
	offsX := float64(path.dst.X) - float64(path.src.X)
	if (path.dst.X == path.src.X) {
		offsX = 0.01
	}
	if (path.dst.Y == path.src.Y) {
		offsY = 0.01
	}
	a := offsY / offsX
	b := float64(path.src.Y) + 0.5 - a * (float64(path.src.X) + 0.5)

	minX, maxX, minY, maxY := getMinMaxXY(path)

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if (!isHitToSquare(a, b, float64(x), float64(y))) {
				continue
			}
			xy.putPoint(GameMapPosition{x, y, 0}, MacroMapTypeRoad)
		}
	}
}

//x,yを起点とする1x1の四角形にy=ax+bの直線が重なるならtrue
func isHitToSquare(a float64, b float64, x float64, y float64) bool {
	//下の辺
	{
		yy1 := x * a + b
		yy2 := (x + 1.0) * a + b
		if (yy1 <= y && yy2 > y) {
			return true
		}
		if (yy1 > y && yy2 <= y) {
			return true
		}
	}
	//左の辺
	{
		xx1 := (y - b) / a
		xx2 := (y + 1.0 - b) / a
		if (xx1 <= x && xx2 > x) {
			return true
		}
		if (xx1 > x && xx2 <= x) {
			return true
		}
	}
	//上の辺
	{
		yy1 := x * a + b
		yy2 := (x + 1.0) * a + b
		if (yy1 <= (y + 1.0) && yy2 > (y + 1.0)) {
			return true
		}
		if (yy1 > (y + 1.0) && yy2 <= (y + 1.0)) {
			return true
		}
	}
	//右の辺
	{
		xx1 := (y - b) / a
		xx2 := (y + 1.0 - b) / a
		if (xx1 <= (x + 1.0) && xx2 > (x + 1.0)) {
			return true
		}
		if (xx1 > (x + 1.0) && xx2 <= (x + 1.0)) {
			return true
		}
	}
	return false
}

func getMinMaxXY(path PathPosition) (minX int, maxX int, minY int, maxY int) {
	if (path.src.X < path.dst.X) {
		minX = path.src.X
		maxX = path.dst.X
	}else {
		minX = path.dst.X
		maxX = path.src.X
	}
	if (path.src.Y < path.dst.Y) {
		minY = path.src.Y
		maxY = path.dst.Y
	}else {
		minY = path.dst.Y
		maxY = path.src.Y
	}
	return minX, maxX, minY, maxY
}

func (xy *xymap) printMapForDebug() {
	fmt.Printf("\n")
	for y := (xy.mapSize.MaxY - 1); y >= 0; y-- {
		fmt.Printf("%02d ", y)
		for x := 0; x < xy.mapSize.MaxX; x++ {
			switch xy.matrix[y][x] {
			case MacroMapTypeCantEnter:
				fmt.Print("#")
			case MacroMapTypeRoad:
				fmt.Print(".")
			case MacroMapTypeRough:
				fmt.Print(";")
			case MacroMapTypeWall:
				fmt.Print("=")
			case MacroMapTypeAllyPoint:
				fmt.Print("A")
			case MacroMapTypeEnemyPoint:
				fmt.Print("E")
			case MacroMapTypeSwampWater: //水系
				fallthrough
			case MacroMapTypeSwampRava:
				fallthrough
			case MacroMapTypeSwampPoison:
				fallthrough
			case MacroMapTypeSwampHeal:
				fmt.Print("~")
			}
		}
		fmt.Print("   ")
		for x := 0; x < xy.mapSize.MaxX; x++ {
			fmt.Printf("%2d", xy.high[y][x])
		}
		fmt.Print("\n")
	}
	fmt.Printf("\n")
	for y := (xy.mapSize.MaxY - 1); y >= 0; y-- {
		fmt.Printf("%02d ", y)
		for x := 0; x < xy.mapSize.MaxX; x++ {
			switch xy.matrix[y][x] {
			case MacroMapTypeCantEnter:
				fmt.Print("#")
			case MacroMapTypeRoad:
				fmt.Print(".")
			case MacroMapTypeRough:
				fmt.Print(";")
			case MacroMapTypeWall:
				fmt.Print("=")
			case MacroMapTypeAllyPoint:
				fmt.Print("A")
			case MacroMapTypeEnemyPoint:
				fmt.Print("E")
			case MacroMapTypeSwampWater: //水系
				fallthrough
			case MacroMapTypeSwampRava:
				fallthrough
			case MacroMapTypeSwampPoison:
				fallthrough
			case MacroMapTypeSwampHeal:
				fmt.Print("~")
			}
		}
		fmt.Print("   ")
		for x := 0; x < xy.mapSize.MaxX; x++ {
			fmt.Printf("%2d", xy.areaId[y][x])
		}
		fmt.Print("\n")
	}
}

//水地形(池)を生成
func (xy *xymap) makeSwamp(geo Geographical) {
	//ゾーニング
	zones := xy.zoningForMakeSwamp()
	for _, val := range zones {
		DDDlog("ahongoS:%d \n", len(val))
	}

	actualZones := trimZoneSizeForSwamp(zones, geo)

	//どれくらい水を配置するのか？
	mounts := mountsForSwamp(actualZones, geo)
	DDDlog("amounts:%d \n", mounts)

	//どのidxに配置するのか？
	alreadyPutsIdx := []int{}
	i := 0
	for ; i < mounts; {
		idx := lottery.GetRandomInt(0, len(actualZones))
		if (idx > (len(actualZones) - 1)) {
			idx = len(actualZones) - 1
		}
		contain := false
		for _, val := range alreadyPutsIdx {
			if idx == val {
				contain = true
			}
		}
		if contain {
			continue
		}
		alreadyPutsIdx = append(alreadyPutsIdx, idx)
		i++
	}

	//池を生成
	for _, idx := range alreadyPutsIdx {
		targetZone := actualZones[idx]
		xy.putSwamp(targetZone, geo)
	}
}

//池を生成
func (xy *xymap) putSwamp(targetZone []GameMapPosition, geo Geographical) {
	var swampType MacroMapType
	swampType = MacroMapTypeSwampWater
	switch geo {
	case GeographicalStep:
	case GeographicalCave:
	case GeographicalSnow:
	case GeographicalPoison:
		swampType = MacroMapTypeSwampPoison
	case GeographicalFire:
		swampType = MacroMapTypeSwampRava
	default:
	}
	for _, pos := range targetZone {
		//harfか？
		high := xy.high[pos.Y][pos.X]
		if (high % 2) == 0 {
			//harfじゃない
			continue
		}
		//エッジ部分には水置かない
		DDDlog("unko2000")
		if (isEdgeInZone(targetZone, pos)) {
			continue
		}
		DDDlog("unko2001")
		xy.matrix[pos.Y][pos.X] = swampType
	}

}

//posはzoneの中のエッジ部分ならtrue
func isEdgeInZone(zone []GameMapPosition, pos GameMapPosition) bool {
	NotEdgeTop := false
	NotEdgeBottom := false
	NotEdgeRight := false
	NotEdgeLeft := false

	for _, posInZone := range zone {
		if ((posInZone.X == pos.X - 1) && (posInZone.Y == pos.Y)) {
			NotEdgeLeft = true
		}
		if ((posInZone.X == pos.X + 1) && (posInZone.Y == pos.Y)) {
			NotEdgeRight = true
		}
		if ((posInZone.Y == pos.Y + 1) && (posInZone.X == pos.X)) {
			NotEdgeTop = true
		}
		if ((posInZone.Y == pos.Y - 1) && (posInZone.X == pos.X)) {
			NotEdgeBottom = true
		}
	}

	NotEdge := NotEdgeTop && NotEdgeBottom && NotEdgeRight && NotEdgeLeft
	return !NotEdge
}

//zonesのなかから、池を張ってもいいくらいの大きさの広場を抽出する
func trimZoneSizeForSwamp(zones [][]GameMapPosition, geo Geographical) [][]GameMapPosition {
	filtered := [][]GameMapPosition{}
	i := 0
	for _, zone := range zones {
		//エッジを取り除いても大丈夫か？
		if (HaveWithOutEdge(zone)) {
			filtered = append(filtered, []GameMapPosition{})
			filtered[i] = append(filtered[i], zone...)
			i++
		}
	}

	return filtered
}

//エッジを取り除いても面積があるか
func HaveWithOutEdge(zone []GameMapPosition) bool {
	areaEdge := 0
	areaAll := len(zone)
	for _, pos := range zone {
		if isEdgeInZone(zone, pos) {
			areaEdge++
		}
	}
	if (areaAll - areaEdge) > 0 {
		return true
	}
	return false
}

//いくつSwampを作るか？
func mountsForSwamp(zones [][]GameMapPosition, geo Geographical) int {
	zoneCounts := len(zones)
	mountsPercents := 0
	switch geo {
	case GeographicalStep:
		mountsPercents = lottery.GetRandomInt(0, 30)
	case GeographicalCave:
		mountsPercents = lottery.GetRandomInt(0, 50)
	case GeographicalSnow:
		mountsPercents = lottery.GetRandomInt(0, 30)
	case GeographicalPoison:
		mountsPercents = lottery.GetRandomInt(0, 70)
	case GeographicalFire:
		mountsPercents = lottery.GetRandomInt(0, 30)
	default:
	}
	mounts := int(float64(zoneCounts * mountsPercents) / 100.0)
	return mounts //TODO 数が少ない？
}

//勾配を生成
//ルール: x,yが大きいほど高い
func (xy *xymap) makeGradient(geo Geographical) {
	geta := 0
	rowestHigh := 0
	//険しさ(勾配の範囲)
	//	steepness := 0
	//まず地形でだいたいの高さ
	coefficient := 3
	switch geo {
	case GeographicalStep:
		base := lottery.GetRandomInt(2, 4)
		rowestHigh = base * coefficient
		break
	case GeographicalCave:
		rowestHigh = 10 * coefficient
		break
	case GeographicalRemain:
		base := lottery.GetRandomInt(2, 6)
		rowestHigh = base * coefficient
		break

	case GeographicalPoison:
		base := lottery.GetRandomInt(2, 5)
		rowestHigh = base * coefficient
		break
	case GeographicalFire:
		rowestHigh = 4 * coefficient
		break
	case GeographicalSnow:
		base := lottery.GetRandomInt(2, 4)
		rowestHigh = base * coefficient
		break

	case GeographicalJozen:
		rowestHigh = 2 * coefficient
		break
	case GeographicalCastle:
		rowestHigh = 2 * coefficient
		break
	}

	high := rowestHigh + geta

	//手前から道の高さを調整していく（歩けるように高さ調整していく）
	xy.makeGradientRoad(high)

	//だんだんと段差になっていく
	//パーリンノイズかける
	xy.makeGradientRough(high)

	//それ以外のところは2個あがったり下がったりさせる
	//ラフ
}

//手前から道の高さを調整していく（歩けるように高さ調整していく）
func (xy *xymap) makeGradientRoad(rowestHigh int) {
	currentHeight := rowestHigh
	for y := 0; y < xy.mapSize.MaxY; y++ {
		for x := 0; x < xy.mapSize.MaxX; x++ {
			switch xy.matrix[y][x] {
			case MacroMapTypeRoad:
				fallthrough
			case MacroMapTypeAllyPoint:
				fallthrough
			case MacroMapTypeEnemyPoint:
				xy.high[y][x] = currentHeight
			}
		}
	}
}

func (xy *xymap) makeGradientRough(rowestHigh int) {
	critHeight := rowestHigh
	coefficient := 0.1
	for y := 0; y < xy.mapSize.MaxY; y++ {
		for x := 0; x < xy.mapSize.MaxX; x++ {
			val := opensimplex.NewWithSeed(0).Eval2(float64(x) * coefficient, float64(y) * coefficient)
			floatHeight := float64(critHeight) * (val + 1.0) / 2.0
			height := int(floatHeight)
			if (height < 1) {
				height = 1
			}
			switch xy.matrix[y][x] {
			case MacroMapTypeRoad:
				fallthrough;
			case MacroMapTypeAllyPoint:
				fallthrough;
			case MacroMapTypeEnemyPoint:
				break;
			default:
				xy.matrix[y][x] = MacroMapTypeRough
				xy.high[y][x] = height
			}
		}
	}
}


//x,yを含むつながっている歩行可能領域を返す
//restrictedNumの数で打ち切る,0なら打ち切らない
func (xy *xymap) getNearHeightPanels(x int, y int, opened *[]GameMapPosition, restrictedNum int) *[]GameMapPosition {
	return xy.openNearPanel(x, y, opened, restrictedNum)
}

//上下左右のパネルをオープンし、歩行可能なら返す
//x, y, openedはすでにオープンした累積
//restrictedNumの数で打ち切る,0なら打ち切らない
func (xy *xymap) openNearPanel(x int, y int, opened *[]GameMapPosition, restrictedNum int) *[]GameMapPosition {
	currentHight := xy.high[y][x]
	DDDlog("open0 x:%d, y:%d, z:%d\n", x, y, currentHight);

	if (xy.shouldOpen(currentHight, x, y, *opened, restrictedNum)) {
		*opened = append(*opened, GameMapPosition{X:x, Y:y, Z:currentHight})
		DDDlog("opened1:%d\n", len(*opened));
	}
	if (xy.shouldOpen(currentHight, x + 1, y, *opened, restrictedNum)) {
		opens := xy.openNearPanel(x + 1, y, opened, restrictedNum)
		*opened = append(*opened, *opens...)
		trim(opened)
		DDDlog("opened5:%d\n", len(*opened));
	}
	if (xy.shouldOpen(currentHight, x, y + 1, *opened, restrictedNum)) {
		opens := xy.openNearPanel(x, y + 1, opened, restrictedNum)
		*opened = append(*opened, *opens...)
		trim(opened)
		DDDlog("opened3:%d\n", len(*opened));
	}
	if (xy.shouldOpen(currentHight, x, y - 1, *opened, restrictedNum)) {
		opens := xy.openNearPanel(x, y - 1, opened, restrictedNum)
		*opened = append(*opened, *opens...)
		DDDlog("opened20:%d\n", len(*opened));
		trim(opened)
		DDDlog("opened2 :%d\n", len(*opened));
	}
	if (xy.shouldOpen(currentHight, x - 1, y, *opened, restrictedNum)) {
		opens := xy.openNearPanel(x - 1, y, opened, restrictedNum)
		*opened = append(*opened, *opens...)
		trim(opened)
		DDDlog("opened4:%d\n", len(*opened));
	}
	DDDlog("close x:%d, y:%d, z:%d\n", x, y, currentHight);
	//	trim(opened);
	return opened
}

func trim(opened *[]GameMapPosition) {
	newOpenArray := []GameMapPosition{}
	newOpened := map[string]GameMapPosition{}
	for _, pos := range *opened {
		key := pos.X + pos.Y * 100 + pos.Z * 10000
		newOpened[strconv.Itoa(key)] = GameMapPosition{}
	}
	for key, _ := range newOpened {
		val, _ := strconv.Atoi(key)
		z := val / 10000
		y := (val % 10000) / 100
		x := val % 100
		newOpenArray = append(newOpenArray, GameMapPosition{X:x, Y:y, Z:z})
	}
	*opened = newOpenArray
	return;
}

//openすべきならtrue
func (xy *xymap) shouldOpen(currentHigh int, x int, y int, opened []GameMapPosition, restrictedNum int) bool {
	if (x >= xy.mapSize.MaxX || y >= xy.mapSize.MaxY || x < 0 || y < 0) {
		//マップ領域外
		return false;
	}
	if (xy.high[y][x] < (currentHigh - 2) || xy.high[y][x] > (currentHigh + 1)) {
		//高さ的にアウト
		return false;
	}
	for _, pos := range opened {
		if pos.X == x && pos.Y == y {
			//すでにオープンしていた
			DDDlog("alreadyed x:%d, y:%d\n", x, y);
			return false;
		}
	}
	if (restrictedNum != 0 && restrictedNum < len(opened)) {
		return false
	}
	return true;
}

//通れない/ハマり地形を正す
func (xy *xymap) validate() {

	retry:

	//ゾーニング
	zones := xy.zoningForValidate()
	for _, val := range zones {
		DDDlog("ahongoF:%d \n", len(val))
	}

	done := true
	//隣り合うゾーンとつなげる
	for idx, zone := range zones {
		for idx2, neighboughZone := range zones {
			if (idx == idx2) {
				continue
			}
			isNeighbough, zone1edge, zone2edge := xy.isNeighbough(zone, neighboughZone, false)
			if (!isNeighbough) {
				continue
			}
			//階段でつなげる
			done1 := xy.addStirsBetweenZones(zone, neighboughZone, zone1edge, zone2edge)
			if !done1 {
				done = false
			}
		}
	}
	if (!done) {
		goto retry;
	}
	DDDlog("unko50000")
}

//zone1,zone2が隣り合っていればtrue,
//zone1の隣り合っている座標を返す
//zone2の隣り合っている座標を返す
//isHighNear高さも考慮するならtrue
func (xy *xymap) isNeighbough(zone1 []GameMapPosition, zone2 []GameMapPosition, isHighNear bool) (bool,
GameMapPosition, GameMapPosition) {
	matrix := make([][]int, xy.mapSize.MaxY)
	for y := 0; y < xy.mapSize.MaxY; y++ {
		matrix[y] = make([]int, xy.mapSize.MaxX)
	}
	for x := 0; x < xy.mapSize.MaxX; x++ {
		for y := 0; y < xy.mapSize.MaxY; y++ {
			matrix[y][x] = 0
		}
	}
	for _, pos := range zone1 {
		matrix[pos.Y][pos.X] = 1
	}
	for _, pos := range zone2 {
		matrix[pos.Y][pos.X] = 2
	}
	beforeId := 0
	beforePos := GameMapPosition{}
	for x := 0; x < xy.mapSize.MaxX; x++ {
		beforeId = 0
		for y := 0; y < xy.mapSize.MaxY; y++ {
			if (isHighNear && !xy.isNearHigh(beforePos, GameMapPosition{X:x, Y:y})) {
				beforeId = matrix[y][x]
				beforePos = GameMapPosition{X:x, Y:y}
				continue
			}
			if (beforeId == 2 && matrix[y][x] == 1) {
				return true, GameMapPosition{X:x, Y:y}, beforePos
			}
			if (beforeId == 1 && matrix[y][x] == 2) {
				return true, beforePos, GameMapPosition{X:x, Y:y}
			}
			beforeId = matrix[y][x]
			beforePos = GameMapPosition{X:x, Y:y}
		}
	}
	beforeId = 0
	beforePos = GameMapPosition{}
	for y := 0; y < xy.mapSize.MaxY; y++ {
		beforeId = 0
		for x := 0; x < xy.mapSize.MaxX; x++ {
			if (isHighNear && !xy.isNearHigh(beforePos, GameMapPosition{X:x, Y:y})) {
				beforeId = matrix[y][x]
				beforePos = GameMapPosition{X:x, Y:y}
				continue
			}
			if (beforeId == 2 && matrix[y][x] == 1) {
				return true, GameMapPosition{X:x, Y:y}, beforePos
			}
			if (beforeId == 1 && matrix[y][x] == 2) {
				return true, beforePos, GameMapPosition{X:x, Y:y}
			}
			beforeId = matrix[y][x]
			beforePos = GameMapPosition{X:x, Y:y}
		}
	}
	return false, GameMapPosition{}, GameMapPosition{}
}

//2つのposが高さ的に差1以内ならtrue
func (xy *xymap) isNearHigh(pos1 GameMapPosition, pos2 GameMapPosition) bool {
	high1 := xy.high[pos1.Y][pos1.X]
	high2 := xy.high[pos2.Y][pos2.X]
	diff := high1 - high2
	if (diff > 1){
		return false
	}else if(diff < -1) {
		return false
	}
	return true
}

//zone1,zone2を階段でつなげる
//無事全てをつなぎ終えたらtrue
func (xy *xymap) addStirsBetweenZones(zone1 []GameMapPosition, zone2 []GameMapPosition,
zone1edge GameMapPosition, zone2edge GameMapPosition) bool {
	DDDlog("edge1:%+v edge2:%+v", zone1edge, zone2edge)

	beginingHeight := 0

	//低い土地を盛り上げる
	zone1high := xy.high[zone1edge.Y][zone1edge.X]
	zone2high := xy.high[zone2edge.Y][zone2edge.X]
	replaceForZone := []GameMapPosition{}
	replaceForEdge := GameMapPosition{}
	if (zone1high > zone2high) {
		replaceForZone = zone2
		replaceForEdge = zone2edge
		beginingHeight = zone1high
	}else {
		replaceForZone = zone1
		replaceForEdge = zone1edge
		beginingHeight = zone2high
	}
	//直線的にいく、その過程でぶつかれば曲げる

	//直線の幅は1-3のランダム(zone面積に応じて比例)
	pipeWidth := int(len(replaceForZone)) / 20 + 1
	if (pipeWidth > 3) {
		pipeWidth = 3
	}

	type Direction int
	const (
		DirectionLeft = 1 + iota
		DirectionRight
		DirectionUp
		DirectionDown
	)

	//方向
	var toDirection Direction
	if (zone1edge.X > zone2edge.X) {
		if (zone1high > zone2high) {
			toDirection = DirectionLeft
		}else {
			toDirection = DirectionRight
		}
	}else if (zone1edge.X < zone2edge.X) {
		if (zone1high > zone2high) {
			toDirection = DirectionRight
		}else {
			toDirection = DirectionLeft
		}
	}else if (zone1edge.Y > zone2edge.Y) {
		if (zone1high > zone2high) {
			toDirection = DirectionDown
		}else {
			toDirection = DirectionUp
		}
	}else if (zone1edge.Y < zone2edge.Y) {
		if (zone1high > zone2high) {
			toDirection = DirectionUp
		}else {
			toDirection = DirectionDown
		}
	}

	proceed := 0
	doneConnected := false
	for {
		nextX := 0
		nextY := 0
		switch toDirection {
		case DirectionLeft:
			nextX = replaceForEdge.X - proceed
			nextY = replaceForEdge.Y
		case DirectionRight:
			nextX = replaceForEdge.X + proceed
			nextY = replaceForEdge.Y
		case DirectionUp:
			nextX = replaceForEdge.X
			nextY = replaceForEdge.Y + proceed
		case DirectionDown:
			nextX = replaceForEdge.X + proceed
			nextY = replaceForEdge.Y
		}
		replaceToHeight := beginingHeight - proceed - 1
		if (!containsInZone(nextX, nextY, replaceForZone)) {
			//行き詰った
			break;
		}
		if (xy.high[nextY][nextX] == replaceToHeight) {
			//つなぎ終えた
			doneConnected = true
			break
		}
		if (replaceToHeight < 1) {
			//行き詰った
			break;
		}
		xy.high[nextY][nextX] = replaceToHeight
		if (doneConnected) {
			break
		}
		proceed++
	}

	return doneConnected
}

//zoneにx,yが含まれていればtrue
func containsInZone(x int, y int, zone []GameMapPosition) bool {
	for _, pos := range zone {
		if pos.X == x && pos.Y == y {
			return true
		}
	}
	return false
}

//各タイルを検証していって、ゾーニングする
func (xy *xymap) zoningForValidate() [][]GameMapPosition {
	zones := [][]GameMapPosition{}
	totalOpened := []GameMapPosition{}
	i := 0
	for y := 0; y < xy.mapSize.MaxY; y++ {
		for x := 0; x < xy.mapSize.MaxX; x++ {
			alreadyOpened := false
			for _, pos := range totalOpened {
				if pos.X == x && pos.Y == y {
					//すでにオープンしていた
					alreadyOpened = true
				}
			}
			if alreadyOpened {
				continue
			}
			opened := &[]GameMapPosition{}
			zone := xy.getNearHeightPanels(x, y, opened, 0)
			//			fmt.Printf("\nlen%d",len(zone))
			if len(*zone) == 0 {
				continue
			}
			totalOpened = append(totalOpened, *opened...)

			zones = append(zones, []GameMapPosition{})
			zones[i] = append(zones[i], *zone...)

			totalOpened = append(totalOpened, *zone...)
			i++
			DDDlog("ahongo:%d x:%d,y:%d,z%d, opened:%d\n", i, x, y, xy.high[y][x], len(*zone))
		}
	}
	return zones
}


//各タイルを検証していって、ゾーニングする
func (xy *xymap) zoningForMakeSwamp() [][]GameMapPosition {
	zones := [][]GameMapPosition{}
	totalOpened := []GameMapPosition{}
	i := 0
	for y := 0; y < xy.mapSize.MaxY; y++ {
		for x := 0; x < xy.mapSize.MaxX; x++ {
			alreadyOpened := false
			for _, pos := range totalOpened {
				if pos.X == x && pos.Y == y {
					//すでにオープンしていた
					alreadyOpened = true
				}
			}
			if alreadyOpened {
				continue
			}
			opened := &[]GameMapPosition{}
			zone := xy.getNearHeightPanelsSame(x, y, opened)
			//			fmt.Printf("\nlen%d",len(zone))
			if len(*zone) == 0 {
				continue
			}
			totalOpened = append(totalOpened, *opened...)

			zones = append(zones, []GameMapPosition{})
			zones[i] = append(zones[i], *zone...)

			totalOpened = append(totalOpened, *zone...)
			i++
			DDDlog("ahongo:%d x:%d,y:%d,z%d, opened:%d\n", i, x, y, xy.high[y][x], len(*zone))
		}
	}
	return zones
}

//A*はやくするためにエリア分けする
func (xy *xymap)zoningForAstar(restricted int) [][]GameMapPosition {
	zones := [][]GameMapPosition{}
	totalOpened := []GameMapPosition{}
	i := 0
	for y := 0; y < xy.mapSize.MaxY; y++ {
		for x := 0; x < xy.mapSize.MaxX; x++ {
			alreadyOpened := false
			for _, pos := range totalOpened {
				if pos.X == x && pos.Y == y {
					//すでにオープンしていた
					alreadyOpened = true
				}
			}
			if alreadyOpened {
				continue
			}
			opened := &[]GameMapPosition{}
			zone := xy.getNearHeightPanels(x, y, opened, restricted)
			//			fmt.Printf("\nlen%d",len(zone))

			if len(*zone) == 0 {
				continue
			}
			//totalOpened = append(totalOpened, *opened...)
			newZone := []GameMapPosition{}

			DDDlog("zone:%+v \n", *zone)

			shouldCreate := false
			for _, open := range *zone {
				shouldAppend := true
				for _, pos := range totalOpened {
					if pos.X == open.X && pos.Y == open.Y {
						shouldAppend = false
					}
				}
				if !shouldAppend {
					continue
				}
				shouldCreate = true
				newZone = append(newZone, open)
			}

			if (shouldCreate) {
				//1増やす
				zones = append(zones, []GameMapPosition{})
				zones[i] = append(zones[i], newZone...)
				/*
				for _, value := range newZone {
					xy.areaId[value.Y][value.X] = i
				}
				xy.maxAreaId = i
				*/
				i++
			}

			totalOpened = append(totalOpened, *zone...)
		}
	}
	return zones
}

//x,yを含むつながっている同じ高さの領域を返す
func (xy *xymap) getNearHeightPanelsSame(x int, y int, opened *[]GameMapPosition) *[]GameMapPosition {
	return xy.openNearPanelSame(x, y, opened)
}

//上下左右のパネルをオープンし、同一高さなら返す
//x, y, openedはすでにオープンした累積
func (xy *xymap) openNearPanelSame(x int, y int, opened *[]GameMapPosition) *[]GameMapPosition {
	currentHight := xy.high[y][x]
	DDDlog("open1 x:%d, y:%d, z:%d\n", x, y, currentHight);

	if (xy.shouldOpenToSame(currentHight, x, y, *opened)) {
		*opened = append(*opened, GameMapPosition{X:x, Y:y, Z:currentHight})
		DDDlog("opened1:%d\n", len(*opened));
	}
	if (xy.shouldOpenToSame(currentHight, x, y - 1, *opened)) {
		opens := xy.openNearPanelSame(x, y - 1, opened)
		*opened = append(*opened, *opens...)
		DDDlog("opened20:%d\n", len(*opened));
		trim(opened)
		DDDlog("opened2 :%d\n", len(*opened));
	}
	if (xy.shouldOpenToSame(currentHight, x, y + 1, *opened)) {
		opens := xy.openNearPanelSame(x, y + 1, opened)
		*opened = append(*opened, *opens...)
		trim(opened)
		DDDlog("opened3:%d\n", len(*opened));
	}
	if (xy.shouldOpenToSame(currentHight, x - 1, y, *opened)) {
		opens := xy.openNearPanelSame(x - 1, y, opened)
		*opened = append(*opened, *opens...)
		trim(opened)
		DDDlog("opened4:%d\n", len(*opened));
	}
	if (xy.shouldOpenToSame(currentHight, x + 1, y, *opened)) {
		opens := xy.openNearPanelSame(x + 1, y, opened)
		*opened = append(*opened, *opens...)
		trim(opened)
		DDDlog("opened5:%d\n", len(*opened));
	}
	DDDlog("close x:%d, y:%d, z:%d\n", x, y, currentHight);
	//	trim(opened);
	return opened
}

//openすべきならtrue
func (xy *xymap) shouldOpenToSame(currentHigh int, x int, y int, opened []GameMapPosition) bool {
	if (x >= xy.mapSize.MaxX || y >= xy.mapSize.MaxY || x < 0 || y < 0) {
		//マップ領域外
		return false;
	}
	if (xy.high[y][x] != currentHigh) {
		//高さ的にアウト
		return false;
	}
	for _, pos := range opened {
		if pos.X == x && pos.Y == y {
			//すでにオープンしていた
			DDDlog("alreadyed x:%d, y:%d\n", x, y);
			return false;
		}
	}
	return true;
}


//zonesで離れた地点を分離
func (xy *xymap)validateForZone(game_map *GameMap, zones [][]GameMapPosition) [][]GameMapPosition {

	newZones := [][]GameMapPosition{}

	//ゾーニング
	for _, zone := range zones {
		zones := zoningForValidateForAstar(game_map, zone)
		newZones = append(newZones, zones...)
	}
	return newZones
}

//小さいzone(restricted未満のサイズ)を他のzoneにくっつける
//小さいzoneが存在しなくなったらtrueを返す
func (xy *xymap)roundZones(zonesSrc [][]GameMapPosition, restricted int) ([][]GameMapPosition, bool) {
	newZones := [][]GameMapPosition{}
	alreadyAddedIdxs := []int{}
	i := 0
	done := true

	for idx, zone := range zonesSrc {
		if (contains(alreadyAddedIdxs, idx)) {
			continue
		}
		if (len(zone) >= (restricted)) {
			//十分な大きさがある
			newZones = append(newZones, []GameMapPosition{})
			newZones[i] = append(newZones[i], zone...)
			alreadyAddedIdxs = append(alreadyAddedIdxs, idx)
			i++
			continue
		}
		tmpZone := []GameMapPosition{}
		tmpZone = append(tmpZone, zone...)
		for idx2, zone2 := range zonesSrc {
			if (idx == idx2) {
				continue
			}
			if (contains(alreadyAddedIdxs, idx2)) {
				continue
			}
			if (len(zone2) >= (restricted)) {
				//十分な大きさがある
				continue
			}
			//隣接、高さチェック
			isNeighbough, _, _ := xy.isNeighbough(zone, zone2, true)
			if !isNeighbough {
				continue
			}
			//隣り合っており、高さ的にも結合して問題ない
			tmpZone = append(tmpZone, zone2...)
			alreadyAddedIdxs = append(alreadyAddedIdxs, idx)
			alreadyAddedIdxs = append(alreadyAddedIdxs, idx2)
		}
		newZones = append(newZones, []GameMapPosition{})
		newZones[i] = append(newZones[i], tmpZone...)
		i++
	}

	for _, zoneN := range newZones {
		if len(zoneN) < restricted {
			done = false
		}
	}
	return newZones, done
}

func remove(numbers []GameMapPosition, search GameMapPosition) []GameMapPosition {
	result := []GameMapPosition{}
	for _, num := range numbers {
		if num != search {
			result = append(result, num)
		}
	}
	return result
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

//ゾーンをグラフ化
//グラフの辺はゲーム内で計算(A*でずれをださないため)
func (xy *xymap)makeGraphForAstar(zones [][]GameMapPosition) {
	line := make([][]int, xy.maxAreaId)
	for i := 0; i < xy.maxAreaId; i++ {
		//		line[i] = make([]int, 0)
		line[i] = []int{}
	}

	//隣り合うゾーンとつなげる
	for idx, zone := range zones {
		Dlog("zone %d:%+v\n", idx, zone)
		areaId := xy.areaId[zone[0].Y][zone[0].X]
		for idx2, neighboughZone := range zones {
			if (idx == idx2) {
				continue
			}
			isNeighbough, _, _ := xy.isNeighbough(zone, neighboughZone, true)
			if (!isNeighbough) {
				continue
			}
			//隣り合っている
			neighboughAreaId := xy.areaId[neighboughZone[0].Y][neighboughZone[0].X]
			//すでに追加されている
			already := false
			for _, val := range line[areaId] {
				if val == neighboughAreaId {
					already = true
				}
			}
			if already {
				continue;
			}
			line[areaId] = append(line[areaId], neighboughAreaId)
		}
	}

	for i := 0; i < xy.maxAreaId; i++ {
		Dlog("Lines areaId:%d : %+v\n", i, line[i])
	}
	xy.areaPath = line
}



//引数のゾーンをさらにゾーニングする
func zoningForValidateForAstar(game_map *GameMap, input []GameMapPosition) [][]GameMapPosition {
	zones := [][]GameMapPosition{}
	totalOpened := []GameMapPosition{}
	xymap := NewXYMap(game_map.Size)
	i := 0

	for _, va := range input {
		xymap.high[va.Y][va.X] = 10
	}

	for _, val := range input {
		x := val.X
		y := val.Y
		alreadyOpened := false
		for _, pos := range totalOpened {
			if pos.X == x && pos.Y == y {
				//すでにオープンしていた
				alreadyOpened = true
			}
		}
		if alreadyOpened {
			continue
		}
		opened := &[]GameMapPosition{}
		zone := xymap.getNearHeightPanels(x, y, opened, 0)
		//			fmt.Printf("\nlen%d",len(zone))
		if len(*zone) == 0 {
			continue
		}
		totalOpened = append(totalOpened, *opened...)

		zones = append(zones, []GameMapPosition{})
		zones[i] = append(zones[i], *zone...)

		totalOpened = append(totalOpened, *zone...)
		i++
		//		DDDlog("ahongo:%d x:%d,y:%d,z%d, opened:%d\n", i, x, y, xy.high[y][x], len(*zone))
	}
	return zones
}

