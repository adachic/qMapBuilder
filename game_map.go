package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/adachic/lottery"
	"image"
	"image/png"
	"os"
	"github.com/satori/go.uuid"
	"image/color"
	"encoding/json"
	"sort"
)

//マップの大きさ
type GameMapSize struct {
	MaxX int
	MaxY int
	MaxZ int
}

//面積
func (s GameMapSize) area() int {
	return s.MaxY * s.MaxX
}

//マップの難度
type Difficult int
const (
	easy   Difficult = 10
	normal Difficult = 3
	hard   Difficult = 2
	exhard Difficult = 1
)

//長方形の形式
type RectForm int
const (
	horizontalLong RectForm = 6 //横長
	verticalLong   RectForm = 5
	square         RectForm = 2
)

//地形
type Geographical int
const (
	GeographicalStep     Geographical = 14 + 10
	GeographicalCave     Geographical = 8 + 10
	GeographicalRemain     Geographical = 7 + 10

	GeographicalPoison  Geographical = 6 + 10
	GeographicalFire     Geographical = 9 + 10
	GeographicalSnow     Geographical = 5 + 10

	GeographicalJozen  Geographical = 3 + 10
	GeographicalCastle   Geographical = 4 + 10
)

//マップメタ(ここから詳細なパーツを決定)
type MacroMapType int
const (
	MacroMapTypeRoad = 1 + iota
	MacroMapTypeRough
	MacroMapTypeWall
	MacroMapTypeCantEnter //進入不可地形 TODO:これつかってんだっけ？

	MacroMapTypeSwampWater //水系
	MacroMapTypeSwampRava
	MacroMapTypeSwampPoison
	MacroMapTypeSwampHeal

	MacroMapTypeOther //他
	MacroMapTypeAllyPoint
	MacroMapTypeEnemyPoint
)

//alreadyに登録されてなくて自分を除く最も近いポイントを返す
//errがtrueなら見つからなかった
func (pos GameMapPosition) searchNearPositionWithOutMe(positions []GameMapPosition, alreadys []PathPosition) (nearPos GameMapPosition, err bool) {
	err = true;
	minDistance := 10000
	for i := 0; i < len(positions); i++ {
		tgtPos := positions[i]
		if (pos.equalXYTo(tgtPos)) {
			continue
		}
		if (containsPath(alreadys, PathPosition{pos, tgtPos})) {
			continue
		}
		distance := pos.distance(tgtPos)
		if (distance < minDistance) {
			minDistance = distance
			nearPos = tgtPos
			err = false
		}
	}
	return nearPos, err
}

func (pos GameMapPosition) equalXYTo(another GameMapPosition) bool {
	return (pos.X == another.X) && (pos.Y == another.Y)
}

func (pos GameMapPosition) distance(another GameMapPosition) int {
	absx := pos.X - another.X
	absy := pos.Y - another.Y
	if (absx < 0) {
		absx *= -1
	}
	if (absy < 0) {
		absy *= -1
	}
	return absx + absy
}

//座標
type GameMapPosition struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
}

//確率を返す
func (d Difficult) Prob() int {
	return int(d)
}

func (d Category) Prob() int {
	return int(d)
}

func (d RectForm) Prob() int {
	return int(d)
}

func (d Geographical) Prob() int {
	return int(d)
}

func NewGameMap(condition GameMapCondition) *GameMap {
	game_map := &GameMap{}
	return game_map.init(condition)
}

func (game_map *GameMap) init(condition GameMapCondition) *GameMap {
	//難易度を初期化
	game_map.initMapDifficult()
	fmt.Printf("difficult: %+v\n", game_map.Difficult)

	//カテゴリを初期化
	game_map.initMapCategory()
	fmt.Printf("category: %+v\n", game_map.Category)

	//道の舗装度を初期化
	game_map.initMapPavement()

	//マップのサイズを初期化
	game_map.initMapSize()
	fmt.Printf("mapSize: %+v\n", game_map.Size)

	//大まかな地形を初期化
	game_map.initMapGeographical()

	//味方開始ポイントを初期化
	game_map.initAllyStartPoint()
	fmt.Printf("allyStartPoint: %+v\n", game_map.AllyStartPoint)

	//敵開始ポイントを決定
	game_map.initEnemyStartPoints()
	for _, enemyStartPoint := range game_map.EnemyStartPoints { // キーは使われません
		fmt.Printf("enemyStartPoint: %+v\n", enemyStartPoint)
	}

	//2次元マップ生成
	{
		xymap := NewXYMap(game_map.Size)

		//広場配置
		xymap.putPlazas(game_map.Difficult, game_map.AllyStartPoint, game_map.EnemyStartPoints)

		//道配置
		xymap.putRoads(game_map.Difficult, game_map.AllyStartPoint, game_map.EnemyStartPoints)

		//壁配置

		//ラフ配置

		//味方、敵ポイント
		xymap.putPoint(game_map.AllyStartPoint, MacroMapTypeAllyPoint)
		for i := 0; i < len(game_map.EnemyStartPoints); i++ {
			xymap.putPoint(game_map.EnemyStartPoints[i], MacroMapTypeEnemyPoint)
		}

		//勾配を生成
		xymap.makeGradient(game_map.Geographical)

		//バリデーション
		xymap.validate()

		//水、毒沼配置
		xymap.makeSwamp(game_map.Geographical)

		//alloc/init
		game_map.allocToJungleGym()

		//xymap情報をコピー
		game_map.copyFromXY(xymap)

		//dump
		xymap.printMapForDebug()
	}
	return game_map
}


//マップ難易度の抽選結果を返す
func (game_map *GameMap) initMapDifficult() {
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
	game_map.Difficult = difficults[result].(Difficult)
}

//雪原かどうかの抽選結果を返す
func (game_map *GameMap) initMapSnow() {
	dice := lottery.GetRandomInt(0, 100);
	game_map.IsSnow = dice < 2
}

//舗装度の抽選結果を返す
func (game_map *GameMap) initMapPavement() {
	dice := lottery.GetRandomInt(0, 100);
	switch{
	case dice < 20:
		game_map.PavementLv = 1;
	case dice < 40:
		game_map.PavementLv = 2;
	case dice < 60:
		game_map.PavementLv = 3;
	case dice < 80:
		game_map.PavementLv = 4;
	case dice < 100:
		game_map.PavementLv = 5;
	}
}

//マップカテゴリの抽選結果を返す
func (game_map *GameMap) initMapCategory() {
	lot := lottery.New(rand.New(rand.NewSource(time.Now().UnixNano())))
	categories := []lottery.Interface{
		CategoryStep,
		CategoryCave,
		CategoryRemains,

		CategoryFire,
		CategoryPoison,
		CategorySnow,

		CategoryJozen,
		CategoryCastle,
	}
	result := lot.Lots(
		categories...,
	)
	game_map.Category = categories[result].(Category)

	switch(game_map.Category){
	case CategoryPoison:
		fallthrough
	case CategoryFire:
		fallthrough
	case CategorySnow:
		fallthrough
	case CategoryJozen:
		game_map.Category = CategoryStep
	}
}


//マップサイズの抽選結果を返す
func (game_map *GameMap) initMapSize() {
	//横長か縦長か
	rectForm := CreateRectForm()
	//x/yアスペクト比
	aspect := CreateAspectOfRectFrom(rectForm)
	//面積
	area := CreateArea(game_map.Difficult)

	fmt.Printf("form: %+v\n", rectForm)
	fmt.Printf("area: %+v\n", area)
	fmt.Printf("aspect: %f\n", aspect)
	yy := float32(area) / float32(aspect)
	fmt.Printf("yy: %f\n", yy)
	y := int(math.Sqrt(float64(yy)))
	if y < 1 {
		y = 1
	}
	x := int(area / y)

	game_map.Size = GameMapSize{x, y, 30}
}

//地形の抽選結果を返す
func (game_map *GameMap) initMapGeographical() {
	//カテゴリに応じて、地形を対応させる
	var geo Geographical
	switch game_map.Category {
	case CategoryStep:
		dice := lottery.GetRandomInt(0, 3)
		switch dice{
		case 0:
			geo = GeographicalStep
		case 1:
			geo = GeographicalPoison
		case 2:
			geo = GeographicalSnow
		case 3:
			geo = GeographicalJozen
		}
	case CategoryFire:
		geo = GeographicalFire
	case CategoryCave:
		geo = GeographicalCave
	case CategoryRemains:
		geo = GeographicalRemain
	case CategoryJozen:
		geo = GeographicalJozen
	case CategorySnow:
		geo = GeographicalSnow
	case CategoryCastle:
		geo = GeographicalCastle
	}

	/*
	GeographicalPoison
	GeographicalSnow
	GeographicalJozen
	GeographicalFire
	*/
	game_map.Geographical = geo
}

//味方の出撃座標を返す
func (game_map *GameMap) initAllyStartPoint() {
	var seed Range
	//難度が高いと中央寄りになる
	switch game_map.Difficult{
	case easy:
		seed = Range{70, 100}
		break
	case normal:
		seed = Range{30, 80}
		break
	case hard:
		seed = Range{0, 50}
		break
	case exhard:
		seed = Range{0, 20}
		break
	}
	fmt.Print("seed", seed.Min, " ", seed.Max)
	distanceFrom := lottery.GetRandomInt(seed.Min, seed.Max)
	game_map.AllyStartPoint, _ = CreateRandomPositionInMap(
		game_map.Size,
		GameMapPosition{game_map.Size.MaxX / 2, game_map.Size.MaxY / 2, 0},
		distanceFrom)
}

//敵出現座標の一覧を返す
func (game_map *GameMap) initEnemyStartPoints() {
	type rangeFromAlly struct {
		Min int
		Max int
	}
	//味方と敵との大体の距離感
	var rangeFrom rangeFromAlly
	//敵出撃座標の数
	var sattyPointNum int

	switch game_map.Difficult{
	case easy:
		sattyPointNum = lottery.GetRandomInt(1, 3)
		rangeFrom = rangeFromAlly{50, 100}
		break
	case normal:
		sattyPointNum = lottery.GetRandomInt(3, 6)
		rangeFrom = rangeFromAlly{30, 100}
		break
	case hard:
		sattyPointNum = lottery.GetRandomInt(5, 10)
		rangeFrom = rangeFromAlly{13, 100}
		break
	case exhard:
		sattyPointNum = lottery.GetRandomInt(10, 20)
		rangeFrom = rangeFromAlly{10, 100}
		break
	}
	var sattyPoints []GameMapPosition
	var radians []float64
	for sattyPointNum > 0 {
		sattyPointNum--
		//味方ポイントからの距離
		distance := lottery.GetRandomInt(rangeFrom.Min, rangeFrom.Max)
		sattyPoint, radian := CreateRandomPositionInMap(game_map.Size, game_map.AllyStartPoint, distance)
		sattyPoints = append(sattyPoints, sattyPoint)
		radians = append(radians, radian)
	}
	//角度をソートする
	sort.Float64s(radians)

	//エッジ分を追加
	var sattyPointsEdge []GameMapPosition
	sattyPointsEdge = CreateEdgePositionInMap(game_map.Size, game_map.AllyStartPoint, radians)
	sattyPoints = append(sattyPoints, sattyPointsEdge...)

	game_map.EnemyStartPoints = sattyPoints
}

func (game_map *GameMap)fillJungleGymToEmpty() {
	for z := 0; z <= game_map.Size.MaxZ; z++ {
		for y := 0; y < game_map.Size.MaxY; y++ {
			for x := 0; x < game_map.Size.MaxX; x++ {
				game_map.JungleGym[z][y][x] = GameParts{IsEmpty:true};
			}
		}
	}
}

func (game_map *GameMap) allocToJungleGym() {
	game_map.JungleGym = make([][][]GameParts, game_map.Size.MaxZ + 1)
	for z := 0; z <= game_map.Size.MaxZ; z++ {
		game_map.JungleGym[z] = make([][]GameParts, game_map.Size.MaxY)
		for y := 0; y < game_map.Size.MaxY; y++ {
			game_map.JungleGym[z][y] = make([]GameParts, game_map.Size.MaxX)
		}
	}
	game_map.fillJungleGymToEmpty();

	game_map.MacroMapTypes = make([][][]MacroMapType, game_map.Size.MaxZ + 1)
	for z := 0; z <= game_map.Size.MaxZ; z++ {
		game_map.MacroMapTypes[z] = make([][]MacroMapType, game_map.Size.MaxY)
		for y := 0; y < game_map.Size.MaxY; y++ {
			game_map.MacroMapTypes[z][y] = make([]MacroMapType, game_map.Size.MaxX)
		}
	}

	game_map.High = make([][]int, game_map.Size.MaxY)
	for y := 0; y < game_map.Size.MaxY; y++ {
		game_map.High[y] = make([]int, game_map.Size.MaxX)
	}
}

//ここでxyをjungleGymへ移行
func (game_map *GameMap) copyFromXY(xy *xymap) {
	for x := 0; x < xy.mapSize.MaxX; x++ {
		for y := 0; y < xy.mapSize.MaxY; y++ {
			macro := xy.getMatrix(x, y);
			high := xy.getHigh(x, y);
			game_map.High[y][x] = high;
			for z := 0; z < high; z++ {
				game_map.MacroMapTypes[z][y][x] = macro;
			}
		}
	}
}

//パーツとのひも付け
//失敗ならfalse
func (game_map *GameMap) bindToGameParts(gamePartsDict map[string]GameParts) bool {
	/*選定パーツのゾーニング*/
	//1.主幹パーツの決定:道・ラフ・その他

	idsRoadFull := GetIdsRoad(game_map, gamePartsDict, false)
	idsRoughFull := GetIdsRough(game_map, gamePartsDict, false)
	idsWallFull := GetIdsWall(game_map, gamePartsDict, false)

	idsRoadHalf := GetIdsRoad(game_map, gamePartsDict, true)
	idsRoughHalf := GetIdsRough(game_map, gamePartsDict, true)
	idsWallHalf := GetIdsWall(game_map, gamePartsDict, true)

	if (len(idsRoadHalf) == 0 || len(idsRoughHalf) == 0 || len(idsWallHalf) == 0) {
		if (len(idsRoadHalf) == 0) {
			if (len(idsRoughHalf) != 0) {
				idsRoadHalf = idsRoughHalf
			}
			if (len(idsWallHalf) != 0) {
				idsRoadHalf = idsWallHalf
			}
		}
		if (len(idsRoughHalf) == 0) {
			if (len(idsRoadHalf) != 0) {
				idsRoughHalf = idsRoadHalf
			}
			if (len(idsWallHalf) != 0) {
				idsRoughHalf = idsWallHalf
			}
		}
		if (len(idsWallHalf) == 0) {
			if (len(idsRoadHalf) != 0) {
				idsWallHalf = idsRoughHalf
			}
			if (len(idsRoughHalf) != 0) {
				idsWallHalf = idsRoughHalf
			}
		}
		if (len(idsRoadHalf) == 0 || len(idsRoughHalf) == 0 || len(idsWallHalf) == 0) {
			return false
		}
	}

	//2.パーツ割当
	for x := 0; x < game_map.Size.MaxX; x++ {
		for y := 0; y < game_map.Size.MaxY; y++ {
			high := game_map.High[y][x];
			for z := 0; z < high; z++ {

				shouldHalf := ((high - 1) == z && (high % 2) > 0 )

				/*
				//half段目か？であれば、halfとし、それ以外は非half
				if((high - 1) == z && high%2 ){
					//half確定
				}else{
				}
				*/

				macro := game_map.MacroMapTypes[z][y][x]
				//1.土
				if (z < high - 1) {
					parts := GetGamePartsFoundation(idsWallFull, idsRoughFull, idsRoadFull, gamePartsDict, x, y, z);
					fmt.Printf("found  : %2d,%2d,%2d id:%s \n", z, y, x, parts.Id)
					if (shouldHalf) {
						//halfにコンバートする
						before := parts.Id
						parts = GetHalfParts(idsRoadHalf, idsRoughHalf, idsWallHalf, parts, gamePartsDict, macro, x, y, z)
						fmt.Printf("converted from %s to %s \n", before, parts.Id)
					}
					game_map.JungleGym[z][y][x] = parts;
					continue;
				}
				//2.表層(道,ラフ,壁)
				parts := GetGamePartsSurface(idsWallFull, idsRoughFull, idsRoadFull, gamePartsDict, macro, x, y, z);
				fmt.Printf("surface: %2d,%2d,%2d id:%s macro[%v]\n", z, y, x, parts.Id, macro)
				if (shouldHalf) {
					//halfにコンバートする
					before := parts.Id
					parts = GetHalfParts(idsRoadHalf, idsRoughHalf, idsWallHalf, parts, gamePartsDict, macro, x, y, z)
					fmt.Printf("converted from %s to %s \n", before, parts.Id)
				}
				game_map.JungleGym[z][y][x] = parts;
			}
		}
	}
	return true
}

//tileのimageを得る
func imageTile32(tile Tile) *image.RGBA {
	file, err := os.Open("./assets/" + tile.FilePath)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	img, err := png.Decode(file)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	//	tileRect := image.Rect(tile.X, tile.Y, tile.Width, tile.Height)
	yy := img.Bounds().Max.Y - tile.Y - tile.Height;

	outputRect := image.Rect(0, 0, tile.Width, tile.Height)
	outputImg := image.NewRGBA(outputRect)
	for x := 0; x < outputRect.Max.X; x++ {
		for y := 0; y < outputRect.Max.Y; y++ {
			outputImg.Set(x, y, img.At(tile.X + x, yy + y))
		}
	}
	return outputImg
}

//tileのimageを得る
func imageTile64(tile Tile) *image.RGBA {
	file, err := os.Open("./assets/" + tile.FilePath)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	img, err := png.Decode(file)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	//	tileRect := image.Rect(tile.X, tile.Y, tile.Width, tile.Height)
	yy := img.Bounds().Max.Y - tile.Y - tile.Height;

	outputRect := image.Rect(0, 0, tile.Width / 2, tile.Height / 2)
	outputImg := image.NewRGBA(outputRect)
	for x := 0; x < outputRect.Max.X; x++ {
		for y := 0; y < outputRect.Max.Y; y++ {
			outputImg.Set(x, y, img.At(tile.X + x * 2, yy + y * 2))
		}
	}
	return outputImg
}

//srcImgのsrcRectを dstImgのdstRectにコピ-
func clipAfromB(srcImg *image.RGBA, srcRect image.Rectangle, dstImg *image.RGBA, dstRect image.Rectangle) {
	for x := 0; x < srcRect.Max.X; x++ {
		for y := 0; y < srcRect.Max.Y; y++ {
			srcRGBA := srcImg.At(x, y)
			_, _, _, a := srcRGBA.RGBA()
			if (a == 0) {
				//0は透明
				continue
			}
			dstImg.Set(dstRect.Min.X + x, dstRect.Min.Y + y, srcRGBA)
		}
	}
}

//x,y,zの起点座標の計算
func (game_map *GameMap) targetDrawPoint(x int, y int, z int) image.Point {
	xx := (x) * 16 + (y) * 16
	yy := (game_map.Size.MaxZ - z) * 16 / 2 +
	(game_map.Size.MaxY - y) * 8 +
	x * 8;
	//	yy := (8 * x * -1) + (8 * y) + (16 * z)
	return image.Point{xx, yy}
}

/*
	CGFloat xOrigin = self.aspectX / 2.0f * matrix.x +
	self.aspectX / 2.0f * matrix.y;
	CGFloat yOrigin =
	self.aspectY / 2.0f * matrix.x * -1.0f +
	self.aspectY / 2.0f * matrix.y +
	self.aspectT * matrix.z;
	CGFloat yAid = [self aidOfZ0Position];
	return CGPointMake(xOrigin, yOrigin + yAid);
*/

func (game_map *GameMap) updateMaxZ() {
	maxHigh := 0
	for y := 0; y < game_map.Size.MaxY; y++ {
		for x := 0; x < game_map.Size.MaxX; x++ {
			if (maxHigh < game_map.High[y][x]) {
				maxHigh = game_map.High[y][x]
			}
		}
	}
	game_map.Size.MaxZ = maxHigh
}

//png生成
func (game_map *GameMap) createPng(gamePartsDict map[string]GameParts) {
	game_map.updateMaxZ();

	outputWidth := 16 + game_map.Size.MaxX * 16 + (game_map.Size.MaxY - 1) * 16
	outputHeight := 16 + game_map.Size.MaxX * 8 + game_map.Size.MaxY * 8 + game_map.Size.MaxZ * 16
	outputRect := image.Rect(0, 0, outputWidth, outputHeight)

	fmt.Printf("aho1:%+v\n", outputRect)

	//出力するイメージ
	outputImg := image.NewRGBA(outputRect)

	//塗りつぶす
	for x := 0; x < outputRect.Max.X; x++ {
		for y := 0; y < outputRect.Max.Y; y++ {
			//fmt.Printf("aho111:%d,%d/,%d,%d\n",x,y,outputRect.Max.X,outputRect.Max.Y)
			outputImg.Set(x, y, color.Black)
		}
	}

	fmt.Printf("aho2:rendering...\n")
	for z := 0; z <= game_map.Size.MaxZ; z++ {
		for y := (game_map.Size.MaxY - 1); y >= 0; y-- {
			for x := 0; x < game_map.Size.MaxX; x++ {
				cube := game_map.JungleGym[z][y][x]
				if (cube.IsEmpty) {
					//fmt.Printf("gomi:%d,%d,%d\n", z, y, x)
					continue
				}
				if (z % 2 > 0) {
					//レンダリングしなくて良い
					continue
				}
				if (game_map.shouldLightening(x, y, z)) {
					//肉抜き
					continue
				}
				if (!cube.Harf) {
					cube = game_map.JungleGym[z + 1][y][x]
				}
				//fmt.Printf("unko: %d,%d,%d\n", z, y, x)
				//fmt.Printf("cube: %+v\n", cube)

				//切り出す
				tile := cube.Tiles[0]
				var tileImage *image.RGBA
				if (cube.RezoType == RezoTypeRect64) {
					tileImage = imageTile64(tile)
				}else {
					tileImage = imageTile32(tile)
				}
				tileRect := image.Rect(0, 0, tile.Width, tile.Height)

				//バッファへ貼り付け
				dstPoint := game_map.targetDrawPoint(x, y, z)
				//dstPoint := targetDrawPoint(x, y, z)
				dstRect := image.Rect(dstPoint.X, dstPoint.Y, outputWidth, outputHeight)
				clipAfromB(tileImage, tileRect, outputImg, dstRect)
			}
		}
	}
	fmt.Printf("aho3:drawed\n")

	//ファイル出力
	fileName := uuid.NewV4().String()
	game_map.Filename = fileName
	file, err := os.Create("./output/" + fileName + ".png")
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = png.Encode(file, outputImg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("aho4:png outputed\n")
}

//マップ
type GameMap struct {
	JungleGym        [][][]GameParts
	MacroMapTypes    [][][]MacroMapType
	High             [][]int

	Size             GameMapSize

	AllyStartPoint   GameMapPosition
	EnemyStartPoints []GameMapPosition

	Difficult        Difficult
	Geographical     Geographical
	Category         Category
	IsSnow           bool
	PavementLv       int

	Filename         string
}

type JsonPanel struct {
	X  int `json:"x"`
	Y  int `json:"y"`
	Z  int `json:"z"`
	Id string `json:"id"`
}

//マップ
type JsonGameMap struct {
	MaxX             int `json:"maxX"`
	MaxY             int `json:"maxY"`
	MaxZ             int `json:"maxZ"`
	AspectX          int `json:"aspectX"`
	AspectY          int `json:"aspectY"`
	AspectT          int `json:"aspectT"`
	JungleGym        []JsonPanel `json:"jungleGym"`
	GameParts        []GameParts `json:"gameParts"`

	AllyStartPoint   GameMapPosition `json:"allyStartPoint"`
	EnemyStartPoints []GameMapPosition `json:"enemyStartPoints"`
	Category         Category `json:"category"`
}

//Json生成
func (game_map *GameMap) createJson(gamePartsDict map[string]GameParts) {
	fmt.Printf("==output json==\n")

	game_map.AllyStartPoint.Z = game_map.High[game_map.AllyStartPoint.Y][game_map.AllyStartPoint.X] / 2 // - 1
	for i, enemyStartPoint := range game_map.EnemyStartPoints { // キーは使われません
		game_map.EnemyStartPoints[i].Z = game_map.High[game_map.EnemyStartPoints[i].Y][game_map.EnemyStartPoints[i].X] / 2  //- 1
		fmt.Printf("enemyStartPoint: %+v\n", enemyStartPoint)
	}

	jsonStub := JsonGameMap{
		MaxX:game_map.Size.MaxX,
		MaxY:game_map.Size.MaxY,
		MaxZ:game_map.Size.MaxZ,
		AspectX:32,
		AspectY:16,
		AspectT:16,
		AllyStartPoint:game_map.AllyStartPoint,
		EnemyStartPoints:game_map.EnemyStartPoints,
		Category:game_map.Category,
	}

	var flags map[string]GameParts
	for z := 0; z < game_map.Size.MaxZ; z++ {
		for y := 0; y < game_map.Size.MaxY; y++ {
			for x := 0; x < game_map.Size.MaxX; x++ {
				cube := game_map.JungleGym[z][y][x]
				if (cube.IsEmpty) {
					continue
				}
				if (z % 2 > 0) {
					//レンダリングしなくて良い
					continue
				}
				if (!cube.Harf) {
					cube = game_map.JungleGym[z + 1][y][x]
				}
				if (game_map.shouldLightening(x, y, z)) {
					//肉抜き
					continue
				}
				jsonStub.JungleGym = append(jsonStub.JungleGym, JsonPanel{X:x, Y:y, Z:z / 2, Id:cube.Id})
				_, ok := flags[cube.Id]
				if (!ok) {
					jsonStub.GameParts = append(jsonStub.GameParts, cube)
				}
			}
		}
	}

	//	fmt.Printf("%+v\n", jsonStub)

	bytes, json_err := json.Marshal(jsonStub)
	if json_err != nil {
		fmt.Println("Json Encode Error: ", json_err)
	}

	//	fmt.Printf("bytes:%+v\n", string(bytes))

	file, err := os.Create("./output/" + game_map.Filename + ".json")
	_, err = file.Write(bytes)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
}

//該当パネルが肉抜きできるならtrue
func (game_map *GameMap) shouldLightening(x int, y int, z int) bool {
	//上のパネルがあるz+1
	//上のパネルがあるz+2
	//手前のパネルがあって、harfではないx+1,y-1
	if (y - 1 < 0) {
		return false;
	}
	if (x + 1 >= game_map.Size.MaxX) {
		return false;
	}
	if (z + 1 >= game_map.Size.MaxZ) {
		return false;
	}
	if (z + 2 >= game_map.Size.MaxZ) {
		return false;
	}
	cube := game_map.JungleGym[z + 1][y][x]
	cube2 := game_map.JungleGym[z][y - 1][x]
	cube3 := game_map.JungleGym[z][y][x + 1]
	cube4 := game_map.JungleGym[z + 2][y][x]
	if (cube.IsEmpty || cube2.IsEmpty || cube3.IsEmpty || cube4.IsEmpty) {
		return false
	}
	if (cube2.Harf || cube3.Harf) {
		return false
	}
	return true
}


//地形に適さないパーツ判定、除外ならtrue
func (game_map *GameMap)isExcludedByGeography(parts GameParts) bool {
	switch game_map.Geographical {
	case GeographicalStep:
		if parts.Snow > 0 {
			return true
		}
	case GeographicalSnow:
		if parts.Snow == 0 {
			return true
		}
	}

	return false
}

