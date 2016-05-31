package model

import (
	"errors"
)

type MetalColor struct {
	Id        int
	Name      string
	ShortName string
}

type Metal struct {
	Id        int
	Name      string
	ShortName string
	Karat     string
}

const NUM_METALS = 4
const NUM_METAL_COLORS = 3
const NUM_METAL_VARIANTS = 10 //3 types of gold * 3 gold colors + silver
const METAL_SHORTNAME_LEN = 3
const (
	METAL_WHITE = iota
	METAL_YELLOW
	METAL_ROSE
)
const (
	GOLD_22KT = iota
	GOLD_18KT
	GOLD_14KT
	SILVER
)

var MetalColorList []MetalColor = []MetalColor{
	METAL_WHITE:  MetalColor{Id: METAL_WHITE, Name: "White", ShortName: "W"},
	METAL_YELLOW: MetalColor{Id: METAL_YELLOW, Name: "Yellow", ShortName: "Y"},
	METAL_ROSE:   MetalColor{Id: METAL_ROSE, Name: "Rose", ShortName: "P"},
}
var MetalColorReverseMap map[string]int = map[string]int{
	"W": METAL_WHITE,
	"Y": METAL_YELLOW,
	"P": METAL_ROSE,
}
var MetalList []Metal = []Metal{
	GOLD_22KT: Metal{Id: GOLD_22KT, Name: "Gold", ShortName: "G22", Karat: "22K"},
	GOLD_18KT: Metal{Id: GOLD_18KT, Name: "Gold", ShortName: "G18", Karat: "18K"},
	GOLD_14KT: Metal{Id: GOLD_14KT, Name: "Gold", ShortName: "G14", Karat: "14K"},
	SILVER:    Metal{Id: SILVER, Name: "Silver", ShortName: "SLV", Karat: ""},
}
var MetalReverseMap map[string]int = map[string]int{
	"G22": GOLD_22KT,
	"G18": GOLD_18KT,
	"G14": GOLD_14KT,
	"SLV": SILVER,
}

//TODO: Fetch the gold price dynamically everyday
var TODAYS_GOLD_PRICE float64 = 3000
var TODAYS_SILVER_PRICE float64 = 50

func IsValidMetal(metalId int) bool {
	if metalId >= 0 && metalId < NUM_METALS {
		return true
	}
	return false
}

func GetMetalList() []Metal {
	return MetalList
}

func GetMetal(metalId int) Metal {
	return MetalList[metalId]
}

//Weight is in grams
func GetMetalPrice(metalId int, weight float64) (float64, error) {
	var priceMultiplier, price float64
	if weight <= 0.0 {
		return price, errors.New("Invalid weight value passed for metal")
	}

	if metalId == GOLD_22KT {
		priceMultiplier = 0.935 * TODAYS_GOLD_PRICE
	} else if metalId == GOLD_18KT {
		priceMultiplier = 0.76 * TODAYS_GOLD_PRICE
	} else if metalId == GOLD_14KT {
		priceMultiplier = 0.6 * TODAYS_GOLD_PRICE
	} else if metalId == SILVER {
		priceMultiplier = TODAYS_SILVER_PRICE
	} else {
		return price, errors.New("Pricing data not found for metal")
	}

	price = weight * priceMultiplier
	if price <= 0 {
		return price, errors.New("Metal price coming out to be <= 0")
	}

	return price, nil
}
