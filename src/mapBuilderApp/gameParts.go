package main

type Tile struct {

}

type WaterType struct {

}

type GameParts struct {
	var tiles []Tile
	var id int

	var walkable bool
	var harf bool
	var waterType WaterType
}

var gamePartsDict []GameParts

