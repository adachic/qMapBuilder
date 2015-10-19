package MapBuilder

import (
	"github.com/adachic/lottery"
	"time"
"math/rand"
)

type GameMap struct {
	MaxX int
	MaxY int
	MaxZ int
	JungleGym [][][]GameParts
}

type Difficult int
const (
	easy Difficult = 10
	normal Difficult  = 3
	hard Difficult = 1
	exhard Difficult = 1
)

func (d Difficult) Prob() int {
	return int(d)
}

func decisionDifficult() Difficult {
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

func CreateGameMap(gamePartsDict map[string]GameParts) GameMap{
	//難易度を決定
	decisionDifficult()

	//マップのサイズを決定

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

func createJsonFromMap(){

}

func createPngFromMap(){

}


